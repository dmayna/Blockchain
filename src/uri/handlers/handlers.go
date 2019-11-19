package handlers

import (
	"blockchain"
	"blockchain/block"
	"data"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"node"
	"strconv"

	"github.com/gorilla/mux"
)

var Bc = blockchain.InitBlockchain()
var Miner node.Node
var PeerList data.PeerList

// FOR SERVER
func HandleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Bc, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func HandlePeers(w http.ResponseWriter, r *http.Request) {
	// Register miner if not already registered
	if !stringInSlice(Miner.Id, PeerList.PeerIds) {
		PeerList.PeerIds = append(PeerList.PeerIds, Miner.Id)
		//Miner.PeerList = PeerList.PeerIds
	}
	fmt.Fprintln(w, PeerList.PeerIds)
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(Bc); err != nil {
		panic(err)
	}
}

func HandleBlocks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockHeight := vars["height"]
	blockHash := vars["hash"]
	height, err := strconv.ParseInt(blockHeight, 10, 32)
	//hash, err := strconv.Atoi(blockHash)
	var b = block.Block{}
	for i := range Bc.Chain[int32(height)] {
		if Bc.Chain[int32(height)][i].Header.Hash == blockHash {
			b = Bc.Chain[int32(height)][i]
		}
	}
	if b == (block.Block{}) {
		w.WriteHeader(204)
	} else if err != nil {
		w.WriteHeader(500)
	} else {
		fmt.Fprintln(w, "Block: ", b)
	}
}

func HandleShow(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, Bc.Show())
}

func HandleHeartbeatReceive(w http.ResponseWriter, r *http.Request) {
	var b block.Block
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &b); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	Bc.Insert(b)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(b); err != nil {
		panic(err)
	}
}

func HandleStart(w http.ResponseWriter, r *http.Request) {
	go func() {
		Miner.StartTryingNonces(*Bc)
		fmt.Fprintln(w, Bc.Show())
	}()
}
