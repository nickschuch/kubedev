package colorpicker

import (
	"strconv"
	"sync"

	"github.com/mgutz/ansi"
)

var (
	lock    sync.RWMutex
	matches map[string]string
)

func init() {
	matches = make(map[string]string)
}

// Wrap the text in repeatable colors based on the input.
func Wrap(text string) string {
	lock.RLock()
	defer lock.RUnlock()

	if val, ok := matches[text]; ok {
		return ansi.Color(text, val)
	}

	// Ensure that we don't go over 256.
	if len(matches) > 256 {
		return ansi.Color(text, "default")
	}

	// Save it for later.
	matches[text] = strconv.Itoa(len(matches) + 1)

	return ansi.Color(text, matches[text])
}
