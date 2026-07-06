package general

type GeneralRequest struct {
	Userid   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}
