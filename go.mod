module github.com/okex/mainnet

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/ethereum/go-ethereum v1.9.25
	github.com/okex/okexchain v0.16.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.9
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/okex/cosmos-sdk v0.39.2-okexchain2
	github.com/tendermint/iavl => github.com/okex/iavl v0.14.1-okexchain1
	github.com/tendermint/tendermint => github.com/okex/tendermint v0.33.9-okexchain1
)
