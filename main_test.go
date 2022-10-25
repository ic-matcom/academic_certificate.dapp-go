package main

import (
	"dapp/lib"
	"dapp/repo"
	"dapp/schema/dto"
	"dapp/service"
	"dapp/service/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"reflect"
	"sort"

	//"reflect"
	"runtime"
	"strings"

	"dapp/schema"
	"os"
	"testing"
)

// use a single instance of Validate, it caches struct info

var app *iris.Application
var appConf *utils.SvcConfig
var validate *validator.Validate
var err error

func TestNewApp(t *testing.T) {
	// set environment variable
	if runtime.GOOS == "windows" {
		_ = os.Setenv(schema.EnvConfigPath, "./conf/conf.sample.windows.yaml")
	} else {
		_ = os.Setenv(schema.EnvConfigPath, "./conf/conf.sample.unix.yaml")
	}
	app, appConf = newApp()
	e := httptest.New(t, app)

	// check server status
	e.GET("/status").Expect().Status(httptest.StatusOK)

	// init user repo
	userRepo := repo.NewRepoUser(appConf)
	userServ := service.NewSvcUserReqs(userRepo)

	// get users list
	expectedUsers, _ := userServ.GetUsersSvc()

	// without basic auth
	usersBody := e.GET("/api/v1/users").Expect().Body().Raw()
	var gotUsers []dto.UserResponse
	err := jsoniter.UnmarshalFromString(usersBody, &gotUsers)
	if err != nil {
		t.Errorf("unmarshal error: %v", err)
	}

	// sort arrays to avoid deeply equal error
	sort.Slice(gotUsers, func(p, q int) bool {
		return gotUsers[p].Username < gotUsers[q].Username
	})

	sort.Slice(*expectedUsers, func(p, q int) bool {
		return (*expectedUsers)[p].Username < (*expectedUsers)[q].Username
	})

	if !reflect.DeepEqual(expectedUsers, &gotUsers) {
		t.Errorf("got users { %+v } did not match expected users { %+v }", &gotUsers, expectedUsers)
	}

	// with valid JWT auth
	cred := dto.UserCredIn{
		Username: (*expectedUsers)[0].Username,
		Password: "password1",
	}

	e.POST("/api/v1/auth").WithJSON(cred).Expect().Status(httptest.StatusOK)

	bearer := fmt.Sprintf("Bearer %s", fakeAccessToken((*expectedUsers)[0].Username))
	e.GET("/api/v1/auth/profile").WithHeader("Authorization", bearer).Expect().Status(httptest.StatusOK)

	// with invalid JWT auth
	cred = dto.UserCredIn{
		Username: (*expectedUsers)[0].Username,
		Password: "fakepasswd",
	}

	e.POST("/api/v1/auth").WithJSON(cred).Expect().Status(httptest.StatusUnauthorized)

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
	err = validate.Struct(userValid)
	if err != nil {
		t.Errorf("user %s must be valid", userValid.Username)
	}

	userInvalid := dto.User{
		Username:   "fakeusername",
		Passphrase: lib.GenerateUUIDStr(),
		FirstName:  "Fake Name",
		LastName:   "Fake LastName",
		Email:      "fake-email.com", // invalid email
	}
	// validate user fields
	err = validate.Struct(userInvalid)
	if err == nil {
		t.Errorf("user %s must be invalid", userInvalid.Username)
	}
}

func TestListUsers(t *testing.T) {

}

func fakeAccessToken(username string) []byte {
	tokenData := dto.AccessTokenData{Scope: strings.Fields("dapp.fabric"), Claims: dto.InjectedParam{ID: username, Username: username}}
	accessToken, _ := lib.MkAccessToken(&tokenData, []byte(appConf.JWTSignKey), appConf.TkMaxAge)
	return accessToken
}
