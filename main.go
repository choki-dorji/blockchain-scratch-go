package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("use: balances | transfer FROM TO AMOUNT | hash")
	}
	g, balances, err := loadState()
	if err != nil {
		return err
	}
	switch args[0] {
	case "hash":
		if len(args) != 1 {
			return fmt.Errorf("use: hash")
		}
		return printHistoryHash()

	case "balances":
		if len(args) != 1 {
			return fmt.Errorf("use: balances")
		}
		fmt.Println("Chain:", g.ChainID)
		names := make([]string, 0, len(balances))
		for name := range balances {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Printf("%s: %d\n", name, balances[name])
		}
		return nil
	case "transfer":
		if len(args) != 4 {
			return fmt.Errorf("use: transfer FROM TO AMOUNT")
		}
		amount, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return fmt.Errorf("amount must be a whole number")
		}
		tx := Transaction{
			From: args[1], To: args[2], Value: amount,
		}
		if err := apply(balances, tx); err != nil {
			return err
		}
		if err := appendTransaction(tx); err != nil {
			return fmt.Errorf("save failed; inspect history: %w", err)
		}
		fmt.Printf("Saved: %s sent %d to %s\n",
			tx.From, tx.Value, tx.To)
		return nil
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}
