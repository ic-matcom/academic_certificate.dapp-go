package service

import (
	"dapp/lib"
	"dapp/repo"
	"dapp/schema"
	"dapp/schema/dto"
	"dapp/schema/mapper"
	"dapp/schema/models"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
)

// region ======== SETUP =================================================================

// ISvcDapp Dapp request service interface
type ISvcDapp interface {
	Query(query dto.Transaction, did string) (interface{}, *dto.Problem)
	Invoke(req dto.Transaction, did string) (interface{}, *dto.Problem)
	GetAsset(id string, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem)
	GetAssetsByState(status int, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem)
	GetAssetsByAccredited(accredited string, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem)
	CreateAsset(req *dto.CreateAsset, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem)
	UpdateAsset(req *dto.Asset, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem)
	ValidateAsset(req *dto.SignAsset, userParam *dto.InjectedParam, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem)
	InvalidateAsset(req *dto.InvalidateAsset, userParam *dto.InjectedParam, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem)
	DeleteAsset(id string, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem)
}

type svcDapp struct {
	repoDapp *repo.RepoDapp
}

// endregion =============================================================================

// NewSvcDappReqs instantiate the Dapp request services
func NewSvcDappReqs(repoDapp *repo.RepoDapp) ISvcDapp {
	return &svcDapp{repoDapp}
}

// region ======== METHODS ======================================================

func (s *svcDapp) Query(query dto.Transaction, did string) (interface{}, *dto.Problem) {
	// requesting blockchain ledger
	raw, e := s.repoDapp.Query(query, did)
	if e != nil {
		return nil, lib.NewProblem(iris.StatusBadGateway, schema.ErrBlockchainTxs, e.Error())
	}

	result := mapper.DecodePayload(raw)
	return result, nil
}

func (s *svcDapp) Invoke(req dto.Transaction, did string) (interface{}, *dto.Problem) {
	// requesting blockchain ledger
	result, e := (*s.repoDapp).Invoke(req, did)
	if e != nil {
		return nil, lib.NewProblem(iris.StatusBadGateway, schema.ErrBlockchainTxs, e.Error())
	}
	dPayload := mapper.DecodePayload(result)
	qResult := dto.TxReceipt{
		ReplyCommon: dto.ReplyCommon{Headers: dto.ReplyHeaders{
			CommonHeaders: req.Headers.CommonHeaders,
			Received:      "",
			Elapsed:       0,
			ReqOffset:     "",
			ReqID:         "",
		}},
		ResponsePayload: dPayload,
	}

	return qResult, nil
}

func (s *svcDapp) GenericGetAssets(payload interface{}, funcName string, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem) {
	var b map[string]interface{}
	err := mapstructure.Decode(payload, &b)
	if err != nil {
		return nil, lib.NewProblem(iris.StatusBadRequest, schema.ErrJsonParse, err.Error())
	}

	tx := dto.Transaction{
		RequestCommon: dto.RequestCommon{Headers: dto.RequestHeaders{CommonHeaders: dto.CommonHeaders{
			PayloadType:  "object",
			Signer:       queryParams.Signer,
			ChannelID:    queryParams.Channel,
			ChaincodeID:  queryParams.Chaincode,
			ContractName: "",
		}}},
		Function:   funcName,
		Payload:    b,
		StrongRead: false,
	}

	// requesting blockchain ledger
	result, e := (*s.repoDapp).Query(tx, did)
	if e != nil {
		return nil, lib.NewProblem(iris.StatusBadGateway, schema.ErrBlockchainTxs, e.Error())
	}
	dPayload := mapper.DecodePayload(result)
	qResult := dto.TxReceipt{
		ReplyCommon: dto.ReplyCommon{Headers: dto.ReplyHeaders{
			CommonHeaders: tx.Headers.CommonHeaders,
		}},
		ResponsePayload: dPayload,
	}
	return qResult, nil
}

func (s *svcDapp) GetAsset(id string, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem) {
	return s.GenericGetAssets(&dto.GetRequestCC{ID: id}, schema.ReadAsset, did, queryParams)
}

