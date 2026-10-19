package logos

import "github.com/Cimmerial/tps/internal/config"

// Preset describes a title-card logo choice in the options menu.
type Preset struct {
	ID    string
	Label string
	Art   string
}

const logoDefault = `      ___           ___           ___     
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

const logoWire = "        ,----,                        \n" +
	"      ,/   .`|,-.----.                \n" +
	"    ,`   .'  :\\    /  \\    .--.--.    \n" +
	"  ;    ;     /|   :    \\  /  /    '.  \n" +
	".'___,/    ,' |   |  .\\ :|  :  /`. /  \n" +
	"|    :     |  .   :  |: |;  |  |--`   \n" +
	";    |.';  ;  |   |   \\ :|  :  ;_     \n" +
	"`----'  |  |  |   : .   / \\  \\    `.  \n" +
	"    '   :  ;  ;   | |`-'   `----.   \\ \n" +
	"    |   |  '  |   | ;      __ \\  \\  | \n" +
	"    '   :  |  :   ' |     /  /`--'  / \n" +
	"    ;   |.'   :   : :    '--'.     /  \n" +
	"    '---'     |   | :      `--'---'   \n" +
	"              `---'.|                 \n" +
	"                `---`                  "

const logoOBird = `  ____o__ __o____   o__ __o       o__ __o    
  /   \   /   \   <|     v\     /v     v\   
       \o/        / \     <\   />       <\  
        |         \o/     o/  _\o____       
       < >         |__  _<|/       \_\_o__ 
        |          |                     \  
        o         <o>          \         /  
       <|          |            o       o   
       / \        / \           <\__ __/>   `

const logoSSS = `_____ ____  ____ 
/__ __Y  __\/ ___\
  / \ |  \/||    \
  | | |  __/\___ |
  \_/ \_/   \____/`

const logoDollar = `$$$$$$$$\ $$$$$$$\   $$$$$$\  
\__$$  __|$$  __$$\ $$  __$$\ 
   $$ |   $$ |  $$ |$$ /  \__|
   $$ |   $$$$$$$  |\$$$$$$\  
   $$ |   $$  ____/  \____$$\ 
   $$ |   $$ |      $$\   $$ |
   $$ |   $$ |      \$$$$$$  |
   \__|   \__|       \______/ `

const logoTPSBlocks = `  _______ _____   _____ 
 |__   __|  __ \ / ____|
    | |  | |__) | (___  
    | |  |  ___/ \___ \ 
    | |  | |     ____) |
    |_|  |_|    |_____/ `

const logoTPSGradient = `░▒▓████████▓▒░▒▓███████▓▒░ ░▒▓███████▓▒░ 
   ░▒▓█▓▒░   ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░        
   ░▒▓█▓▒░   ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░        
   ░▒▓█▓▒░   ░▒▓███████▓▒░ ░▒▓██████▓▒░  
   ░▒▓█▓▒░   ░▒▓█▓▒░             ░▒▓█▓▒░ 
   ░▒▓█▓▒░   ░▒▓█▓▒░             ░▒▓█▓▒░ 
   ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓███████▓▒░  `

const logoTPSSlash = ".------..------..------.\n" +
	"|T.--. ||P.--. ||S.--. |\n" +
	"| :/\\: || :/\\: || :/\\: |\n" +
	"| (__) || (__) || :\\/: |\n" +
	"| '--'T|| '--'P|| '--'S|\n" +
	"`------'`------'`------'"

const logoPixel = `▗▄▄▄▖▗▄▄▖  ▗▄▄▖
  █  ▐▌ ▐▌▐▌   
  █  ▐▛▀▘  ▝▀▚▖
  █  ▐▌   ▗▄▄▞▘`

const logoOO = "(`-')      _  (`-') (`-').-> \n" +
	"( OO).->   \\-.(OO ) ( OO)_   \n" +
	"/    '._   _.'    \\(_)--\\_)  \n" +
	"|'--...__)(_...--''/    _ /  \n" +
	"`--.  .--'|  |_.' |\\_..`--.  \n" +
	"   |  |   |  .___.'.-._)   \\ \n" +
	"   |  |   |  |     \\       / \n" +
	"   `--'   `--'      `-----'  "

const logoSlashes = `__/\\\\\\\\\\\\\\\\\\\\\\__/\\\\\\\\\\\\\\\\\\\\\\_______/\\\\\\\\\\\\\\\\\\\\\___        
 _\///////\\\/////__\/\\\/////////\\\___/\\\/////////\\\_       
  _______\/\\\_______\/\\\_______\/\\\__\//\\\______\///__      
   _______\/\\\_______\/\\\\\\\\\\\\\\\\\\/____\////\\\_________     
    _______\/\\\_______\/\\\/////////_________\////\\\______    
     _______\/\\\_______\/\\\_____________________\////\\\___   
      _______\/\\\_______\/\\\______________/\\\______\//\\\__  
       _______\/\\\_______\/\\\_____________\///\\\\\\\\\\\\\/___ 
        _______\///________\///________________\///////////_____`

const logoSmallAlt = `    ___       ___       ___   
   /\  \     /\  \     /\  \  
   \:\  \   /::\  \   /::\  \ 
   /::\__\ /::\:\__\ /\:\:\__\
  /:/\/__/ \/\::/  / \:\:\/__/
  \/__/       \/__/   \::/  / 
                      \/__/  `

var catalog = []Preset{
	{ID: "default", Label: "default (geometric TPS)", Art: logoDefault},
	{ID: "classic", Label: "classic (compact slashes)", Art: logoSmallAlt},
	{ID: "wire", Label: "wire block letters", Art: logoWire},
	{ID: "o_bird", Label: "o bird", Art: logoOBird},
	{ID: "sss", Label: "SSS", Art: logoSSS},
	{ID: "dollar", Label: "dollar", Art: logoDollar},
	{ID: "tps_blocks", Label: "TPS block letters", Art: logoTPSBlocks},
	{ID: "tps_gradient", Label: "TPS gradient blocks", Art: logoTPSGradient},
	{ID: "tps_slash", Label: "TPS slash cards", Art: logoTPSSlash},
	{ID: "pixel", Label: "pixel TPS", Art: logoPixel},
	{ID: "oo", Label: "OO ascii", Art: logoOO},
	{ID: "slashes", Label: "big slashes", Art: logoSlashes},
}

// Presets returns all selectable logo presets (stable order).
func Presets() []Preset {
	return catalog
}

// Art returns ASCII art for id, or default when unknown.
func Art(id string) string {
	id = config.NormalizeLogoID(id)
	for _, p := range catalog {
		if p.ID == id {
			return p.Art
		}
	}
	return logoDefault
}

// Label returns a human label for id.
func Label(id string) string {
	id = config.NormalizeLogoID(id)
	for _, p := range catalog {
		if p.ID == id {
			return p.Label
		}
	}
	if id == "" || id == "default" {
		return "default (geometric TPS)"
	}
	return id
}
