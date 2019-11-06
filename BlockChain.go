package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
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
	Nonce      string
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
	//verify Nonce
	h := sha256.New()
	hash := block.Header.ParentHash + block.Header.Nonce + block.Header.Hash
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
	for b := range bc.Chain[block.Header.Height] {
		if bc.Chain[block.Header.Height][b].Header.Hash == block.Header.Hash {
			return
		}
	}
	// checking if parentHash matches blockchain
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
	var b Blockchain
	for block := range inData {
		b.Insert(*DecodeBlockFromJson(inData[block]))
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

func handlePeers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Todo Index!")
}

func run() error {
	r := mux.NewRouter()
	r.HandleFunc("/", handleGetBlockchain).Methods("GET")
	r.HandleFunc("/peer", handlePeers)

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

type Node struct {
	PeerList []string
}

var blockchain = InitBlockchain()

func main() {
	go func() {
		// You can see blockchain in the terminal or in browser (easiest)
		// at http://localhost:8080/
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
		blockchain.Insert(*Initial(11, "c2e4f5eebed86504fbb982f9f2f876386cea9e73d84a989930824172cb8967b5b6c93f00487c0b303418ab95ad67b3c8d3972ddbd86f5227deb47e98a4912c25", "97313"))
		blockchain.PrintChain()
		fmt.Println(DecodeBlockFromJson(b3.EncodeToJson()))
		fmt.Println(blockchain.EncodeToJson())
		Chain2 := DecodeBlockchainFromJson(blockchain.EncodeToJson())
		fmt.Println(Chain2)
		//fmt.Println(blockchain.EncodeToJson())
		//fmt.Println(len(blockchain.EncodeToJson()))
		//spew.Dump(blockchain.EncodeToJson())
	}()
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				fmt.Println(ip)
			case *net.IPAddr:
				ip = v.IP
				fmt.Println(ip)
			}
			// process IP address
			fmt.Println(ip)
			log.Fatal(run())
		}
	}
}
