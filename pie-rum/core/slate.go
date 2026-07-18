package pierum

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	rumpaint "pie-rum-sdk/paint"
	"strings"
	"time"
)

type ISlate struct {
	usage      map[string]int64
	lastUpdate map[string]time.Time // updated during toggle system
	metadata   *IMetadata
}

func NewSlate() *ISlate {
	return &ISlate{
		metadata:   NewMetadata(buffers),
		usage:      make(map[string]int64),
		lastUpdate: make(map[string]time.Time),
	}
}

// RecordChange updates the lastUpdate time for a given component
func (s *ISlate) RecordChange(name string) {
	s.lastUpdate[name] = time.Now()
}

// RecordUsage increments the usage count for a given component
func (s *ISlate) RecordUsage(name string) {
	s.usage[name]++
}

// type Slate[In any] struct {
// 	Profile []ISequence[In]

// 	DeActivate *bool
// 	Activate   *bool
// 	Remove     *bool
// }

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
