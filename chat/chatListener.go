package main

import (
    "fmt"
    "github.com/pwrlabs/pwrgo/pwrgo"
    "encoding/hex"
    "time"
)

type Listener struct{}

func NewListener() *Listener {
    return &Listener{}
}

func (l *Listener) Listen() {
    go func() {
        blockNumber := pwrgo.BlocksCount() - 1
        for {
            latestBlockNumber := pwrgo.BlocksCount() - 1

            if blockNumber < latestBlockNumber {
                blockNumber++
                block := pwrgo.GetBlock(blockNumber)
                for _, txn := range block.Transactions {
					if txn.Type == "VM Data" {
						txAppId := txn.VmId
						if txAppId != appId {
							continue
						}
						
						data, _ := hex.DecodeString(txn.Data)
					
						fmt.Printf("Message From %s: %s\n\n> ", txn.From, string(data))
					}
                }
            }

            time.Sleep(10 * time.Millisecond)
        }
    }()
}
