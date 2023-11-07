package main

import (
    "fmt"
    "pwr/pwrgo"
	"log"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Import wallet by private key
	pHex := "E83385AF76B2B1997326B567461FB73DD9C27EAB9E1E86D26779F4650C5F2B75"
    fmt.Printf("Private key hex: %s\n", pHex)
    pk, err := crypto.HexToECDSA(pHex)
    if err != nil {
        log.Fatal(err.Error())
    }

	// Random wallet for testing (does not match above Private key)
	var address = "0x4097a04a9a8fef9ffb4a64e460193e6eb0b557c8"
	fmt.Println("Running nonce/balance test for: ", address)

	// Get nonce for address
	var nonce = pwrgo.NonceOfUser(address)
	fmt.Println("Nonce: ", nonce)

	// Get PWR balance of address
	var balance = pwrgo.BalanceOf(address)
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
	var transferTx = pwrgo.TransferPWR(address, "1", nonce + 1, pk) // send 1 PWR to address, given nonce and private key bytes
    fmt.Println("Transfer tx : ", transferTx)
}