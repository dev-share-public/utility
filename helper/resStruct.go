package helper

import (
	"utility/config/envi"

	"github.com/gofiber/fiber/v2"
)

type StructMasterResponse struct {
	Status     bool   `json:"status"`
	VersionApi string `json:"version_api"`
	StatusCode string `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func (r StructMasterResponse) StructMasterResponseFinal() StructMasterResponse {
	if r.Data == nil {
		r.Data = map[string]any{}
	}

	if r.VersionApi == "" {
		r.VersionApi = envi.GetEnv("MASTER_VERSION_API", "")
	}

	return r
}

type StructMasterErrorResponse struct {
	Status     bool   `json:"status"`
	VersionApi string `json:"version_api"`
	StatusCode string `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func (r StructMasterErrorResponse) StructMasterErrorResponseFinal() StructMasterErrorResponse {
	if r.Data == nil {
		r.Data = map[string]any{}
	}

	if r.VersionApi == "" {
		r.VersionApi = envi.GetEnv("MASTER_VERSION_API", "")
	}

	return r
}

// ErrorCode เก็บรายละเอียดข้อผิดพลาด รองรับหลายภาษา
type ErrorCodeDescription struct {
	RefCode     string            `json:"refcode,omitempty"`
	Description map[string]string `json:"description,omitempty"`
}

var ErrorDefaultMessages = map[string]ErrorCodeDescription{
	"BadRequest":          {"400", map[string]string{"en": "Bad request error", "th": "คำขอไม่ถูกต้อง"}},
	"NotFound":            {"404", map[string]string{"en": "Resource not found", "th": "ไม่พบทรัพยากร"}},
	"Unauthorized":        {"401", map[string]string{"en": "Unauthorized access", "th": "ไม่มีสิทธิ์เข้าถึง"}},
	"BadRequestParamiter": {"2001", map[string]string{"en": "BadRequest paramiter", "th": "คำขอ paramiter ไม่ถูกต้อง"}},
	"BadRequestHeader":    {"2001", map[string]string{"en": "BadRequest header", "th": "คำขอ header ไม่ถูกต้อง"}},
}

func GetErrorByKeyV1(ctx *fiber.Ctx, key string, err_cust ...map[string]interface{}) StructMasterErrorResponse {
	var validationErrors []map[string]interface{}

	for _, m := range err_cust {
		item := make(map[string]interface{})

		if v, ok := m["field"]; ok {
			if fn, ok := v.(func() string); ok {
				item["field"] = fn()
			} else {
				item["field"] = v
			}
		}

		if v, ok := m["tag"]; ok {
			if fn, ok := v.(func() string); ok {
				item["tag"] = fn()
			} else {
				item["tag"] = v
			}
		}

		validationErrors = append(validationErrors, item)
	}

	map_err := make(map[string]interface{})
	if len(validationErrors) > 0 {
		if envi.GetEnv("IS_DISPLAY_VALIDATESTRUCT", "false") == "true" {
			map_err["validation_errors"] = validationErrors
		}
		// ถ้าไม่มี error เลย จะไม่ใส่ key นี้
	}

	lang := ctx.Get("Accept-Language", "en")
	if err, exists := ErrorDefaultMessages[key]; exists {
		msg := err.Description[lang]
		if msg == "" {
			msg = err.Description["en"]
		}
		return StructMasterErrorResponse{
			Status:     false,
			StatusCode: err.RefCode,
			Message:    msg,
			Data:       map_err,
		}.StructMasterErrorResponseFinal()
	}
	return StructMasterErrorResponse{
		Status:     false,
		StatusCode: "500",
		Message:    "Unknown error",
		Data:       map_err,
	}.StructMasterErrorResponseFinal()
}
