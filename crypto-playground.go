/*
Crypto Playground is just me playing around with crypto/hashing concepts
from class to make sure I understand them...
 */

package main

import (
   "github.com/cg-/crypto-playground/rsa"
   "github.com/op/go-logging"
   "os"
)

var log = logging.MustGetLogger("crypto-playground")
var format = logging.MustStringFormatter(
   `%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)


func init(){
   backend := logging.NewLogBackend(os.Stderr, "", 0)
   backendFormatter := logging.NewBackendFormatter(backend, format)
   //backendLeveled := logging.AddModuleLevel(backend)
   //backendLeveled.SetLevel(logging.INFO, "")
   logging.SetBackend(backendFormatter)
   log.Info("Finished main init function.")
}


func main() {
   log.Info("Starting main program.")
   rsa.RSAGen()
}

func check(e error){
   if e != nil {
      log.Fatalf(e.Error())
      os.Exit(5)
   }
}
