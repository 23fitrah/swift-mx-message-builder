package general

type GeneralRequest struct {
	Userid   string `json:"user_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}
