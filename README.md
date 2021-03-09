# OKExChain mainnet

This repo collects the genesis and configuration files for the various OKExChain
mainnet. It exists so the [OKExChain repo](https://github.com/okex/okexchain)
does not get bogged down with large genesis files and status updates.

## Getting Started

To get started with the latest mainnet, see the
[docs](https://okexchain-docs.readthedocs.io/en/latest/getting-start/join-okexchain-mainnet.html).

## Startup an okexchain full node by the genesis.json file
- Build okexchaind by [the latest released version v0.16.3](https://github.com/okex/okexchain/releases/tag/v0.16.3)

- Start mainnet with the okexchaind binary
```
okexchaind init your_custom_moniker --chain-id okexchain-66 --home ~/.okexchaind

wget https://raw.githubusercontent.com/okex/mainnet/main/genesis.json -O ~/.okexchaind/config/genesis.json

export OKEXCHAIN_SEEDS="e926c8154a2af4390de02303f0977802f15eafe2@3.16.103.80:26656,7fa5b1d1f1e48659fa750b6aec702418a0e75f13@35.177.8.240:26656,c8f32b793871b56a11d94336d9ce6472f893524b@18.167.16.85:26656"

okexchaind start --chain-id okexchain-66 --home ~/.okexchaind --p2p.seeds $OKEXCHAIN_SEEDS
```

- Sanity check the [genesis file](https://raw.githubusercontent.com/okex/mainnet/main/genesis.json)

```bash
$ shasum -a 256 ~/.okexchaind/config/genesis.json
1705b40f65f9f77083658a12e557e3225ecba529ec1328dcb08c0df1d4e42125  ~/.okexchaind/config/genesis.json
```

## Startup an okexchain full node by the mainnet snapshot data

```

okexchaind init your_custom_moniker --chain-id okexchain-66 --home ~/.okexchaind

rm -rf ~/.okexchaind/data
cd ~/.okexchaind

wget https://ok-public-hk.oss-cn-hongkong.aliyuncs.com/cdn/okexchain/snapshot/okexchain-v0.16.3-mainnet-20210127-height_275913.tar.gz

tar -zxvf okexchain-v0.16.3-mainnet-20210127-height_275913.tar.gz

export OKEXCHAIN_SEEDS="e926c8154a2af4390de02303f0977802f15eafe2@3.16.103.80:26656,7fa5b1d1f1e48659fa750b6aec702418a0e75f13@35.177.8.240:26656,c8f32b793871b56a11d94336d9ce6472f893524b@18.167.16.85:26656"

okexchaind start --chain-id okexchain-66 --home ~/.okexchaind --p2p.seeds $OKEXCHAIN_SEEDS
```
