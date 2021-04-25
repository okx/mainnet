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
#### 2.1 Export by exchaind
```
   1. git clone -b v0.16.3.1 https://github.com/okex/exchain.git
   2. make install
   3. exchaind export --home ${EXCHAIND_PATH} --height=2322600 --for-zero-height --log_level evm:debug --log_file ./export.log --log_stdout=false > ${EXCHAIND_PATH}/config/genesis_no_migrate.json
   4. git checkout v0.18.0
   5. make GenesisHeight=2322600 install
   6.exchaind migrate v0.18 ${EXCHAIND_PATH}/config/genesis_no_migrate.json --chain-id=exchain66 > genesis.json
```
#### 2.2 Download from official
```
   wget https://raw.githubusercontent.com/okex/mainnet/main/genesis.json -O ${EXCHAIND_PATH}/config/genesis.json
```
Note: it needs to check genesis.json no matter which way is used
```
$ shasum -a 256 ~/.okexchaind/config/genesis.json
${genesis_of_2322600_shasum}  ${EXCHAIND_PATH}/config/genesis.json
```
### 4. Build exchaind binary
Build exchaind by [the latest released version v0.18.0](https://github.com/okex/exchain/releases/tag/v0.18.0)

### 5. Reset data
`exchaind unsafe-reset-all --home /data/okexchaind`
### 6. Start
`exchaind start --chain-id exchain-66 --home /data/okexchaind`


## Startup an okexchain full node by the genesis.json file
### 1. Build exchaind by [the latest released version v0.18.0](https://github.com/okex/exchain/releases/tag/v0.18.0)


### 2. Start mainnet with the okexchaind binary

```
exchaind init your_custom_moniker --chain-id exchain-66 --home ${EXCHAIND_PATH}

wget https://raw.githubusercontent.com/okex/mainnet/main/genesis.json -O ${EXCHAIND_PATH}/config/genesis.json

export OKEXCHAIN_SEEDS="e926c8154a2af4390de02303f0977802f15eafe2@3.16.103.80:26656,7fa5b1d1f1e48659fa750b6aec702418a0e75f13@35.177.8.240:26656,c8f32b793871b56a11d94336d9ce6472f893524b@18.167.16.85:26656"

exchaind start --chain-id exchain-66 --home ${EXCHAIND_PATH} --p2p.seeds $OKEXCHAIN_SEEDS
```

Note: it needs to check the [genesis file](https://raw.githubusercontent.com/okex/mainnet/main/genesis.json)

```bash
$ shasum -a 256 ${EXCHAIND_PATH}/config/genesis.json
1705b40f65f9f77083658a12e557e3225ecba529ec1328dcb08c0df1d4e42125  ${EXCHAIND_PATH}/config/genesis.json
```
