package pwrgo

import (
   "crypto/ecdsa"
   "github.com/ethereum/go-ethereum/crypto"
   "encoding/hex"
   "math/big"
)

var RPC_ENDPOINT = "https://pwrrpc.pwrlabs.io"


func txBytes(txType int, nonce int, amount *big.Int, recipient string) ([]byte, error) {
   typeByte := decToBytes(txType, 1)
   nonceBytes := decToBytes(nonce, 4)

   amountBytes := amount.Bytes()
   recipientBytes, err := hex.DecodeString(recipient[2:])
   if err != nil {
      return nil, err
   }
   
   paddedNonce := make([]byte, 4)
   copy(paddedNonce[4-len(nonceBytes):], nonceBytes)

   paddedAmount := make([]byte, 8)
   copy(paddedAmount[8-len(amountBytes):], amountBytes)
   
   paddedRecipient := make([]byte, 20)
   copy(paddedRecipient[20-len(recipientBytes):], recipientBytes)
   
   var txnBytes []byte
   txnBytes = append(txnBytes, typeByte...)
   txnBytes = append(txnBytes, paddedNonce...)
   txnBytes = append(txnBytes, paddedAmount...)
   txnBytes = append(txnBytes, paddedRecipient...)
   
   return txnBytes, nil
}

func vmDataBytes(vmId *big.Int, nonce int, data []byte) ([]byte, error) {
   typeByte := decToBytes(5, 1)
   nonceBytes := decToBytes(nonce, 4)
   vmIdBytes := vmId.Bytes()

   paddedNonce := make([]byte, 4)
   copy(paddedNonce[4-len(nonceBytes):], nonceBytes)
   
   var txnBytes []byte
   txnBytes = append(txnBytes, typeByte...)
   txnBytes = append(txnBytes, paddedNonce...)
   txnBytes = append(txnBytes, vmIdBytes...)
   txnBytes = append(txnBytes, data...)
   
   return txnBytes, nil
}

func decToBytes(value, length int) []byte {
   result := make([]byte, length)
   for i := 0; i < length; i++ {
      result[length-1-i] = byte(value >> (8 * i))
   }
   return result
}

//func decToBytes(value, length int) []byte {
//   result := make([]byte, length)
//   for i := 0; i < length; i++ {
//      result[i] = byte(value >> (8 * i))
//   }
//   return result
//}

func signMessage(message []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
   messageHash := crypto.Keccak256(message)
   signature, err := crypto.Sign(messageHash, privateKey)
   if err != nil {
      return nil, err
   }
   
   if signature[64] == 0 || signature[64] ==  1 {
     signature[64] += 27
   } 
   
   return signature, nil
}