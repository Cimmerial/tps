package cli

import (
	"fmt"

	"github.com/Cimmerial/tps/internal/ui"
)

func printConfigExamples() {
	fmt.Printf("%s\n", ui.Header("config examples"))
	fmt.Println()
	sections := []struct {
		title string
		body  string
	}{
		{
			title: "link a project (config.yaml)",
			body: `projects:
  - id: <uuid>
    name: Clastone
    path: /Users/you/dev/Clastone
    git_remote: git@github.com:you/Clastone.git
    commands:`,
		},
		{
			title: "basic command",
			body: `      - alias: run
        abbreviation: r
        description: dev server
        shell: npm run dev
        detach: true`,
		},
		{
			title: "tunnel + capture",
			body: `      - alias: tunnel
        abbreviation: cf
        shell: cloudflared tunnel --url http://localhost:5173
        capture: trycloudflare
        capture_actions: [print, clipboard, sms]
        detach: true`,
		},
		{
			title: "tags (menu folders)",
			body: `      - alias: merge to main
        abbreviation: mtm
        tags: [git]
        shell: git checkout main && git merge dev && git push`,
		},
		{
			title: "precursor (run another command first)",
			body: `      - alias: inspect classes
        abbreviation: ic
        precursor: fas
        shell: npm run inspect:library-classes`,
		},
		{
			title: "shell inputs ([INPUT1], [INPUT2], …)",
			body: `      - alias: inspect classes
        abbreviation: ic
        shell: npm run inspect:library-classes -- --username [INPUT1]
        # interactive: ic -Cimmerial   (first -value → INPUT1, second → INPUT2, …)`,
		},
		{
			title: "multi shell (parallel terminals)",
			body: `      - alias: stack
        shell1: npm run dev
        shell2: npx cloudflared tunnel --url http://localhost:5173
        detach: true`,
		},
		{
			title: "platform-specific shells",
			body: `      - alias: build
        linux_shell: ./scripts/build.sh
        darwin_shell: ./scripts/build-mac.sh
        windows_shell: powershell -File scripts/build.ps1`,
		},
		{
			title: "interactive: nc",
			body: `  nc -n "my cmd" -a x -c "npm test" -desc "run tests"
  nc -n tunnel -a cf -c "cloudflared …" -t -m trycloudflare -actions print,clipboard,sms
  nc -n admin -a adm -c "./script.sh" -tags admin -precursor fas
  nc -n build -a b -mac ./scripts/build-mac.sh -linux ./scripts/build.sh -win "powershell -File scripts/build.ps1"
  nc -n stack -a st -c "npm run dev" -c "npx cloudflared tunnel --url http://localhost:5173" -t
  nc -n once -a once -c ./scripts/seed.sh -t -close
  nc -n repl -a repl -c python -H
  nc -help`,
		},
		{
			title: "detach + close tab (macOS)",
			body: `      - alias: once
        shell: ./scripts/seed.sh
        detach: true
        close_on_complete: true`,
		},
		{
			title: "TPS.md + TPS_SPECS.md (linked app repos)",
			body: `  # Auto-written at repo root on setup / project add / u (gitignored)
  # Cheat sheet: internal/instance/README.md
  # Command/script spec: internal/instance/SPECS.md`,
		},
		{
			title: "setup / export / import",
			body: `  setup              # link cwd, refresh TPS.md + TPS_SPECS.md
  setup -n Clastone --path ~/dev/Clastone
  export             # menu: options, a project, or everything (save dialog)
  export -a          # everything in one file (options + all commands)
  import             # pick file, then destination
  import backup.tps.yaml`,
		},
	}
	for i, sec := range sections {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s\n", ui.SubHeader(sec.title))
		fmt.Println(ui.Dim(sec.body))
	}
	fmt.Println()
}
