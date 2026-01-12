package validationstruct

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
	"utility/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var MyValidate *validator.Validate

func SetupValidate() {
	MyValidate = validator.New()
	MyValidate.RegisterValidation("date_ymd", func(fl validator.FieldLevel) bool {
		_, err := time.Parse("2006-01-02", fl.Field().String())
		return err == nil
	})

	MyValidate.RegisterValidation("e164prefix", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		if !strings.HasPrefix(val, "+") || len(val) <= 1 {
			return false
		}
		for _, r := range val[1:] {
			if !unicode.IsDigit(r) {
				return false
			}
		}
		return true
	})

	MyValidate.RegisterValidation("password_complex", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		if len(password) < 8 {
			return false
		}

		var (
			hasUpper  = false
			hasLower  = false
			hasNumber = false
			hasSymbol = false
		)

		for _, c := range password {
			switch {
			case unicode.IsUpper(c):
				hasUpper = true
			case unicode.IsLower(c):
				hasLower = true
			case unicode.IsDigit(c):
				hasNumber = true
			case unicode.IsPunct(c) || unicode.IsSymbol(c):
				hasSymbol = true
			}
		}

		return hasUpper && hasLower && hasNumber && hasSymbol
	})

	MyValidate.RegisterValidation("engnumdashuscoredot", func(fl validator.FieldLevel) bool {
		pattern := `^[a-zA-Z0-9_.-]+$`
		matched, _ := regexp.MatchString(pattern, fl.Field().String())
		return matched
	})

	// enum validator: WEB, ANDROID, IOS
	MyValidate.RegisterValidation("platform_enum", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		switch val {
		case "WEB", "ANDROID", "IOS":
			return true
		}
		return false
	})

	fmt.Println("✅ Setup MyValidate Struct")

}

type ValidateStructRes struct {
	Field   string      `json:"field"` // by passing alt name to ReportError like below
	Tag     string      `json:"tag"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
	Param   interface{} `json:"param"`
}

func ValidateStruct(c *fiber.Ctx, r interface{}) (helper.StructMasterErrorResponse, error) {
	var errors []*ValidateStructRes
	var errorMaps []map[string]interface{}
	err := MyValidate.Struct(r)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var ele ValidateStructRes
			ele.Field = err.Field()
			ele.Tag = err.Tag()
			ele.Message = err.Error()
			ele.Param = err.Param()
			ele.Value = err.Value()
			// fmt.Println(ele.Tag)
			errors = append(errors, &ele)
			errorMaps = append(errorMaps, map[string]interface{}{
				"field":   err.Field,
				"tag":     err.Tag,
				"message": ele.Message,
				"param":   ele.Param,
				"value":   ele.Value,
			})
		}
	}
	return helper.GetErrorByKeyV1(c, "BadRequest", errorMaps...), err
}
