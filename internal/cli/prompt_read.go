package cli

import (
	"bufio"
	"fmt"

	"github.com/Cimmerial/tps/internal/ui"
)

func readLine(reader *bufio.Reader, running int, tick int) (string, error) {
	fmt.Print(ui.PromptBusy(running, tick))
	return reader.ReadString('\n')
}
