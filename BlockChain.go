package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sort"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
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

func Initial(height int32, parentHash string, value string) *Block {
	header := Header{Height: height, Timestamp: time.Now().Unix(), ParentHash: parentHash, Size: 32}
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

type Blockchain struct {
	Chain  map[int32][]Block
	Length int32
}

func (bc *Blockchain) Get(height int32) []Block {
	return bc.Chain[height]
}

func GenesisBlock() *Block {
	return Initial(1, "Genesis Block", "Genesis")
}

func InitBlockchain() *Blockchain {
	bc := Blockchain{}
	bc.Chain = make(map[int32][]Block)
	bc.Chain[1] = append(bc.Chain[1], *GenesisBlock())
	return &bc
}

func (bc *Blockchain) Insert(block Block) {
	for b := range bc.Chain[block.Header.Height] {
		if bc.Chain[block.Header.Height][b].Header.Hash == block.Header.Hash {
			return
		}
	}
	for b := range bc.Chain[block.Header.Height-1] {
		if bc.Chain[block.Header.Height-1][b].Header.Hash == block.Header.ParentHash {
			bc.Chain[block.Header.Height] = append(bc.Chain[block.Header.Height], block)
			if block.Header.Height > bc.Length {
				bc.Length = block.Header.Height
			}
		} else {
			// change this to an error later
			fmt.Println("Parent Hash does not Match")
		}
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
		//fmt.Println("Height: ", k+1, " Blocks: ", bc.Chain[int32(k+1)])
		spew.Dump(bc.Chain[int32(k+1)])
	}
}

// FOR SERVER
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func run() error {
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/", handleGetBlockchain).Methods("GET")

	httpAddr := flag.String("http", "8080", "http listen address")
	log.Println("Listening on ", *httpAddr)
	s := &http.Server{
		Addr:           ":" + *httpAddr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// PROOF OF WORK
const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}
	return pow
}

var blockchain = InitBlockchain()

func main() {
	go func() {
		b3 := Initial(3, "816da3c677f40fce32377aa69f1f488ea7787f0a1f11bb024ee4d246b68f6bd4adb8fbf32e51745e83d47566d40e6c6274be1a6cd79a88269c1e541191d54822", "18964")
		b4 := Initial(4, "d97e47faaf24e025c133ae913d231403b84996616230113777fd46b9311e815134e2d76eb148b6e0541e0692e44bf7addd4320c58e6bfe15895b5bbbdaf4d33e", "89064")
		b5 := Initial(4, "d97e47faaf24e025c133ae913d231403b84996616230113777fd46b9311e815134e2d76eb148b6e0541e0692e44bf7addd4320c58e6bfe15895b5bbbdaf4d33e", "33567")
		b6 := Initial(5, "2d2eabeff80606d2a458bef8bb8e29a16477e9a2a0c85c77d848a576b776f0c1bcd824f4ac535649014d119f3d16c595a61a44bf6bc45133bcf76682b8cd629a", "26795")
		b7 := Initial(6, "343frwgrgregw3", "97313")
		b8 := Initial(7, "mbnotuwh4tg47g", "30783")

		blockchain.Insert(*Initial(1, "aq3r3r32rer232a", "33345"))
		blockchain.Insert(*Initial(2, "bfcb2310fccfda2c6fc35d157694817d94d1fa3066547972555f89ab5d1c66ac5a5757a7659f6fddb45d14d232997ea4718737215a0e3beb5a63cadf4f4a051f", "12976"))
		blockchain.Insert(*b3)
		blockchain.Insert(*b4)
		blockchain.Insert(*b5)
		blockchain.Insert(*b6)
		blockchain.Insert(*b7)
		blockchain.Insert(*b8)
		blockchain.PrintChain()
		//Chain2 := DecodeBlockchainFromJson(Chain.EncodeToJson())
		//fmt.Println(Chain2)
		//fmt.Println(blockchain.EncodeToJson())
		//fmt.Println(len(blockchain.EncodeToJson()))
		//spew.Dump(blockchain)
	}()

	log.Fatal(run())
}
