package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	echov4 "github.com/labstack/echo/v4"

	authdomain "github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/domain"
	authoutbound "github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound"
	authmock "github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound/mock"
	userdomain "github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
	usermock "github.com/Genedevelop/7solutions-backend-challenge/internal/user/port/outbound/mock"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/logger"
)

type sharedUsers struct {
	repo   *usermock.UserRepository
	hashes map[string]string
}

func (s *sharedUsers) Create(ctx context.Context, a *authoutbound.Account) error {
	u := &entity.User{Name: a.Name, Email: a.Email, CreatedAt: a.CreatedAt}
	if err := s.repo.Create(ctx, u); err != nil {
		return err
	}
	a.ID = u.ID
	s.hashes[u.ID] = a.PasswordHash
	return nil
}

func (s *sharedUsers) GetByEmail(ctx context.Context, email string) (*authoutbound.Credential, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Email == email {
			return &authoutbound.Credential{UserID: u.ID, PasswordHash: s.hashes[u.ID]}, nil
		}
	}
	return nil, errs.ErrCredentialNotFound
}

type prefixTokens struct{}

func (prefixTokens) Parse(token string) (string, error) {
	userID, ok := strings.CutPrefix(token, "token:")
	if !ok || userID == "" {
		return "", errors.New("invalid token")
	}
	return userID, nil
}

type testAPI struct {
	t    *testing.T
	srv  http.Handler
	repo *usermock.UserRepository
	logs *bytes.Buffer
}

func newTestAPI(t *testing.T) *testAPI {
	t.Helper()
	repo := usermock.NewUserRepository()
	users := userdomain.NewUserService(repo)
	shared := &sharedUsers{repo: repo, hashes: map[string]string{}}
	auth, err := authdomain.NewAuthService(shared, shared, authmock.PlainHasher{}, authmock.Tokens{})
	if err != nil {
		t.Fatal(err)
	}
	logs := &bytes.Buffer{}
	logger.Init(logs, "info")
	return &testAPI{t: t, srv: newServer(users, auth, prefixTokens{}), repo: repo, logs: logs}
}

type userResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type errorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

func (a *testAPI) do(method, path, token, body string) *httptest.ResponseRecorder {
	a.t.Helper()
	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, r)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	a.srv.ServeHTTP(rec, req)
	return rec
}

func (a *testAPI) register(email string) userResponse {
	a.t.Helper()
	rec := a.do(http.MethodPost, "/auth/register", "", `{"name":"Test","email":"`+email+`","password":"password123"}`)
	if rec.Code != http.StatusCreated {
		a.t.Fatalf("register %s: %d %s", email, rec.Code, rec.Body)
	}
	var u userResponse
	decode(a.t, rec, &u)
	return u
}

func decode(t *testing.T, rec *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), v); err != nil {
		t.Fatalf("decode %q: %v", rec.Body, err)
	}
}

func expectStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("status = %d, want %d, body %s", rec.Code, want, rec.Body)
	}
}

func TestRegisterAndLogin(t *testing.T) {
	api := newTestAPI(t)

	u := api.register("Alice@Example.com")
	if u.ID == "" || u.Email != "alice@example.com" {
		t.Errorf("got %+v", u)
	}

	rec := api.do(http.MethodPost, "/auth/register", "", `{"name":"x","email":"alice@example.com","password":"password123"}`)
	expectStatus(t, rec, http.StatusConflict)

	rec = api.do(http.MethodPost, "/auth/register", "", `{"name":"","email":"bad","password":"1"}`)
	expectStatus(t, rec, http.StatusBadRequest)
	var verr errorResponse
	decode(t, rec, &verr)
	if len(verr.Fields) != 3 {
		t.Errorf("fields = %v", verr.Fields)
	}

	expectStatus(t, api.do(http.MethodPost, "/auth/register", "", `{not json`), http.StatusBadRequest)

	rec = api.do(http.MethodPost, "/auth/login", "", `{"email":"alice@example.com","password":"password123"}`)
	expectStatus(t, rec, http.StatusOK)
	var login loginResponse
	decode(t, rec, &login)
	if login.AccessToken != "token:"+u.ID || login.TokenType != "Bearer" {
		t.Errorf("login = %+v", login)
	}

	expectStatus(t, api.do(http.MethodPost, "/auth/login", "", `{"email":"alice@example.com","password":"wrong-pass"}`), http.StatusUnauthorized)
}

func TestPasswordNeverReturned(t *testing.T) {
	api := newTestAPI(t)
	u := api.register("a@example.com")
	rec := api.do(http.MethodGet, "/users/"+u.ID, "token:"+u.ID, "")
	expectStatus(t, rec, http.StatusOK)
	if strings.Contains(strings.ToLower(rec.Body.String()), "password") {
		t.Errorf("response leaks password: %s", rec.Body)
	}
}

