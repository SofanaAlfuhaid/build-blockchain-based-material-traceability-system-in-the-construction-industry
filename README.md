#Build Blockchain-based material traceability platform

This project show how blockchain can be used to improve material traceability in the construction industry. The use case here is the typical construction supply chain. In this project, we have three participants, construction company, supplier and regulator. Each of the participants represents a peer in the network, where they can gain their identity and deploy the smart contract.  The Construction company make the purchase order, received the order, store it, and consume it. The Supplier provides the construction company with requested materials.  The regulator monitors materials movement also inspects to check if materials as per requirements and regulations.

Note: This code runs locally, and it has frontend, backend implementation, and blockchain Explorer in another repos. 

Prerequisites:

A sample Hyperledger Fabric Binaries and Docker Images is downloaded, and the developer builds on top of them and customize a solution based on requirements.

Hardware Requirements:

–	PC – this project uses MacBook Pro, processor:  2.7 GHz Intel Core i5,  Memory:  8 GB 1867 MHz DDR3

Software Requirements:

–	Hyperledger Fabric v1.4.2

In order to develop or operate Hyperledger Fabric, the following prerequisites must be installed in the platform operating system:

–	Docker and Docker Compose – v19.03.8
–	cURL - latest
–	NPM – latest
–	nvm - latest
–	Node.js - latest
–	Python  - v2.7.x
– Go - v1.13

Also, to develop and test the platform and the smart contract:

–	Code editor - Visual Studio Code version 1.28, or higher.

For building user interfaces:

– React - v16.13.1

For the middleware:

–	Nodes.js - v8.16.0
–	Postman – API client
–	MongoDB - v4.0.18

For Hyperledger Explorer:

–	PostgreSQL - v12.3
– jq - v1.6


Here some useful instructions to run the blockchain network, more explanation and a demo will be added soon.

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
