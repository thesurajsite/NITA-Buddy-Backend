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

	// Profile
	r.HandleFunc("/profile", authHandler.GetUserProfile).Methods("GET")
	r.HandleFunc("/profile/{id}", authHandler.GetUserProfileFromID).Methods("GET")

	// orders
	r.HandleFunc("/order", orderHandler.PlaceOrder).Methods("POST")
	r.HandleFunc("/allOrders", orderHandler.FetchOtherOrders).Methods("GET")
	r.HandleFunc("/myOrders", orderHandler.FetchMyOrders).Methods("GET")
	r.HandleFunc("/cancelMyOrder/{id}", orderHandler.CancelMyOrder).Methods("DELETE")
	r.HandleFunc("/acceptOrder/{id}", orderHandler.AcceptOrder).Methods("PUT")
}
