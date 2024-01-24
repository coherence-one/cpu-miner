package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"time"
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

func calculateContractAddress(salt *big.Int) {
	start := time.Now()
	localSalt := new(big.Int).Set(salt)
	nonce := uint64(1)

	// keccak := crypto.Keccak256Hash(([]byte{0xff}, b.Bytes(), salt[:], vbg)[12:])
	futureIntemediateFactoryAddress := crypto.CreateAddress2(mainFactoryAddress, bigIntToByteArray(localSalt), hashByteCode)
	futureAddress := crypto.CreateAddress(futureIntemediateFactoryAddress, nonce)
	elapsed := time.Since(start)
	fmt.Printf("It took %s %s\n\n", elapsed, futureAddress)

	fmt.Printf("%d - %s\n", localSalt, futureAddress.String())

}

func main() {
	salt := big.NewInt(int64(786686378459217))
	calculateContractAddress(salt)
}
