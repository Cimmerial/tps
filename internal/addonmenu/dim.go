package addonmenu

import "strings"

func DimReason(item MenuItem) string {
	if item.DimWhen == nil {
		return ""
	}
	d := item.DimWhen
	if r := strings.TrimSpace(d.Reason); r != "" {
		return r
	}
	if b := strings.TrimSpace(d.RequireBinary); b != "" {
		if !BinaryOnPath(b) {
			return b + " not found on PATH"
		}
	}
	if f := strings.TrimSpace(d.MagickRead); f != "" {
		if !MagickAvailable() {
			return "ImageMagick (magick) not found on PATH"
		}
		if !MagickCanRead(f) {
			return "cannot read " + strings.ToUpper(f) + " with ImageMagick"
		}
	}
	if f := strings.TrimSpace(d.MagickWrite); f != "" {
		if !MagickAvailable() {
			return "ImageMagick (magick) not found on PATH"
		}
		if !MagickCanWrite(f) {
			return "cannot write " + strings.ToUpper(f) + " with ImageMagick"
		}
	}
	return ""
}

func IsDimmed(item MenuItem) bool {
	return DimReason(item) != ""
}
