# pwrgo
Golang library for PWR Chain RPC interactions

# Run the example

In the repo root, run

```
go run .
```

# Example output

```
Running test for:  0x4097a04a9a8fef9ffb4a64e460193e6eb0b557c8
Nonce:  1
Balance:  199992000028406
Blocks count:  2218
Validators count:  20
Latest block:  {"data":{"block":{"blockHash":"0x68dd58de587ad45299f83a7b068f92d9ad00013fcfc99e8634447f025e34659d","success":true,"blockNumber":2217,"blockReward":8100,"transactionCount":1,"transactions":[{"positionInTheBlock":0,"nonceOrValidationHash":"7","size":81,"data":"686579","vmId":99,"fee":8100,"from":"0x586f776b04f52651aa185eee3f77cf6160d61b78","to":"VM: 99","txnFee":8100,"type":"VM Data","hash":"0x1db617f14fc5f70cf50241e21786b930605ed59c1d4c43c01eb52ccea5df0bcf"}],"blockSubmitter":"0x3e1fa3b7f1dcf20890604c50a01b79ef79a33a5f","blockSize":200,"timestamp":1699045566}},"status":"success"}
```

# To-do:

- wrap as a package or module
- signing and broadcasting transactions