func (s *svcDapp) GetAssetsByState(status int, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem) {
	payload := dto.PayloadGetByStatus{
		QueryString: dto.QueryStringGetByStatus{
			Selector: dto.SelectorGetByStatus{
				DocType:           schema.DocType,
				CertificateStatus: status,
			},
		},
		PageSize: queryParams.PageLimit,
		Bookmark: queryParams.Bookmark,
	}
	return s.GenericGetAssets(payload, schema.QueryAssetsWithPag, did, queryParams)
}

func (s *svcDapp) GetAssetsByAccredited(accredited string, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem) {
	payload := dto.PayloadGetByAccredited{
		QueryString: dto.QueryStringGetByAccredited{
			Selector: dto.SelectorGetByAccredited{
				DocType:    schema.DocType,
				Accredited: accredited,
			},
		},
		PageSize: queryParams.PageLimit,
		Bookmark: queryParams.Bookmark,
	}
	return s.GenericGetAssets(payload, schema.QueryAssetsWithPag, did, queryParams)
}

func (s *svcDapp) CreateAsset(req *dto.CreateAsset, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem) {
	asset := mapper.MapCreateAsset2Asset(req)
	// Assets have ID <Code>+<Year>+<Month>+<Day>+<Hour>+<Minute>+<Second> of the time when were created
	asset.ID = fmt.Sprintf("%s%s", schema.DocType, time.Now().Format("20060102150405"))
	asset.Status = dto.New
	b, _ := lib.ToMap(asset, "json")

	tx := dto.Transaction{
		RequestCommon: dto.RequestCommon{Headers: dto.RequestHeaders{CommonHeaders: dto.CommonHeaders{
			PayloadType:  "object",
			Signer:       queryParams.Signer,
			ChannelID:    queryParams.Channel,
			ChaincodeID:  queryParams.Chaincode,
			ContractName: "",
		}}},
		Function:   schema.CreateAsset,
		Payload:    b,
		StrongRead: false,
	}
	// requesting blockchain ledger
	result, e := (*s.repoDapp).Invoke(tx, did)
	if e != nil {
		return nil, lib.NewProblem(iris.StatusBadGateway, schema.ErrBlockchainTxs, e.Error())
	}
	dPayload := mapper.DecodePayload(result)
	qResult := dto.TxReceipt{
		ReplyCommon: dto.ReplyCommon{Headers: dto.ReplyHeaders{
			CommonHeaders: tx.Headers.CommonHeaders,
		}},
		ResponsePayload: dPayload,
	}
	return qResult, nil
}

func (s *svcDapp) UpdateAsset(req *dto.Asset, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem) {
	req.DocType = schema.DocType
	b, _ := lib.ToMap(req, "json")

	tx := dto.Transaction{
		RequestCommon: dto.RequestCommon{Headers: dto.RequestHeaders{CommonHeaders: dto.CommonHeaders{
			PayloadType:  "object",
			Signer:       queryParams.Signer,
			ChannelID:    queryParams.Channel,
			ChaincodeID:  queryParams.Chaincode,
			ContractName: "",
		}}},
		Function:   schema.UpdateAsset,
		Payload:    b,
		StrongRead: false,
	}
	// requesting blockchain ledger
	result, e := (*s.repoDapp).Invoke(tx, did)
	if e != nil {
		return nil, lib.NewProblem(iris.StatusBadGateway, schema.ErrBlockchainTxs, e.Error())
	}
	dPayload := mapper.DecodePayload(result)
	qResult := dto.TxReceipt{
		ReplyCommon: dto.ReplyCommon{Headers: dto.ReplyHeaders{
			CommonHeaders: tx.Headers.CommonHeaders,
		}},
		ResponsePayload: dPayload,
	}
	return qResult, nil
}

