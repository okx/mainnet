# Background

Simple ways to withdraw OKT or destroy validator

## Table of Contents

| Section | Description | Target Users | Remark |
|---------|-------------|-------------|-------|
| [1. How to Withdraw OKT](#1-how-to-withdraw-okt) | Methods to withdraw OKT tokens | Stakers | |
| &nbsp;&nbsp;&nbsp;&nbsp;[1.1 Using Command Line](#11-using-command-line) | Command line withdrawal process | Stakers | Or you can also use wallet (1.2) |
| &nbsp;&nbsp;&nbsp;&nbsp;[1.2 Using Wallet and Contract](#12-using-wallet-and-contract) | Wallet and smart contract withdrawal | Stakers | [Get the wallet](https://web3.okx.com/) |
| [2. How to Destroy Validator](#2-how-to-destroy-validator) | Process to deregister a validator | Validators | |
| [3. How to Build exchaincli](#3-how-to-build-exchaincli) | Build and setup exchaincli tool | All Users | |

# 1. How to Withdraw OKT

## 1.1 Using Command Line
```
// Withdraw an amount of OKT and the corresponding shares from all validators.
// You will have to wait 14 days before your OKTs are fully unlocked and transferrable. This will also trigger a passive reward, which will automatically distribute the rewards to its own account (or reward withdrawal account)
// e.g., <amountToWithdraw>=1024okt, <gasPrice>=0.00000001okt

exchaincli tx staking withdraw <amountToWithdraw> --from <delegatorKeyName> --gas auto --gas-adjustment 1.5 --gas-prices <gasPrice>
```

If you encounter the following error:
```
failed. destroyed validator xxxxxx isn't allowed to add shares to. please get rid of it from the shares adding list by adding shares to other validators again or unbond all delegated tokens
```

Please first execute the following voting operation, then run the above `withdraw` command again:
```
exchaincli tx staking add-shares <validator-addr1, validator-addr2, validator-addr3, ... validator-addrN> --from <delegatorKeyName> --gas auto --gas-adjustment 1.5 --gas-prices <gasPrice>
```

## 1.2 Using Wallet and Contract
Through the wallet, call the `genWithdrawMsg` interface on the UI to assemble and generate parameters, then call invoke to execute `Withdraw`:
```
function genWithdrawMsg() public view returns
```

If you encounter the following error:
```
failed. destroyed validator xxxxxx isn't allowed to add shares to. please get rid of it from the shares adding list by adding shares to other validators again or unbond all delegated tokens
```

Please first execute the following `AddShares` interface to assemble and generate parameters, call invoke to execute, then run the above `genWithdrawMsg` process again:
```
function genAddSharesMsg() public view returns
```

# 2. How to Destroy Validator

Deregister a validator and withdraw staked assets. You will have to wait 14 days before your OKTs are fully unlocked and transferrable. This will also trigger a passive reward, which will automatically distribute the rewards to its own account.

```
exchaincli tx staking destroy-validator --from ex1sxxxxx --gas auto --gas-prices 0.0000000001okt --gas-adjustment 1.3 -y  
```

# 3. How to Build exchaincli
```
git clone -b v1.7.0.6 https://github.com/okx/exchain.git
cd exchain; make install
# Import mnemonic:
exchaincli keys add --recover delegator1 -m "mnemonic" -y
```
