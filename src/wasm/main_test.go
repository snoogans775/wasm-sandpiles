package main

import (
  "testing"
  "reflect"
)

func testInitSandpile(t *testing.T) {
  newPile = initSandpiles(100, 1000)
  type = reflect.TypeOf(newPile).String()
  targetType := "int[][]"  

  if  type != targetType {
    t.Errorf("got %q, wanted %q", type, targetType) 
  } 
}
