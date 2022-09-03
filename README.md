# WASM with Go
## Objective
Build an implementation of the sandpiles algorithm with a graphical display using Go as the data processing layer and the browser as the presentation layer
## Open Questions
What is the fastest way to send data between Go and JS?
Is Go fast at updating the DOM?
How hard is it to build/deploy a Go-wasm app?
## 1st step: config
A bit tricky and magical.
The wasm_exec.js file is the magical glue that makes the .wasm output interface with the browser.
M A G I C A L G L U E
Not too hard, make sure GOOS = js and GOARCH = wasm
## 2nd step: Write in Go
Go is fun, Go is fast.
```
func initSandpilesWrapper() js.Func {
	initFunc := js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
            // Stuff that happens
            // return interface{}
        }
    )
    return initFunc
}
js.Global().Set("initSandpiles", initSandpilesWrapper())
```
The break down:
- js.Func is a function that can be accessed by JS
- js.FuncOf() turns a Go function into a js.Func
- js.Value is important because we need to Go-ify anything that comes in
- we would like to have the global 'this' and local 'args' please

## Issues
### Return types on Go js.FuncOf()
Not all Go data structures can be returned via the interface{} return type in js.FuncOf()
For example:
- No multidimensional slices
- In fact, no slices at all D:`
### I kinda suck at Go
My weak dynamically typed brain demands satisfaction.
I just made the variable in this function, it's not changing while we run, I swear.
Really please let me use it as a const.
Pretty please!?
Wait, okay then, Go devs must turn slices into arrays all the time.
No? What, why is this so hard to google?
```copy(varLead.Magic[:], someSlice[0:4])```
Does this operation really merit the word 'Magic'?
Okay I can just loop and load an array, but not dynamically, we'll try the magic thing later.
Only nerds use []byte, I'm sticking with []int
## Time for dependencies!
I am going to try @markfarnan/go-canvas
**behold the monstrosity of typecasting**
```
s.piles[row+1][col] += int(math.Floor(float64(s.toppleThreshold) / 4))
```
Performance is not great on canvas...
## Fully functioning!
Web Canvas API is the bottleneck. There is mention of using goroutines to avoid this.
Go can log like a champ. Logging in JS can really load up the event loop
.Truthy() is js derpiness leaking in to Go interop