func TestUsersRequireToken(t *testing.T) {
	api := newTestAPI(t)
	for _, header := range []string{"", "garbage"} {
		rec := api.do(http.MethodGet, "/users", header, "")
		expectStatus(t, rec, http.StatusUnauthorized)
		if rec.Header().Get("WWW-Authenticate") == "" {
			t.Error("missing WWW-Authenticate header")
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Basic abc")
	rec := httptest.NewRecorder()
	api.srv.ServeHTTP(rec, req)
	expectStatus(t, rec, http.StatusUnauthorized)
}

func TestUserCRUD(t *testing.T) {
	api := newTestAPI(t)
	a := api.register("a@example.com")
	token := "token:" + a.ID

	b := api.register("bob@example.com")

	expectStatus(t, api.do(http.MethodPost, "/users", token, `{"name":"x","email":"x@example.com","password":"password123"}`), http.StatusMethodNotAllowed)

	rec := api.do(http.MethodGet, "/users", token, "")
	expectStatus(t, rec, http.StatusOK)
	var list []userResponse
	decode(t, rec, &list)
	if len(list) != 2 {
		t.Errorf("list = %d users", len(list))
	}

	expectStatus(t, api.do(http.MethodGet, "/users/"+b.ID, token, ""), http.StatusOK)
	expectStatus(t, api.do(http.MethodGet, "/users/nope", token, ""), http.StatusNotFound)

	rec = api.do(http.MethodPatch, "/users/"+a.ID, token, `{"name":"Alice"}`)
	expectStatus(t, rec, http.StatusOK)
	var updated userResponse
	decode(t, rec, &updated)
	if updated.Name != "Alice" || updated.Email != "a@example.com" {
		t.Errorf("updated = %+v", updated)
	}
	expectStatus(t, api.do(http.MethodPatch, "/users/"+a.ID, token, `{"email":"bob@example.com"}`), http.StatusConflict)
	expectStatus(t, api.do(http.MethodPatch, "/users/"+a.ID, token, `{}`), http.StatusBadRequest)
	expectStatus(t, api.do(http.MethodPatch, "/users/"+b.ID, token, `{"name":"x"}`), http.StatusForbidden)

	expectStatus(t, api.do(http.MethodDelete, "/users/"+b.ID, token, ""), http.StatusForbidden)
	expectStatus(t, api.do(http.MethodDelete, "/users/"+a.ID, token, ""), http.StatusNoContent)
	expectStatus(t, api.do(http.MethodGet, "/users/"+a.ID, "token:"+b.ID, ""), http.StatusNotFound)
}

func TestDatabaseErrorIsHidden(t *testing.T) {
	api := newTestAPI(t)
	a := api.register("a@example.com")
	api.repo.Err = errors.New("mongo: connection refused to 10.0.0.5")

	rec := api.do(http.MethodGet, "/users", "token:"+a.ID, "")
	expectStatus(t, rec, http.StatusInternalServerError)
	if strings.Contains(rec.Body.String(), "10.0.0.5") {
		t.Errorf("internal error leaked: %s", rec.Body)
	}
	for _, want := range []string{`"level":"ERROR"`, "10.0.0.5", `"error_code":"INTERNAL"`} {
		if !strings.Contains(api.logs.String(), want) {
			t.Errorf("error log missing %q:\n%s", want, api.logs)
		}
	}
	var body errorResponse
	decode(t, rec, &body)
	if body.Code != "INTERNAL" || body.Message != "internal server error" {
		t.Errorf("body = %+v", body)
	}
}

func TestRequestLogger(t *testing.T) {
	api := newTestAPI(t)
	api.do(http.MethodGet, "/healthz", "", "")
	api.do(http.MethodGet, "/users", "", "")

	logs := api.logs.String()
	for _, want := range []string{`"method":"GET"`, `"path":"/healthz"`, `"status":200`, `"path":"/users"`, `"status":401`, `"duration"`, `"request_id"`} {
		if !strings.Contains(logs, want) {
			t.Errorf("log missing %q:\n%s", want, logs)
		}
	}
}

func TestStrictJSONBody(t *testing.T) {
	api := newTestAPI(t)
	a := api.register("a@example.com")
	for _, body := range []string{`{"name":"Y"} garbage`, `{"name":"Y"}{"name":"Z"}`, ``} {
		rec := api.do(http.MethodPatch, "/users/"+a.ID, "token:"+a.ID, body)
		expectStatus(t, rec, http.StatusBadRequest)
	}
}

func TestOwnIDInAnyLetterCase(t *testing.T) {
	api := newTestAPI(t)
	a := api.register("a@example.com")
	rec := api.do(http.MethodPatch, "/users/"+strings.ToUpper(a.ID), "token:"+a.ID, `{"name":"Upper"}`)
	expectStatus(t, rec, http.StatusOK)
}

func TestServerTimeouts(t *testing.T) {
	srv := newTestAPI(t).srv.(*echov4.Echo).Server
	if srv.ReadHeaderTimeout == 0 || srv.ReadTimeout == 0 || srv.WriteTimeout == 0 || srv.IdleTimeout == 0 {
		t.Errorf("server timeouts not set: %+v", srv)
	}
}
