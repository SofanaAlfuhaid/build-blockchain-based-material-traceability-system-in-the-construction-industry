Build Blockchain-based material traceability platform

This project show how blockchain can be used to improve material traceability in the construction industry. The use case here is the typical construction supply chain. In this project, we have three participants, construction company, supplier and regulator. Each of the participants represents a peer in the network, where they can gain their identity and deploy the smart contract.  The Construction company make the purchase order, received the order, store it, and consume it. The Supplier provides the construction company with requested materials.  The regulator monitors materials movement also inspects to check if materials as per requirements and regulations.



### channel list
peer channel list

### create the mychannel.tx file
```sh
export CHANNEL_NAME=mychannel && configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME
```
### create orderer genesis block
```sh
../bin/configtxgen -profile TwoOrgsOrdererGenesis -channelID byfn-sys-channel -outputBlock ./channel-artifacts/genesis.block

```

### running the network
```sh

docker-compose -f docker-compose-cli.yaml up -d

```

### check running containers
```sh

docker ps

```
### create the channel iteslf

```sh
docker exec -it cli_cst bash
export CHANNEL_NAME=chentity && peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/$CHANNEL_NAME.tx --tls --cafile $ORDERER_CA
```

### joining the channel

```sh
peer channel join -b mychannel.block

```
