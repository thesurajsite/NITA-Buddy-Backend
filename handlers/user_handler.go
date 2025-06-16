package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *AuthHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {

	// Get userId from token
	userID, err := h.GetUserIDFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Invalid or missing token",
			"user":    nil,
		})
		return
	}

	// Get user details from token
	user, err := h.userModel.GetUserByID(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "user not found",
			"user":    nil,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "user details fetched successfully",
		"user":    user,
	})
}

func (h *AuthHandler) GetUserProfileFromID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userIDstr := vars["id"]
	userID, err := primitive.ObjectIDFromHex(userIDstr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "invalid user ID format",
			"user":    nil,
		})
		return
	}

	// Get user details from token
	user, err := h.userModel.GetUserByID(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "user not found",
			"user":    nil,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "user details fetched successfully",
		"user":    user,
	})
}
