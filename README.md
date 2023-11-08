# pwrgo
Golang library for PWR Chain RPC interactions

# Run the example

In the repo root, run

```
go run example.go
```

# Functions

- NonceOfUser(address string) int
- BalanceOf(address string) int
- BlocksCount() int
- ValidatorsCount() int
- GetBlock(blockNumber int) int
- TransferPWR(address string, amount string, nonce int, PrivateKey) string
- FromPrivateKey(privateKey string) PWRWallet
- NewWallet() PWRWallet
- SendVMDataTx(vmId string, data []byte, nonce int, PrivateKey) string

# Example output

```
Public key: 0x04369d83469a66920f31e4cf3bd92cb0bc20c6e88ce010dfa43e5f08bc49d11da87970d4703b3adbc9a140b4ad03a0797a6de2d377c80c369fe76a0f45a7a39d3f
Address: 0x2605c1Ad496F428aB2b700Edd257f0a378f83750
Nonce:  1
Balance:  99999990100
Blocks count:  49
Validators count:  2
Latest block:  {"data":{"block":{"blockHash":"0x5b8ea364240abe31266d7856519102ba0a0e84cb7bc1556b3e5af40317f5aef9","success":true,"blockNumber":48,"blockReward":40900,"transactionCount":1,"transactions":[{"positionInTheBlock":0,"nonceOrValidationHash":"6","size":409,"data":"f9014801843b9aca00828dfd94759ff85f2a3f8ebd38bd79bda3af56d499fa4b1f80b8e4f7742d2f00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000008dbdfb3ae060fb00000000000000000000000000000000000000000000000008ac7230489e80000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001900000000000000000000000000000000000000000000000000000000000000018363534613665336436373466363331356366643363656430000000000000000062a084608d71ecacf858e00b26363ceee8c5a550e5a51704538f6e0c233c05d504eaa032cd9078b3ab18eeeb2f0d8aec7848291fb2a58ac38d8af713d77bbe67957dac","vmId":31,"fee":40900,"from":"0xc7e9fbfbd60df0cfcf3e8a3c5b30ed0def178c57","to":"VM: 31","txnFee":40900,"type":"VM Data","hash":"0xdbb16e09e89aa16cf69992ee08beb1d66a37b8b4264ccd6adeb74d166026293a"}],"blockSubmitter":"0x61bd8fc1e30526aaf1c4706ada595d6d236d9883","blockSize":528,"timestamp":1699376702}},"status":"success"}
Transfer tx :  {"data":{"message":"Txn broadcasted to validator nodes"},"status":"success"}
New wallet address:  0x2EdcE756c2b332DF96EeE3f8b71252901bB43274
New wallet private key:  0x6013906878601d0cb8bd79f751e67c35bc611c8ae663c45b9bd10832e85d1141
New wallet public key:  0x04110631961a78de320903c7747f444677e651fe16e9828fa2fd1825b697da47a9a343e93e8e387abf3952b1c5ce29c664711cb26df3e1dc4f55a56dab8e668006
```

# To-do:

- wrap as a package or module
- sign and broadcast VM data