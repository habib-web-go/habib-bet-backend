package forms

import "github.com/habib-web-go/habib-bet-backend/models"

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Coins    uint   `json:"coins"`
}

func CreateUserResponse(u *models.User) *UserResponse {
	return &UserResponse{
		Id:       u.ID,
		Username: u.Username,
		Coins:    u.Coins,
	}
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type IncreaseCoinsRequest struct {
	Amount uint `json:"amount" binding:"required"`
}
