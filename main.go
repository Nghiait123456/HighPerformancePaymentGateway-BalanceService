package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/high-performance-payment-gateway/balance-service/balance"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/config/env"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/error_handle"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/log_init"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/validate"
	"os"
)

func handle(c *fiber.Ctx) error {
	panic("have error")
	return fiber.NewError(500, "have erorssss")
	return c.SendString("Hello, World!")

}

type Person struct {
	Name uint64 `json:"name" xml:"name" form:"name"`
	//Pass interface{} `json:"pass" xml:"pass" form:"pass"`
}

func main() {
	balanceModule := balance.NewModule()
	balanceModule.Start()
}

func response(c *fiber.Ctx) error {
	return c.SendString("Success !!!!!!!!!!!!!!!111")
}

func temp() {

	log_init.Init(log_init.Log{
		TypeFormat: log_init.TYPE_FORMAT_TEXT,
		TypeOutput: log_init.TYPE_OUTPUT_FILE,
		LinkFile:   "balance/infrastructure/log/log_file/log.log",
	})
	//
	//gCf := TestEnvGlobal()
	//TestAutoChangeSecret(gCf.AuthInternalServiceConfig())

	// Create a new fiber instance with custom config
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: error_handle.CustomMessageError,
	})

	//config error handle
	eH := error_handle.ErrorHandle{
		App: app,
	}
	eH.Init()

	// config alert panic notify
	alertAc := error_handle.AlertAc{
		Dsn:              "https://4ea29cdaaa424266a113571ac407c88a@o1406092.ingest.sentry.io/6739879",
		Repanic:          true,
		Debug:            true,
		AttachStacktrace: true,
	}
	alertPanic := error_handle.NewPanicHandle(&alertAc, app)
	alertPanic.Init()

	//app.Use(fibersentry.New(fibersentry.Config{
	//	Repanic:         true,
	//	WaitForDelivery: true,
	//}))
	//
	//enhanceSentryEvent := func(c *fiber.Ctx) error {
	//	if hub := fibersentry.GetHubFromContext(c); hub != nil {
	//		hub.Scope().SetTag("someRandomTag", "maybeYouNeedIt")
	//	}
	//	return c.Next()
	//}

	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("panic0999hbhbhb" + "333333333333333333333")
	})

	//app.All("/", func(c *fiber.Ctx) error {
	//	if hub := fibersentry.GetHubFromContext(c); hub != nil {
	//		hub.WithScope(func(scope *sentry.Scope) {
	//			scope.SetExtra("unwantedQuery", "someQueryDataMaybe")
	//			hub.CaptureMessage("User provided unwanted query string, but we recovered just fine")
	//		})
	//	}
	//	return c.SendStatus(fiber.StatusOK)
	//})

	errApp := app.Listen(":3000")
	if errApp != nil {
		errMApp := fmt.Sprintf("Init app error, message: %s", errApp.Error())
		panic(errMApp)
		os.Exit(0)
	}

}
func TestJsonHamas() {
	ListES := validate.ListErrorsDefaultShow{}

	oneErr := validate.OneErrorDefaultShow{
		Field:     "name",
		Rule:      "require",
		Message:   "field is require",
		ParamRule: "required",
	}

	ListES["name"] = oneErr

	showE := validate.DefaultShowError{
		ListError: ListES,
	}

	rs, errCV := json.Marshal(showE)
	if errCV != nil {
		fmt.Println("error")
	}

	fmt.Println(string(rs))
}

func TestValidate() {

	// User contains user information
	type User struct {
		FirstName   string `json:"fname" validate:"maxlengError,alpha,required"`
		LastName    string `json:"lname" validate:"alpha"`
		Age         uint8  `validate:"gte=20,lte=65"`
		Email       string `json:"e-mail" validate:"required,email"`
		JoiningDate string `validate:"datetime"`
	}

	// use a single instance of validate, it caches struct info
	var validate *validator.Validate

	validate = validator.New()
	validate.RegisterValidation("maxlengError", func(fl validator.FieldLevel) bool {
		fmt.Println("ffffff", fl.Field())
		return len(fl.Field().String()) > 6
	})

	user := &User{
		FirstName:   "",
		LastName:    "Test",
		Age:         75,
		Email:       "Badger.Smith@",
		JoiningDate: "005-25-10",
	}

	err := validate.Struct(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		fmt.Println("------ List of tag fields with error ---------")

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.StructField())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println("---------------")
		}
		return
	}
}

func TestEnvGlobal() env.GlobalConfigInterface {
	os.Setenv(env.ENV_ENVIRONMENT, env.ENV_DEV)
	os.Setenv(env.AWS_SECRET_NAME_GLOBAL_KEY, "payment-balance-service-qwedjfndasndajndn12")
	os.Setenv(env.AWS_REGION_GLOBAL_KEY, "ap-southeast-1")
	os.Setenv(env.AWS_VERSION_STATE_GLOBAL_KEY, "AWSCURRENT")

	gCf := env.NewGlobalConfig()
	gCf.Init(env.ConfigEnv{
		FilePatchLogInLocal: "balance/infrastructure/config/env/.env",
	})

	return gCf
}

func TestAutoChangeSecret(ai env.AuthInternalServiceConfigInterface) {
	//firstInit
	//err := auth_internal_service.FirstTimeContructSecret()
	//if err != nil {
	//	fmt.Println("contrcut secret firsttime erro : %s", err.Error())
	//} else {
	//	fmt.Println("contruct secret in firsttime success")
	//}

	//auto change secret
	//auto := auth_internal_service.NewAutoChangeSecret(ai)
	//err := auto.Init()
	//if err != nil {
	//	fmt.Println("init error : %s", err.Error())
	//} else {
	//	fmt.Println("init success")
	//}

	//test auth
	//auth := auth_internal_service.NewAuthInternalService(ai)
	//auth.Init()
	//fmt.Println(auth.Auth("42zPQ0sIAmjtHSTRANmc00000000000000000004"))

	// test integrated

	//auto change secret
	//go func() {
	//	auto := auth_internal_service.NewAutoChangeSecret(ai)
	//	err := auto.Init()
	//	if err != nil {
	//		fmt.Println("init error : %s", err.Error())
	//	} else {
	//		fmt.Println("init success")
	//	}
	//}()
	//
	////init auth
	//auth := auth_internal_service.NewAuthInternalService(ai)
	//auth.Init()
	//go func() {
	//	for {
	//		fmt.Println(auth.Auth("K38xvC2XthUC6fHSxz6W00000000000000000006"))
	//		time.Sleep(5 * time.Second)
	//	}
	//}()
}
