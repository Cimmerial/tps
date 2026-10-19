package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/Cimmerial/tps/internal/jobs"
	"github.com/Cimmerial/tps/internal/ui"
)

func parseTerminalCheck(fields []string) (MenuInput, error) {
	start := 0
	if len(fields) >= 2 && strings.EqualFold(fields[0], "terminal") && strings.EqualFold(fields[1], "check") {
		start = 2
	} else if len(fields) >= 1 && strings.EqualFold(fields[0], "tc") {
		start = 1
	} else {
		return MenuInput{}, fmt.Errorf("tc: expected tc or terminal check")
	}
	in := MenuInput{Action: ActionTerminalCheck}
	for i := start; i < len(fields); i++ {
		switch fields[i] {
		case "-k", "--k", "-kill", "--kill":
			if i+1 >= len(fields) {
				in.TcKillAll = true
				continue
			}
			switch {
			case strings.EqualFold(fields[i+1], "all"):
				in.TcKillAll = true
				i++
			default:
				id, err := parseTcKillID(fields[i+1])
				if err != nil {
					return errInput(err)
				}
				in.TcKillID = id
				i++
			}
		default:
			return MenuInput{}, fmt.Errorf("tc: unknown flag %q", fields[i])
		}
	}
	return in, nil
}

func errInput(err error) (MenuInput, error) {
	return MenuInput{}, err
}

func parseTcKillID(s string) (int, error) {
	var id int
	_, err := fmt.Sscanf(s, "%d", &id)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("tc: invalid job number %q", s)
	}
	return id, nil
}

func printTerminalJobs(projectID string) error {
	list, err := jobs.Reap(projectID)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Println(ui.Dim("no terminals running (tc -k N to stop one)"))
		return nil
	}
	fmt.Println(ui.Bold("TERMINALS:"))
	for _, j := range list {
		age := time.Since(j.Started).Round(time.Second)
		fmt.Printf("  %d) %s  pid %d  %s  %s\n",
			j.ID,
			ui.Abbrev(j.Alias),
			j.PID,
			ui.Dim(age.String()),
			ui.Dim(truncateRunes(j.Shell, 40)),
		)
	}
	fmt.Println(ui.Dim("tc -k · tc -k 1"))
	return nil
}

func runTerminalCheck(projectID string, in MenuInput) error {
	if in.TcKillAll {
		n, err := jobs.KillAll(projectID)
		if err != nil {
			return err
		}
		fmt.Println(ui.Green(fmt.Sprintf("stopped %d terminal(s)", n)))
		return nil
	}
	if in.TcKillID > 0 {
		n, err := jobs.Kill(projectID, in.TcKillID)
		if err != nil {
			return err
		}
		fmt.Println(ui.Green(fmt.Sprintf("stopped %d terminal(s)", n)))
		return nil
	}
	return printTerminalJobs(projectID)
}
