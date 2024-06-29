package common

import (
	"fmt"
)

// FormatPort returns a formatted port for httpServer consumption
func FormatPort(port string) string {
	return fmt.Sprintf(":%v", port)
}
