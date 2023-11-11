package main

import (
	"pwr/sssa"
	"fmt"
	"math/rand"
	"github.com/ethereum/go-ethereum/crypto"
    "encoding/hex"
	ecies "github.com/ecies/go/v2"
	"log"
    "pwr/pwrgo"
	"golang.org/x/crypto/sha3"
    "crypto/ecdsa"
	"strings"
    "net/http"
    "io/ioutil"
	"encoding/json"
	"bytes"
    "time"
)

var appId int64 = 222 // VM ID for Social Recovery
var stashSecretsTxType = 1 // Type byte added to Stash tx's
var trusteeRecoverSecret = 2 // Type byte added to tx's where Trustee is recovering their secret share for a given addr && a given recovery addr.
							 // The trustee will decrypt their SSS share off-chain and Encrypt it on-chain with the Recovery accounts Public Key.
							 // Once 3/5 (ie) shares are Decrypted by trustee and Encrypted for recovery addr, then the recovery account may
							 // decrypt the shares and recover their secret (key/phrase/etc)
var pwrExplorerAPI_URL = "https://pwrexplorerbackend.pwrlabs.io"
var pwrFaucetAPI_URL = "https://pwrfaucet.pwrlabs.io/"

type Trustee struct {
	Address string
	PublicKey string
    PrivateKey string
    ECDSAPrivateKey *ecdsa.PrivateKey
	EncryptedSecretShare []byte
}

type PWRTransaction struct {
	TxHash string
	BlockNumber int
	Signature string
}

type Transaction struct {
	TxnHash            string `json:"txnHash"`
	TimeStamp          int    `json:"timeStamp"`
	ValueInUSD         int    `json:"valueInUsd"`
	NonceOrValidation  string `json:"nonceOrValidationHash"`
	Block              int    `json:"block"`
	TxnType            string `json:"txnType"`
	From               string `json:"from"`
	To                 string `json:"to"`
	TxnFeeInUSD        float64 `json:"txnFeeInUsd"`
	TxnFee             int     `json:"txnFee"`
	Value              int     `json:"value"`
}

type Metadata struct {
	TotalItems    int `json:"totalItems"`
	StartIndex    int `json:"startIndex"`
	PreviousPage  int `json:"previousPage"`
	ItemsPerPage  int `json:"itemsPerPage"`
	EndIndex      int `json:"endIndex"`
	NextPage      int `json:"nextPage"`
	TotalPages    int `json:"totalPages"`
	CurrentPage   int `json:"currentPage"`
}

type Data struct {
	Metadata           Metadata      `json:"metadata"`
	HashOfFirstTxnSent string        `json:"hashOfFirstTxnSent"`
	HashOfLastTxnSent  string        `json:"hashOfLastTxnSent"`
	Transactions       []Transaction `json:"transactions"`
	TimeOfFirstTxnSent int           `json:"timeOfFirstTxnSent"`
	TimeOfLastTxnSent  int           `json:"timeOfLastTxnSent"`
}

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


func post(url string, jsonStr string) string {
    var jsonBytes = []byte(jsonStr)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    return string(body)
}

func hexToEncryptedSecrets(hexSt string) []string {
	byteSlice, err := hex.DecodeString(hexSt)
	if err != nil {
		log.Fatal("Error decoding hex string:", err.Error())
	}
	combinedSecrets := string(byteSlice)
	encryptedSecrets := strings.Fields(combinedSecrets)
	//fmt.Println("Original Strings:", encryptedSecrets)
	return encryptedSecrets
}

func encryptedSecretsToHex(secretsBytes [][]byte) string {
	//fmt.Println("encryptedSecretsToHex, secretBytes:", secretsBytes)
	//fmt.Println("encryptedSecretsToHex, len secretBytes:", len(secretsBytes))
	var delimiter = ","
	var secrets []string
	for i:=0; i<len(secretsBytes) -1 ;i++ {
		secrets = append(secrets, hex.EncodeToString(secretsBytes[i]))
	}
	combinedSecrets := strings.Join(secrets, delimiter)
	secretsCombinedBytes := []byte(combinedSecrets)
	secretsHex := hex.EncodeToString(secretsCombinedBytes)
	//fmt.Println("Hex String:", hexString)
	return secretsHex
}

func keccak256(input string) []byte {
	data := []byte(input)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}

