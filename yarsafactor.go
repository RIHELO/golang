package main

import (
 "fmt"
 "math/big"
 "runtime"
 "sync"
 "os"
)

var wg sync.WaitGroup

func isEven(bgNumber *big.Int) bool {
  bigstr := bgNumber.String()
  last := bigstr[len(bigstr)-1:]
  switch last {
    case "0", "2", "4", "5", "6", "8":
      return true
    case "1", "3", "7", "9":
      return false
  }
  return false
}

func isPerfectSquare(bgNumber *big.Int) bool {
  bigstr := bgNumber.String()
  last := bigstr[len(bigstr)-1:]
  switch last {
    case "0", "1", "4", "5", "6", "9":
      sqrt := big.NewInt(1).Sqrt(bgNumber)
      test := big.NewInt(1).Mod(bgNumber,sqrt)
      if test.Cmp(big.NewInt(0)) == 0 {
         return true
      } else {
         return false
      }
    case "2", "3", "7", "8":
      return false
  }
  return false
}

func factor(process int, rsaNumber *big.Int, rng *big.Int, initialValue *big.Int) {
  defer wg.Done()
  fmt.Println("Process:",process, "initialValue:", initialValue, " range:", rng)
  innerSquareSide := initialValue
  innerSquare  := big.NewInt(0).Exp(innerSquareSide, big.NewInt(2), nil)
  quarterSquare := big.NewInt(0).Add(innerSquare,rsaNumber)
  foundPerfectSquare:=isPerfectSquare(quarterSquare)
  for foundPerfectSquare == false && innerSquareSide.Cmp(rng) <= 0{ 
    innerSquareSide.Add(innerSquareSide, big.NewInt(1))
    innerSquare = big.NewInt(0).Exp(innerSquareSide, big.NewInt(2), nil)
    quarterSquare = big.NewInt(0).Add(innerSquare,rsaNumber)
    foundPerfectSquare = isPerfectSquare(quarterSquare)
  }
  if foundPerfectSquare {
    finalSquare := big.NewInt(0).Mul(big.NewInt(4),quarterSquare)
    externalSquareSizeLength := new(big.Int).Sqrt(finalSquare)
    innerSquareSide4 := big.NewInt(0).Sqrt(new(big.Int).Mul(innerSquare,big.NewInt(4)))
    difference := big.NewInt(0).Sub(externalSquareSizeLength,innerSquareSide4)
    sideP := big.NewInt(0).Div(difference,big.NewInt(2))
    if isEven(sideP) == false {
      sideQ := big.NewInt(0).Add(sideP,big.NewInt(1).Mul(innerSquareSide,big.NewInt(2)))
      fmt.Println("RSA:", rsaNumber,"Process:",process, "initialValue:",initialValue, " range:", rng,"innerSquareSide: ", innerSquareSide, "p: ", sideP, "q: ", sideQ)
      os.Exit(0)
    }
  }
}

func main() {
  fmt.Println("Version: ", runtime.Version(), "NumCPU:", runtime.NumCPU(), "GOMAXPROCS",runtime.GOMAXPROCS(0))
  var rsa = os.Args[1] 
  rsaNumber , _:= big.NewInt(0).SetString(rsa,10)
  //sqrtRsaNumber := big.NewInt(0).Sqrt(rsaNumber)
  numProcesses := int64(8)
  initialValue, _ := big.NewInt(0).SetString(os.Args[2],10) 
  rangePerProcessor, _ := big.NewInt(0).SetString(os.Args[3],10) 
  rng := big.NewInt(0).Add(initialValue,rangePerProcessor)
  for i := 0; i < int(numProcesses); i++ {
    wg.Add(1)
    go factor(i, rsaNumber, rng, initialValue)
    rng = big.NewInt(0).Add(rng,rangePerProcessor)
    initialValue = big.NewInt(0).Add(initialValue,rangePerProcessor)
  }
  wg.Wait()
}

