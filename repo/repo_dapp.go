package repo

import (
	"dapp/lib"
	"dapp/schema"
	"dapp/schema/dto"
	"dapp/service/utils"
	"errors"
	"fmt"
	"github.com/cloudflare/cfssl/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	jsoniter "github.com/json-iterator/go"
	"path/filepath"
	"sync"
)

// region ======== SETUP =================================================================

// the ccpClientWrapper implements the RPCClient interface using the fabric-sdk-go implementation
// based on the static network description provided via the CCP yaml
type ccpClientWrapper struct {
	channelClient   *channel.Client
	channelProvider context.ChannelProvider
	signer          *msp.IdentityIdentifier
}

type RepoDapp struct {
	ChannelName       string          // ChannelName HLF channel name
	CppPath           string          // CppPath path to the connection profile
	WalletPath        string          // WalletPath path to the wallets folders
	Wallet            *gateway.Wallet // Wallet with admin privilege identity for admins ops on the network
	DappIdentityUser  string          // DappIdentityUser dapp user identity to authenticate normal dapp ops in the HLF network
	DappIdentityAdmin string          // DappIdentityAdmin dapp admin identity to authenticate admin dapp ops in the HLF network

	configProvider core.ConfigProvider
	sdk            *fabsdk.FabricSDK
	channelCreator channelCreator
	channelClient  ccpClientWrapper
}

var singleton *RepoDapp

// using Go sync package to invoke a method exactly only once
var onceRepoDapp sync.Once

// endregion =============================================================================

func NewRepoDapp(svcConf *utils.SvcConfig) *RepoDapp {
	onceRepoDapp.Do(func() {
		wallet, err := gateway.NewFileSystemWallet(filepath.Join(svcConf.WalletFolder, schema.WalletStr))
		if err != nil {
			panic(schema.ErrDetWalletProc + " ." + err.Error())
		}
		exist, err := lib.FileExists(svcConf.CppPath)
		if err != nil || !exist {
			panic(fmt.Errorf("fabric network profile config not found"))
		}

		configProvider := config.FromFile(filepath.Clean(svcConf.CppPath))
		sdk, err := fabsdk.New(configProvider)
		if err != nil {
			panic(schema.ErrDetSDKInit + " ." + err.Error())
		}

		exist, err = lib.FileExists(svcConf.WalletFolder)
		if err != nil || !exist {
			panic(fmt.Errorf("wallet folder not found, check in the dapp configuration the parameter \"WalletFolder\""))
		}

		singleton = &RepoDapp{
			CppPath:           svcConf.CppPath,
			WalletPath:        svcConf.WalletFolder,
			Wallet:            wallet,
			DappIdentityUser:  svcConf.DappIdentityUser,
			DappIdentityAdmin: svcConf.DappIdentityAdmin,

			configProvider: configProvider,
			sdk:            sdk,
			channelCreator: createChannelClient,
			channelClient:  ccpClientWrapper{},
		}
	})
	return singleton
}

// region ======== METHODS ===============================================================

func (r *RepoDapp) Query(query dto.Transaction, did string) ([]byte, error) {
	var err error
	var args_ []string

	if query.IsSchema {
		// if a isSchema is true, the payload property in the body must be a JSON structure
		argsMap, ok := query.Payload.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("the \"payload\" property must be JSON if a isSchema property is true")
		}

		res, err := jsoniter.MarshalToString(argsMap)
		if err != nil {
			return nil, err
		}
		args_ = append(args_, res)
	} else {
		argVals, ok := query.Payload.([]interface{})
		if !ok {
			return nil, fmt.Errorf("no payload schema is specified in the payload's \"headers\", the \"args\" property must be an array of strings")
		}

		args_ = make([]string, len(argVals))
		for i, v := range query.Payload.([]interface{}) {
			args_[i] = v.(string)
		}
	}

	if query.StrongRead {
		// getting bc components instance
		_, contract, err := r.getSDKComponents(query, false)
		if err != nil {
			return nil, err
		}

		// Getting the identities list
		res, err := contract.EvaluateTransaction(query.Function, args_...)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	peerEndpoint, err := getFirstPeerEndpointFromConfig(r.configProvider)
	if err != nil {
		return nil, err
	}

	req := channel.Request{
		ChaincodeID: query.Headers.ContractName,
		Fcn:         query.Function,
		Args:        convert(args_),
	}

	//TODO: then move to struct
	// prepare contexts
	// TODO: move to getChannelClient func (uncreated)
	channelContext := r.sdk.ChannelContext(query.Headers.ChannelID, fabsdk.WithUser(r.DappIdentityUser))//, fabsdk.WithOrg("Org1MSP"))

	cClient, _ := r.channelCreator(channelContext)
	if err != nil {
		return nil, err
	}
	// TODO: move to getChannelClient func (uncreated)
	r.channelClient.channelClient = cClient
	r.channelClient.channelProvider = channelContext

	result, err := r.channelClient.channelClient.Query(req, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(peerEndpoint))
	if err != nil {
		log.Errorf("Failed to send query [%s:%s:%s]. %s", query.Headers.ChannelID, query.Headers.ContractName, query.Function, err)
		return nil, err
	}

	return result.Payload, nil
}

