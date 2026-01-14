package logic

import (
	"utility/controller/card/req"
	"utility/controller/card/res"
	"utility/helper"

	"github.com/gofiber/fiber/v2"
)

func C_GetBrandCard(c *fiber.Ctx, req *req.ReqGetBrandCard) (helper.StructMasterResponse, helper.StructMasterErrorResponse, error) {
	GetCreditCardBrand := helper.GetCreditCardBrand(req.CardNumber)
	return helper.StructMasterResponse{
		Status:     true,
		StatusCode: "C0001",
		Message:    "Success",
		Data:       res.ResGetBrandCard{CardBrand: GetCreditCardBrand},
	}.StructMasterResponseFinal(), helper.StructMasterErrorResponse{}.StructMasterErrorResponseFinal(), nil
}
