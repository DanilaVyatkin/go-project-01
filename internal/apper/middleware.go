package apper

import (
	"github.com/pkg/errors"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func MiddleWare(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppErr
		err := h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				w.Header().Set("Content-Type", "application/json")
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				}
				err = err.(*AppErr)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(appErr.Marshal())
				return
			}

			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).Marshal())
		}
	}
}
