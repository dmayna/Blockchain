package block

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Header struct {
	Height     int32
	Timestamp  int64
	Hash       string
	ParentHash string
	Size       int32
	Nonce      string
}

type Block struct {
	Header Header
	Value  string
}

func Initial(height int32, parentHash string, value string) *Block {
	header := Header{Height: height, Timestamp: time.Now().Unix(), ParentHash: parentHash, Size: 32, Nonce: ""}
	b := Block{Header: header, Value: value}
	h := sha512.New()
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
