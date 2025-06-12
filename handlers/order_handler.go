package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/suraj/nitabuddy/models"
	"github.com/suraj/nitabuddy/utils"
)

type OrderHandler struct {
	orderModel *models.OrderModel
}

func NewOrderHandler(orderModel *models.OrderModel) *OrderHandler {
	return &OrderHandler{
		orderModel: orderModel,
	}
}

func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {

	// Check for valid user
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Unauthorized: " + err.Error(),
			"status":  false,
		})
		return
	}

	var input struct {
		Store        string `json:"store"`
		OrderDetails string `json:"order_details"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
			"status":  false,
		})
		return
	}

	_, err = h.orderModel.CreateOrder(input.Store, input.OrderDetails, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
			"status":  false,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Order Placed Successfully",
		"status":  true,
	})

}

func (h *OrderHandler) FetchMyOrders(w http.ResponseWriter, r *http.Request) {

	userID, err := utils.ExtractUserIDFromToken(r)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized " + err.Error(),
			"orders":  nil,
		})
		return
	}

	// Fetch orders from DB
	orders, err := h.orderModel.GetOrdersByUserID(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Failed to fetch Orders " + err.Error(),
			"orders":  nil,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "orders fetched",
		"orders":  orders,
	})
}
