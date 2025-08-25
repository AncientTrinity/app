// Filename: cmd/api/main.go

package main

import (
	"fmt"
)

func printUB() string {
    return "Hello, UB!"
}

func main() {
    greeting := printUB()
    fmt.Println(greeting)
}
