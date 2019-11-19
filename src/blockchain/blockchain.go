package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"blockchain/block"

	"github.com/davecgh/go-spew/spew"
)

const Difficulty = 4

type Blockchain struct {
	Chain  map[int32][]block.Block
	Length int32
}

func (bc *Blockchain) Get(height int32) []block.Block {
	return bc.Chain[height]
}

func GenesisBlock() *block.Block {
	return block.Initial(1, "Genesis Block", "Genesis")
}

func InitBlockchain() *Blockchain {
	bc := Blockchain{}
	bc.Chain = make(map[int32][]block.Block)
	bc.Chain[1] = append(bc.Chain[1], *GenesisBlock())
	bc.Length = 1
	return &bc
}

func (bc *Blockchain) Insert(b block.Block) {
	//verify Nonce
	h := sha256.New()
	hash := b.Header.ParentHash + b.Header.Nonce + b.Value
	h.Write([]byte(hash))
	pow_answer := hex.EncodeToString(h.Sum(nil))
	runes := []rune(pow_answer)
	for i := 0; i <= Difficulty; i++ {
		if string(runes[i]) != "0" {
			fmt.Printf("Rune %v is '%c'\n", i, runes[i])
			return
		}
	}
	// we dont store duplicate blocks
	for i := range bc.Chain[b.Header.Height] {
		if bc.Chain[b.Header.Height][i].Header.Hash == b.Header.Hash {
			return
		}
	}
	// checking if parentHash matches blockchain
	for i := range bc.Chain[b.Header.Height-1] {
		if bc.Chain[b.Header.Height-1][i].Header.Hash == b.Header.ParentHash {
			bc.Chain[b.Header.Height] = append(bc.Chain[b.Header.Height], b)
			if b.Header.Height > bc.Length {
				bc.Length = b.Header.Height
			}
		} else {
			// change this to an error later
			fmt.Println("Parent Hash does not Match")
			fmt.Println(bc.Chain[b.Header.Height-1][0].Header.Hash)
			fmt.Println(b.Header.ParentHash)
		}
	}
}

type BlocksJson []string

func (bc *Blockchain) EncodeToJson() []string {

	blocksJson := []string{}
	for k, _ := range bc.Chain {
		for block := range bc.Chain[k] {
			blocksJson = append(blocksJson, bc.Chain[k][block].EncodeToJson())
		}
	}
	return blocksJson
}

func DecodeBlockchainFromJson(inData []string) *Blockchain {
	var bc Blockchain
	for b := range inData {
		bc.Insert(*block.DecodeBlockFromJson(inData[b]))
	}

	return &bc
}

func (bc *Blockchain) PrintChain() {
	var Keys []int32
	for k := range bc.Chain {
		Keys = append(Keys, k)
	}
	sort.Slice(Keys, func(i, j int) bool { return Keys[i] < Keys[j] })
	for k := range Keys {
		//fmt.Println("Height: ", k+1, " Blocks: ", bc.Chain[int32(k+1)])
		spew.Dump(bc.Chain[int32(k+1)])
	}
}

func (bc *Blockchain) GetLatestBlocks() []block.Block {
	return bc.Get(bc.Length)
}

func (bc *Blockchain) GetParentBlock(b block.Block) *block.Block {
	//potentialParentBlocks := bc.Get(block.Header.Height - 1)
	for i := range bc.Get(b.Header.Height - 1) {
		if bc.Chain[b.Header.Height-1][i].Header.Hash == b.Header.ParentHash {
			return &bc.Chain[b.Header.Height-1][i]
		}
	}
	return &block.Block{}
}

func (blockchain *Blockchain) Show() string {
	rs := ""
	var idList []int
	for id := range blockchain.Chain {
		idList = append(idList, int(id))
	}
	sort.Ints(idList)
	for _, id := range idList {
		var hashs []string
		for _, block := range blockchain.Chain[int32(id)] {
			hashs = append(hashs, block.Header.Hash+"<="+block.Header.ParentHash)
		}
		sort.Strings(hashs)
		rs += fmt.Sprintf("%v: ", id)
		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}
		rs += "\n"
	}
	sum := sha256.Sum256([]byte(rs))
	rs = fmt.Sprintf("This is the BlockChain: %s\n", hex.EncodeToString(sum[:])) + rs
	return rs
}
