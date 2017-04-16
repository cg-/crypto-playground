package rsa

import (
	"fmt"
	"github.com/op/go-logging"
	"math/big"
	"math/rand"
	"time"
)

type KeyPair struct {
	Size    int
	PubKey  Key
	PrivKey Key
}

func (k *KeyPair) Pad(data []byte) []byte {
	paddedData := make([]byte, k.Size/8)

	for i := range paddedData {
		flip := len(paddedData) - 1 - i
		if i < len(data) {
			paddedData[flip] = data[i]
		} else {
			paddedData[flip] = 0
		}
	}

	log.Infof("Padded to %v", paddedData)
	return paddedData
}

func (k *Key) Encrypt(data []byte) []byte {
	log.Infof("Encrypting %v", data)

	inputInt := big.NewInt(0)
	inputInt.SetBytes(data)

	log.Infof("Input int is: %v", inputInt)
	inputInt.Exp(inputInt, k.Key, nil)
	log.Infof("After Exp int is: %v", inputInt)
	inputInt.Mod(inputInt, k.Mod)
	log.Infof("After Mod int is: %v", inputInt)

	return inputInt.Bytes()
}

func (kp *KeyPair) String() string {
	s1 := ""
	s1 += fmt.Sprintf("PubKey START\n%s\nPubKey END", kp.PubKey.String())
	s1 += fmt.Sprintln()
	s1 += fmt.Sprintf("PrivKey START\n%s\nPrivKey END", kp.PrivKey.String())
	return s1
}

type Key struct {
	Size int
	Key  *big.Int
	Mod  *big.Int
}

func (k *Key) String() string {
	s1 := ""
	s1 += fmt.Sprintf("Key: [%v]\n", k.Key)
	s1 += fmt.Sprintf("Mod: [%v]\n", k.Mod)
	return s1
}

var log = logging.MustGetLogger("crypto-playground")

func RSAGen(size int) *KeyPair {
	p := getPrime(time.Now().UnixNano(), size)
	q := getPrime(time.Now().UnixNano(), size)

	n := big.NewInt(0)
	n.Mul(p, q)

	pMinus1 := big.NewInt(0)
	pMinus1.Sub(p, big.NewInt(1))

	qMinus1 := big.NewInt(0)
	qMinus1.Sub(p, big.NewInt(1))

	toilent := big.NewInt(0)
	toilent.Mul(pMinus1, qMinus1)

	d := big.NewInt(0)
	d.ModInverse(big.NewInt(65537), n)

	toRet := &KeyPair{
		Size: size,
		PubKey: Key{
			Size: size,
			Mod:  n,
			Key:  big.NewInt(65537),
		},
		PrivKey: Key{
			Size: size,
			Mod:  d,
			Key:  big.NewInt(65537),
		},
	}
	return toRet
}

/*
getPrime generates a large prime number by making a bunch of 64 bit numbers and appending
them, ensures it's odd, then tests it for primality probabilistically. If it's not prime,
we add 2 to get another odd number, and repeat.
*/
func getPrime(seed int64, size int) *big.Int {
	rand.Seed(seed)
	randArray := make([]byte, size/8)
	rand.Read(randArray)
	bigV := big.NewInt(0)
	bigV.SetBytes(randArray)
	log.Infof("Generated big int: %v bit length: %v", bigV, bigV.BitLen())

	evenTest := big.NewInt(2)

	log.Info("Testing if it's even...")
	evenTest.Mod(bigV, evenTest)
	if evenTest.Int64() == 0 {
		log.Info("It's even. Adding 1.")
		bigV = bigV.Add(bigV, big.NewInt(1))
	}

	log.Infof("Proceeding with: %v", bigV)

	for !bigV.ProbablyPrime(10) {
		log.Infof("Testing for prime: %v", bigV)
		log.Info("Probably not prime! Adding 2.")
		bigV = bigV.Add(bigV, big.NewInt(2))
	}

	log.Infof("%v is probably prime.", bigV)
	/*
		log.Infof("Let's make sure...")
		if ensurePrime(bigV) {
			return bigV
		} else {
			return getPrime(seed, size)
		}
	*/
	return bigV
}

func ensurePrime(int *big.Int) bool {
	i := big.NewInt(1)
	j := big.NewInt(1)

	sqrt := big.NewInt(0)
	sqrt.Sqrt(int)

	res := big.NewInt(0)
	outerGo := true
	innerGo := true
	for outerGo {
		log.Infof("Outer with i = %v, j = %v", i, j)
		for innerGo {
			log.Infof("Inner with i = %v, j = %v", i, j)
			res.Mul(i, j)
			if res == int {
				return false
			}
			j.Add(j, big.NewInt(1))

			if j.Cmp(sqrt) > 0 {
				j = big.NewInt(1)
				break
			}
		}
		i.Add(i, big.NewInt(1))
		if i.Cmp(sqrt) > 0 {
			break
		}
	}
	return true
}
