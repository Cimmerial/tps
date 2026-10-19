package pick

import "fmt"

// ErrCancelled means the user dismissed the picker or declined to type a path.
var ErrCancelled = fmt.Errorf("pick cancelled")
