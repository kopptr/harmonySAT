package main

import (
   "fmt"
   "rand"
   "time"
   "flag"
)

var (
   seed = flag.Int64("seed", time.Nanoseconds(), "seed for random number generator")
)

func main() {

   flag.Parse()
   rand.Seed(*seed)

   fmt.Print("Hello, world!\n")
}
