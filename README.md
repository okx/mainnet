# OKExChain mainnet

This repo collects the genesis and configuration files for the various OKExChain
mainnet. It exists so the [OKExChain repo](https://github.com/okex/okexchain)
does not get bogged down with large genesis files and status updates.

## Getting Started

To get started with the latest mainnet, see the
[docs](https://okexchain-docs.readthedocs.io/en/latest/getting-start/join-okexchain-mainnet.html).

## mainnet Status
Source Code: [latest released version](https://github.com/okex/okexchain/releases/tag/v0.16.0)

⚠️ Latest mainnet: [okexchain v0.16.0](https://github.com/okex/okexchain/releases/tag/v0.16.0) ⚠️
* *Jan 6, 2021 11:19 UTC* - okexchain-v0.16

Download the [genesis file](https://raw.githubusercontent.com/okex/mainnet/main/genesis.json)

```bash
$ shasum -a 256 genesis.json
8379fd587486f0fd108b16099b5f36274563e28722e252d177c343c35aaf7ddb  genesis.json
```
Please read [GENESIS.md](GENESIS.md) for details on how it was generated and
to recompute it for yourself.

Seed nodes:
```
e926c8154a2af4390de02303f0977802f15eafe2@3.16.103.80:26656
7fa5b1d1f1e48659fa750b6aec702418a0e75f13@35.177.8.240:26656
c8f32b793871b56a11d94336d9ce6472f893524b@18.167.16.85:26656
```