func (r *RepoDapp) Invoke(query dto.Transaction, did string) ([]byte, error) {
	var err error
	var args_ []string

	if query.IsSchema {
		// if a isSchema is true, the payload property in the body must be a JSON structure
		argsMap, ok := query.Payload.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("the \"payload\" property must be JSON if a isSchema property is true")
		}

		res, err := jsoniter.MarshalToString(argsMap)
		if err != nil {
			return nil, err
		}
		args_ = append(args_, res)
	} else {
		argVals, ok := query.Payload.([]interface{})
		if !ok {
			return nil, fmt.Errorf("no payload schema is specified in the payload's \"headers\", the \"args\" property must be an array of strings")
		}

		args_ = make([]string, len(argVals))
		for i, v := range query.Payload.([]interface{}) {
			args_[i] = v.(string)
		}
	}

	if query.StrongRead {
		// getting bc components instance
		_, contract, err := r.getSDKComponents(query, false)
		if err != nil {
			return nil, err
		}

		// invoking the contract
		res, err := contract.SubmitTransaction(query.Function, args_...)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	peerEndpoint, err := getFirstPeerEndpointFromConfig(r.configProvider)
	if err != nil {
		return nil, err
	}

	req := channel.Request{
		ChaincodeID: query.Headers.ContractName,
		Fcn:         query.Function,
		Args:        convert(args_),
	}

	//TODO: then move to struct
	// prepare contexts
	// TODO: move to getChannelClient func (uncreated)
	channelContext := r.sdk.ChannelContext(query.Headers.ChannelID, fabsdk.WithUser(r.DappIdentityUser))//, fabsdk.WithOrg("Org1MSP"))

	cClient, _ := r.channelCreator(channelContext)
	if err != nil {
		return nil, err
	}
	// TODO: move to getChannelClient func (uncreated)
	r.channelClient.channelClient = cClient
	r.channelClient.channelProvider = channelContext

	result, err := r.channelClient.channelClient.Query(req, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(peerEndpoint))
	if err != nil {
		log.Errorf("Failed to send query [%s:%s:%s]. %s", query.Headers.ChannelID, query.Headers.ContractName, query.Function, err)
		return nil, err
	}

	return result.Payload, nil

}

type channelCreator func(context.ChannelProvider) (*channel.Client, error)

func createChannelClient(channelProvider context.ChannelProvider) (*channel.Client, error) {
	return channel.New(channelProvider)
}

func (r *RepoDapp) getSDKComponents(query dto.Transaction, withAdminIdentity bool) (*gateway.Network, *gateway.Contract, error) {
	chID := query.Headers.ChannelID
	cName := query.Headers.ContractName

	var identityLabel = r.DappIdentityUser

	if withAdminIdentity {
		identityLabel = r.DappIdentityAdmin
	}

	if !r.Wallet.Exists(identityLabel) {
		return nil, nil, fmt.Errorf("the %s identity not exist in wallet: %s", identityLabel, r.WalletPath)
	}
	// trying to get an instance of HLF SDK network gateway, from the connection profile
	gw, err := gateway.Connect( // gt = gateway
		gateway.WithConfig(r.configProvider),
		gateway.WithIdentity(r.Wallet, identityLabel),
		gateway.WithSDK(r.sdk),
		)
	if err != nil {
		return nil, nil, err
	}
	defer gw.Close()

	// trying to get an instance of the gateway network
	nt, e := gw.GetNetwork(chID) // nt == network
	if e != nil {
		return nil, nil, e
	}
	// trying to get the contract
	contract := nt.GetContractWithName(cName, "")
	if contract == nil {
		return nil, nil, errors.New(schema.ErrDetContractNotFound)
	}

	// so far so good, returning the instance pointers
	return nt, contract, nil
}

func getOrgFromConfig(config core.ConfigProvider) (string, error) {
	configBackend, err := config()
	if err != nil {
		return "", err
	}
	if len(configBackend) != 1 {
		return "", fmt.Errorf("invalid config file")
	}

	cfg := configBackend[0]
	value, ok := cfg.Lookup("client.organization")
	if !ok {
		return "", fmt.Errorf("no client organization defined in the config")
	}

	return value.(string), nil
}

func getFirstPeerEndpointFromConfig(config core.ConfigProvider) (string, error) {
	org, err := getOrgFromConfig(config)
	if err != nil {
		return "", err
	}
	configBackend, _ := config()
	cfg := configBackend[0]
	value, ok := cfg.Lookup(fmt.Sprintf("organizations.%s.peers", org))
	if !ok {
		return "", fmt.Errorf("no peers list found in the organization %s", org)
	}
	peers := value.([]interface{})
	if len(peers) < 1 {
		return "", fmt.Errorf("peers list for organization %s is empty", org)
	}
	return peers[0].(string), nil
}

func convert(args []string) [][]byte {
	result := [][]byte{}
	for _, v := range args {
		result = append(result, []byte(v))
	}
	return result
}
// region ======== Dapp ======================================================


// endregion =============================================================================
