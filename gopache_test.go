package gopache

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	file, _ := os.Open("./data/httpd.conf")

	defer file.Close()

	root, err := Parse(file)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, 39, len(root.Children))
}

func TestIsRootNode(t *testing.T) {
	file, _ := os.Open("./data/httpd.conf")

	defer file.Close()

	root, _ := Parse(file)

	assert.Equal(t, true, root.isRootNode())
}

func TestFindOne(t *testing.T) {
	file, _ := os.Open("./data/httpd.conf")
	defer file.Close()

	root, _ := Parse(file)

	node, err := root.FindOne("LoadModule")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "LoadModule", node.Name)

	assert.Equal(t, "alias_module modules/mod_alias.so", node.Content)
}

func TestFind(t *testing.T) {
	file, _ := os.Open("./data/httpd.conf")
	defer file.Close()

	root, _ := Parse(file)

	matches, err := root.Find("LoadModule")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 10, len(matches))

	for _, m := range matches {
		assert.Equal(t, "LoadModule", m.Name)
	}
}
