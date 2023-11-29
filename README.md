# pwrgo
Golang library for PWR Chain RPC interactions

# Installation

To install the pwrgo library, run

```
go get github.com/pwrlabs/pwrgo@v0.0.3
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
Nonce:  11
Balance:  99999895793
Blocks count:  32328
Validators count:  2
Latest block hash:  0x5973e0dbe796cc8f64587db20e74d42029b3e74b1c1bb4e313e18170f5003d9c
Latest block timestamp:  1701298200
Latest block tx count:  7
Latest block submitter:  0x61bd8fc1e30526aaf1c4706ada595d6d236d9883
[Block #32328] Transfer tx hash: 0xe5f4f06a59119636dd9a4bab9874b2ae54eb4ae955fe8530cf3e32511fcade45
New wallet address:  0xECD1908715fcA3538818307e77bc1Dc682763A08
New wallet private key:  0xdafc07dfb1c52691045328cd4faad4886fb6e013e028f355edb29a94ef67385c
New wallet public key:  0x046dd883f70cfe3e046244c05a71c94b5860722662aedeeded557b768ba21d6e2cdeb3d6a780b14dc7280a9c1640c97e58606037b2f0b3baa8e14b3ca879aadfb5
[Block #32328] VM data tx hash: 0xfd9a16fc66f3a936c39b55de5a7a1a1058f4275d740ccefdec5a5c95cb9dfe51
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

- Social Recovery VM: implement call to new RPC method to get raw TX bytes
- Social Recovery VM: Recover Signature from TX bytes, and recover public key from signature + TX data