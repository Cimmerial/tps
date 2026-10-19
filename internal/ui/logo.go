package ui

import (
	"fmt"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui/logos"
)

const SmallLogo = `      ___           ___           ___     
     /\  \         /\  \         /\  \    
     \:\  \       /::\  \       /::\  \   
      \:\  \     /:/\:\  \     /:/\ \  \  
      /::\  \   /::\~\:\  \   _\:\~\ \  \ 
     /:/\:\__\ /:/\:\ \:\__\ /\ \:\ \ \__\
    /:/  \/__/ \/__\:\/:/  / \:\ \:\ \/__/
   /:/  /           \::/  /   \:\ \:\__\  
   \/__/             \/__/     \:\/:/  /  
                                \::/  /   
                                 \/__/    `

func PrintLogo() {
	fmt.Println()
	fmt.Println()
	fmt.Println(ColorizeLines(logoArt(), titleCardColorCode()))
}

// PrintLogoPreview renders one title-card preset (resolves random to a sample pick).
func PrintLogoPreview(logoID string) {
	id := config.NormalizeLogoID(logoID)
	if id == "random" {
		pick := logos.PickRandom()
		fmt.Println(Dim("example this refresh: " + pick))
		fmt.Println()
		fmt.Println()
		fmt.Println(ColorizeLines(logos.Art(pick), titleCardColorCode()))
		return
	}
	fmt.Println()
	fmt.Println(ColorizeLines(logos.Art(id), titleCardColorCode()))
}
