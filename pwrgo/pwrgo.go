package pwrgo

import (
   "encoding/json"
   "io/ioutil"
   "log"
   "net/http"
   //"fmt"
   "strconv"
)

var RPC_ENDPOINT = "https://pwrrpc.pwrlabs.io"

func get(url string) (response string) {
   resp, err := http.Get(url)
   if err != nil {
      log.Fatalln(err)
   }

   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
   response = string(body)
   return
}

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