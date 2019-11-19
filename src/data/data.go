package data

import (
	"blockchain"
)

type HeartBeatData struct {
	Id          string // sender id
	Addr        string // senders addr
	BlocksJson  string
	PeerMapJson string
}

type PeerList struct {
	SelfId  string
	PeerIds []string
	Length  int
}

type RegisterData struct {
	AssignedId  string
	PeerMapJson string
}

type SyncBlockChain struct {
	Bc blockchain.Blockchain
}
