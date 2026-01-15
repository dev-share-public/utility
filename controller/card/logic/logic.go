package logic

import (
	"utility/controller/card/req"
	"utility/controller/card/res"
	"utility/helper"
	lib_card "utility/library/card"

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

func L_GetCardDetail(c *fiber.Ctx, req *req.ReqGetCardDetail) (helper.StructMasterResponse, helper.StructMasterErrorResponse, error) {
	GetCardDetail, err := lib_card.GetCardDetail(req.CardNumber)
	if err != nil {
		return helper.StructMasterResponse{}.StructMasterResponseFinal(), helper.StructMasterErrorResponse{
			StatusCode: "C0002",
			Message:    err.Error(),
		}.StructMasterErrorResponseFinal(), err
	}
	return helper.StructMasterResponse{
		Status:     true,
		StatusCode: "C0001",
		Message:    "Success",
		Data:       GetCardDetail,
	}.StructMasterResponseFinal(), helper.StructMasterErrorResponse{}.StructMasterErrorResponseFinal(), nil
}
