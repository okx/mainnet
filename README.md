# ExChain mainnet

This repo collects the genesis and configuration files for the various ExChain
mainnet. It exists so the [ExChain repo](https://github.com/okex/exchain)
does not get bogged down with large genesis files and status updates.


## Startup an exchain full node by the exchaind binary

### 1. Build exchaind by [the latest released version](https://github.com/okex/exchain/releases/latest)
```
git clone -b latest_version https://github.com/okex/exchain.git  # latest_version refers to https://github.com/okex/exchain/releases/latest
cd exchain
make mainnet
```

### 2. Start full node

```
export EXCHAIND_PATH=~/.exchaind (You can also specify other directory)

exchaind init your_custom_moniker --chain-id exchain-66 --home ${EXCHAIND_PATH}
```

#### 2.1. Using genesis.json file
If you want to synchronize blocks from Genesis block, use the genesis.json file.
```
wget https://raw.githubusercontent.com/okex/mainnet/main/genesis.json -O ${EXCHAIND_PATH}/config/genesis.json

exchaind start --chain-id exchain-66 --home ${EXCHAIND_PATH}
```

Note: it needs to check the [genesis file](https://raw.githubusercontent.com/okex/mainnet/main/genesis.json)

```bash
$ shasum -a 256 ${EXCHAIND_PATH}/config/genesis.json
0958b6c9f5f125d1d6b8f56e042fa8a71b1880310227b8b2f27ba93ff7cd673b  ${EXCHAIND_PATH}/config/genesis.json
```
#### 2.2. Using snapshots
If you want to quickly sync to the latest block, use the snapshot.
[Snapshot link](https://forum.okt.club/d/380-the-latest-snapshots-of-mainnet-and-testnet)

```
# download snapshot
rm -rf ${EXCHAIND_PATH}/data
cd ${EXCHAIND_PATH}
wget https://okg-pub-hk.oss-cn-hongkong.aliyuncs.com/cdn/okc/snapshot/mainnet-$version-$date-$height-rocksdb.tar.gz
tar -zxvf mainnet-$version-$date-$height-rocksdb.tar.gz
```
#### 2.3. Start the node
```
exchaind start --chain-id exchain-66 --home ${EXCHAIND_PATH}
```

## Startup an exchain full node with docker
### 1. make the data dir
```shell
mkdir -p ~/.exchaind/data
echo '{\n"height": "0",\n"round": "0",\n"step": 0\n}' > ~/.exchaind/data/priv_validator_state.json
```

### 2. run docker image
```shell
docker run -d --name exchain-mainnet-fullnode -v ~/.exchaind/data:/root/.exchaind/data/ -p 8545:8545 -p 26656:26656 okexchain/fullnode-mainnet:latest
```

### 3. check log
```shell
docker logs --tail 100 -f exchain-mainnet-fullnode
```

### 4. stop and remove the docker container
```shell
docker rm -f exchain-mainnet-fullnode
```

### 5. restart
You can restart in the previous data dir
```shell
docker run -d --name exchain-mainnet-fullnode -v ~/.exchaind/data:/root/.exchaind/data/ -p 8545:8545 -p 26656:26656 okexchain/fullnode-mainnet:latest
```


## Upgrade an exchain full node to latest

### 1. Stop exchain  full node
### 2. Build exchaind binary
Build exchaind by [the latest released version](https://github.com/okex/exchain/releases/latest)
```
git clone -b latest_version https://github.com/okex/exchain.git  # latest_version refers to https://github.com/okex/exchain/releases/latest
cd exchain
make mainnet
```
### 3. Start
```
export EXCHAIND_PATH=~/.exchaind (You can also specify other directory)
exchaind start --chain-id exchain-66 --home ${EXCHAIND_PATH}
```

