package endpoints

import (
	"dapp/lib"
	"dapp/repo"
	"dapp/schema"
	"dapp/schema/dto"
	"dapp/schema/models"
	"dapp/service"
	"dapp/service/utils"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/hero"
)

// DappHandler  endpoint handler struct for Dapp
type DappHandler struct {
	response *utils.SvcResponse
	service  *service.ISvcDapp
	validate *validator.Validate // handle validations for structs and individual fields based on tags
	uTrans   *ut.UniversalTranslator
}

// NewDappHandler create and register the handler for Dapp
//
// - app [*iris.Application] ~ Iris App instance
//
// - MdwAuthChecker [*context.Handler] ~ Authentication checker middleware
//
// - svcR [*utils.SvcResponse] ~ GrantIntentResponse service instance
//
// - svcC [utils.SvcConfig] ~ Configuration service instance
func NewDappHandler(app *iris.Application, mdwAuthChecker *context.Handler, svcR *utils.SvcResponse, svcC *utils.SvcConfig, validate *validator.Validate, uT *ut.UniversalTranslator) DappHandler { // --- VARS SETUP ---
	repoDapp := repo.NewRepoDapp(svcC)
	svc := service.NewSvcDappReqs(repoDapp)
	// registering protected / guarded router
	h := DappHandler{svcR, &svc, validate, uT}

	// --- DEPENDENCIES ---
	hero.Register(lib.DepObtainUserDid)

	// Simple group: v1
	v1 := app.Party("/api/v1")
	{
		publicAPI := v1.Party("/dapp")
		{
			publicAPI.Get("/certificates_by_state/{state: int}", hero.Handler(h.getCertificatesByState))
			publicAPI.Get("/certificates_by_accredited/{accredited: string}", hero.Handler(h.getCertificatesByAccredited))
			publicAPI.Get("/certificates/{id: string}", hero.Handler(h.getAssetById))
		}
		// registering protected / guarded router
		protectedAPI := v1.Party("/dapp")
		{
			// --- GROUP / PARTY MIDDLEWARES ---
			protectedAPI.Use(*mdwAuthChecker)

			protectedAPI.Post("/query", hero.Handler(h.postQuery))
			protectedAPI.Post("/transaction", hero.Handler(h.postTransaction))
			protectedAPI.Post("/certificates", hero.Handler(h.postCreateAsset))
			protectedAPI.Put("/certificates", hero.Handler(h.putUpdateAsset))
			protectedAPI.Put("/validate_certificate", hero.Handler(h.putValidateCertificate))
			protectedAPI.Put("/invalidate_certificate", hero.Handler(h.putInvalidateCertificate))
			protectedAPI.Delete("/certificates/{id: string}", hero.Handler(h.deleteAssetById))
		}
	}
	return h
}

// endregion ======== Dapp ======================================================

// postQuery Performs a query in blockchain
// @Summary Performs a query in blockchain
// @description.markdown Query
// @Tags DApp
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	true 	"Insert access token" default(Bearer <Add access token here>)
// @Param 	Query		    body 	dto.Transaction 	true	"Data as a JSON object"
// @Success 200 {object} dto.QueryResult "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 502 {object} dto.Problem "err.bad_gateway"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /dapp/query [post]
func (h DappHandler) postQuery(ctx *context.Context, params dto.InjectedParam) {
	// getting query data
	var query dto.Transaction

	// unmarshalling the json and check
	if err := ctx.ReadJSON(&query); err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}
	bcRes, problem := (*h.service).Query(query, params.Username)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}

	(*h.response).ResOKWithData(bcRes, &ctx)
}

