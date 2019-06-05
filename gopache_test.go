package gopache

import (
	"os"
	"strconv"
	"strings"
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

func TestFindOneVHost(t *testing.T) {
	file, _ := os.Open("./data/vhost.conf")
	defer file.Close()

	root, _ := Parse(file)

	match, err := root.FindOne("VirtualHost")
	if err != nil {
		panic(err)
	}
	tokens := strings.Split(match.Content, ":")
	port, _ := strconv.ParseUint(tokens[1], 10, 32)

	a := assert.New(t)
	a.Equal("172.20.30.40:4000", match.Content)
	a.Equal(uint(4000), uint(port))
	a.Equal("172.20.30.40", tokens[0])

	docRoot, err := root.FindOne("DocumentRoot")
	if err != nil {
		panic(err)
	}
	a.Equal("/www/subdomain/sub2", strings.Trim(docRoot.Content, "\""))
	admin, err := root.FindOne("ServerAdmin")
	if err != nil {
		panic(err)
	}
	a.Equal("name@email.corporation.com", admin.Content)
}
