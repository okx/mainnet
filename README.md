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

wget https://raw.githubusercontent.com/okex/mainnet/main/genesis.json -O ${EXCHAIND_PATH}/config/genesis.json

export EXCHAIN_SEEDS="e926c8154a2af4390de02303f0977802f15eafe2@3.16.103.80:26656,7fa5b1d1f1e48659fa750b6aec702418a0e75f13@35.177.8.240:26656,c8f32b793871b56a11d94336d9ce6472f893524b@18.167.16.85:26656"

exchaind start --chain-id exchain-66 --mempool.sort_tx_by_gp --home ${EXCHAIND_PATH} --p2p.seeds $EXCHAIN_SEEDS
```

Note: it needs to check the [genesis file](https://raw.githubusercontent.com/okex/mainnet/main/genesis.json)

```bash
$ shasum -a 256 ${EXCHAIND_PATH}/config/genesis.json
0958b6c9f5f125d1d6b8f56e042fa8a71b1880310227b8b2f27ba93ff7cd673b  ${EXCHAIND_PATH}/config/genesis.json
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
exchaind start --chain-id exchain-66 --mempool.sort_tx_by_gp --home ${EXCHAIND_PATH}
```

