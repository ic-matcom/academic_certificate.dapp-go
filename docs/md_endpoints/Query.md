# example requests

## ReadAsset transaction example
> Query ReadAsset transaction of basic chaincode `https://github.com/kmilodenisglez/fabric-testnet-nano-without-syschannel/tree/main/chaincodes-external/cc-assettransfer-go` 
```json
{
  "func": "ReadAsset",
  "headers": {
    "channel": "mychannel",
    "contractName": "basic",
    "signer": ""
  },
  "payload": ["1"],
  "isSchema": false,
  "strongRead": false
}
```