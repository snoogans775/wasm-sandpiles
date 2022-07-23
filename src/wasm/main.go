package main

import (
  "fmt"
  "syscall/js"
)

func initSandpiles(size int, initialPile int) ([][]int, error) {
  piles := make([][]int, size)
  for i:= range piles {
    piles[i] = make([]int, size)
  }
  piles[size / 2][size / 2] = initialPile


  return piles, nil
}

func updateSandpiles(piles [][]int) {
  for i:= range piles {
    for j:= range piles[i] {
      if piles[i][j] >= 4 {
        // init this pile
        piles[i][j] -= 4

        //update neighbors
        piles[i+1][j]++
        piles[i-1][j]++
        piles[i][j-1]++
        piles[i][j+1]++

      }
    }
  }
}

func initSandpilesWrapper() js.Func {
  initFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
      return "Invalid number of arguments passed"
    }
    pileSize := args[0].Int()
    initPile := args[1].Int()

    fmt.Println("input pileSize %s\n", pileSize)
    fmt.Println("input initPile %s\n", initPile)

    sandpiles, err := initSandpiles(pileSize, initPile)
    if err != nil {
      fmt.Printf("Unable to init sandpiles %s\n", err)
      return err.Error()
    }
    return sandpiles
  })

  return initFunc
}

func main() {
  fmt.Println("Sandpiles Functions loaded")
  //js.Global().Set("initSandpiles", initSandpilesWrapper())
  //js.Global().Set("updateSandpiles", updateSandpiles)

  <-make(chan bool)
}
