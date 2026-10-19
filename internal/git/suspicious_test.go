package git

import (
	"testing"
)

func TestPathLooksSuspicious(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{".env", true},
		{"config/.env", true},
		{".env.local", true},
		{"dir/.env.production", true},
		{"secrets/server.pem", true},
		{"certs/client.PFX", true},
		{"credentials", true},
		{"config/credentials.json", true},
		{"my-secret-key.txt", true},
		{"internal/private/cache.go", true},
		{"README.md", false},
		{"internal/config/public.go", false},
		{"main.go", false},
		{"vendor/foo/secret-sauce/recipe.md", true},
	}
	for _, tc := range cases {
		got := pathLooksSuspicious(tc.path)
		if got != tc.want {
			t.Errorf("pathLooksSuspicious(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestSuspiciousPathsDedupes(t *testing.T) {
	got := SuspiciousPaths([]string{".env", ".env", "main.go", "foo.pem"})
	want := []string{".env", "foo.pem"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
