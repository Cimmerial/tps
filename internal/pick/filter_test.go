package pick

import "testing"

func TestFileFilterValidate(t *testing.T) {
	f := FileFilter{Extensions: []string{"pdf"}}
	if err := f.Validate("/x/doc.pdf"); err != nil {
		t.Fatal(err)
	}
	if err := f.Validate("/x/doc.png"); err == nil {
		t.Fatal("expected reject")
	}
}
