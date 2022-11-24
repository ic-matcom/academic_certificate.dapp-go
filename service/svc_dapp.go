package service

import (
	"dapp/lib"
	"dapp/repo"
	"dapp/schema"
	"dapp/schema/dto"
	"dapp/schema/mapper"
	"github.com/kataras/iris/v12"
)

// region ======== SETUP =================================================================

// ISvcDapp Dapp request service interface
type ISvcDapp interface {
	Query(query dto.Transaction, did string) (interface{}, *dto.Problem)
	Invoke(req dto.Transaction, did string) (interface{}, *dto.Problem)
	CreateAsset(req dto.Asset, did, channel, chaincode, signer string) (interface{}, *dto.Problem)
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

func (s *svcDapp) CreateAsset(req dto.Asset, did, channel, chaincode, signer string) (interface{}, *dto.Problem) {

	b, _ := lib.ToMap(&req, "json")

	tx := dto.Transaction{
		RequestCommon: dto.RequestCommon{Headers: dto.RequestHeaders{CommonHeaders: dto.CommonHeaders{
			PayloadType:  "object",
			Signer:       signer,
			ChannelID:    channel,
			ChaincodeID:  chaincode,
			ContractName: "",
		}}},
		Function:   "CreateAsset", //TODO: esto si puede quedar anclado en el codigo, pero recomiendo que muevas todas TxName al schema/constants.go o un fichero .go similar
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
