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

func (k *Key) Encrypt(b []byte) []byte {
	log.Infof("Encrypting %v", b)
	toRet := []byte{}
	res := big.NewInt(0)
	for i := range b {
		m := int8(b[i])
		log.Infof("m: %v", m)
		log.Infof("m 64: %v", int64(m))
		res.Exp(big.NewInt(int64(m)), k.Exponent, nil)
		res.Mod(res, k.Mod)
		log.Info(res)
	}
	return toRet
}

func (kp *KeyPair) String() string {
	s1 := ""
	s1 += fmt.Sprintf("PubKey START\n%s\nPubKey END", kp.PubKey.String())
	s1 += fmt.Sprintln()
	s1 += fmt.Sprintf("PrivKey START\n%s\nPrivKey END", kp.PrivKey.String())
	return s1
}

type Key struct {
	Size     int
	Exponent *big.Int
	Mod      *big.Int
}

func (k *Key) String() string {
	s1 := ""
	s1 += fmt.Sprintf("Exponent: [%v]\n", k.Exponent)
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
			Size:     size,
			Mod:      n,
			Exponent: big.NewInt(65537),
		},
		PrivKey: Key{
			Size:     size,
			Mod:      d,
			Exponent: big.NewInt(65537),
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
		log.Info("Testing for prime: %v", bigV)
		log.Info("Probably not prime! Adding 2.")
		bigV = bigV.Add(bigV, big.NewInt(2))
	}

	log.Infof("%v is probably prime.", bigV)
	return bigV
}
