package greet

import "fmt"

// Hello greets someone.
func Hello(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}
