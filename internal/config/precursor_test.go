package config

import "testing"

func TestCommandRunChain(t *testing.T) {
	p := &Project{
		Commands: []Command{
			{Alias: "tunnel", Abbrev: "cf", Shell: "cloudflared"},
			{Alias: "run", Abbrev: "r", Shell: "npm run dev", Precursor: "cf"},
		},
	}
	chain, err := p.CommandRunChain("run")
	if err != nil {
		t.Fatal(err)
	}
	if len(chain) != 2 || chain[0] != "tunnel" || chain[1] != "run" {
		t.Fatalf("chain: %v", chain)
	}
}

func TestCommandRunChainCycle(t *testing.T) {
	p := &Project{
		Commands: []Command{
			{Alias: "a", Shell: "a", Precursor: "b"},
			{Alias: "b", Shell: "b", Precursor: "a"},
		},
	}
	_, err := p.CommandRunChain("a")
	if err == nil {
		t.Fatal("expected cycle error")
	}
}
