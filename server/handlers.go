package server

import (
	"context"
	"encoding/json"
	"github.com/eine-doodka/twoStage/customerrors"
	"github.com/eine-doodka/twoStage/logic"
	"github.com/google/uuid"
	"net/http"
)

type Handlers struct {
	logic logic.Logic
}

func NewHandlers(logicImpl logic.Logic) *Handlers {
	return &Handlers{
		logic: logicImpl,
	}
}

func (h *Handlers) handleInit() http.Handler {
	type response struct {
		OpID uuid.UUID `json:"operation_id"`
		Code string    `json:"code"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		opid, randn, err := h.logic.HandleInit(ctx)
		if err != nil {
			h.error(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{
			OpID: opid,
			Code: randn,
		})
	})
}
func (h *Handlers) handleCommit() http.Handler {
	type request struct {
		OperationId string `json:"id"`
		UserCode    string `json:"code"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		rq := &request{}
		err := json.NewDecoder(r.Body).Decode(rq)
		if err != nil {
			h.error(w, err, http.StatusBadRequest)
			return
		}
		id, err := uuid.Parse(rq.OperationId)
		if err != nil {
			h.error(w, err, http.StatusBadRequest)
			return
		}

		err = h.logic.HandleCommit(ctx, id, rq.UserCode)

		if err == customerrors.ErrNotFound {
			h.error(w, customerrors.ErrNotFound, http.StatusNotFound)
			return
		} else if err == customerrors.ErrNotMatch {
			h.error(w, customerrors.ErrNotMatch, http.StatusExpectationFailed)
		} else if err != nil {
			h.error(w, err, http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusFound)
			return
		}
	})
}

func (h *Handlers) error(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
	return
}
