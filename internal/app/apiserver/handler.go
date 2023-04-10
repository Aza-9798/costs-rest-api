package apiserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Aza-9798/costs-rest-api/internal/app/model"
	"github.com/Aza-9798/costs-rest-api/internal/app/store"

	"github.com/gorilla/mux"
)

func (s *server) handleUserCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleAccountGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, _ := strconv.Atoi(mux.Vars(r)["id"])
		a, err := s.store.Account().Find(accountId)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, a)
	}
}
func (s *server) handleAccountGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		res, err := s.store.Account().GetAllByUser(u.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		} else {
			s.respond(w, r, http.StatusOK, res)
		}
	}
}

func (s *server) handleAccountCreate() http.HandlerFunc {
	type request struct {
		Name        string  `json:"name"`
		Type        string  `json:"type"`
		Description string  `json:"description"`
		Balance     float64 `json:"balance"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := r.Context().Value(ctxKeyUser).(*model.User)
		acc := &model.Account{
			Name:        req.Name,
			Type:        req.Type,
			Description: req.Description,
			Balance:     req.Balance,
			User:        u.ID,
		}
		if err := s.store.Account().Create(acc); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, acc)
	}
}

func (s *server) handleAccountDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, _ := strconv.Atoi(mux.Vars(r)["id"])
		err := s.store.Account().Delete(accountId)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleAccountUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a := &model.Account{}
		if err := json.NewDecoder(r.Body).Decode(a); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}
		accountId, _ := strconv.Atoi(mux.Vars(r)["id"])
		a.ID = accountId
		err := s.store.Account().Save(a)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleTransactionCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &model.TransactionJSON{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		SourceAcc, err := s.store.Account().Find(req.Source)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		DestinationAcc, err := s.store.Account().Find(req.Destination)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		t := &model.TransactionDB{
			TransactionDate: req.TransactionDate,
			Source:          SourceAcc,
			Destination:     DestinationAcc,
			Type:            req.Type,
			Amount:          req.Amount,
			Description:     req.Description,
		}

		if err := s.store.Transaction().Create(t); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusCreated, t.ToJSON())
	}
}

func (s *server) handleTransactionGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		t, err := s.store.Transaction().Find(id)
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, t)
	}
}

func (s *server) handleTransactionGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID, _ := strconv.Atoi(mux.Vars(r)["accountID"])
		u := r.Context().Value(ctxKeyUser).(*model.User)
		if ok, err := s.store.Account().IsAccountBelongsUser(accountID, u.ID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		} else if !ok {
			s.error(w, r, http.StatusUnauthorized, nil)
			return
		}
		tr, err := s.store.Transaction().GetAllByAccount(accountID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, tr)
	}
}

func (s *server) handleTransactionDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		t, err := s.store.Transaction().Find(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		if err := s.store.Transaction().Delete(t); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleTransactionUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		t := &model.TransactionJSON{}
		if err := json.NewDecoder(r.Body).Decode(t); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}
		source, err := s.store.Account().Find(t.Source)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		destination, err := s.store.Account().Find(t.Destination)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		tDB := &model.TransactionDB{
			ID:              id,
			TransactionDate: t.TransactionDate,
			Source:          source,
			Destination:     destination,
			Amount:          t.Amount,
			Type:            t.Type,
			Description:     t.Description,
		}
		if err := s.store.Transaction().Save(tDB); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusAccepted, tDB)
	}
}

func (s *server) handleSummaryGet() http.HandlerFunc {
	type request struct {
		DateStart time.Time `json:"date_start"`
		DateEnd   time.Time `json:"date_end"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}
		u := r.Context().Value(ctxKeyUser).(*model.User)
		res, err := s.store.Transaction().GetSummary(u.ID, req.DateStart, req.DateEnd)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusOK, res)
	}
}

func (s *server) handleSummaryAccountGet() http.HandlerFunc {
	type request struct {
		DateStart time.Time `json:"date_start"`
		DateEnd   time.Time `json:"date_end"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}
		accountID, _ := strconv.Atoi(mux.Vars(r)["accountID"])
		acc, err := s.store.Account().Find(accountID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		res, err := model.GetSummaryByAccount(acc)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		if err := res.SetPeriod(req.DateStart, req.DateEnd); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		transactions, err := s.store.Transaction().GetAllByAccountAndPeriod(accountID, req.DateStart, req.DateEnd)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		res.CalculateSummary(transactions)
		s.respond(w, r, http.StatusOK, res)
	}
}
