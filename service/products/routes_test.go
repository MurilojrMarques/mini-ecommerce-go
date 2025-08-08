package product

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MurilojrMarques/mini-ecommerce-go/types"
	"github.com/gorilla/mux"
)

func TestProductServiceHandlers(t *testing.T) {
	productStore := &mockProductStore{}
	userStore := &mockUserStore{}
	handler := NewHandler(productStore, userStore)

	t.Run("Receber os protudos", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products", handler.handleGetProducts).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Esperado status %d, mas recebeu %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("Falhar se o ID do produto não for número", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/products/abc", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products/{productID}", handler.handleGetProduct).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Esperado status %d, mas recebeu %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("Manipular o recebimento do produto por id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/products/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products/{productID}", handler.handleGetProduct).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Esperado status %d, mas recebeu %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("Falhar se na criação de um produto a payload estiver incorreto", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products", handler.handleCreateProduct).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Esperado status %d, mas recebeu %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("Realizar a criação do produto", func(t *testing.T) {
		payload := types.CreateProductPayload{
			Name:        "test",
			Price:       100,
			Image:       "test.jpg",
			Description: "test description",
			Quantity:    10,
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products", handler.handleCreateProduct).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Esperado status %d, mas recebeu %d", http.StatusCreated, rr.Code)
		}
	})
}

type mockProductStore struct{}

func (m *mockProductStore) GetProductByID(productID int) (*types.Product, error) {
	return &types.Product{}, nil
}

func (m *mockProductStore) GetProducts() ([]*types.Product, error) {
	return []*types.Product{}, nil
}

func (m *mockProductStore) CreateProduct(product types.CreateProductPayload) error {
	return nil
}

func (m *mockProductStore) UpdateProduct(product types.Product) error {
	return nil
}

func (m *mockProductStore) GetProductsByID(ids []int) ([]types.Product, error) {
	return []types.Product{}, nil
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByID(userID int) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) CreateUser(user types.User) error {
	return nil
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return &types.User{}, nil
}
