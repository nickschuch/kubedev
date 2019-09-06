package log

import (
	"fmt"

	"github.com/nickschuch/kubedev/internal/colorpicker"
)

// Infoln wraps fmt.Println with color.
func Infoln(name, message string) {
	fmt.Println(colorpicker.Wrap(name), message)
}

// Info wraps fmt.Printf with color.
func Info(name, message string) {
	fmt.Printf("%s %s", colorpicker.Wrap(name), message)
}

// Error wraps fmt.Println with color and supports errors.
func Error(name, message string, err error) {
	fmt.Println(colorpicker.Wrap(name), message, err)
}
