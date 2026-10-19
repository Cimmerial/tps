package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/jobs"
	"github.com/Cimmerial/tps/internal/ui"
)

func parseTerminalCheck(fields []string) (MenuInput, error) {
	start := 0
	if len(fields) >= 2 && strings.EqualFold(fields[0], "terminal") && strings.EqualFold(fields[1], "check") {
		start = 2
	} else if len(fields) >= 1 && (strings.EqualFold(fields[0], "t") || strings.EqualFold(fields[0], "tc")) {
		start = 1
	} else {
		return MenuInput{}, fmt.Errorf("t: expected t or terminal check")
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
		case "-r", "--r", "-restart", "--restart":
			in.TcRestart = true
		default:
			return MenuInput{}, fmt.Errorf("t: unknown flag %q", fields[i])
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
		return 0, fmt.Errorf("t: invalid job number %q", s)
	}
	return id, nil
}

func printTerminalJobs(projectID string) error {
	list, err := jobs.Reap(projectID)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Println(ui.Dim("no terminals running (tk · t -k 1)"))
		return nil
	}
	fmt.Println(ui.Header("TERMINALS:"))
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
	fmt.Println(ui.Dim("tk · t -k 1 · tr"))
	return nil
}

func runTerminalCheck(p *config.Project, in MenuInput) error {
	if in.TcRestart {
		var started []jobs.Job
		if err := runWithBraille("restarting terminals…", func() error {
			var rerr error
			started, rerr = jobs.RestartAll(p)
			return rerr
		}); err != nil {
			return err
		}
		for _, job := range started {
			watchCaptureIfNeeded(p, job)
		}
		fmt.Println(ui.Green(fmt.Sprintf("restarted %d terminal(s)", len(started))))
		return nil
	}
	if in.TcKillAll {
		var n int
		if err := runWithBraille("stopping terminals…", func() error {
			var kerr error
			n, kerr = jobs.KillAll(p.ID)
			return kerr
		}); err != nil {
			return err
		}
		fmt.Println(ui.Green(fmt.Sprintf("stopped %d terminal(s)", n)))
		return nil
	}
	if in.TcKillID > 0 {
		var n int
		if err := runWithBraille("stopping terminal…", func() error {
			var kerr error
			n, kerr = jobs.Kill(p.ID, in.TcKillID)
			return kerr
		}); err != nil {
			return err
		}
		fmt.Println(ui.Green(fmt.Sprintf("stopped %d terminal(s)", n)))
		return nil
	}
	return printTerminalJobs(p.ID)
}

func watchCaptureIfNeeded(p *config.Project, job jobs.Job) {
	cmd := p.FindCommand(job.Alias)
	if cmd == nil || cmd.Capture == "" {
		return
	}
	logPath, err := jobs.LogFile(p.ID, job.ID)
	if err != nil {
		return
	}
	pattern := config.ResolveCapturePattern(cmd.Capture)
	jobs.WatchCapture(logPath, pattern, 90*time.Second)
}

func runWithBraille(label string, fn func() error) error {
	var err error
	ui.BrailleWait(label, func() {
		err = fn()
	})
	return err
}
