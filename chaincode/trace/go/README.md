## 22-April-2020
### channel list
peer channel list

### channel create
```sh
export CHANNEL_NAME=mychannel && peer channel create -o orderer.track.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/$CHANNEL_NAME.tx --tls --cafile $ORDERER_CA

```

### channel join
```sh
peer channel join -b mychannel.block

```


### chaincode install
```sh

peer chaincode install -n chaincode_name -v 1.0 -p github.com/chaincode/trace/go/

```
### chaincode instantiate
```sh

peer chaincode instantiate -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C channel_name -n chaincode_name -v version -c '{"Args":["init"]}' -P "OR ('CSTMSP.peer','SUPMSP.peer', 'REGMSP.peer')"

```

### chaincode upgrade 
```sh

peer chaincode upgrade -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C channel_name -n chaincode_name -v version -c '{"Args":["init"]}' -P "OR ('CSTMSP.peer','SUPMSP.peer', 'REGMSP.peer')"

```

### check running containers
```sh

docker ps

```
### check running container logs

```sh
docker logs container_name -f --tail 1000
```

### chaincode functions

PurchaseOrder
```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n track -c '{"args":["purchaseOrder","{\"po\":\"QT041111111111111151\",\"posts\":\"create\",\"itemno\":\"111101\",\"itemname\":\"Cement\",\"desc\":\"poatlad cement 40 kg\",\"quan\":8,\"uprice\":\"100\",\"addr\":\"c-30\",\"delvry\":\"10-06-2020\",\"buyerid\":\"B0001\",\"suppid\":\"S0001\",\"cts\":\"456789\",\"uts\":\"456787678\",\"amt\":\"800\"}"]}'
```
supplierRecOrderSts
```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n track -c '{"args":["supplierRecOrderSts","{\"po\":\"QT041111111111111151\",\"posts\":\"inProgress\",\"uts\":\"456787678\"}"]}'
```

supplierOrderCorrection
```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n track -c '{"args":["supplierOrderCorrection","{\"po\":\"QT041111111111111151\",\"posts\":\"inProgress\",\"uts\":\"456787678\"}"]}'
```


createOrderBySupplier
```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n mytrack -c '{"args":["createOrderBySupplier","{\"po\":\"QT041111111111111151\",\"cid\":\"S0001\",\"shid\":\"DO0001\",\"trno\":\"T001\",\"regid\":\"R001\",\"dosts\":\"expecting confirmation from regulator\",\"gtin\":\"10614141999996\",\"uts\":\"456787678\"}"]}'
```

logisticApproval
```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n mytrack -c '{"args":["logisticApproval","{\"po\":\"QT041111111111111151\",\"dosts\":\"shipped\", \"uts\":\"456787678\"}"]}'
```
inventoryManagerReceipt
```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n mytrack -c '{"args":["inventoryManagerReceipt","{\"po\":\"QT041111111111111151\",\"invmngid\":\"IM001\",\"expdate\":\"10-06-2020\", \"gis\":\"zone1\",\"grept\":\"GRxxxxxx\",\"grsts\":\"pending\", \"uts\":\"456787678\"}"]}'
```

inventoryApproval
```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n trackk -c '{"args":["inventoryApproval","{\"po\":\"QT041111111111112222\",\"dosts\":\"arrived\",\"posts\":\"inStock\", \"grsts\":\"received\", \"uts\":\"456787678\", \"stanbathweght\":320, \"itemname\":\"Cement\"}"]}'
```

FormenConsumption
```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n track -c '{"args":["formenConsumption","{\"po\":\"QT041111111111111131\",\"foreid\":\"F0001\",\"purps\":\"for concreate making\",\"fdec\":\"same as desc\", \"ccorder\":\"created\",\"conum\":\"0004\",\"pquanty\":4, \"futs\":\"456787678\"}"]}'
```

stockRelease
```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n track -c '{"args":["stockRelease","{\"po\":\"QT041111111111111131\",\"ccorder\":\"expecting confirmation from regulator\", \"futs\":\"456787678\",\"batchid\":\"b001\",\"conum\":\"0001\"}"]}'
```

consumptionApproval
```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n track -c '{"args":["consumptionApproval","{\"po\":\"QT041111111111111131\",\"ccorder\":\"ready to use\",\"bweght\":\"1000\", \"futs\":\"456787678\",\"conum\":\"0001\"}"]}'
```
displayOrderStatus
```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n track -c '{"args":["displayOrderStatus","{\"po\":\"QT041111111111111151\",\"ccorder\":\"expecting confirmation from regulator for pouring\", \"futs\":\"456787678\",\"conum\":\"0001\"}"]}'
```

consumptionApprovalForPouring
```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n track -c '{"args":["consumptionApprovalForPouring","{\"po\":\"QT041111111111111151\",\"ccorder\":\"ready to be poured\",\"density\":\"2.03\", \"futs\":\"456787678\",\"conum\":\"0001\"}"]}'
```

materialQuery

```sh
peer chaincode invoke -o orderer.track.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n track -c '{"args":["materialQuery","10614141999996"]}'
```


