package main

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func printHistoryHash() error {
	data, err := os.ReadFile("database/tx.db")
	if err != nil {
		return err
	}
	sum := sha256.Sum256(data)
	fmt.Printf("%x\n", sum)
	return nil
}
