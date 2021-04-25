# ExChain mainnet

This repo collects the genesis and configuration files for the various ExChain
mainnet. It exists so the [ExChain repo](https://github.com/okex/exchain)
does not get bogged down with large genesis files and status updates.

## Getting Started

To get started with the latest mainnet, see the
[docs](https://okexchain-docs.readthedocs.io/en/latest/getting-start/join-okexchain-mainnet.html).


## Upgrade an exchain full node base on v0.16.3

### 1. Stop exchain  full node ,and start exchain with 'halt-height = 2322600'
### 2. two ways to get the genesis.json file
```
export EXCHAIND_PATH=~/.okexchaind (If your directory is not ~/.okexchaind, specify your own directory)
```
#### 2.1 Export by exchaind
```
   cd ${EXCHAIND_PATH}
   git clone -b v0.16.3.1 https://github.com/okex/exchain.git
   cd exchain
   make install
   okexchaind export --home ${EXCHAIND_PATH} --height=2322600 --for-zero-height --log_level evm:debug --log_file ./export.log --log_stdout=false > ${EXCHAIND_PATH}/config/genesis_no_migrate.json
   git checkout v0.18.1
   make GenesisHeight=2322600 install
   exchaind migrate v0.18 ${EXCHAIND_PATH}/config/genesis_no_migrate.json --chain-id=exchain66 > ${EXCHAIND_PATH}/config/genesis.json
```
#### 2.2 Download genesis.json
```
   wget https://raw.githubusercontent.com/okex/mainnet/main/genesis.json -O ${EXCHAIND_PATH}/config/genesis.json
```
Note: it needs to check genesis.json no matter which way is used
```
$ shasum -a 256 ${EXCHAIND_PATH}/config/genesis.json
0958b6c9f5f125d1d6b8f56e042fa8a71b1880310227b8b2f27ba93ff7cd673b  ${EXCHAIND_PATH}/config/genesis.json
```
### 3. Build exchaind binary
Build exchaind by [the latest released version v0.18.1](https://github.com/okex/exchain/releases/tag/v0.18.1)
```
   cd ${EXCHAIND_PATH}
   git clone -b v0.18.1 https://github.com/okex/exchain.git
   cd exchain
   make GenesisHeight=2322600 install
```
### 4. Reset data
`exchaind unsafe-reset-all --home ${EXCHAIND_PATH}`
### 5. Start
`exchaind start --chain-id exchain-66 --home ${EXCHAIND_PATH}`


## Startup an exchain full node by the genesis.json file
### 1. Build exchaind by [the latest released version v0.18.1](https://github.com/okex/exchain/releases/tag/v0.18.1)
```
export EXCHAIND_PATH=~/.okexchaind (You can also specify other directory)
```

### 2. Start mainnet with the exchaind binary

```
exchaind init your_custom_moniker --chain-id exchain-66 --home ${EXCHAIND_PATH}

wget https://raw.githubusercontent.com/okex/mainnet/main/genesis.json -O ${EXCHAIND_PATH}/config/genesis.json

export EXCHAIN_SEEDS="e926c8154a2af4390de02303f0977802f15eafe2@3.16.103.80:26656,7fa5b1d1f1e48659fa750b6aec702418a0e75f13@35.177.8.240:26656,c8f32b793871b56a11d94336d9ce6472f893524b@18.167.16.85:26656"

exchaind start --chain-id exchain-66 --home ${EXCHAIND_PATH} --p2p.seeds $EXCHAIN_SEEDS
```

Note: it needs to check the [genesis file](https://raw.githubusercontent.com/okex/mainnet/main/genesis.json)

```bash
$ shasum -a 256 ${EXCHAIND_PATH}/config/genesis.json
0958b6c9f5f125d1d6b8f56e042fa8a71b1880310227b8b2f27ba93ff7cd673b  ${EXCHAIND_PATH}/config/genesis.json
```
