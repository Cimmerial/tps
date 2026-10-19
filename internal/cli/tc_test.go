package cli

import "testing"

func TestParseTkInput(t *testing.T) {
	in, err := parseTkInput([]string{"tk"})
	if err != nil || !in.TcKillAll || in.TcKillAllProjects {
		t.Fatalf("tk: %+v err=%v", in, err)
	}
	in, err = parseTkInput([]string{"tk", "-a"})
	if err != nil || !in.TcKillAll || !in.TcKillAllProjects {
		t.Fatalf("tk -a: %+v err=%v", in, err)
	}
	in, err = parseTkInput([]string{"tk", "-all"})
	if err != nil || !in.TcKillAllProjects {
		t.Fatalf("tk -all: %+v err=%v", in, err)
	}
}

func TestParseTerminalCheckKillAllProjects(t *testing.T) {
	in, err := parseTerminalCheck([]string{"t", "-k", "-a"})
	if err != nil || !in.TcKillAll || !in.TcKillAllProjects {
		t.Fatalf("t -k -a: %+v err=%v", in, err)
	}
}
