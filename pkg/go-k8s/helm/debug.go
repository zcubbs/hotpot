package helm

import (
	"fmt"
	"log"
)

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	err := log.Output(2, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Printf("error: %s", err)
	}
}

func getDebug(enable bool) func(format string, v ...interface{}) {
	if enable {
		return debug
	}
	return func(format string, v ...interface{}) {
		// do nothing
	}
}
