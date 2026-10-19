package logos

// WalkArtRunes calls fn for every rune in embedded title-card art (all presets).
func WalkArtRunes(fn func(rune, string)) {
	if fn == nil {
		return
	}
	for _, p := range readmeCatalog() {
		src := "logo/" + p.ID
		for _, r := range p.Art {
			fn(r, src)
		}
	}
}
