package endpoints

import (
	"dapp/lib"
	"dapp/repo"
	"dapp/schema"
	"dapp/schema/dto"
	"dapp/service"
	"dapp/service/utils"
	"encoding/json"
	"fmt"

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
		// registering protected / guarded router
		protectedAPI := v1.Party("/dapp")
		{
			// --- GROUP / PARTY MIDDLEWARES ---
			protectedAPI.Use(*mdwAuthChecker)

			protectedAPI.Post("/query", hero.Handler(h.postQuery))
			protectedAPI.Post("/transaction", hero.Handler(h.postTransaction))

			protectedAPI.Get("/certificates_by_state/{state: int}", hero.Handler(h.getCertificatesByState))
			protectedAPI.Get("/certificates_by_accredited/{accredited: string}", hero.Handler(h.getCertificatesByAccredited))
			protectedAPI.Post("/create_asset", hero.Handler(h.postCreateAsset))
			protectedAPI.Post("/validate_asset", hero.Handler(h.postValidateAsset))
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
// @Summary Send transaction to peers
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
// @Router /dapp/create_asset [post]
func (h DappHandler) postCreateAsset(ctx iris.Context, params dto.InjectedParam) {
	queryParams := new(dto.QueryParamChaincode)
	lib.ParamsToStruct(ctx, queryParams)

	var requestData dto.CreateAsset
	// unmarshalling the json and check
	if err := ctx.ReadJSON(&requestData); err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}
	// trying to submit the transaction
	bcRes, problem := (*h.service).CreateAsset(requestData, params.Username, queryParams)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}

	(*h.response).ResOKWithData(bcRes, &ctx)
}

// postValidateAsset Validate Asset in ledger
// @Summary Send transaction to peers
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
// @Router /dapp/validate_asset [post]
func (h DappHandler) postValidateAsset(ctx iris.Context, params dto.InjectedParam) {
	queryParams := new(dto.QueryParamChaincode)
	lib.ParamsToStruct(ctx, queryParams)

	var requestData dto.SignAsset
	// unmarshalling the json and check
	if err := ctx.ReadJSON(&requestData); err != nil {
		(*h.response).ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}
	// trying to submit the transaction
	bcRes, problem := (*h.service).ValidateAsset(requestData, &params, queryParams)
	if problem != nil {
		(*h.response).ResErr(problem, &ctx)
		return
	}

	(*h.response).ResOKWithData(bcRes, &ctx)
}

// getCertificatesByState Performs a query in blockchain for certificates with some state
// @Summary Performs a query in blockchain for certificates with some state
// @Tags Certificate
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	true 	"Insert access token" default(Bearer <Add access token here>)
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
func (h DappHandler) getCertificatesByState(ctx *context.Context, params dto.InjectedParam) {
	state := ctx.Params().GetIntDefault("state", -1)
	if state == -1 {
		h.response.ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: schema.ErrDetInvalidField}, &ctx)
		return
	}
	qp := new(dto.QueryParamChaincode)
	lib.ParamsToStruct(ctx, qp)

	queryJSON := fmt.Sprintf(`
	{
		"func": "common:QueryAssetsWithPagination",
		"headers": {
		  "chaincode": "%s",
		  "channel": "%s",
		  "contractName": "",
		  "payloadType": "object",
		  "signer": "%s"
		},
		"payload": {"queryString":{"selector":{"docType":"CERT", "certificate_status":%d}}, "pageSize":%d, "bookmark":"%s"},
		"strongRead": false
	}`, qp.Chaincode, qp.Channel, qp.Signer, state, qp.PageLimit, qp.Bookmark)
	var query dto.Transaction

	// unmarshalling the json and check
	if err := json.Unmarshal([]byte(queryJSON), &query); err != nil {
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

// getNewCertificates Performs a query in blockchain for certificates with some state
// @Summary Performs a query in blockchain for certificates with some state
// @Tags Certificate
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	true 	"Insert access token" default(Bearer <Add access token here>)
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
func (h DappHandler) getCertificatesByAccredited(ctx *context.Context, params dto.InjectedParam) {
	accredited := ctx.Params().GetDefault("accredited", "")
	if accredited == "" {
		h.response.ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: schema.ErrDetInvalidField}, &ctx)
		return
	}

	qp := new(dto.QueryParamChaincode)
	lib.ParamsToStruct(ctx, qp)

	queryJSON := fmt.Sprintf(`
	{
		"func": "common:QueryAssetsWithPagination",
		"headers": {
		  "chaincode": "%s",
		  "channel": "%s",
		  "contractName": "",
		  "payloadType": "object",
		  "signer": "%s"
		},
		"payload": {"queryString":{"selector":{"docType":"CERT", "accredited":%d}}, "pageSize":%d, "bookmark":"%s"},
		"strongRead": false
	}`, qp.Chaincode, qp.Channel, qp.Signer, accredited, qp.PageLimit, qp.Bookmark)

	var query dto.Transaction
	// unmarshalling the json and check
	if err := json.Unmarshal([]byte(queryJSON), &query); err != nil {
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

// region ======== LOCAL DEPENDENCIES ====================================================

// endregion =============================================================================
