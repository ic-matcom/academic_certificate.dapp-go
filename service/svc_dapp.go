package service

import (
	"dapp/lib"
	"dapp/repo"
	"dapp/schema"
	"dapp/schema/dto"
	"dapp/schema/mapper"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
)

// region ======== SETUP =================================================================

// ISvcDapp Dapp request service interface
type ISvcDapp interface {
	Query(query dto.Transaction, did string) (interface{}, *dto.Problem)
	Invoke(req dto.Transaction, did string) (interface{}, *dto.Problem)
	CreateAsset(req dto.CreateAsset, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem)
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

func (s *svcDapp) CreateAsset(req dto.CreateAsset, did string, queryParams *dto.QueryParamChaincode) (interface{}, *dto.Problem) {
	asset := mapper.MapCreateAsset2Asset(req)
	// Assets have ID <Code>+<Year>+<Month>+<Day>+<Hour>+<Minute>+<Second> of the time when were created
	asset.ID = fmt.Sprintf("%s%s", schema.DocType, time.Now().Format("20060102150405"))
	asset.Status = dto.New
	b, _ := lib.ToMap(&asset, "json")

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
