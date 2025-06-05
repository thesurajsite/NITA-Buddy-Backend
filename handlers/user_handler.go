package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *AuthHandler) GetUserDetails(w http.ResponseWriter, r *http.Request) {

	// Get userId from token
	userID, err := h.GetUserIDFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid or missing token",
			"status":  false,
			"user":    nil,
		})
		return
	}

	// Get user details from token
	user, err := h.userModel.GetByID(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "user not found",
			"status":  false,
			"user":    nil,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "user details fetched successfully",
		"status":  true,
		"user":    user,
	})
}
