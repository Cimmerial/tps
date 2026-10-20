//go:build !darwin

package pick

func ChooseFile(prompt string, filter FileFilter) (string, error) {
	path, err := readPathLine(prompt)
	if err != nil {
		return "", err
	}
	if err := filter.Validate(path); err != nil {
		return "", err
	}
	return path, nil
}

func ChooseFiles(prompt string, filter FileFilter) ([]string, error) {
	paths, err := readPathsLine(prompt)
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		if err := filter.Validate(path); err != nil {
			return nil, err
		}
	}
	return paths, nil
}

func ChooseFolder(prompt string) (string, error) {
	return readPathLine(prompt)
}
