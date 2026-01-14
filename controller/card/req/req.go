package req

type ReqGetBrandCard struct {
	CardNumber string `form:"card_number" json:"card_number" validate:"required,number,max=20"`
}
