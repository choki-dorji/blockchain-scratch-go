package main

import "fmt"

func apply(balances map[string]int64, tx Transaction) error {
	if tx.From == "" || tx.To == "" {
		return fmt.Errorf("sender and receiver are required")
	}
	if tx.From == tx.To {
		return fmt.Errorf("sender and receiver must differ")
	}
	if tx.Value <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if balances[tx.From] < tx.Value {
		return fmt.Errorf("insufficient balance for %s", tx.From)
	}
	const maxBalance int64 = 1<<63 - 1
	if balances[tx.To] > maxBalance-tx.Value {
		return fmt.Errorf("receiver balance would overflow")
	}
	balances[tx.From] -= tx.Value
	balances[tx.To] += tx.Value
	return nil
}
