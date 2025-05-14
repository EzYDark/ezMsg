package fonts

import (
	"os"

	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/text"
	"github.com/rs/zerolog/log"
)

// appendFont loads a font from a byte slice and adds it to the collection.
func AppendFont(collection []text.FontFace, fontProps font.Font, fontBytes []byte) []text.FontFace {
	face, err := opentype.Parse(fontBytes) // Use opentype.ParseCollection for .otc files
	if err != nil {
		log.Fatal().Msgf("failed to parse font %s: %v", fontProps.Typeface, err)
	}
	return append(collection, text.FontFace{Font: fontProps, Face: face})
}

// loadFontFile reads a font file from disk.
func LoadFontFile(path string) []byte {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Msgf("Failed to read font file %s: %v", path, err)
	}
	return fontBytes
}
