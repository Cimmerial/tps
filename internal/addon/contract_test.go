package addon

import "testing"

func TestCheckContractOK(t *testing.T) {
	m := &Manifest{ContractIteration: CurrentContractIteration}
	c := CheckContract(m)
	if !c.OK || len(c.Warnings) != 0 {
		t.Fatalf("expected ok: %+v", c)
	}
}

func TestCheckContractMissing(t *testing.T) {
	c := CheckContract(&Manifest{})
	if c.OK || len(c.Warnings) == 0 {
		t.Fatalf("expected warning: %+v", c)
	}
}

func TestCheckContractBehind(t *testing.T) {
	if CurrentContractIteration < 2 {
		t.Skip("no prior iteration to test")
	}
	m := &Manifest{ContractIteration: CurrentContractIteration - 1}
	c := CheckContract(m)
	if c.OK || len(c.Warnings) == 0 {
		t.Fatalf("expected warning: %+v", c)
	}
}

func TestCheckContractAhead(t *testing.T) {
	m := &Manifest{ContractIteration: CurrentContractIteration + 1}
	c := CheckContract(m)
	if c.OK || len(c.Warnings) == 0 {
		t.Fatalf("expected warning: %+v", c)
	}
}