func privateKeyToAddress (privateKey *ecdsa.PrivateKey) string {
   publicKey := &privateKey.PublicKey
   address := crypto.PubkeyToAddress(*publicKey)
   return address.Hex()
}

// Public Key recovery from signature + hash //
func bytesToECIESPubKey(publicKeyBytes []byte) (*ecies.PublicKey, error) {
	fmt.Println("bytesToECIESPubKey, publicKeyBytes: ", publicKeyBytes)
	eciesPublicKey, err := ecies.NewPublicKeyFromBytes(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	return eciesPublicKey, nil
}

func recoverPublicKeyFromSignature(signature, messageHash []byte) (*ecdsa.PublicKey, error) {
	pubKey, err := crypto.SigToPub(messageHash, signature)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}
//////////////////////////////////////////////////

func stashSecretsBytes(stasher string, secrets []byte) ([]byte) {
   typeByte := decToBytes(stashSecretsTxType, 1)
   addrBytes := keccak256(stasher)

   var txnBytes []byte
   txnBytes = append(txnBytes, typeByte...)
   txnBytes = append(txnBytes, addrBytes...)
   txnBytes = append(txnBytes, secrets...)
   
   return txnBytes
}

func decToBytes(value, length int) []byte {
   result := make([]byte, length)
   for i := 0; i < length; i++ {
      result[length-1-i] = byte(value >> (8 * i))
   }
   return result
}

func stashSecrets(secrets [][]byte, nonce int, privateKey *ecdsa.PrivateKey) { // stashes ecies encrypted secrets. Not plaintext secrets
    // To stash secrets on Social Recovery VM 222: 
	// Data template
	//
	// first byte is the Tx Type (0 = Register Public Key, 1 = Stash Secrets, 2 = TrusteeRecoverShare)
	// Bytes 2 - 34 are a Keccak256 hash of the stashers address
	// Bytes 35 and after are Hex encoded bytes of the encrypted secret shares delimited by commas
	// 

	var stashersAddr = privateKeyToAddress(privateKey)
	var secretsHex = encryptedSecretsToHex(secrets)
	var secretBytes,_ = hex.DecodeString(secretsHex)
	var secretsBytesStashed = stashSecretsBytes(stashersAddr, secretBytes)

    var dataTx = pwrgo.SendVMDataTx(appId, secretsBytesStashed, nonce, privateKey)
	if dataTx.Success {
		fmt.Println("Stashed secrets: ", dataTx.TxHash)
	} else {
		log.Fatal("Error stashing secrets ", dataTx.Error)
	}
}

func getFirstTxByAddress(address string) *PWRTransaction {
	apiResp := get(pwrExplorerAPI_URL + "/transactionHistory/?address=" + address + "&count=10&page=1")

	var data Data
	err := json.Unmarshal([]byte(apiResp), &data)
	pwrTx := new(PWRTransaction)
	if err != nil {
		fmt.Println("Error:", err)
		return  pwrTx
	}
	if len(data.Transactions) == 0 {
		return pwrTx // empty PWRtx object
	}
	firstTx := data.Transactions[0] // TO-DO: ensure this first tx is actually sent from address. Endpoint returns tx's to and from

	pwrTx.TxHash = firstTx.TxnHash
	pwrTx.BlockNumber = firstTx.Block
	return pwrTx
}

func separateSignatureFromTxnData(transactionData []byte) ([]byte, []byte) {
	signatureStartIndex := len(transactionData)
	for i := len(transactionData) - 1; i >= 0; i-- {
		if transactionData[i] < 128 {
			signatureStartIndex = i
			break
		}
	}
	bufferData := transactionData[:signatureStartIndex]
	signature := transactionData[signatureStartIndex:]

	//bufferStr := hex.EncodeToString(bufferData)
	//signatureStr := hex.EncodeToString(signature)
	//
	//fmt.Println("Buffer Data:", bufferStr)
	//fmt.Println("Signature:", signatureStr)
	return bufferData, signature
}

func main() {
    var privateKeyStr = "0x8e5d3ea16c6a9c73b4b5f49de3cb2c9b57cb6b17fecbd13b0b0b7e745307a4d9"
    var wallet = pwrgo.FromPrivateKey(privateKeyStr)
	var nonce = pwrgo.NonceOfUser(wallet.Address)

    var secretStr = "carpet cat flower chair foot river make image amazing three say shoe" // or just secretStr := privateKeyStr to make priv_key the secret

    // Generate 5 new Address/Public/Private Keys
	var trustees []Trustee

	var trusteeWallet1 = pwrgo.NewWallet()
	var trusteeWallet2 = pwrgo.NewWallet()
	var trusteeWallet3 = pwrgo.NewWallet()
	var trusteeWallet4 = pwrgo.NewWallet()
	var trusteeWallet5 = pwrgo.NewWallet()

	trustees = append(trustees, Trustee{Address: trusteeWallet1.Address, ECDSAPrivateKey: trusteeWallet1.PrivateKey})
	trustees = append(trustees, Trustee{Address: trusteeWallet2.Address, ECDSAPrivateKey: trusteeWallet2.PrivateKey})
	trustees = append(trustees, Trustee{Address: trusteeWallet3.Address, ECDSAPrivateKey: trusteeWallet3.PrivateKey})
	trustees = append(trustees, Trustee{Address: trusteeWallet4.Address, ECDSAPrivateKey: trusteeWallet4.PrivateKey})
	trustees = append(trustees, Trustee{Address: trusteeWallet5.Address, ECDSAPrivateKey: trusteeWallet5.PrivateKey})

   fmt.Println("Trustees: ", trustees)

   var minShareRequired = 3
   var totalShares = 5
   fmt.Printf("Creating Shamirs Shared Secret with %d/%d shares required for recovery", minShareRequired, totalShares)
   fmt.Println("Shared secret plaintext: " + secretStr)

   var shares,_ = sssa.Create(minShareRequired, totalShares, secretStr)
   
   fmt.Println("Shares: ", shares)

	//////////////////// Encrypt Shares ////////////////////////
    if privateKeyStr[0:2] == "0x" {
       privateKeyStr = privateKeyStr[2:]
    }
    
    fmt.Println("Private key str: ", privateKeyStr)
    privateKey, err := crypto.HexToECDSA(privateKeyStr)
    if err != nil {
        log.Fatal(err.Error())
    }
    
    fmt.Println("Private key: ", privateKey)
    fmt.Println("Trustees count: ", len(trustees))

	//pubKeyStr, _ := privateKeyToWallet(privateKey)

    //fmt.Println("Public key str: ", pubKeyStr)

	var sharesEncrypted [][]byte

	for i := 0; i < len(trustees); i++ {
		// POST to PWR Faucet for 100 PWR tokens
		var faucetResult = post(pwrFaucetAPI_URL + "claimPWR/?userAddress=" + trustees[i].Address, `{"userAddress":"`+trustees[i].Address+`"}`)
		fmt.Printf("Trustee %d faucet result: %s\n", i, faucetResult)

		var waitTimeSecs = time.Duration(10)
		fmt.Printf("Waiting %d seconds...\n", waitTimeSecs)
        time.Sleep(waitTimeSecs * time.Second) 

		pwrgo.ReturnBlockNumberOnTx = true // calls blocksCount/ endpoint before broadcasting tx, returns BlockNumber in Response

		// Burn 1 PWR token (so the trustees public key is recoverable from on-chain signature). Any tx should suffice
		var transferTx = pwrgo.TransferPWR("0x0000000000000000000000000000000000000000", "1", 0, trustees[i].ECDSAPrivateKey) // burn 1 PWR
		if transferTx.Success {
			fmt.Printf("[Block #%d] Trustee%d Transfer tx: %s\n", transferTx.BlockNumber, i, transferTx.TxHash)
		} else {
			log.Fatal("Error for Trustee%d Transfer tx: %s\n", i, transferTx.Error)
		}

		fmt.Printf("Waiting %d seconds...\n", waitTimeSecs)
        time.Sleep(waitTimeSecs * time.Second) 

		// //////// Recover trustee public key from transaction hash + signature /////////
		// var pwrTx = getFirstTxByAddress(trustees[i].Address)
		// if pwrTx.BlockNumber == 0 {
	    // 	log.Fatal("Error: Trustee has no transactions on PWR chain (cannot extract public key from signature + hash) for " + trustees[i].Address)
		// }
		// 
		fmt.Println("First txn: ", transferTxHash) //pwrTx.TxHash)

		var txnData []byte
		block := pwrgo.GetBlock(transferTx.BlockNumber)
        for _, txn := range block.Transactions {
			if txn.Hash == transferTxHash {
				// fmt.Println("Found first tx for trustee: ", txn.Hash)
				txnDataBytes, _ := hex.DecodeString(txn.Data)
				txnData = txnDataBytes // TO-DO: instead of reading bytes of TX data, use new RPC call to get tx bytes
			}
        }
		// fmt.Println("First tx for trustee: ", txnData)

		if len(txnData) == 0 {
			log.Fatal("Error: Could not find trustee TX data for ", trustees[i].Address, " on Block number ", pwrTx.BlockNumber)
		}

		txnDataHex, signatureHex := separateSignatureFromTxnData(txnData)

		
		fmt.Println("TxData Buffer: ", txnDataHex)
		fmt.Println("signatureHex: ", signatureHex)

		var ecdsaPubKey,err = recoverPublicKeyFromSignature(signatureHex, txnDataHex)
		if err != nil {
			log.Fatal(err.Error())
		}
		ecdsaPubKeyBytes := crypto.FromECDSAPub(ecdsaPubKey)
		var eciesPubKey,_ = bytesToECIESPubKey(ecdsaPubKeyBytes) // passed empty bytes ERR
		trustees[i].PublicKey = eciesPubKey.Hex(false)
		/////////////////////////////////////////////////////////////////////


        pubKeyStr := trustees[i].PublicKey
    
        if pubKeyStr[0:2] == "0x" {
          pubKeyStr = pubKeyStr[2:]
        }
	    
        pubKey, err := ecies.NewPublicKeyFromHex(pubKeyStr) 
	    if err != nil {
	    	log.Fatal(err.Error())
	    }
	    fmt.Println("Share ", i)
	    fmt.Println("Share bytes: ", []byte(shares[i]))
	    cipherBytes, err := ecies.Encrypt(pubKey, []byte(shares[i]))
	    if err != nil {
	    	log.Fatal(err.Error())
	    }
        fmt.Println("Cipher bytes: ", cipherBytes)
        cipherHex := hex.EncodeToString(cipherBytes)
        fmt.Println("Cipher hex: ", cipherHex)
        fmt.Println("--------------------------")
		// sharesEncrypted[i] = cipherHex
		sharesEncrypted = append(sharesEncrypted, cipherBytes)
		trustees[i].EncryptedSecretShare = cipherBytes
    }

    fmt.Println("Shares encrypted : ", sharesEncrypted)
    fmt.Println("--------------------------")

	stashSecrets(sharesEncrypted, nonce, privateKey)

   //////////////////////////////////////////////////////////////////////////////////////////////////

	/// Remove 2 Trustees (secret shares) randomly ////
	index1 := rand.Intn(len(trustees))
	index2 := rand.Intn(len(trustees))

	for index2 == index1 {
		index2 = rand.Intn(len(trustees))
	}

	var majorityTrustees []Trustee

	for i := 0; i < len(trustees); i++ {
		if i != index1 && i != index2 {
			majorityTrustees = append(majorityTrustees, trustees[i])
		}
	}
   
   fmt.Println("Trustees (2 removed): ", majorityTrustees)
   //////////////////////////////////////////////////////////////


	/////////////// Decryption of 3 shares ///////////////////

	var decryptedSecretShares []string

	for i := 0; i < len(majorityTrustees); i++ {
		privKeyStr := majorityTrustees[i].PrivateKey
		if privKeyStr[0:2] == "0x" {
          privKeyStr = privKeyStr[2:]
        }
		pk, err := ecies.NewPrivateKeyFromHex(privKeyStr)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("[!] pk: ",pk)
		plaintextBytes, _ := ecies.Decrypt(pk, majorityTrustees[i].EncryptedSecretShare)
		fmt.Println("Plaintext bytes: ", plaintextBytes)
		fmt.Println("Plaintext: ", string(plaintextBytes[:]))
		decryptedSecretShares = append(decryptedSecretShares, string(plaintextBytes[:]))
	}
	
    fmt.Println("Decrypted shares: ", decryptedSecretShares)

	///////////////////////////////////


   // Secret Recovery
   var secret,_ = sssa.Combine(decryptedSecretShares)
   fmt.Println("Recovered: ", secret)
}


