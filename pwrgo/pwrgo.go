package pwrgo

import (
   "encoding/json"
   "log"
   "strconv"
   "crypto/ecdsa"
   "github.com/ethereum/go-ethereum/common/hexutil"
   "math/big"
)

type Response struct {
   Data Data `json:"data"`
   Status string `json:"status"`
}

type Data struct {
  Nonce int `json:"nonce,omitempty"`
  Balance int `json:"balance,omitempty"`
  BlocksCount int `json:"blocksCount,omitempty"`
  ValidatorsCount int `json:"validatorsCount,omitempty"`
}

func parseResponse(responseStr string) (response Response) {
    err := json.Unmarshal([]byte(responseStr), &response)
    if err != nil {
        log.Fatalf("Error unmarshaling %s", err)
    }
	return
}

func NonceOfUser(address string) (int) {
	var response = get(RPC_ENDPOINT + "/nonceOfUser/?userAddress=" + address)
	var resp = parseResponse(response)
	return resp.Data.Nonce
}

func BalanceOf(address string) (int) {
	var response = get(RPC_ENDPOINT + "/balanceOf/?userAddress=" + address)
	var resp = parseResponse(response)
	return resp.Data.Balance
}

func BlocksCount() (int) {
	var response = get(RPC_ENDPOINT + "/blocksCount/")
	var resp = parseResponse(response)
	return resp.Data.BlocksCount
}

func ValidatorsCount() (int) {
	var response = get(RPC_ENDPOINT + "/validatorsCount/")
	var resp = parseResponse(response)
	return resp.Data.ValidatorsCount
}

func GetBlock(blockNumber int) (string) {
	var blockNumberStr = strconv.Itoa(blockNumber)
	var response = get(RPC_ENDPOINT + "/block/?blockNumber=" + blockNumberStr)
	return response
}

func TransferPWR(to string, amount string, nonce int, privateKey *ecdsa.PrivateKey) (string) {
	if len(to) != 42 {
		return "Invalid address"
	}
	if nonce < 0 {
		return "Nonce cannot be negative"
	}

	amt := new(big.Int)
    amt.SetString(amount, 10)
	var buffer []byte
	buffer, err := txBytes(0, nonce, amt, to)
	if err != nil {
		return "Failed to get tx bytes"
	}

	signature, err := signMessage(buffer, privateKey)
	if err != nil {
		return "Failed to sign message"
	}

	finalTxn := append(buffer, signature...)
	var transferTx = hexutil.Encode(finalTxn)
	var transferTxn = `{"txn":"` + transferTx[2:] + `"}`
	var result = post(RPC_ENDPOINT + "/broadcast/", transferTxn)
	return result
}
