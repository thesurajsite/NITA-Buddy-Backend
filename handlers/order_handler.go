package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suraj/nitabuddy/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHandler struct {
	orderModel  *models.OrderModel
	authHandler *AuthHandler // Add this field
}

func NewOrderHandler(orderModel *models.OrderModel, authHandler *AuthHandler) *OrderHandler {
	return &OrderHandler{
		orderModel:  orderModel,
		authHandler: authHandler, // Initialize it
	}
}

func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {

	// Check for valid user - USE authHandler instead of utils
	userID, err := h.authHandler.GetUserIDFromToken(r)
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
		"message": "Request created Successfully",
		"status":  true,
	})

}

func (h *OrderHandler) FetchOtherOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authHandler.GetUserIDFromToken(r)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized " + err.Error(),
			"orders":  []interface{}{},
		})
		return
	}

	// Fetch orders from DB
	orders, err := h.orderModel.GetOtherIncompleteOrders(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Failed to fetch Requests " + err.Error(),
			"orders":  []interface{}{},
		})
		return
	}

	if orders == nil {
		orders = []models.Order{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Requests fetched",
		"orders":  orders,
	})

}

func (h *OrderHandler) FetchMyOrders(w http.ResponseWriter, r *http.Request) {

	userID, err := h.authHandler.GetUserIDFromToken(r)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized " + err.Error(),
			"orders":  []interface{}{},
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
			"message": "Failed to fetch Requests " + err.Error(),
			"orders":  []interface{}{},
		})
		return
	}

	if orders == nil {
		orders = []models.Order{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Requests fetched",
		"orders":  orders,
	})
}

func (h *OrderHandler) CancelMyOrder(w http.ResponseWriter, r *http.Request) {

	userID, err := h.authHandler.GetUserIDFromToken(r)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized: " + err.Error(),
		})
		return
	}

	vars := mux.Vars(r)
	orderIDstr := vars["id"]
	orderID, err := primitive.ObjectIDFromHex(orderIDstr)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Invalid Request ID",
		})
		return
	}

	err = h.orderModel.CancelOrder(userID, orderID)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Could not cancel Request: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Request cancelled and Deleted",
	})
}

func (h *OrderHandler) AcceptOrder(w http.ResponseWriter, r *http.Request) {

	userID, err := h.authHandler.GetUserIDFromToken(r)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized: " + err.Error(),
		})
		return
	}

	vars := mux.Vars(r)
	orderIDstr := vars["id"]
	orderId, err := primitive.ObjectIDFromHex(orderIDstr)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Invalid Request ID",
		})
		return
	}

	err = h.orderModel.AcceptOrder(userID, orderId)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "can't Accept Request: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Request Accepted",
	})
}

func (h *OrderHandler) FetchAcceptedOrders(w http.ResponseWriter, r *http.Request) {

	userID, err := h.authHandler.GetUserIDFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized" + err.Error(),
			"orders":  []interface{}{},
		})
		return
	}

	orders, err := h.orderModel.GetAcceptedOrders(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Failed to fetch Requests " + err.Error(),
			"orders":  []interface{}{},
		})
		return
	}

	if orders == nil {
		orders = []models.Order{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Requests fetched",
		"orders":  orders,
	})
}

func (h *OrderHandler) CompleteOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := h.authHandler.GetUserIDFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized: " + err.Error(),
		})
		return
	}

	var input struct {
		OrderID string `json:"order_id"`
		OTP     string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Invalid input: " + err.Error(),
		})
		return
	}

	orderObjectID, err := primitive.ObjectIDFromHex(input.OrderID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Invalid Order ID",
		})
		return
	}

	err = h.orderModel.CompleteOrder(userID, orderObjectID, input.OTP)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Order completed successfully. Rewards updated.",
	})
}
