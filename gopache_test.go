package gopache

import (
	"fmt"
	"log"
	"os"
	"testing"
)

// TODO: Add unit tests
func TestParse(t *testing.T) {
	file, err := os.Open("./data/httpd.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	root, err := Parse(file)

	if err != nil {
		log.Fatal(err)
	}
	if len(root.Children) != 39 {
		t.Errorf("len(root.Children) = %d; want 39", len(root.Children))
	}

	// Debug string
	fmt.Print(root.String(0))
}
