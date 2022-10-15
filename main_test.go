package main

import (
	"dapp/lib"
	"dapp/schema/dto"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/httptest"
	"runtime"

	"dapp/schema"
	"os"
	"testing"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func TestNewApp(t *testing.T) {
	// set environment variable
	if runtime.GOOS == "windows" {
		_ = os.Setenv(schema.EnvConfigPath, "./conf/conf.sample.windows.yaml")
	} else {
		_ = os.Setenv(schema.EnvConfigPath, "./conf/conf.sample.unix.yaml")
	}
	app, _ := newApp()
	e := httptest.New(t, app)
	// check server status
	e.GET("/status").Expect().Status(httptest.StatusOK)

	// without basic auth
	//e.GET("/api/v1/drones").Expect().Status(httptest.StatusUnauthorized)
	//e.GET("/api/v1/medications").Expect().Status(httptest.StatusUnauthorized)

	// with valid JWT auth
	//cred := dto.UserCredIn{
	//	Username: "richard.sargon@meinermail.com",
	//	Password: "password1",
	//}
	//
	//_ = e.POST("/api/v1/auth").WithJSON(cred).Expect().Status(httptest.StatusOK)
	//
	//// with invalid JWT auth
	//cred = dto.UserCredIn{
	//	Username: "noexist@meinermail.com",
	//	Password: "fakepasswd",
	//}
	//
	//_ = e.POST("/api/v1/auth").WithJSON(cred).Expect().Status(httptest.StatusUnauthorized)

	// user valid
	userValid := dto.User{
		Username:   "fakeusername",
		Passphrase: lib.GenerateUUIDStr(),
		FirstName:  "Fake Name",
		LastName:   "Fake LastName",
		Email:      "fake@email.com",
	}

	validate = validator.New()

	// validate drone fields
	err := validate.Struct(userValid)
	if err != nil {
		t.Errorf("user %s must be valid", userValid.Username)
	}

	userInvalid := dto.User{
		Username:   "fakeusername",
		Passphrase: lib.GenerateUUIDStr(),
		FirstName:  "Fake Name",
		LastName:   "Fake LastName",
		Email:      "fake-email.com",   // invalid email
	}
	// validate drone fields
	err = validate.Struct(userInvalid)
	if err == nil {
		t.Errorf("user %s must be invalid", userInvalid.Username)
	}
}
