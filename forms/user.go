package forms

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Coins    uint   `json:"coins"`
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type IncreaseCoinsRequest struct {
	Amount uint `json:"amount" binding:"required"`
}
