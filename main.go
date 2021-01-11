package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/okex/okexchain/app"
	"github.com/okex/okexchain/app/codec"
	ethermint "github.com/okex/okexchain/app/types"
	swap "github.com/okex/okexchain/x/ammswap"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/dex"
	"github.com/okex/okexchain/x/distribution"
	farm "github.com/okex/okexchain/x/farm/types"
	"github.com/okex/okexchain/x/genutil"
	"github.com/okex/okexchain/x/gov"
	"github.com/okex/okexchain/x/order"
	"github.com/okex/okexchain/x/params"
	"github.com/okex/okexchain/x/staking"
	"github.com/okex/okexchain/x/token"
	"github.com/okex/okexchain/x/token/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	amino "github.com/tendermint/go-amino"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	genesisTemplate = "params/okexchain_genesis_template.json"
	genTxPath       = "gentx"
	tokensPath      = "tokens"
	accountDir      = "accounts"
	tokenSchema     = "token.json"
	genesisFile     = "okexchain-genesis.json"

	denomination      = common.NativeToken
	timeGenesisString = "2021-01-10 04:00:00 -0000 UTC"
)

// constants but can't use `const`
var (
	timeGenesis time.Time
)

func loadTokens() (res []types.Token) {
	fs, err := ioutil.ReadDir(tokensPath)
	if err != nil {
		panic(err)
	}
	for _, f := range fs {
		name := f.Name()
		if name == "README.md" {
			continue
		}
		bz, err := ioutil.ReadFile(path.Join(tokensPath, name+"/"+tokenSchema))
		if err != nil {
			panic(err)
		}
		var token types.Token
		err = json.Unmarshal(bz, &token)
		if err != nil {
			panic(err)
		}
		res = append(res, token)
	}
	return
}

func loadAccounts(tokenName string) map[string]float64 {
	path := tokensPath + "/" + tokenName + "/" + accountDir

	fs, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	accounts := make(map[string]float64)
	for _, f := range fs {
		name := f.Name()
		if err := accumulateContributors(path+"/"+name, accounts); err != nil {
			panic(err)
		}
	}
	return accounts
}
func sortAccounts(accountMap map[string]authexported.GenesisAccount) []authexported.GenesisAccount {
	var genesisAccounts []authexported.GenesisAccount
	for _, acc := range accountMap {
		genesisAccounts = append(genesisAccounts, acc)
	}
	// sort the accounts
	sort.SliceStable(genesisAccounts, func(i, j int) bool {
		return strings.Compare(
			genesisAccounts[i].GetAddress().String(),
			genesisAccounts[j].GetAddress().String(),
		) < 0
	})
	return genesisAccounts
}

// initialize the times!
func init() {
	var err error
	timeLayoutString := "2006-01-02 15:04:05 -0700 MST"
	timeGenesis, err = time.Parse(timeLayoutString, timeGenesisString)
	if err != nil {
		panic(err)
	}
}

