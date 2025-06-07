package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/suraj/nitabuddy/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	userModel *models.UserModel
	jwtSecret []byte
}

func NewAuthHandler(userModel *models.UserModel, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		userModel: userModel,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		Name       string `json:"name"`
		Enrollment string `json:"enrollment"`
		Phone      string `json:"phone"`
		Hostel     string `json:"hostel"`
		Branch     string `json:"branch"`
		Year       string `json:"year"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
			"status":  false,
			"token":   "",
		})
		return
	}

	user, err := h.userModel.Create(input.Email, input.Password, input.Name, input.Enrollment, input.Phone, input.Hostel, input.Branch, input.Year)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
			"status":  false,
			"token":   "",
		})
		return
	}

	tokenString, err := h.generateJWT(user.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Failed to generate token",
			"status":  false,
			"token":   "",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User Registered Successfully",
		"status":  true,
		"token":   tokenString,
	})

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
			"status":  false,
			"token":   "",
		})
		return
	}

	user, err := h.userModel.GetByEmail(input.Email)
	if err != nil || !h.userModel.VerifyPassword(user, input.Password) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid Credentials",
			"status":  false,
			"token":   "",
		})
		return
	}

	tokenString, err := h.generateJWT(user.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Failed to generate token",
			"status":  false,
			"token":   "",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login Successful",
		"status":  true,
		"token":   tokenString,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// JWT is stateless: Just ask client to discard the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Logout Successful",
		"status":  true,
	})
}

func (h *AuthHandler) generateJWT(userID primitive.ObjectID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.Hex(),
		"iat":     time.Now().Unix(), // gives new token at every login
		// No expiration - Token Never Expires
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwtSecret)
}

func (h *AuthHandler) GetUserIDFromToken(r *http.Request) (primitive.ObjectID, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return primitive.ObjectID{}, http.ErrNoCookie
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return primitive.ObjectID{}, fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("invalid claims")
	}

	userIDHex, ok := claims["user_id"].(string)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("user_id not found in token")
	}

	return primitive.ObjectIDFromHex(userIDHex)
}
