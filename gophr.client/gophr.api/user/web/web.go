package web

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jayvib/golog"
	"gophr/v2/gophr.api/user"
	"net/http"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Method  string      `json:"method,omitempty"`
}

type Parameters struct {
	UserService user.Service `validation:"required"`
}

func (p *Parameters) Validate() error {
	validation := validator.New()
	return validation.Struct(p)
}

func New(param *Parameters) *Handler {
	if err := param.Validate(); err != nil {
		panic(err)
	}

	return &Handler{
		svc: param.UserService,
	}
}

type Handler struct {
	svc user.Service
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	usr, _ := h.svc.GetByID(r.Context(), userId)
	h.writeResponse(w, r, "OK", usr, http.StatusOK)
}

func (h *Handler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	golog.Debug("Email:", v["email"])
	email := v["email"]
	res, _ := h.svc.GetByEmail(r.Context(), email) // TODO: Handle the error
	h.writeResponse(w, r, "OK", res, http.StatusOK)
}

func (h *Handler) writeResponse(w http.ResponseWriter, r *http.Request, message string, data interface{}, code int) {
	resp := &Response{
		Data:    data,
		Message: message,
		Method:  r.Method,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		golog.Errorf("error while encoding json to the http writer: %v", err.Error())
		return
	}
	w.WriteHeader(code)
}

func (h *Handler) GetByUsername(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) Save(w http.ResponseWriter, r *http.Request)          {}
func (h *Handler) Login(w http.ResponseWriter, r *http.Request)         {}
