package main

import (
	"net/http"
	"time"

	echov4 "github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	authecho "github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/adapter/inbound/echo"
	authinbound "github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/inbound"
	userecho "github.com/Genedevelop/7solutions-backend-challenge/internal/user/adapter/inbound/echo"
	userinbound "github.com/Genedevelop/7solutions-backend-challenge/internal/user/port/inbound"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/middleware"
)

func newServer(users userinbound.IUserUseCase, auth authinbound.IAuthUseCase, tokens middleware.TokenParser) *echov4.Echo {
	// Step 1: create Echo with the central error handler and server timeouts so slow clients cannot hold connections open
	e := echov4.New()
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = middleware.ErrorHandler()
	e.Server.ReadHeaderTimeout = 5 * time.Second
	e.Server.ReadTimeout = 15 * time.Second
	e.Server.WriteTimeout = 15 * time.Second
	e.Server.IdleTimeout = 60 * time.Second

	// Step 2: request id first so every log line of a request shares it, then logging, panic recovery and a body limit
	e.Use(echomw.RequestID())
	e.Use(middleware.RequestLogger())
	e.Use(echomw.Recover())
	e.Use(echomw.BodyLimit("1M"))

	// Step 3: liveness probe for docker and load balancers
	e.GET("/healthz", func(c echov4.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Step 4: each bounded context registers its own routes, a new context only adds one line here
	authecho.NewAuthHandler(auth).Register(e)
	userecho.NewUserHandler(users).Register(e, middleware.JWTAuth(tokens))

	return e
}
