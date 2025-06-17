package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/suraj/nitabuddy/models"
	"github.com/suraj/nitabuddy/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type RewardsHandler struct {
	rewardsModel *models.RewardsModel
}

func NewRewardsHandler(rewardsModel *models.RewardsModel) *RewardsHandler {
	return &RewardsHandler{
		rewardsModel: rewardsModel,
	}
}

func (h *RewardsHandler) FetchRewardsByID(w http.ResponseWriter, r *http.Request) {

	userID, err := utils.ExtractUserIDFromToken(r)
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
