package pwrgo

import (
   "encoding/json"
   "log"
   "strconv"
   "crypto/ecdsa"
   "github.com/ethereum/go-ethereum/common/hexutil"
   "github.com/ethereum/go-ethereum/crypto"
   "math/big"
)

var ReturnBlockNumberOnTx = false

type Transaction struct {
	PositionInTheBlock int `json:"positionInTheBlock"`
	NonceOrValidationHash string `json:"nonceOrValidationHash"`
	Size int `json:"size"`
	Data string `json:"data"`
	VmId int `json:"vmId"`
	Fee int `json:"fee"`
	From string `json:"from"`
	To string `json:"to"`
	TxnFee int `json:"txnFee"`
	Type string `json:"type"`
	Hash string `json:"hash"`
}

type Block struct {
	BlockHash string `json:"blockHash"`
	Success bool `json:"success"`
	BlockNumber int `json:"blockNumber"`
	BlockReward int `json:"blockReward"`
	TransactionCount int `json:"transactionCount"`
	Transactions   []Transaction `json:"transactions"`
	BlockSubmitter string `json:"blockSubmitter"`
	BlockSize int `json:"blockSize"`
	Timestamp int `json:"timestamp"`
}

type Data struct {
  Nonce int `json:"nonce,omitempty"`
  Balance int `json:"balance,omitempty"`
  BlocksCount int `json:"blocksCount,omitempty"`
  ValidatorsCount int `json:"validatorsCount,omitempty"`
  Block Block `json:"block,omitempty"`
  Message string `json:"message,omitempty"`
}

type Response struct {
   Data Data `json:"data"`
   Status string `json:"status"`
   Success bool
   TxHash string
   BlockNumber int
   Error string
}

func parseResponse(responseStr string) (response Response) {
    err := json.Unmarshal([]byte(responseStr), &response)
    if err != nil {
        log.Fatalf("Error unmarshaling %s", err)
    }
	if response.Status == "success" {
		response.Success = true
	} else {
		response.Success = false
		response.Error = response.Data.Message
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

func GetBlock(blockNumber int) (Block) {
    var blockNumberStr = strconv.Itoa(blockNumber)
    var response = get(RPC_ENDPOINT + "/block/?blockNumber=" + blockNumberStr)
	var resp = parseResponse(response)
    return resp.Data.Block
}

func TransferPWR(to string, amount string, nonce int, privateKey *ecdsa.PrivateKey) (Response) {
    if len(to) != 42 {
        log.Fatalf("Invalid address ", to)
    }
    if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
    }

    amt := new(big.Int)
    amt.SetString(amount, 10)
    var buffer []byte
    buffer, err := txBytes(0, nonce, amt, to)
    if err != nil {
        log.Fatalf("Failed to get tx bytes ", err.Error())
    }

    signature, err := signMessage(buffer, privateKey)
    if err != nil {
        log.Fatalf("Failed to sign message ", err.Error())
    }

	var blockNumber = 0
	if ReturnBlockNumberOnTx {
		blockNumber = BlocksCount()
	}

    finalTxn := append(buffer, signature...)
    var transferTx = hexutil.Encode(finalTxn)
    var transferTxn = `{"txn":"` + transferTx[2:] + `"}`
    var result = post(RPC_ENDPOINT + "/broadcast/", transferTxn)
	hash := crypto.Keccak256Hash(finalTxn)

	transferResponse := parseResponse(result)
	transferResponse.TxHash = hash.Hex()
	transferResponse.BlockNumber = blockNumber
    return transferResponse
}

func SendVMDataTx(vmId int64, data []byte, nonce int, privateKey *ecdsa.PrivateKey) (Response) {
    if nonce < 0 {
        log.Fatalf("Nonce cannot be negative ", nonce)
    }

    var buffer []byte
    buffer, err := vmDataBytes(vmId, nonce, data)
    if err != nil {
        log.Fatalf("Failed to get vm data bytes ", err.Error())
    }

    signature, err := signMessage(buffer, privateKey)
    if err != nil {
        log.Fatalf("Failed to sign message ", err.Error())
    }

	var blockNumber = 0
	if ReturnBlockNumberOnTx {
		blockNumber = BlocksCount()
	}

    finalTxn := append(buffer, signature...)
    var dataTx = hexutil.Encode(finalTxn)
    var dataTxn = `{"txn":"` + dataTx[2:] + `"}`
    var result = post(RPC_ENDPOINT + "/broadcast/", dataTxn)
	hash := crypto.Keccak256Hash(finalTxn)

	vmDataTxResponse := parseResponse(result)
	vmDataTxResponse.TxHash = hash.Hex()
	vmDataTxResponse.BlockNumber = blockNumber
    return vmDataTxResponse
}
