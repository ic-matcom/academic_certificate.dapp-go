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
		ReplyCommon:     dto.ReplyCommon{Headers: dto.ReplyHeaders{
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