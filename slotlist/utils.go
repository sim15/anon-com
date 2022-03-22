// package slotlist

// import (
// 	"crypto/aes"
// 	"encoding/binary"
// 	"fmt"
// 	"math/rand"
// )

// // TODO: IMPORTANT-- DPF OUTPUT SHOULD BE 128 BITS, NOT 64
// func PRGe(seed int64, size uint64) []byte {
// 	rand.Seed(seed)
// 	res := make([]byte, size)
// 	rand.Read(res)
// 	return res
// 	// return []byte{122, 65}

// }

// // add prg data to clients and servers. store its size and initial state
// func PRG(seed int64, size uint64) []byte {

// 	// bc, err := aes.NewCipher([]byte("key3456789012345"))
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// fmt.Printf("The block size is %d\n", bc.BlockSize())
// 	// TODO: FIX FIX
// 	// var dst = make([]byte, 32)
// 	// var src = []byte("sensitive1234567")

// 	// for i := 0; i < 20; i++ {

// 	// 	bc.Encrypt(dst, src)
// 	// }

// 	// bc.Encrypt(dst, src)

// 	seedb := make([]byte, 16)
// 	binary.LittleEndian.PutUint64(seedb, uint64(seed))
// 	bc, err := aes.NewCipher(seedb)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	var dst = make([]byte, 32)

// 	for i := uint64(0); i < (size / 16); i++ {

// 		bc.Encrypt(dst, seedb)
// 	}

// 	return dst

// }

package slotlist

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"math/rand"
)

// TODO: IMPORTANT-- DPF OUTPUT SHOULD BE 128 BITS, NOT 64
func PRGe(seed int64, size uint64) []byte {
	rand.Seed(seed)
	res := make([]byte, size)
	rand.Read(res)
	return res
	// return []byte{122, 65}

}

// add prg data to clients and servers. store its size and initial state
// TODO: take a 128 bit int; two 64 bit
func PRG(seed int64, size uint64) []byte {
	// TODO: send ctr preassigned
	seedb := make([]byte, 16)
	binary.LittleEndian.PutUint64(seedb, uint64(seed))

	bcc, _ := aes.NewCipher(seedb)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// ivVal := make([]byte, len())

	// using the seed itself as the initial value. TODO: is that ok?
	ctr := cipher.NewCTR(bcc, []byte(seedb))
	dest, src := make([]byte, size), make([]byte, size)

	ctr.XORKeyStream(dest, src)

	return dest

}
