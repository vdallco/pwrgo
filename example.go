package main

import (
    "fmt"
    "github.com/pwrlabs/pwrgo/pwrgo"
)

func main() {
    // Import wallet by private key
    privateKeyHex := "0x9d4428c6e0638331b4866b70c831f8ba51c11b031f4b55eed4087bbb8ef0151f"
    var wallet = pwrgo.FromPrivateKey(privateKeyHex)
    
    fmt.Printf("Public key: %s\n", wallet.PublicKey)
    fmt.Printf("Address: %s\n", wallet.Address)
    
    // Get nonce for address
    var nonce = pwrgo.NonceOfUser(wallet.Address)
    fmt.Println("Nonce: ", nonce)
    
    // Get PWR balance of address
    var balance = pwrgo.BalanceOf(wallet.Address)
    fmt.Println("Balance: ", balance)
    
    // Get total blocks count
    var blocksCount = pwrgo.BlocksCount()
    fmt.Println("Blocks count: ", blocksCount)

    // Get total validators count
    var validatorsCount = pwrgo.ValidatorsCount()
    fmt.Println("Validators count: ", validatorsCount)
    
    // Get block info by Block Number
    var latestBlock = pwrgo.GetBlock(blocksCount - 1)
    fmt.Println("Latest block hash: ", latestBlock.BlockHash)
    fmt.Println("Latest block timestamp: ", latestBlock.Timestamp)
    fmt.Println("Latest block tx count: ", latestBlock.TransactionCount)
    fmt.Println("Latest block submitter: ", latestBlock.BlockSubmitter)
    
	// Transfer PWR
	pwrgo.ReturnBlockNumberOnTx = true // automatically calls blocksCount from RPC and returns BlockNumber on tx response
    var transferTx = pwrgo.TransferPWR("0x61bd8fc1e30526aaf1c4706ada595d6d236d9883", "1", nonce, wallet.PrivateKey) // send 1 PWR
	if transferTx.Success {
		fmt.Printf("[Block #%d] Transfer tx hash: %s\n", transferTx.BlockNumber, transferTx.TxHash)
		nonce = nonce + 1 // increment nonce since we just Transferred PWR
	} else {
		fmt.Println("Error sending Transfer tx: ", transferTx.Error)
		fmt.Println("Error sending ", transferTx.TxHash)
	}

    // Create new wallet and print address and keys
    var newWallet = pwrgo.NewWallet()
    fmt.Println("New wallet address: ", newWallet.Address)
    fmt.Println("New wallet private key: ", newWallet.PrivateKeyStr)
    fmt.Println("New wallet public key: ", newWallet.PublicKey)
    
    // Send data to VM 1337
    var data = []byte("Hello world")
    var vmTxResponse = pwrgo.SendVMDataTx(1337, data, nonce, wallet.PrivateKey)
	if vmTxResponse.Success {
		fmt.Printf("[Block #%d] VM data tx hash: %s", vmTxResponse.BlockNumber, vmTxResponse.TxHash)
	} else {
		fmt.Println("Error sending VM data tx: ", vmTxResponse.Error)
		fmt.Println("Error sending ", vmTxResponse.TxHash)
	}
	
}