// postTransaction Send transaction to peers
// @Summary Send transaction to peers
// @description.markdown Transaction
// @Tags DApp
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	true 	"Insert access token" default(Bearer <Add access token here>)
// @Param 	Transaction		body 	dto.Transaction	true	"Transaction Data"
// @Success 202 {object} dto.Transaction "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 502 {object} dto.Problem "err.bad_gateway"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /dapp/transaction [post]
func (h DappHandler) postTransaction(ctx iris.Context, params dto.InjectedParam) {
	if params.Role != models.Role_CertificateAdmin {
		(*h.response).ResUnauthorized(&ctx)
		return
	}
	// getting data from client
	var requestData dto.Transaction

	// unmarshalling the json and check
	if err := ctx.ReadJSON(&requestData); err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}
	// trying to submit the transaction
	bcRes, problem := (*h.service).Invoke(requestData, params.Username)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}

	(*h.response).ResOKWithData(bcRes, &ctx)
}

// postCreateAsset Create Asset in ledger
// @Summary Create Certificate
// @description.markdown CreateAsset
// @Tags Certificate
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	        true  "Insert access token" default(Bearer <Add access token here>)
// @Param   channel         query   string          true  "Insert channel" default(mychannel)"
// @Param   chaincode       query   string          true  "Insert chaincode id" default(certificate)"
// @Param   signer          query   string          true  "Insert signer" default(User1)"
// @Param 	Transaction		body 	dto.CreateAsset	true  "Transaction Data"
// @Success 202 {object} dto.Asset "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 502 {object} dto.Problem "err.bad_gateway"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /dapp/certificates [post]
func (h DappHandler) postCreateAsset(ctx iris.Context, params dto.InjectedParam) {
	if params.Role != models.Role_CertificateAdmin {
		(*h.response).ResUnauthorized(&ctx)
		return
	}
	queryParams := new(dto.QueryParamChaincode)
	err := lib.ParamsToStruct(ctx, queryParams)
	if err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}

	var requestData dto.CreateAsset
	// unmarshalling the json and check
	if err := ctx.ReadJSON(&requestData); err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}
	// trying to submit the transaction
	bcRes, problem := (*h.service).CreateAsset(&requestData, params.Username, queryParams)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}

	(*h.response).ResOKWithData(bcRes, &ctx)
}

// putUpdateAsset Update Asset in ledger with given data
// @Summary Update Certificate
// @Description Update data from certificate with specified ID
// @Tags Certificate
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	        true  "Insert access token" default(Bearer <Add access token here>)
// @Param   channel         query   string          true  "Insert channel" default(mychannel)"
// @Param   chaincode       query   string          true  "Insert chaincode id" default(certificate)"
// @Param   signer          query   string          true  "Insert signer" default(User1)"
// @Param 	Transaction		body 	dto.Asset    	true  "Transaction Data"
// @Success 202 {object} dto.Asset "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 502 {object} dto.Problem "err.bad_gateway"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /dapp/certificates [put]
func (h DappHandler) putUpdateAsset(ctx iris.Context, params dto.InjectedParam) {
	if params.Role != models.Role_CertificateAdmin {
		(*h.response).ResUnauthorized(&ctx)
		return
	}
	queryParams := new(dto.QueryParamChaincode)
	err := lib.ParamsToStruct(ctx, queryParams)
	if err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}

	var requestData dto.Asset
	// unmarshalling the json and check
	if err := ctx.ReadJSON(&requestData); err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}
	// trying to submit the transaction
	bcRes, problem := (*h.service).UpdateAsset(&requestData, params.Username, queryParams)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}

	(*h.response).ResOKWithData(bcRes, &ctx)
}

// getAssetById Get Asset from ledger with specified ID
// @Summary Get Certificate
// @Description Get Certificate data from ledger with specified ID
// @Tags Certificate
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param 	id		    	path 	string     true	 "Cartificate ID"
// @Param   channel         query   string     true  "Insert channel" default(mychannel)"
// @Param   chaincode       query   string     true  "Insert chaincode id" default(certificate)"
// @Param   signer          query   string     true  "Insert signer" default(User1)"
// @Success 202 {object} dto.Asset "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 502 {object} dto.Problem "err.bad_gateway"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /dapp/certificates/{id} [get]
func (h DappHandler) getAssetById(ctx iris.Context) {
	id := ctx.Params().GetString("id")
	queryParams := new(dto.QueryParamChaincode)
	lib.ParamsToStruct(ctx, queryParams)

	bcRes, problem := (*h.service).GetAsset(id, schema.GuestUser, queryParams)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}

	(*h.response).ResOKWithData(bcRes, &ctx)
}

