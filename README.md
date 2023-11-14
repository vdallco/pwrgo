# pwrgo
Golang library for PWR Chain RPC interactions

# Installation

To install the pwrgo library, run

```
go get github.com/pwrlabs/pwrgo@v0.0.2
```

# Using the library

Import the library

```
import(
   "github.com/pwrlabs/pwrgo/pwrgo"
)
```

# Run the example

In the repo root, run

```
go run example.go
```

# Functions

## General
- NonceOfUser(address string) int
- BalanceOf(address string) int
- BlocksCount() int
- ValidatorsCount() int
- GetBlock(blockNumber int) Block

## Wallet operations
- FromPrivateKey(privateKey string) PWRWallet
- NewWallet() PWRWallet

## Transactions
- TransferPWR(address string, amount string, nonce int, PrivateKey) Response
- SendVMDataTx(vmId string, data []byte, nonce int, PrivateKey) Response

# example.go output

```
Public key: 0x040cd999a20b0eba1cf86362c738929671902c9b337ab1370d2ba790be68b01227cab9fa9096b87651686bf898acf11857906907ba7fca4f5f5d9513bdd16e0a52
Address: 0xA4710E3D79E1ED973af58E0f269e9b21dD11Bc64
Nonce:  15
Balance:  99999868051
Blocks count:  29706
Validators count:  4
Latest block hash:  0xf514d6859d12eedb94f56d56f7119c8b45a4bd203e28b8b083b200ed7924ab92
Latest block timestamp:  1699684646
Latest block tx count:  1
Latest block submitter:  0x61bd8fc1e30526aaf1c4706ada595d6d236d9883
[Block #29706] Transfer tx hash: 0x2c65c38b6fb6bc5a9444b29791994ed6eb31aceb7c0b72128d89c6309c7522d6
New wallet address:  0xF0bc82Df249E93cF847E130eeFf314d80B05160f
New wallet private key:  0x51f2240c0de430a0d0b7c658653968902f2cc4a5d2ea85b445d1fb866b165050
New wallet public key:  0x04bd6f3c8402b9f58c8b5c432f43871c91733f2eef5b8cfa2f54d0d01c70f9d44f4aca8d7e5ce858cedcbbcfe36b66eb99b24b76862ac3d5fab9d98a1e54cadb75
```

# Proof-of-concept apps/VMs

## Chat

See the Go-Messaging-App repo: https://github.com/pwrlabs/Go-Messaging-App

## Social Recovery

To store Shamir's Secret Sharing (SSS) shares encrypted on-chain for trustees, run

```
go run .\sssExample.go
```

# To-do:

- wrap as a package or module
- Social Recovery VM: implement call to new RPC method to get raw TX bytes
- Social Recovery VM: Recover Signature from TX bytes, and recover public key from signature + TX data