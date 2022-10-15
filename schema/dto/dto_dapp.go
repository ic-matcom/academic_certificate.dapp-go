package dto

type StatusMsg struct {
	OK bool `json:"ok"`
}

// CommonHeaders are common to all messages
type CommonHeaders struct {
	//ID           string `json:"id,omitempty" example:""`
	//MsgType      string `json:"type,omitempty" example:""`
	Signer       string `json:"signer,omitempty" example:""`
	ChannelID    string `json:"channel" example:"mychannel"`
	ChaincodeID  string `json:"chaincodeID" example:"basic"`
	ContractName string `json:"contractName,omitempty" example:""`
}

// ReplyHeaders are common to all replies
type ReplyHeaders struct {
	CommonHeaders
	Received  string  `json:"timeReceived"`
	Elapsed   float64 `json:"timeElapsed"`
	ReqOffset string  `json:"requestOffset"`
	ReqID     string  `json:"requestId"`
}

// RequestHeaders are common to all requests
type RequestHeaders struct {
	CommonHeaders
}

// RequestCommon is a common interface to all requests
type RequestCommon struct {
	Headers RequestHeaders `json:"headers"`
}

// ReplyCommon is a common interface to all replies
type ReplyCommon struct {
	Headers ReplyHeaders `json:"headers"`
}

type DevPopulateOut struct {
	TxId     string
	Identity string
}

// TxDataRequest standard fronted transaction request structure
type TxDataRequest struct {
	Signature string                 `json:"signature"`
	Payload   map[string]interface{} `json:"payload"`
}

type MapStruct struct {
	RequestCommon
}

type Transaction struct {
	RequestCommon
	Function   string `json:"func" validate:"required"`
	Payload    any    `json:"payload" binding:"required"`
	IsSchema   bool   `json:"isSchema" binding:"required,boolean"`
	StrongRead bool   `json:"strongRead" binding:"required,boolean"`
}

type QueryResult struct {
	ReplyCommon
	Result interface{} `json:"result"`
}

type ChaincodeSpec struct {
	Type    int    `json:"type"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type TxReceipt struct {
	// BlockNumber     uint64              `json:"blockNumber"`
	// SignerMSP       string              `json:"signerMSP"`
	// Signer          string              `json:"signer"`
	// ChaincodeSpec   ChaincodeSpec       `json:"chaincode"`
	// TransactionID   string              `json:"transactionID"`
	// Status          pb.TxValidationCode `json:"status"`
	// SourcePeer      string              `json:"peer"`
	ReplyCommon
	ResponsePayload any `json:"responsePayload"`
}
