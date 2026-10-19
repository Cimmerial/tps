package config

import "testing"

func TestResolveCapturePattern(t *testing.T) {
	if ResolveCapturePattern("url") != captureURLPattern {
		t.Fatal("url preset")
	}
	if ResolveCapturePattern("cf") != ResolveCapturePattern("trycloudflare") {
		t.Fatal("cf alias")
	}
}

func TestExtractCapture(t *testing.T) {
	log := "INF tunnel https://foo.trycloudflare.com ready\n"
	got, err := ExtractCapture(log, ResolveCapturePattern("trycloudflare"))
	if err != nil || got != "https://foo.trycloudflare.com" {
		t.Fatalf("trycloudflare: %q err=%v", got, err)
	}
	log2 := "see https://127.0.0.1:3000) for local"
	got2, err := ExtractCapture(log2, ResolveCapturePattern("url"))
	if err != nil || got2 != "https://127.0.0.1:3000" {
		t.Fatalf("url trim: %q err=%v", got2, err)
	}
}

func TestParseCaptureActions(t *testing.T) {
	acts := ParseCaptureActions([]interface{}{"print", "sms", "bogus", "clipboard"})
	if len(acts) != 3 || acts[0] != "print" {
		t.Fatalf("%v", acts)
	}
	c := Command{Capture: "url", CaptureActions: []string{"sms"}}
	if !c.HasCaptureAction("sms") || c.HasCaptureAction("clipboard") {
		t.Fatal("per-command actions")
	}
	empty := Command{Capture: "url"}
	def := empty.EffectiveCaptureActions()
	if len(def) != 2 || def[1] != "clipboard" {
		t.Fatalf("default %v", def)
	}
}
