package handle

import (
	"net/http"
	"utility/config/validationstruct"
	"utility/controller/card/logic"
	"utility/controller/card/req"
	"utility/helper"

	"github.com/gofiber/fiber/v2"
)

func H_GetBrandCard(c *fiber.Ctx) error {
	req := req.ReqGetBrandCard{}
	if err := c.BodyParser(&req); err != nil {
		ResBadRequestParamiter := helper.GetErrorByKeyV1(c, "BadRequestParamiter")
		return c.Status(http.StatusOK).JSON(ResBadRequestParamiter)
	}

	if validator, err := validationstruct.ValidateStruct(c, req); err != nil {
		return c.Status(http.StatusOK).JSON(validator)
	}

	resp, errMessage, err := logic.C_GetBrandCard(c, &req)
	if err != nil {
		return c.Status(200).JSON(errMessage)
	}

	return c.Status(200).JSON(resp)
}

func H_GetCardDetail(c *fiber.Ctx) error {
	req := req.ReqGetCardDetail{}
	if err := c.BodyParser(&req); err != nil {
		ResBadRequestParamiter := helper.GetErrorByKeyV1(c, "BadRequestParamiter")
		return c.Status(http.StatusOK).JSON(ResBadRequestParamiter)
	}

	if validator, err := validationstruct.ValidateStruct(c, req); err != nil {
		return c.Status(http.StatusOK).JSON(validator)
	}

	resp, errMessage, err := logic.L_GetCardDetail(c, &req)
	if err != nil {
		return c.Status(200).JSON(errMessage)
	}

	return c.Status(200).JSON(resp)
}
