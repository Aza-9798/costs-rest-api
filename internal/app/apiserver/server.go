package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Aza-9798/costs-rest-api/internal/app/store"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "somesessionname"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

type ctxKey int8

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store, logLevel string) (*server, error) {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetLevel(lvl)
	s := &server{
		router:       mux.NewRouter(),
		logger:       logger,
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/user", s.handleUserCreate()).Methods("POST")
	s.router.HandleFunc("/session", s.handleSessionCreate()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)

	private.HandleFunc("/summary", s.handleSummaryGet()).Methods("POST")
	//счета
	private.HandleFunc("/account", s.handleAccountCreate()).Methods("POST")
	private.HandleFunc("/account/{id:[0-9]+}", s.handleAccountGet()).Methods("GET")
	private.HandleFunc("/account/{id:[0-9]+}", s.handleAccountDelete()).Methods("DELETE")
	private.HandleFunc("/account/{id:[0-9]+}", s.handleAccountUpdate()).Methods("PUT")
	private.HandleFunc("/account/all", s.handleAccountGetAll()).Methods("GET")
	private.HandleFunc("/account/{accountID:[0-9]+}/all_transactions", s.handleTransactionGetAll()).Methods("GET")
	private.HandleFunc("/account/{accountID:[0-9]+}/summary", s.handleSummaryAccountGet()).Methods("POST")
	//переводы
	private.HandleFunc("/transaction", s.handleTransactionCreate()).Methods("POST")
	private.HandleFunc("/transaction/{id:[0-9]+}", s.handleTransactionGet()).Methods("GET")
	private.HandleFunc("/transaction/{id:[0-9]+}", s.handleTransactionDelete()).Methods("DELETE")
	private.HandleFunc("/transaction/{id:[0-9]+}", s.handleTransactionUpdate()).Methods("PUT")
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
