package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/suraj/nitabuddy/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type RewardsHandler struct {
	rewardsModel *models.RewardsModel
	authHandler  *AuthHandler // Add this field
}

func NewRewardsHandler(rewardsModel *models.RewardsModel, authHandler *AuthHandler) *RewardsHandler {
	return &RewardsHandler{
		rewardsModel: rewardsModel,
		authHandler:  authHandler, // Initialize it
	}
}

func (h *RewardsHandler) FetchRewardsByID(w http.ResponseWriter, r *http.Request) {

	userID, err := h.authHandler.GetUserIDFromToken(r) // Use authHandler instead of utils
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Unauthorized: " + err.Error(),
			"coins":   0,
		})
		return
	}

	reward, err := h.rewardsModel.GetRewardsByUserID(userID)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  false,
				"message": "No rewards found for this user",
				"coins":   0,
			})
			return
		}

		// any other DB error
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Failed to fetch rewards: " + err.Error(),
			"coins":   0,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     true,
		"message":    "Rewards fetched successfully",
		"coins":      reward.Coins,
		"fetched_at": time.Now().Format(time.RFC3339),
	})

}
