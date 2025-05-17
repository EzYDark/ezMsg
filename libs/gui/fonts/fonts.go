package fonts

import (
	"io"       // Added for reading response body
	"net/http" // Added for HTTP requests
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

// LoadFontURL fetches a font file from a URL.
func LoadFontURL(url string) []byte {
	// Why: Make an HTTP GET request to the specified URL.
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal().Msgf("Failed to get font from URL %s: %v", url, err)
	}
	// Why: Ensure the response body is closed after the function finishes.
	defer resp.Body.Close()

	// Why: Check if the HTTP request was successful.
	if resp.StatusCode != http.StatusOK {
		log.Fatal().Msgf("Failed to download font from URL %s: status code %d", url, resp.StatusCode)
	}

	// Why: Read all data from the response body.
	fontBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Msgf("Failed to read font data from URL %s: %v", url, err)
	}
	return fontBytes
}
