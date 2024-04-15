package main

import "os"

func main() {
	os.Exit(0) // want "shouldn't call os.Exit in main function"
}
