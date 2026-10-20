//go:build !darwin

package pick

func ChooseFile(prompt string) (string, error) {
	return readPathLine(prompt)
}

func ChooseFiles(prompt string) ([]string, error) {
	return readPathsLine(prompt)
}

func ChooseFolder(prompt string) (string, error) {
	return readPathLine(prompt)
}
