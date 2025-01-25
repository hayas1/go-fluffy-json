[![Build Status](https://github.com/hayas1/go-fluffy-json/actions/workflows/go.yml/badge.svg)](https://github.com/hayas1/go-fluffy-json/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayas1/go-fluffy-json)](https://goreportcard.com/report/github.com/hayas1/go-fluffy-json)

# fluffyjson
fluffyjson can deal with JSON fluffily.
- Compatible with `encode/json` and better interface than `any`(`interface{}`).
- Useful methods to handle JSON value such as cast, access, visit, and so on.
- Pure Go implementation.

# Usage
## Install
TODO

## Unmarshal JSON
```go
import fluffyjson "github.com/hayas1/go-fluffy-json"

target := `{"hello":"world"}`
var v struct {
    Hello fluffyjson.Value `json:"hello"`
}
if err := json.Unmarshal([]byte(target), &v); err != nil {
    panic(err)
}

world, err := v.Hello.AsString()
if err != nil {
    panic(err)
}
fmt.Println(world)
// Output: world
```

## Marshal JSON
TODO

## Cast JSON value as Go struct
TODO

## Access to JSON value
TODO

## Visit JSON as tree
TODO
