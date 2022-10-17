package endpoints

import (
	"dapp/lib"
	"dapp/repo"
	"dapp/schema"
	"dapp/schema/dto"
	"dapp/schema/mapper"
	"dapp/service"
	"dapp/service/auth"
	"dapp/service/utils"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/hero"
)

type HAuth struct {
	response  *utils.SvcResponse
	appConf   *utils.SvcConfig
	providers map[string]bool
	validate  *validator.Validate // handle validations for structs and individual fields based on tags
}

// NewAuthHandler create and register the authentication handlers for the App. For the moment, all the
// auth handlers emulates the Oauth2 "password" grant-type using the "client-credentials" flow.
//
// - app [*iris.Application] ~ Iris App instance
//
// - MdwAuthChecker [*context.Handler] ~ Authentication checker middleware
//
// - svcR [*utils.SvcResponse] ~ GrantIntentResponse service instance
//
// - svcC [utils.SvcConfig] ~ Configuration service instance
func NewAuthHandler(app *iris.Application, mdwAuthChecker *context.Handler, svcR *utils.SvcResponse, svcC *utils.SvcConfig, validate *validator.Validate) HAuth { // --- VARS SETUP ---
	h := HAuth{svcR, svcC, make(map[string]bool), validate}
	// filling providers
	h.providers["dapp_provider"] = true

	repoUser := repo.NewRepoUser(svcC)
	svcAuth := auth.NewSvcAuthentication(h.providers, repoUser) // instantiating authentication Service
	svcUser := service.NewSvcUserReqs(repoUser)

	// Simple group: v1
	v1 := app.Party("/api/v1")
	{
		// registering unprotected router
		authRouter := v1.Party("/auth") // authorize
		{
			// --- GROUP / PARTY MIDDLEWARES ---

			// --- DEPENDENCIES ---
			hero.Register(depObtainUserCred)
			hero.Register(svcAuth) // as an alternative, we can put these dependencies as property in the struct HAuth, as we are doing in the rest of the endpoints / handlers
			hero.Register(svcUser)

			// --- REGISTERING ENDPOINTS ---
			// authRouter.Post("/<provider>")	// provider is the auth provider to be used.
			authRouter.Post("/", hero.Handler(h.authIntent))
		}

		// registering protected router
		guardAuthRouter := v1.Party("/auth")
		{
			// --- GROUP / PARTY MIDDLEWARES ---
			guardAuthRouter.Use(*mdwAuthChecker) // registering access token checker middleware

			// --- DEPENDENCIES ---
			hero.Register(DepObtainUserDid)
			hero.Register(svcUser)

			// --- REGISTERING ENDPOINTS ---
			guardAuthRouter.Get("/logout", h.logout)
			guardAuthRouter.Get("/profile", hero.Handler(h.getUserProfile))
		}

		// User management CRUD
		guardUserManagerRouter := v1.Party("/users")
		{
			guardUserManagerRouter.Get("", hero.Handler(h.getUsers))
			guardUserManagerRouter.Get("/{id:string}", hero.Handler(h.getUserById))
			guardUserManagerRouter.Put("/{id:string}", hero.Handler(h.putUserById))
			guardUserManagerRouter.Post("", hero.Handler(h.postUser))

			// TODO: missing to finish (Create and Delete)
			// POST: /users | createUser()
			// DELETE: /users/:id | deleteUser()

			// --- DEPENDENCIES ---
			hero.Register(repoUser)
		}
	}
	return h
}

// region ======== ENDPOINT HANDLERS =====================================================

// authIntent Intent to grant authentication using the provider user's credentials and the specified  auth provider
// @Summary User authentication
// @description.markdown AuthIntent
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param 	credential 	body 	dto.UserCredIn 	true	"User Login Credential"
// @Success 200 "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.wrong_auth_provider"
// @Failure 504 {object} dto.Problem "err.network"
// @Failure 500 {object} dto.Problem "err.json_parse"
// @Router /auth [post]
func (h HAuth) authIntent(ctx iris.Context, uCred *dto.UserCredIn, svcAuth *auth.SvcAuthentication) {
	// using a provider named 'dapp_provider', also injecting dependencies
	provider := "dapp_provider"

	// ej: Aqui podemos comprobar si la bd esta poblada
	//populate := r.IsPopulateDBSvc()
	//if !populate {
	//	h.response.ResErr(&dto.Problem{Status: iris.StatusInternalServerError, Title: schema.ErrBuntdbNotPopulated, Detail: "The database has not been populated yet"}, &ctx)
	//	return
	//}

	authGrantedData, problem := svcAuth.AuthProviders[provider].GrantIntent(uCred, nil) // requesting authorization to evote (provider) mechanisms in this case
	if problem != nil {                                                                 // check for errors
		h.response.ResErr(problem, &ctx)
		return
	}

	// TODO: pass this to the service
	// if so far so good, we are going to create the auth token
	tokenData := mapper.ToAccessTokenDataV(authGrantedData)
	accessToken, err := lib.MkAccessToken(tokenData, []byte(h.appConf.JWTSignKey), h.appConf.TkMaxAge)
	if err != nil {
		h.response.ResErr(&dto.Problem{Status: iris.StatusInternalServerError, Title: schema.ErrJwtGen, Detail: err.Error()}, &ctx)
		return
	}

	h.response.ResOKWithData(string(accessToken), &ctx)
}

