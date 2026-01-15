package req

type ReqGetBrandCard struct {
	CardNumber string `form:"card_number" json:"card_number" validate:"required,number,max=20"`
}

type ReqGetCardDetail struct {
	CardNumber string `form:"card_number" json:"card_number" validate:"required,number,min=6,max=6"`
}