// putValidateCertificate Validate Asset in ledger
// @Summary Validate Certificate
// @Description Validate certificate with specified ID. The order for validation is: Secretary -> Dean -> Rector
// @Tags Certificate
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	        true  "Insert access token" default(Bearer <Add access token here>)
// @Param   channel         query   string          true  "Insert channel" default(mychannel)"
// @Param   chaincode       query   string          true  "Insert chaincode id" default(certificate)"
// @Param   signer          query   string          true  "Insert signer" default(User1)"
// @Param 	Transaction		body 	dto.SignAsset	true  "Transaction Data"
// @Success 202 {object} dto.Asset "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 502 {object} dto.Problem "err.bad_gateway"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /dapp/validate_certificate [put]
func (h DappHandler) putValidateCertificate(ctx iris.Context, params dto.InjectedParam) {
	if !lib.Contains([]string{models.Role_Secretary, models.Role_Dean, models.Role_Rector}, params.Role) {
		(*h.response).ResUnauthorized(&ctx)
		return
	}
	queryParams := new(dto.QueryParamChaincode)
	err := lib.ParamsToStruct(ctx, queryParams)
	if err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}

	var requestData dto.SignAsset
	// unmarshalling the json and check
	if err := ctx.ReadJSON(&requestData); err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}
	// trying to submit the transaction
	bcRes, problem := (*h.service).ValidateAsset(&requestData, &params, queryParams)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}

	(*h.response).ResOKWithData(bcRes, &ctx)
}

// putInvalidateCertificate Invalidate Asset in ledger
// @Summary Invalidate Certificate
// @Description Invalidate Certificate with specified ID, Must be provided details about that invalidation.
// @Tags Certificate
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	             true  "Insert access token" default(Bearer <Add access token here>)
// @Param   channel         query   string               true  "Insert channel" default(mychannel)"
// @Param   chaincode       query   string               true  "Insert chaincode id" default(certificate)"
// @Param   signer          query   string               true  "Insert signer" default(User1)"
// @Param 	Transaction		body 	dto.InvalidateAsset	 true  "Transaction Data"
// @Success 202 {object} dto.Asset "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 502 {object} dto.Problem "err.bad_gateway"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /dapp/invalidate_certificate [put]
func (h DappHandler) putInvalidateCertificate(ctx iris.Context, params dto.InjectedParam) {
	if !lib.Contains([]string{models.Role_Secretary, models.Role_Dean, models.Role_Rector, models.Role_CertificateAdmin}, params.Role) {
		(*h.response).ResUnauthorized(&ctx)
		return
	}
	queryParams := new(dto.QueryParamChaincode)
	err := lib.ParamsToStruct(ctx, queryParams)
	if err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}

	var requestData dto.InvalidateAsset
	// unmarshalling the json and check
	if err := ctx.ReadJSON(&requestData); err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}
	// trying to submit the transaction
	bcRes, problem := (*h.service).InvalidateAsset(&requestData, &params, queryParams)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}

	(*h.response).ResOKWithData(bcRes, &ctx)
}

