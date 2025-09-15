package main

import (
	"os"
	"wget/wget"
)

func main() {
	os.Exit(wget.CLI(os.Args[1:]))
}
