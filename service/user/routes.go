package user

import (
	"fmt"
	"net/http"

	"github.com/MurilojrMarques/mini-ecommerce-go/config"
	"github.com/MurilojrMarques/mini-ecommerce-go/service/auth"
	"github.com/MurilojrMarques/mini-ecommerce-go/types"
	"github.com/MurilojrMarques/mini-ecommerce-go/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriterError(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriterError(w, http.StatusUnauthorized, fmt.Errorf("Email ou senha inválidos"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriterError(w, http.StatusUnauthorized, fmt.Errorf("Email ou senha inválidos"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriterError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriterError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriterError(w, http.StatusBadRequest, fmt.Errorf("Erro de validação: %s", error.Error()))
		return
	}

	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriterError(w, http.StatusConflict, fmt.Errorf("Usuário do email %s já existe", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriterError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriterError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
