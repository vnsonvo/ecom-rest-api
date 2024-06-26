package products

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/vnsonvo/ecom-rest-api/services/auth"
	"github.com/vnsonvo/ecom-rest-api/types"
	"github.com/vnsonvo/ecom-rest-api/utils"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, prefixPath string) {
	mux.HandleFunc(fmt.Sprintf("GET %s/products", prefixPath), h.handleGetProducts)
	mux.HandleFunc(fmt.Sprintf("POST %s/products", prefixPath), auth.JWTAuthMiddleware(h.handlerCreateProduct, h.userStore))
	mux.HandleFunc(fmt.Sprintf("GET %s/products/{productId}", prefixPath), h.handleGetProduct)

}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handlerCreateProduct(w http.ResponseWriter, r *http.Request) {
	var product types.CreateProductPayload
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload: %v", errors))
		return
	}

	if err := h.store.CreateProduct(product); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, product)
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	productIdStr := r.PathValue("productId")
	if productIdStr == "" {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Missing product ID"))
		return
	}

	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Invalid product ID"))
		return
	}

	product, err := h.store.GetProductByID(productId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, product)
}
