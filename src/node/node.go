package node

import (
	"blockchain"
	"blockchain/block"
	"crypto/sha256"
	"crypto/sha512"
	"data"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

const Difficulty = blockchain.Difficulty

type Node struct {
	Id       string
	Port     string
	PeerList data.PeerList
}

func (miner *Node) SendHeartBeat(blockData data.HeartBeatData) {
	// send a heartbeat post request to /heartbeat/receive for every peer in the PeerList
	for peer := range miner.PeerList.PeerIds {
		fmt.Println(peer)
	}
}

func (miner *Node) StartTryingNonces(bc blockchain.Blockchain) {
	for {
		latestBlock := bc.GetLatestBlocks()[0]
		header := block.Header{Height: latestBlock.Header.Height + 1, Timestamp: time.Now().Unix(), ParentHash: latestBlock.Header.Hash, Size: 32, Nonce: ""}
		h := sha512.New()
		// for random value to start finding nonce
		hashTime := string(strconv.FormatInt(time.Now().Unix(), 10)) + miner.Port
		h.Write([]byte(hashTime))
		nonceAttempt := hex.EncodeToString(h.Sum(nil))
		// this code takes the first 16 hexes from random hex value
		runes := []rune(nonceAttempt)
		// ... Convert back into a string from rune slice.
		nonceAttempt = string(runes[0:16])
		fmt.Println("----------------------- RUNE SUBSTRING:", nonceAttempt)
		// check to see if nonce passes difficulty test
		notFound := true
		for notFound {
			h = sha256.New()
			hash := string(header.ParentHash) + string(nonceAttempt) + "value"
			h.Write([]byte(hash))
			hashAttempt := hex.EncodeToString(h.Sum(nil))
			for i, char := range hashAttempt { // change this for hash outout not nonce number itself!!!
				if string(char) != "0" && i < Difficulty {
					//nonce fails leave loop, increment nonce and try again
					break
				}
				if string(char) == "0" && i == Difficulty {
					// nonce found add new block to blockchain
					fmt.Println("nonce found------------------------------------")

					header.Nonce = nonceAttempt
					b := block.Block{Header: header, Value: "value"} // figure out value
					h := sha512.New()
					hash := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash + string(b.Header.Size) + b.Value
					h.Write([]byte(hash))
					b.Header.Hash = hex.EncodeToString(h.Sum(nil))
					bc.Insert(b)
					notFound = false
					break
				}
			}
			if notFound == true {
				// this code is to increment nonce value if wrong
				newNonce, _ := strconv.ParseInt(nonceAttempt, 16, 64)
				newNonce += 1
				if newNonce < 0 {
					newNonce *= -1
				}

				nonceAttempt = strconv.FormatInt(newNonce, 16)
				if nonceAttempt == "7fffffffffffffff" {
					nonceAttempt = "1000000000000000"
				}
			}
		}
	}
}
