# example invoke request

## CreateAsset with params transaction example
> Query CreateAsset transaction using params in basic chaincode `https://github.com/kmilodenisglez/fabric-testnet-nano-without-syschannel/tree/main/chaincodes-external/cc-assettransfer-go`
```json
{
  "func": "CreateAsset",
  "headers": {
    "channel": "mychannel",
    "contractName": "basic",
    "signer": ""
  },
  "payload": ["1","blue","35","tom","1000"],
  "isSchema": false,
  "strongRead": false
}
```

## CreateAsset with json object transaction example
> Query CreateAsset transaction using a json object  in basic chaincode `https://github.com/kmilodenisglez/fabric-testnet-nano-without-syschannel/tree/main/chaincodes-external/cc-assettransfer-go`
```json
{
  "func": "CreateAssetUsingStructParam",
  "headers": {
    "channel": "mychannel",
    "contractName": "basic",
    "signer": ""
  },
  "payload": {
    "ID": "14",
    "color": "green",
    "size": 101,
    "owner": "kmilo",
    "appraisedValue": 90
   },
  "isSchema": true,
  "strongRead": false
}
```