func main() {

	// 0. load tokens and feed accounts
	tokens := loadTokens()
	accountMap := make(map[string]authexported.GenesisAccount, 0)
	for _, token := range tokens {
		makeGenesisAccounts(token.Symbol, accountMap)
	}
	genesisAccounts := sortAccounts(accountMap)

	fmt.Println("-----------")
	fmt.Printf("TOTAL genesis accounts: %d \n", len(genesisAccounts))

	// 2. sanity check totals
	checkTotals(genesisAccounts, tokens)

	// 3. load gentxs, validators
	fs, err := ioutil.ReadDir(genTxPath)
	if err != nil {
		panic(err)
	}

	var genTxs []json.RawMessage
	for _, f := range fs {
		name := f.Name()
		if name == "README.md" {
			continue
		}
		bz, err := ioutil.ReadFile(path.Join(genTxPath, name))
		if err != nil {
			panic(err)
		}
		genTxs = append(genTxs, json.RawMessage(bz))
	}

	fmt.Printf("TOTAL genesis validators: %d\n", len(genTxs))
	fmt.Println("-----------")

	cdc := codec.MakeCodec(app.ModuleBasics)

	// 4. produce the genesis file
	genesisDoc := makeGenesisDoc(cdc, genesisAccounts, genTxs, tokens)
	// write the genesis file
	bz, err := cdc.MarshalJSON(genesisDoc)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer([]byte{})
	err = json.Indent(buf, bz, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(genesisFile, buf.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

func fromBech32(address string) sdk.AccAddress {
	bech32PrefixAccAddr := "okexchain"
	bz, err := sdk.GetFromBech32(address, bech32PrefixAccAddr)
	if err != nil {
		panic(err)
	}
	if len(bz) != sdk.AddrLen {
		panic("Incorrect address length")
	}
	return sdk.AccAddress(bz)
}

// load a map of hex addresses and convert them to bech32
func accumulateContributors(fileName string, contribs map[string]float64) error {
	allocations := ObjToMap(fileName)

	for addr, amt := range allocations {
		if _, ok := contribs[addr]; ok {
			fmt.Println("Duplicate addr", addr)
		}
		contribs[addr] += amt
	}
	return nil
}

func makeGenesisAccounts(tokenName string, accountMap map[string]authexported.GenesisAccount) {

	accounts := loadAccounts(tokenName)
	for addr, amount := range accounts {
		var account authexported.GenesisAccount
		var ok bool
		account, ok = accountMap[addr]
		if !ok {
			account = &ethermint.EthAccount{
				BaseAccount: authtypes.NewBaseAccount(fromBech32(addr), sdk.DecCoins{}, nil, 0, 0),
				CodeHash:    ethcrypto.Keccak256(nil),
			}
		}
		//amountDec := sdk.MustNewDecFromStr(fmt.Sprintf("%f", amount))
		//tokenAmount := sdk.NewDecCoinsFromDec(tokenName, amountDec)
		tokenAmount, err := sdk.ParseDecCoins(fmt.Sprintf("%f%s", amount, tokenName))
		if err != nil {
			panic(err)
		}
		account.SetCoins(account.GetCoins().Add(tokenAmount...))
		accountMap[addr] = account
	}
}

// check total atoms and no duplicates
func checkTotals(genesisAccounts []authexported.GenesisAccount, tokens []types.Token) {
	// check total
	for _, token := range tokens {
		total := sdk.NewDec(0)
		for _, account := range genesisAccounts {
			total = total.Add(account.GetCoins().AmountOf(token.Symbol))
		}
		fmt.Printf("<%s>: total supply <%s>, ", token.Symbol, token.OriginalTotalSupply)
		fmt.Printf("sum(account balance) <%s>\n", total)

		if !total.Equal(token.OriginalTotalSupply) {
			panic(fmt.Sprintf("Failed to check %s Total Supply", token.Symbol))
		}
	}

	// ensure no duplicates
	checkdupls := make(map[string]struct{})
	for _, acc := range genesisAccounts {
		if _, ok := checkdupls[acc.GetAddress().String()]; ok {
			panic(fmt.Sprintf("Got duplicate: %v", acc.GetAddress()))
		}
		checkdupls[acc.GetAddress().String()] = struct{}{}
	}
	if len(checkdupls) != len(genesisAccounts) {
		panic("length mismatch!")
	}
}

func produceAppState(cdc *amino.Codec,
	genesisAccounts []authexported.GenesisAccount,
	genTxs []json.RawMessage, tokens []types.Token) map[string]json.RawMessage {
	appState := app.ModuleBasics.DefaultGenesis()

	appState = genutil.SetGenesisStateInAppState(cdc, appState, genutil.NewGenesisState(genTxs))

	// auth
	// add genesis account to the app state
	authGenState := auth.DefaultGenesisState()
	authGenState.Accounts = genesisAccounts
	authGenState.Accounts = auth.SanitizeGenesisAccounts(authGenState.Accounts)

	authGenStateBz := cdc.MustMarshalJSON(authGenState)
	appState[auth.ModuleName] = authGenStateBz

	// token
	var genesisToken token.GenesisState
	cdc.MustUnmarshalJSON(appState[token.ModuleName], &genesisToken)
	if len(genesisToken.Tokens) != 1 {
		panic(fmt.Errorf("no genesis denom"))
	}
	// use okt_info.json overwrite default okt into
	genesisToken.Tokens = nil
	genesisToken.Tokens = append(genesisToken.Tokens, tokens...)
	genesisToken.Params.FeeIssue = sdk.ZeroFee()
	genesisToken.Params.FeeMint = sdk.ZeroFee()
	genesisToken.Params.FeeBurn = sdk.ZeroFee()
	genesisToken.Params.FeeModify = sdk.ZeroFee()
	genesisToken.Params.FeeChown = sdk.ZeroFee()

	tokenStateBz := cdc.MustMarshalJSON(genesisToken)
	appState[token.ModuleName] = tokenStateBz

	// staking
	var genesisStaking = staking.DefaultGenesisState()
	genesisStaking.Params.BondDenom = denomination
	genesisStaking.Params.MinSelfDelegation = sdk.NewDec(10000)

	stakingStateBz := cdc.MustMarshalJSON(genesisStaking)
	appState[staking.ModuleName] = stakingStateBz

	// bank
	var genesisBank bank.GenesisState
	genesisBank.SendEnabled = false

	genesisBankBz := cdc.MustMarshalJSON(genesisBank)
	appState[bank.ModuleName] = genesisBankBz

	// gov
	var genesisGov = gov.DefaultGenesisState()
	genesisGov.DepositParams.MinDeposit = sdk.SysCoins{sdk.NewDecCoin(denomination, sdk.NewInt(10))}

	genesisGovBz := cdc.MustMarshalJSON(genesisGov)
	appState[gov.ModuleName] = genesisGovBz

	// params
	var genesisParams = params.DefaultGenesisState()
	genesisParams.Params.MinDeposit = sdk.SysCoins{sdk.NewDecCoin(denomination, sdk.NewInt(10))}

	genesisParamsBz := cdc.MustMarshalJSON(genesisParams)
	appState[params.ModuleName] = genesisParamsBz

	// order
	var genesisOrder = order.DefaultGenesisState()
	genesisOrder.Params.TradeFeeRate = sdk.MustNewDecFromStr("0.001")
	genesisOrder.Params.NewOrderMsgGasUnit = 40000
	genesisOrder.Params.CancelOrderMsgGasUnit = 30000

	genesisOrderBz := cdc.MustMarshalJSON(genesisOrder)
	appState[order.ModuleName] = genesisOrderBz

	// dex
	var genesisDex = dex.DefaultGenesisState()
	genesisDex.Params.ListFee = sdk.NewDecCoinFromDec(denomination, sdk.MustNewDecFromStr("1000"))
	genesisDex.Params.RegisterOperatorFee = sdk.ZeroFee()
	genesisDex.Params.TransferOwnershipFee = sdk.ZeroFee()
	genesisDex.Params.DelistMinDeposit = sdk.SysCoins{sdk.NewDecCoin(denomination, sdk.NewInt(10))}

	genesisDexBz := cdc.MustMarshalJSON(genesisDex)
	appState[dex.ModuleName] = genesisDexBz

	// swap
	var genesisSwap = swap.DefaultGenesisState()
	genesisSwap.Params.FeeRate = sdk.MustNewDecFromStr("0.003")

	genesisSwapBz := cdc.MustMarshalJSON(genesisSwap)
	appState[swap.ModuleName] = genesisSwapBz

	// dist
	var genesisDistr = distribution.DefaultGenesisState()
	genesisDistr.Params.CommunityTax = sdk.MustNewDecFromStr("0.02")

	genesisDistrBz := cdc.MustMarshalJSON(genesisDistr)
	appState[distribution.ModuleName] = genesisDistrBz

	// farm
	var genesisFarm = farm.DefaultGenesisState()
	genesisFarm.Params.CreatePoolDeposit = sdk.ZeroFee()
	genesisFarm.Params.CreatePoolFee = sdk.ZeroFee()

	genesisFarmBz := cdc.MustMarshalJSON(genesisFarm)
	appState[farm.ModuleName] = genesisFarmBz

	return appState
}

// json marshal the initial app state (accounts and gentx) and add them to the template
func makeGenesisDoc(cdc *amino.Codec,
	genesisAccounts []authexported.GenesisAccount,
	genTxs []json.RawMessage, tokens []types.Token) *tmtypes.GenesisDoc {
	// 1. read the template with the params
	_, genesisDoc, err := genutil.GenesisStateFromGenFile(cdc, genesisTemplate)
	if err != nil {
		panic(err)
	}

	// 2. set genesis time
	genesisDoc.GenesisTime = timeGenesis

	// 3. read the okexchain state from the generic app state bytes
	// and populate with the accounts and gentxs
	appState := produceAppState(cdc, genesisAccounts, genTxs, tokens)
	// marshal the okexchain app state back to json and update the genesisDoc
	appStateJSON, err := cdc.MarshalJSON(appState)
	if err != nil {
		panic(err)
	}
	genesisDoc.AppState = appStateJSON
	return genesisDoc
}

// ObjToMap Load a JSON object of addr->amt into a map.
// Expects no duplicates!
// TODO: remove this and do everything through lists so duplicates can always be detected?
func ObjToMap(file string) map[string]float64 {
	bz, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	m := make(map[string]float64)
	err = json.Unmarshal(bz, &m)
	if err != nil {
		panic(err)
	}
	return m
}
