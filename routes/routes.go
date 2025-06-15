package routes

import (
	"github.com/gorilla/mux"
	"github.com/suraj/nitabuddy/handlers"
)

// Setup configures all the routes for the application
func Setup(r *mux.Router, authHandler *handlers.AuthHandler, orderHandler *handlers.OrderHandler) {

	//Auth Routes
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	r.HandleFunc("/profile", authHandler.GetUserProfile).Methods("GET")

	// orders
	r.HandleFunc("/order", orderHandler.PlaceOrder).Methods("POST")
	r.HandleFunc("/myOrders", orderHandler.FetchMyOrders).Methods("GET")
	r.HandleFunc("/cancelMyOrder/{id}", orderHandler.CancelMyOrder).Methods("DELETE")
}