// logout this endpoint invalidated a previously granted access token
// @Summary User logout
// @Description This endpoint invalidated a previously granted access token
// @Tags Auth
// @Security ApiKeyAuth
// @Produce  json
// @Param Authorization header string true "Insert access token" default(Bearer <Add access token here>)
// @Success 204 "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 500 {object} dto.Problem "err.generic
// @Router /auth/logout [get]
func (h HAuth) logout(ctx iris.Context) {
	err := ctx.Logout()

	if err != nil {
		h.response.ResErr(&dto.Problem{Status: iris.StatusInternalServerError, Title: schema.ErrGeneric, Detail: err.Error()}, &ctx)
		return
	}

	// so far so good
	h.response.ResOK(&ctx)
}

// getUserProfile Get user from the BD.
// @Summary Get user
// @Tags Auth
// @Security ApiKeyAuth
// @Produce  json
// @Param Authorization header string true "Insert access token" default(Bearer <Add access token here>)
// @Success 200 {object} dto.UserResponse "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 500 {object} dto.Problem "err.generic
// @Router /auth/profile [get]
func (h HAuth) getUserProfile(ctx iris.Context, params dto.InjectedParam, service service.ISvcUser) {
	user, problem := service.GetUserSvc(params.ID)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}
	h.response.ResOKWithData(user, &ctx)
}

// getUsers Get all users from the BD.
// @Summary Get users
// @Tags Users
// @Security ApiKeyAuth
// @Produce  json
// @Param Authorization header string  true "Insert access token" default(Bearer <Add access token here>)
// @Success 200 {object} []dto.UserResponse "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 500 {object} dto.Problem "err.generic
// @Router /users [get]
func (h HAuth) getUsers(ctx iris.Context, service service.ISvcUser) {
	users, problem := service.GetUsersSvc()
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}
	h.response.ResOKWithData(users, &ctx)
}

// getUserById Get user by ID
// @Summary Get user by ID
// @Tags Users
// @Security ApiKeyAuth
// @Produce  json
// @Param Authorization header string  true "Insert access token" default(Bearer <Add access token here>)
// @Param   id          path   string  true "The unique identifier for the user within the account"     Format(string)
// @Success 200 {object} dto.UserResponse "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 500 {object} dto.Problem "err.generic
// @Router /users/{id} [get]
func (h HAuth) getUserById(ctx iris.Context, service service.ISvcUser) {
	// checking param
	userID := ctx.Params().GetString("id")
	if userID == "" {
		h.response.ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: schema.ErrDetInvalidField}, &ctx)
		return
	}

	user, problem := service.GetUserSvc(userID)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}
	h.response.ResOKWithData(user, &ctx)
}

// putUserById Update user.
// @Tags Users
// @Security ApiKeyAuth
// @Produce  json
// @Param Authorization header string  true "Insert access token" default(Bearer <Add access token here>)
// @Param   id          path   string  true "The unique identifier for the user within the account"     Format(string)
// @Param 	Transaction	body   dto.UserUpdateRequest	true	"User Data"
// @Success 200 {object} dto.UserResponse "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 500 {object} dto.Problem "err.generic
// @Router /users/{id} [put]
func (h HAuth) putUserById(ctx iris.Context, service service.ISvcUser) {
	// checking param
	userID := ctx.Params().GetString("id")
	if userID == "" {
		h.response.ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: schema.ErrDetInvalidField}, &ctx)
		return
	}

	// getting data from client
	var requestData dto.UserUpdateRequest

	// unmarshalling the json and check
	if err := ctx.ReadJSON(&requestData); err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}

	response, problem := service.PutUserSvc(userID, requestData)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}
	h.response.ResOKWithData(response, &ctx)
}

// postUser Create user.
// @Tags Users
// @Security ApiKeyAuth
// @Produce  json
// @Param Authorization header string    true  "Insert access token" default(Bearer <Add access token here>)
// @Param 	Transaction	body   dto.User	 true  "User Data"
// @Success 200 {object} dto.UserResponse "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 500 {object} dto.Problem "err.generic
// @Router /users [post]
func (h HAuth) postUser(ctx iris.Context, service service.ISvcUser) {
	// getting data from client
	var requestData dto.User

	// unmarshalling the json and check
	if err := ctx.ReadJSON(&requestData); err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}

	response, problem := service.PostUserSvc(requestData)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}
	h.response.ResOKWithData(response, &ctx)
}

// endregion =============================================================================

// region ======== LOCAL DEPENDENCIES ====================================================

// depObtainUserCred is used as dependencies to obtain / create the user credential from request body (multipart/form-data).
// It returns a dto.UserCredIn struct
func depObtainUserCred(ctx iris.Context) dto.UserCredIn {
	cred := dto.UserCredIn{}

	// Getting data
	cred.Username = ctx.PostValue("username")
	cred.Password = ctx.PostValue("password")

	// TIP: We can do some validation here if we want
	return cred
}

// endregion =============================================================================
