package addonmenu

import (
	"fmt"

	"github.com/Cimmerial/tps/internal/ui"
)

func printAddonMenuPanel(title string, choices []MenuItem) {
	box := ui.NewMenuBox(ui.ChromeChoice)
	box.Title(ui.HeaderChoice(title))
	box.Plain(ui.SectionTrail{LastSection: true, LastRow: false}, ui.Gray("number to choose — Enter, 0, or q = back"))
	for i, item := range choices {
		label := item.Label
		if IsDimmed(item) {
			label = ui.Gray(label)
		}
		box.ChoiceRow(ui.SectionTrail{LastSection: true, LastRow: false}, fmt.Sprintf("%d", i+1), label, "", "")
	}
	box.ChoiceRow(ui.SectionTrail{LastSection: true, LastRow: true}, "0", "back", "q", "")
	box.Flush()
}
