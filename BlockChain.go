package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type Header struct {
	Height     int32
	Timestamp  int64
	Hash       string
	ParentHash string
	Size       int32
}

type Block struct {
	Header Header
	Value  string
}

func initial(height int32, parentHash string, value string) *Block {
	header := Header{Height: height, Timestamp: time.Now().Unix(), ParentHash: parentHash, Size: 32}
	b := Block{Header: header, Value: value}
	h := sha256.New()
	hash := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash + string(b.Header.Size) + b.Value
	h.Write([]byte(hash))
	b.Header.Hash = hex.EncodeToString(h.Sum(nil))
	return &b
}

func (b *Block) print() {
	fmt.Println(b.Header.Height)
	fmt.Println(b.Header.Timestamp)
	fmt.Printf("%x", b.Header.Hash)
	fmt.Println(b.Header.ParentHash)
	fmt.Println(b.Header.Size)
	fmt.Println(b.Value)
}

func DecodeBlockFromJson(inData string) *Block {
	var b Block
	err := json.Unmarshal([]byte(inData), &b)
	if err != nil {
		panic(err)
	}
	return &b
}

func (b *Block) EncodeToJson() string {
	e, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return string(e)
}

type Blockchain struct {
	Chain  map[int32][]Block
	Length int32
}

func (bc *Blockchain) Get(height int32) []Block {
	return bc.Chain[height]
}

func (bc *Blockchain) Insert(block Block) {
	if bc.Chain == nil {
		bc.Chain = make(map[int32][]Block)
	}
	for b := range bc.Chain[block.Header.Height] {
		if bc.Chain[block.Header.Height][b].Header.Hash == block.Header.Hash {
			return
		}
	}
	bc.Chain[block.Header.Height] = append(bc.Chain[block.Header.Height], block)
	if block.Header.Height > bc.Length {
		bc.Length = block.Header.Height
	}
}

func (bc *Blockchain) EncodeToJson() []string {
	blocksJson := []string{}
	for k, _ := range bc.Chain {
		for block := range bc.Chain[k] {
			//fmt.Println(bc.Chain[k][block])
			blocksJson = append(blocksJson, bc.Chain[k][block].EncodeToJson())
		}
	}
	return blocksJson
}

func DecodeBlockchainFromJson(inData string) *Blockchain {
	var b Blockchain
	err := json.Unmarshal([]byte(inData), &b)
	if err != nil {
		panic(err)
	}
	return &b
}

func (bc *Blockchain) PrintChain() {
	var Keys []int32
	for k := range bc.Chain {
		Keys = append(Keys, k)
	}
	sort.Slice(Keys, func(i, j int) bool { return Keys[i] < Keys[j] })
	for k := range Keys {
		fmt.Println("Height: ", k+1, " Blocks: ", bc.Chain[int32(k+1)])
	}
}

func main() {
	//  b := initial(1,"aq3r3r32rer232a","33345")
	//b1 := initial(2,"werweer232werf","12976")
	b3 := initial(3, "afsdfdsfsdaa2f", "18964")
	b4 := initial(4, "sewr4twsdfdsff", "89064")
	b5 := initial(4, "ar3qrqrfdsfccd", "33567")
	b6 := initial(5, "45tvefwefwtrtr", "26795")
	b7 := initial(6, "343frwgrgregw3", "97313")
	b8 := initial(7, "mbnotuwh4tg47g", "30783")
	//b.print()
	//fmt.Println(b.EncodeToJson())
	//  b2 := DecodeBlockFromJson(b.EncodeToJson())
	//b2.print()
	m := make(map[int32][]Block)
	m[1] = append(m[1], *initial(1, "aq3r3r32rer232", "33443"))
	m[2] = append(m[2], *initial(2, "3sgrer85yeff32", "87844"))
	m[3] = append(m[3], *initial(3, "wegq3r3rdgsgr232", "45443"))
	m[2] = append(m[2], *initial(2, "bfgny32rer232", "90441"))
	//Chain := Blockchain{Chain: m, Length: 4}
	blockchain := Blockchain{}
	blockchain.Insert(*initial(1, "aq3r3r32rer232a", "33345"))
	blockchain.Insert(*initial(2, "werweer232werf", "12976"))
	blockchain.Insert(*b3)
	blockchain.Insert(*b4)
	blockchain.Insert(*b5)
	blockchain.Insert(*b6)
	blockchain.Insert(*b7)
	blockchain.Insert(*b8)
	//blockchain.PrintChain()
	//Chain2 := DecodeBlockchainFromJson(Chain.EncodeToJson())
	//fmt.Println(Chain2)
	fmt.Println(blockchain.EncodeToJson())
	fmt.Println(len(blockchain.EncodeToJson()))
}
