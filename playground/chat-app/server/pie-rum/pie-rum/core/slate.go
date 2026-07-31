package pierum

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	rumpaint "pie-rum-sdk/paint"
	"strings"
	"time"
)

//go:fix inline
func boolPtr(t bool) *bool {
	b := t
	return &b
}

// openSSLHex returns the hex value alike openSSLHex rand -hex rang
func openSSLHex(rang int) (string, error) {
	bytes := make([]byte, rang)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func cleanJSONResponse(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)
	return text
}

func convertStringToTime(t string) time.Time {
	x, _ := time.Parse(t, "")
	return x
}

func convertStringToDuration(t string) *time.Duration {
	d, _ := time.ParseDuration(t)
	return &d
}

func printHeader() {
	t := rumpaint.Header(`
██████╗░██╗░░░██╗███╗░░░███╗
██╔══██╗██║░░░██║████╗░████║
██████╔╝██║░░░██║██╔████╔██║
██╔══██╗██║░░░██║██║╚██╔╝██║
██║░░██║╚██████╔╝██║░╚═╝░██║
╚═╝░░╚═╝░╚═════╝░╚═╝░░░░░╚═╝
	`)
	log.Println(t)
}
