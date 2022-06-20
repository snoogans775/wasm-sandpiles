package main

import (
  "fmt"
//  "syscall/js"
)

func initSandpiles(size, initialPile int) [][]int {
  piles := make([][]int, size)
  for i:= range piles {
    piles[i] = make([]int, size)
  }
  piles[size / 2][size / 2] = initialPile

  
  return piles
}

//func jsonWrapper() js.Func {
//  jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
//    if len(args) != 1 {
//      return "Invalid number of arguments passed"
//    }
//    inputJSON := args[0].String()
//    fmt.Printf("input %s\n", inputJSON)
//    pretty, err := prettyJson(inputJSON)
//    if err != nil {
//      fmt.Printf("Unable to convert to JSON %s\n", err)
//      return err.Error()
//    }
//    return pretty
//  })
//  return jsonFunc
//}

func main() {
  fmt.Println("Begin Sandpiles", initSandpiles(100, 1000))
//  js.Global().Set("formatJSON", jsonWrapper())
//  <-make(chan bool)
}
