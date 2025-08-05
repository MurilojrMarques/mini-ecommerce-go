package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MurilojrMarques/mini-ecommerce-go/types"
	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("falhar ao registrar usuário se o payload for inválido", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName:  "test",
			Email:     "asd",
			Password:  "asd",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Esperado status %d, mas recebeu %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Registrar usuário com sucesso", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName:  "test",
			Email:     "email@email.com",
			Password:  "asd",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Esperado status %d, mas recebeu %d", http.StatusCreated, rr.Code)
		}
	})
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("Usuário não encontrado")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(types.User) error {
	return nil
}
