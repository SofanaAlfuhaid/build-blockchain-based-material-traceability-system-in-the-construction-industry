## 22-April-2020
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