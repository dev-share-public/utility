package setup

import (
	"utility/config/envi"
	"utility/config/timezo"
	"utility/config/validationstruct"
	"utility/route"

	"github.com/gofiber/fiber/v2"
)

type AppSetUp struct {
	FiberV2        *fiber.App
	FiberV2RunPort string
}

func FirstSetup() (*AppSetUp, error) {
	if ErrEnvi := envi.InitEnvi(); ErrEnvi != nil {
		// panic("Error Loading .Env file form path : " + ErrEnvi.Error())
		return &AppSetUp{}, ErrEnvi
	}
	if ErrTimeZo := timezo.InitTime("Asia/Bangkok"); ErrTimeZo != nil {
		return &AppSetUp{}, ErrTimeZo
	}

	validationstruct.SetupValidate()
	FIBER_RUN_PORT := envi.GetEnv("FIBER_RUN_PORT", "80")
	app := route.SetUpRoute()
	return &AppSetUp{
		FiberV2:        app,
		FiberV2RunPort: FIBER_RUN_PORT,
	}, nil
}
