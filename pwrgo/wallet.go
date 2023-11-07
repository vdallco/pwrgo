package pwrgo

import (
   "crypto/ecdsa"
   "log"
   "github.com/ethereum/go-ethereum/crypto"
   "github.com/ethereum/go-ethereum/common/hexutil"
)

type PWRWallet struct {
   PrivateKey *ecdsa.PrivateKey
   PublicKey string
   Address string
   
}

func FromPrivateKey(privateKeyStr string) *PWRWallet {
   if privateKeyStr[0:2] == "0x" {
      privateKeyStr = privateKeyStr[2:]
   }

   privateKey, err := crypto.HexToECDSA(privateKeyStr)
   if err != nil {
       log.Fatal(err.Error())
   }

   publicKey := &privateKey.PublicKey
   publicKeyStr := hexutil.Encode(crypto.FromECDSAPub(publicKey))

   address := crypto.PubkeyToAddress(*publicKey)
	
   var wallet = new(PWRWallet)

   wallet.PrivateKey = privateKey
   wallet.PublicKey = publicKeyStr
   wallet.Address = address.Hex()

   return wallet
}