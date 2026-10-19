package ui

import "fmt"

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
	fmt.Println(ColorizeLines(SmallLogo, cyan))
}
