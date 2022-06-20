package greet

import "testing"

func TestHello(t *testing.T) {
	hello := Hello("World")

	if hello != "Hello, World!" {
		t.Fatalf("expected to greet the World")
	}
}
