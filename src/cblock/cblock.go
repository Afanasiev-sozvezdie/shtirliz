package cblock

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

//Get key from passphrase by sha256
func PassKey(passphrase string) [32]byte {
	return sha256.Sum256([]byte(passphrase))
}

//Get key from random numbers by sha256
func Rand() [32]byte {
	var randNum [32]byte
	io.ReadFull(rand.Reader, randNum[:])
	return randNum
}

type Shifr struct {
	key [32]byte
}

func XORSlice(sl1, sl2 []byte) []byte {
	var res_len int
	if len(sl1) > len(sl2) {
		res_len = len(sl2)
	} else {
		res_len = len(sl1)
	}

	xor := make([]byte, res_len)

	for i := 0; i < res_len; i++ {
		xor[i] = sl1[i] ^ sl2[i]
	}

	return xor

}

func MakeShifr(passPhrase string) Shifr {
	var shifr Shifr

	shifr.key = PassKey(passPhrase)

	fmt.Println("Key:", shifr.key)

	return shifr
}

func (sh *Shifr) Encode(msg []byte) []byte {
	msg_len := len(msg)

	acc := Rand()
	code := make([]byte, msg_len+32)
	copy(code[0:], acc[:32])

	offset := 0

	for msg_len > 0 {
		if msg_len >= 32 {
			copy(acc[:], XORSlice(acc[:], msg[offset:offset+32]))
		} else {
			copy(acc[:], XORSlice(acc[:], msg[offset:offset+msg_len]))
		}

		copy(code[offset+32:], XORSlice(acc[:], sh.key[:]))

		msg_len -= 32
		offset += 32
	}

	return code
}

func (sh *Shifr) Decode(coded_msg []byte) []byte {

	coded_msg_len := len(coded_msg)

	if coded_msg_len <= 32 {
		return []byte("")
	}
	code_len := coded_msg_len - 32

	msg := make([]byte, code_len)

	code := coded_msg[32:]
	acc := coded_msg[0:32]

	offset := 0

	for code_len > 0 {
		var nokey []byte
		if code_len >= 32 {
			nokey = XORSlice(code[offset:offset+32], sh.key[:])
			copy(msg[offset:], XORSlice(nokey, acc))
			acc = XORSlice(acc, msg[offset:offset+32])
		} else {
			nokey = XORSlice(code[offset:offset+code_len], sh.key[:])
			cp := XORSlice(nokey, acc)
			copy(msg[offset:], cp[:code_len])
			acc = XORSlice(acc, msg[offset:offset+code_len])
		}

		code_len -= 32
		offset += 32
	}

	return msg
}
