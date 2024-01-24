package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"
	"runtime"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	mainFactoryAddress          = common.HexToAddress("0xB9af59262147673C2016b2b10808411166756ed3")
	IntermediateFactoryBytecode = "6034600d60003960346000f3fe6000548060008114602b573660008037600080366000855af43d6000803e806026573d6000fd5b3d6000f35b6000356000555050"
	limit                       = 6
	dataBytes, _ = hex.DecodeString(IntermediateFactoryBytecode)
  hashByteCode = crypto.Keccak256(dataBytes)
	toFind = "0x0000000"
)

func bigIntToByteArray(n *big.Int) [32]byte {
	var a [32]byte
	copy(a[32-len(n.Bytes()):], n.Bytes())
	return a
}

var (
	requests int
	mu       sync.Mutex
	start    = time.Now()
)

func calculateContractAddress(salt *big.Int, ch chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()

	localSalt := new(big.Int).Set(salt)
	nonce := uint64(1)

	for {
		localSalt.Add(localSalt, big.NewInt(1))
		mu.Lock()
		requests++
		mu.Unlock()
		futureIntemediateFactoryAddress := crypto.CreateAddress2(mainFactoryAddress, bigIntToByteArray(localSalt), hashByteCode)
		futureAddress := crypto.CreateAddress(futureIntemediateFactoryAddress, nonce)
		if strings.HasPrefix(futureAddress.Hex(), toFind) {
			fmt.Printf("%d - %s\n", localSalt, futureAddress.String())
		}

		mu.Lock()
		if requests % 100000000 == 0 {
			rps := float64(requests) / time.Since(start).Seconds()
			log.Printf("RPS: %f\n", rps)
			requests = 0
			start = time.Now()
		}
		mu.Unlock()
	}
}

func main() {
	numCPU := runtime.NumCPU()
	fmt.Printf("Running on %d cores\n\n", numCPU)

	ch := make(chan bool, numCPU)
	var wg sync.WaitGroup

	for i := 0; i < numCPU; i++ {
		myRand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)))
		wg.Add(1)
		salt := big.NewInt(int64(myRand.Intn(100000000000000)))
		go calculateContractAddress(salt, ch, &wg)
	}
	wg.Wait()
}