// deleteAssetById Delete Asset from ledger with specified ID
// @Summary Delete Certificate
// @Description Delete Certificate from ledger with specified ID
// @Tags Certificate
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	   true  "Insert access token" default(Bearer <Add access token here>)
// @Param 	id		    	path 	string     true	 "Cartificate ID"
// @Param   channel         query   string     true  "Insert channel" default(mychannel)"
// @Param   chaincode       query   string     true  "Insert chaincode id" default(certificate)"
// @Param   signer          query   string     true  "Insert signer" default(User1)"
// @Success 202 {object} dto.Asset "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 502 {object} dto.Problem "err.bad_gateway"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /dapp/certificates/{id} [delete]
func (h DappHandler) deleteAssetById(ctx iris.Context, params dto.InjectedParam) {
	if params.Role != models.Role_CertificateAdmin {
		(*h.response).ResUnauthorized(&ctx)
		return
	}
	id := ctx.Params().GetString("id")
	queryParams := new(dto.QueryParamChaincode)
	err := lib.ParamsToStruct(ctx, queryParams)
	if err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}

	bcRes, problem := (*h.service).DeleteAsset(id, params.Username, queryParams)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}

	(*h.response).ResOKWithData(bcRes, &ctx)
}

// getCertificatesByState Performs a query in blockchain for certificates with some state
// @Summary Performs a query in blockchain for certificates with a specified state
// @Description Return Certificates that have the specified state.
// @Tags Certificate
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param 	state		    path 	int 	true	"State of the assets"
// @Param 	page_limit		query 	int 	true	"Amount of assets per page" default(5)"
// @Param   bookmark        query   string  false   "Bookmark to know last asset gotten"
// @Param   channel         query   string  true    "Insert channel" default(mychannel)"
// @Param   chaincode       query   string  true    "Insert chaincode id" default(certificate)"
// @Param   signer          query   string  true    "Insert signer" default(User1)"
// @Success 200 {object} dto.QueryResult "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 502 {object} dto.Problem "err.bad_gateway"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /dapp/certificates_by_state/{state} [get]
func (h DappHandler) getCertificatesByState(ctx *context.Context) {
	state := ctx.Params().GetIntDefault("state", -1)
	pageLimit := ctx.URLParamIntDefault("page_limit", -1)

	qp := new(dto.QueryParamChaincode)
	err := lib.ParamsToStruct(ctx, qp)
	if err != nil || state == -1 || pageLimit == -1 {
		h.response.ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: schema.ErrDetInvalidField}, &ctx)
		return
	}
	qp.PageLimit = pageLimit

	bcRes, problem := (*h.service).GetAssetsByState(state, schema.GuestUser, qp)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}
	(*h.response).ResOKWithData(bcRes, &ctx)
}

// getNewCertificates Performs a query in blockchain for certificates with some state
// @Summary Get Certificates of Accredited
// @Description Get all Certificates that belongs to the specified Accredited
// @Tags Certificate
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param 	accredited  	path 	string 	true	"Person to whom the certificates were emitted"
// @Param 	page_limit		query 	int 	true	"Amount of assets per page" default(5)"
// @Param   bookmark        query   string  false   "Bookmark to know last asset gotten"
// @Param   channel         query   string  true    "Insert channel" default(mychannel)"
// @Param   chaincode       query   string  true    "Insert chaincode id" default(certificate)"
// @Param   signer          query   string  true    "Insert signer" default(User1)"
// @Success 200 {object} dto.QueryResult "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 502 {object} dto.Problem "err.bad_gateway"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /dapp/certificates_by_accredited/{accredited} [get]
func (h DappHandler) getCertificatesByAccredited(ctx *context.Context) {
	accredited := ctx.Params().GetStringDefault("accredited", "")
	pageLimit := ctx.URLParamIntDefault("page_limit", -1)

	qp := new(dto.QueryParamChaincode)
	err := lib.ParamsToStruct(ctx, qp)
	if err != nil || accredited == "" || pageLimit == -1 {
		h.response.ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: schema.ErrDetInvalidField}, &ctx)
		return
	}
	qp.PageLimit = pageLimit

	bcRes, problem := (*h.service).GetAssetsByAccredited(accredited, schema.GuestUser, qp)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}
	(*h.response).ResOKWithData(bcRes, &ctx)
}

// region ======== LOCAL DEPENDENCIES ====================================================

// endregion =============================================================================
