package addon

import (
	"fmt"

	"github.com/Cimmerial/tps/internal/ui"
)

// CurrentContractIteration must stay in sync with ADDONS.md "Current iteration".
const CurrentContractIteration = 1

// ContractCheck describes how an addon manifest relates to the tps contract.
type ContractCheck struct {
	ManifestIteration int
	CurrentIteration  int
	OK                bool
	Warnings          []string
}

func CheckContract(m *Manifest) ContractCheck {
	cur := CurrentContractIteration
	out := ContractCheck{
		CurrentIteration: cur,
	}
	if m == nil {
		out.Warnings = []string{"no addon manifest"}
		return out
	}
	out.ManifestIteration = m.ContractIteration
	if m.ContractIteration == 0 {
		out.Warnings = append(out.Warnings, fmt.Sprintf("tps-addon.yaml does not set contract_iteration (current contract is %d)", cur))
		return out
	}
	if m.ContractIteration < cur {
		out.Warnings = append(out.Warnings, fmt.Sprintf("addon contract_iteration %d is behind tps contract %d — update the addon repo", m.ContractIteration, cur))
		return out
	}
	if m.ContractIteration > cur {
		out.Warnings = append(out.Warnings, fmt.Sprintf("addon contract_iteration %d is ahead of tps contract %d — update tps", m.ContractIteration, cur))
		return out
	}
	out.OK = true
	return out
}

// PrintContractWarnings prints non-fatal contract mismatch lines to stdout.
func PrintContractWarnings(m *Manifest) {
	check := CheckContract(m)
	for _, w := range check.Warnings {
		fmt.Printf("%s %s\n", ui.Yellow("addon contract:"), w)
	}
}
