package user

import (
	"awesomeProject4/internal/apper"
	"awesomeProject4/internal/handlers"
	"awesomeProject4/pkg/logging"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	usersURL = "/users"
	userURL  = "/user/:uuid"
)

type handler struct {
	logger *logging.Logger //Что бы в handler использовать Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

func (h handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apper.MiddleWare(h.GetList))              //Все пользователи
	router.HandlerFunc(http.MethodPost, usersURL, apper.MiddleWare(h.CreateUsers))         //Создать пользователя
	router.HandlerFunc(http.MethodGet, userURL, apper.MiddleWare(h.GetUserByUUID))         //Найти пользователя по ID
	router.HandlerFunc(http.MethodPut, userURL, apper.MiddleWare(h.UpdateUser))            //Обновить пользователя
	router.HandlerFunc(http.MethodPatch, userURL, apper.MiddleWare(h.PartiallyUpdateUser)) //Частично обновить пользователя
	router.HandlerFunc(http.MethodDelete, userURL, apper.MiddleWare(h.DeleteUser))         //Удалить пользователя
}
func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	return apper.ErrNotFound
}

func (h *handler) CreateUsers(w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf("this is API error")
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {
	return apper.NewAppError(nil, "test", "test", "t13")
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("this is list of users"))
	w.WriteHeader(204)

	return nil
}

func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("this is list of users"))
	w.WriteHeader(200)

	return nil
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("this is list of users"))
	w.WriteHeader(204)

	return nil
}
