package rsa

import (
   "math/rand"
   "strconv"
   "github.com/op/go-logging"
   "time"
   "math/big"
)

var log = logging.MustGetLogger("crypto-playground")

func RSAGen() {
   p := getPrime(time.Now().UnixNano())
   q := getPrime(time.Now().UnixNano())

   log.Infof("P: %v Q: %v", p, q)
}

/*
getPrime generates a large prime number by making a bunch of 64 bit numbers and appending
them, ensures it's odd, then tests it for primality probabilistically. If it's not prime,
we add 2 to get another odd number, and repeat.
 */
func getPrime(seed int64) *big.Int {
   rand.Seed(seed)
   v1 := strconv.FormatInt(rand.Int63(), 10)
   v2 := strconv.FormatInt(rand.Int63(), 10)
   v3 := strconv.FormatInt(rand.Int63(), 10)
   v4 := strconv.FormatInt(rand.Int63(), 10)
   v5 := strconv.FormatInt(rand.Int63(), 10)
   v6 := strconv.FormatInt(rand.Int63(), 10)
   v7 := strconv.FormatInt(rand.Int63(), 10)
   v8 := strconv.FormatInt(rand.Int63(), 10)
   v := v1 + v2 + v3 + v4 +v5 +v6 +v7 + v8
   log.Infof("Generated random 64 bit numbers. Appended to: %v", v)
   bigV := big.NewInt(0)
   bigV.SetString(v, 10)
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

