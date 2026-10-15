package main

import "testing"

func TestRejectInvalidTransfers(t *testing.T) {
	cases := []Transaction{
		{From: "pema", To: "dorji", Value: 151},
		{From: "dorji", To: "pema", Value: 0},
		{From: "dorji", To: "pema", Value: -5},
		{From: "dorji", To: "dorji", Value: 10},
		{From: "", To: "pema", Value: 10},
	}
	for _, tx := range cases {
		b := map[string]int64{"dorji": 850, "pema": 150}
		if err := apply(b, tx); err == nil {
			t.Fatalf("expected rejection for %+v", tx)
		}
		if b["dorji"] != 850 || b["pema"] != 150 {
			t.Fatalf("rejected transfer changed balances: %v", b)
		}
	}
}
func TestTransferPreservesSupply(t *testing.T) {
	b := map[string]int64{"dorji": 850, "pema": 150}
	tx := Transaction{From: "dorji", To: "pema", Value: 100}
	if err := apply(b, tx); err != nil {
		t.Fatal(err)
	}
	if b["dorji"] != 750 || b["pema"] != 250 {
		t.Fatalf("unexpected balances: %v", b)
	}
	if b["dorji"]+b["pema"] != 1000 {
		t.Fatal("transfer changed the total supply")
	}
}
func TestOrderedHistory(t *testing.T) {
	b := map[string]int64{"dorji": 1000, "pema": 0}
	history := []Transaction{
		{From: "dorji", To: "pema", Value: 200},
		{From: "pema", To: "dorji", Value: 50},
	}
	for _, tx := range history {
		if err := apply(b, tx); err != nil {
			t.Fatal(err)
		}
	}
	if b["dorji"] != 850 || b["pema"] != 150 {
		t.Fatalf("incorrect replay result: %v", b)
	}
	fresh := map[string]int64{"dorji": 1000, "pema": 0}
	if err := apply(fresh, history[1]); err == nil {
		t.Fatal("Pema cannot spend before receiving tokens")
	}
}
