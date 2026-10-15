package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func loadState() (Genesis, map[string]int64, error) {
	var g Genesis
	data, err := os.ReadFile("database/genesis.json")
	if err != nil {
		return g, nil, err
	}
	if err := json.Unmarshal(data, &g); err != nil {
		return g, nil, err
	}
	if g.ChainID == "" || len(g.Balances) == 0 {
		return g, nil, fmt.Errorf("invalid genesis")
	}
	balances := make(map[string]int64)
	for name, amount := range g.Balances {
		if name == "" || amount < 0 {
			return g, nil, fmt.Errorf("invalid initial balance")
		}
		balances[name] = amount
	}
	if err := replay(balances); err != nil {
		return g, nil, err
	}
	return g, balances, nil
}

func replay(balances map[string]int64) error {
	f, err := os.Open("database/tx.db")
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		line++
		var tx Transaction
		if err := json.Unmarshal(scanner.Bytes(), &tx); err != nil {
			return fmt.Errorf("line %d: %w", line, err)
		}
		if err := apply(balances, tx); err != nil {
			return fmt.Errorf("line %d: %w", line, err)
		}
	}
	return scanner.Err()
}
func appendTransaction(tx Transaction) error {
	data, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	f, err := os.OpenFile("database/tx.db",
		os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	n, writeErr := f.Write(data)
	if writeErr == nil && n != len(data) {
		writeErr = io.ErrShortWrite
	}
	if writeErr == nil {
		writeErr = f.Sync()
	}
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
