# pwrgo
Golang library for PWR Chain RPC interactions

# Run the example

In the repo root, run

```
go run example.go
```

# Functions

- NonceOfUser(address string)
- BalanceOf(address string)
- BlocksCount()
- ValidatorsCount()
- GetBlock(blockNumber int)
- TransferPWR(address string, amount string, nonce int, private_key)


# Example output

```
Private key hex: E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F4650C5F2B75
Running nonce/balance test for:  0x4097a04a9a8fef9ffb4a64e460193e6eb0b557c8
Nonce:  0
Balance:  0
Blocks count:  15
Validators count:  2
Latest block:  {"data":{"block":{"blockHash":"0xc3d3e1bc3838d721987bd983eb6aab2f54c5a2165c9a54920356c39bbb937934","success":true,"blockNumber":14,"blockReward":9800,"transactionCount":1,"transactions":[{"positionInTheBlock":0,"nonceOrValidationHash":"0","size":98,"fee":9800,"from":"0x2605c1ad496f428ab2b700edd257f0a378f83750","to":"0x3e1fa3b7f1dcf20890604c50a01b79ef79a33a5f","txnFee":9800,"type":"Transfer","value":100,"hash":"0x7f57051cc8ba50c7400d89fba0980c27fda6425ec1d8589a5d24f1dc9f2db919"}],"blockSubmitter":"0x61bd8fc1e30526aaf1c4706ada595d6d236d9883","blockSize":217,"timestamp":1699332516}},"status":"success"}
Transfer tx :  {"data":{"message":"Txn broadcasted to validator nodes"},"status":"success"}
```

# To-do:

- wrap as a package or module