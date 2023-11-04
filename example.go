package main

import (
    "fmt"
    "pwr/pwrgo"
)

func main() {
	var address = "0x4097a04a9a8fef9ffb4a64e460193e6eb0b557c8"
	fmt.Println("Running test for: ", address)

	var nonce = pwrgo.NonceOfUser(address)
	fmt.Println("Nonce: ", nonce)

	var balance = pwrgo.BalanceOf(address)
	fmt.Println("Balance: ", balance)
	
	var blocksCount = pwrgo.BlocksCount()
	fmt.Println("Blocks count: ", blocksCount)
	
	var validatorsCount = pwrgo.ValidatorsCount()
	fmt.Println("Validators count: ", validatorsCount)
	
	var latestBlock = pwrgo.GetBlock(blocksCount - 1)
	fmt.Println("Latest block: ", latestBlock)
}