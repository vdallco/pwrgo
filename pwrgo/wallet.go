package pwrgo

import (
   "crypto/ecdsa"
   "log"
   "github.com/ethereum/go-ethereum/crypto"
   "github.com/ethereum/go-ethereum/common/hexutil"
)

type PWRWallet struct {
   PrivateKey *ecdsa.PrivateKey
   PrivateKeyStr string
   PublicKey string
   Address string
}

func privateKeyToWallet (privateKey *ecdsa.PrivateKey) *PWRWallet {
   publicKey := &privateKey.PublicKey
   publicKeyStr := hexutil.Encode(crypto.FromECDSAPub(publicKey))
   privateKeyStr := hexutil.Encode(crypto.FromECDSA(privateKey))
   address := crypto.PubkeyToAddress(*publicKey)
	
   var wallet = new(PWRWallet)
   wallet.PrivateKey = privateKey
   wallet.PublicKey = publicKeyStr
   wallet.Address = address.Hex()
   wallet.PrivateKeyStr = privateKeyStr
   return wallet
}

func FromPrivateKey(privateKeyStr string) *PWRWallet {
   if privateKeyStr[0:2] == "0x" {
      privateKeyStr = privateKeyStr[2:]
   }

   privateKey, err := crypto.HexToECDSA(privateKeyStr)
   if err != nil {
       log.Fatal(err.Error())
   }

   return privateKeyToWallet(privateKey)
}

func NewWallet() *PWRWallet {
   privateKey, err := crypto.GenerateKey()
   if err != nil {
       log.Fatal(err.Error())
   }
   return privateKeyToWallet(privateKey)
}