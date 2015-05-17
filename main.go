package main

import "os"

func main() {
	config := NewConfig()

	if err := config.Load(os.Args); err != nil {
		panic(err)
	}

	NewDating(config).Run()
}
