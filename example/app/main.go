package main

import (
	"fmt"

	"github.com/chunkgar/gokit/app"
)

func main() {
	// Your Go code here
	app.NewApp("test", "test", app.WithNoConfig(), app.WithRunFunc(func(basename string) error {
		fmt.Println("Hello, World!")

		return nil
	})).Run()
}
