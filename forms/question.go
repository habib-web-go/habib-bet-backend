package forms

type QuestionAnswer struct {
	Option string `json:"option" binding:"required"`
}
