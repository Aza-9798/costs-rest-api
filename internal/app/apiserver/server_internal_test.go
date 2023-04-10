package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aza-9798/costs-rest-api/internal/app/model"
	"github.com/Aza-9798/costs-rest-api/internal/app/store/teststore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

var (
	logLevel = "debug"
)

func TestServer_AuthenticateUser(t *testing.T) {
	st := teststore.New()
	u := model.TestUser(t)
	st.User().Create(u)
	testCases := []struct {
		name         string
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
	}
	secretKey := []byte("secret")
	s, err := newServer(st, sessions.NewCookieStore(secretKey), logLevel)
	if err != nil {
		t.Fatal(err)
	}
	sc := securecookie.New(secretKey, nil)
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode(sessionName, tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.authenticateUser(fakeHandler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleUserCreate(t *testing.T) {
	s, err := newServer(teststore.New(), sessions.NewCookieStore([]byte("secret")), logLevel)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "valid-email@example.org",
				"password": "validpassword",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/user", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, rec.Code, tc.expectedCode)
		})
	}
}

func TestServer_HandleSessionCreate(t *testing.T) {
	st := teststore.New()
	u := model.TestUser(t)
	st.User().Create(u)
	s, err := newServer(st, sessions.NewCookieStore([]byte("secret")), logLevel)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    u.Email,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "user not exists",
			payload: map[string]string{
				"email":    "not-found@example.org",
				"password": u.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    u.Email,
				"password": u.Password + "123",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/session", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, rec.Code, tc.expectedCode)
		})
	}
}

func TestServer_HandelAccountCreate(t *testing.T) {
	st := teststore.New()
	u := model.TestUser(t)
	st.User().Create(u)
	svr, err := newServer(st, sessions.NewCookieStore([]byte("secret")), logLevel)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"name":        "test",
				"type":        model.CurrentAccount,
				"description": "test",
				"balance":     0,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "missing name",
			payload: map[string]interface{}{
				"name":        "",
				"type":        model.CurrentAccount,
				"description": "test",
				"balance":     0,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "wrong type",
			payload: map[string]interface{}{
				"name":        "test",
				"type":        "invalid",
				"description": "test",
				"balance":     0,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//Нужно создать сессию иначе нельзя будет создать счет
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(map[string]string{
				"email":    u.Email,
				"password": u.Password,
			})
			req, _ := http.NewRequest(http.MethodPost, "/session", b)
			svr.ServeHTTP(rec, req)
			assert.Equal(t, rec.Code, http.StatusOK)

			cookie := rec.Header().Get("Set-Cookie")
			rec = httptest.NewRecorder()
			b = &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ = http.NewRequest(http.MethodPost, "/private/account", b)
			req.Header.Set("Cookie", cookie)
			svr.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
