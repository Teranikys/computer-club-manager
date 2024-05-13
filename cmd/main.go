package main

import (
	"computer-club-manager/internal/club"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide input file path")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cc, err := club.NewComputerClub(file)
	if err != nil {
		return
	}

	cc.ProcessEvents()
}
