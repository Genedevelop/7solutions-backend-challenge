package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	authmongo "github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/adapter/outbound/mongo"
	authdomain "github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/domain"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/adapter/inbound/job"
	usermongo "github.com/Genedevelop/7solutions-backend-challenge/internal/user/adapter/outbound/mongo"
	userdomain "github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/config"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/crypto"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/jwt"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/logger"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/mongodb"
)

func main() {
	// Step 1: load config and check the JWT secret before any connection opens, so a failure here has nothing to clean up
	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", err)
		os.Exit(1)
	}
	logger.Init(os.Stdout, cfg.LogLevel)
	tokens, err := jwt.NewJWTManager(cfg.JWTSecret, cfg.JWTTTL, cfg.JWTIssuer)
	if err != nil {
		logger.Error("jwt setup", err)
		os.Exit(1)
	}

	// Step 2: ctx is cancelled on Ctrl+C or SIGTERM (docker stop)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// Step 3: connect to Mongo and make sure the unique email index exists
	connectCtx, cancelConnect := context.WithTimeout(ctx, 10*time.Second)
	client, err := mongodb.Connect(connectCtx, cfg.MongoURI)
	if err != nil {
		cancelConnect()
		logger.Error("mongo connect", err)
		os.Exit(1)
	}
	db := client.Database(cfg.MongoDB)
	err = mongodb.EnsureIndexes(connectCtx, db)
	cancelConnect()
	if err != nil {
		_ = client.Disconnect(context.Background())
		logger.Error("mongo indexes", err)
		os.Exit(1)
	}

	// Step 4: wire adapters into each context's domain service
	hasher := crypto.NewBcryptHasher(0)
	userService := userdomain.NewUserService(usermongo.NewUserRepository(db))
	authService, err := authdomain.NewAuthService(authmongo.NewCredentialReader(db), authmongo.NewCredentialWriter(db), hasher, tokens)
	if err != nil {
		_ = client.Disconnect(context.Background())
		logger.Error("auth service", err)
		os.Exit(1)
	}

	// Step 5: start the background user counter
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		job.RunUserCounter(ctx, userService, cfg.UserCountInterval)
	}()

	// Step 6: start HTTP, a listen error ends the process
	errCh := make(chan error, 1)
	httpServer := newServer(userService, authService, tokens)
	go func() {
		logger.Info("http server listening", "addr", cfg.HTTPAddr)
		if err := httpServer.Start(cfg.HTTPAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// Step 7: wait for a signal or a server failure
	var runErr error
	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case runErr = <-errCh:
		logger.Error("http server", runErr)
	}
	stop()

	// Step 8: stop taking new requests and let in-flight ones finish within the timeout
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("http shutdown", err)
	}

	cancelShutdown()

	// Step 9: wait for the counter goroutine, then close Mongo last with its own budget since Shutdown may have used all of its time
	wg.Wait()
	disconnectCtx, cancelDisconnect := context.WithTimeout(context.Background(), 5*time.Second)
	if err := client.Disconnect(disconnectCtx); err != nil {
		logger.Error("mongo disconnect", err)
	}
	cancelDisconnect()

	logger.Info("server stopped")
	if runErr != nil {
		os.Exit(1)
	}
}