func (s *svcDapp) ValidateAsset(req *dto.SignAsset, userParam *dto.InjectedParam, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem) {

	valAsset := dto.ValidateAsset{
		ID:        req.ID,
		Validator: req.SignedBy,
	}
	if userParam.Role == models.Role_Secretary {
		valAsset.ValidatorT = dto.Secretary
	} else if userParam.Role == models.Role_Dean {
		valAsset.ValidatorT = dto.Dean
	} else if userParam.Role == models.Role_Rector {
		valAsset.ValidatorT = dto.Rector
	} else {
		return nil, lib.NewProblem(iris.StatusUnauthorized, "User have no permission to validate the certificate", "")
	}

	b, _ := lib.ToMap(&valAsset, "json")

	tx := dto.Transaction{
		RequestCommon: dto.RequestCommon{Headers: dto.RequestHeaders{CommonHeaders: dto.CommonHeaders{
			PayloadType:  "object",
			Signer:       queryParams.Signer,
			ChannelID:    queryParams.Channel,
			ChaincodeID:  queryParams.Chaincode,
			ContractName: "",
		}}},
		Function:   schema.ValidateAsset,
		Payload:    b,
		StrongRead: false,
	}
	// requesting blockchain ledger
	result, e := (*s.repoDapp).Invoke(tx, userParam.Username)
	if e != nil {
		return nil, lib.NewProblem(iris.StatusBadGateway, schema.ErrBlockchainTxs, e.Error())
	}
	dPayload := mapper.DecodePayload(result)
	qResult := dto.TxReceipt{
		ReplyCommon: dto.ReplyCommon{Headers: dto.ReplyHeaders{
			CommonHeaders: tx.Headers.CommonHeaders,
		}},
		ResponsePayload: dPayload,
	}
	return qResult, nil
}

func (s *svcDapp) InvalidateAsset(req *dto.InvalidateAsset, userParam *dto.InjectedParam, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem) {
	b, _ := lib.ToMap(req, "json")

	tx := dto.Transaction{
		RequestCommon: dto.RequestCommon{Headers: dto.RequestHeaders{CommonHeaders: dto.CommonHeaders{
			PayloadType:  "object",
			Signer:       queryParams.Signer,
			ChannelID:    queryParams.Channel,
			ChaincodeID:  queryParams.Chaincode,
			ContractName: "",
		}}},
		Function:   schema.InvalidateAsset,
		Payload:    b,
		StrongRead: false,
	}
	// requesting blockchain ledger
	result, e := (*s.repoDapp).Invoke(tx, userParam.Username)
	if e != nil {
		return nil, lib.NewProblem(iris.StatusBadGateway, schema.ErrBlockchainTxs, e.Error())
	}
	dPayload := mapper.DecodePayload(result)
	qResult := dto.TxReceipt{
		ReplyCommon: dto.ReplyCommon{Headers: dto.ReplyHeaders{
			CommonHeaders: tx.Headers.CommonHeaders,
		}},
		ResponsePayload: dPayload,
	}
	return qResult, nil
}

func (s *svcDapp) DeleteAsset(id string, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem) {
	b, _ := lib.ToMap(&dto.GetRequestCC{ID: id}, "json")
	tx := dto.Transaction{
		RequestCommon: dto.RequestCommon{Headers: dto.RequestHeaders{CommonHeaders: dto.CommonHeaders{
			PayloadType:  "object",
			Signer:       queryParams.Signer,
			ChannelID:    queryParams.Channel,
			ChaincodeID:  queryParams.Chaincode,
			ContractName: "",
		}}},
		Function:   schema.DeleteAsset,
		Payload:    b,
		StrongRead: false,
	}
	// requesting blockchain ledger
	result, e := (*s.repoDapp).Invoke(tx, did)
	if e != nil {
		return nil, lib.NewProblem(iris.StatusBadGateway, schema.ErrBlockchainTxs, e.Error())
	}
	dPayload := mapper.DecodePayload(result)
	qResult := dto.TxReceipt{
		ReplyCommon: dto.ReplyCommon{Headers: dto.ReplyHeaders{
			CommonHeaders: tx.Headers.CommonHeaders,
		}},
		ResponsePayload: dPayload,
	}
	return qResult, nil
}
