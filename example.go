package main

import (
    "fmt"
    "pwr/pwrgo"
)

func main() {
	// Import wallet by private key
	privateKeyHex := "0xE83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F4650C5F2B75"
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
	fmt.Println("Latest block: ", latestBlock)
	
	// Transfer PWR
	var transferTx = pwrgo.TransferPWR(wallet.Address, "1", nonce + 1, wallet.PrivateKey) // send 1 PWR to address, given nonce and private key bytes
    fmt.Println("Transfer tx : ", transferTx)

    // Create new wallet and print address and keys
    var newWallet = pwrgo.NewWallet()
	fmt.Println("New wallet address: ", newWallet.Address)
    fmt.Println("New wallet private key: ", newWallet.PrivateKeyStr)
    fmt.Println("New wallet public key: ", newWallet.PublicKey)
}