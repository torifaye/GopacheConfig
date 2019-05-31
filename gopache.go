package gopache

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
)

const (
	commentRegex      = "#.*"
	directiveRegex    = `([^\/\s]+)\s*(.+)`
	sectionOpenRegex  = `<([^\/\s>]+)\s*([^>]+)?>`
	sectionCloseRegex = `<\/([^\s>]+)\s*>`
)

var (
	commentMatcher      = regexp.MustCompile(commentRegex)
	directiveMatcher    = regexp.MustCompile(directiveRegex)
	sectionOpenMatcher  = regexp.MustCompile(sectionOpenRegex)
	sectionCloseMatcher = regexp.MustCompile(sectionCloseRegex)
)

// ConfigNode is a recursively defined n-ary tree
type ConfigNode struct {
	Name     string
	Content  string
	Parent   *ConfigNode
	Children []*ConfigNode
}

func newRootNode() *ConfigNode {
	return &ConfigNode{
		Name:     "",
		Content:  "",
		Parent:   nil,
		Children: nil,
	}
}

func (c *ConfigNode) addChild(child *ConfigNode) {
	c.Children = append(c.Children, child)
}

func createChildNode(name, content string, parent *ConfigNode) (*ConfigNode, error) {
	if name == "" {
		return nil, errors.New("Name cannot be empty")
	}
	if content == "" {
		return nil, errors.New("Content cannot be empty")
	}
	if parent == nil {
		return nil, errors.New("Parent cannot be null")
	}

	child := &ConfigNode{
		Name:    name,
		Content: content,
		Parent:  parent,
	}
	parent.addChild(child)

	return child, nil
}

func (c *ConfigNode) isRootNode() bool {
	return c.Parent == nil
}

// String prints out the contents of a Config node in an easy to read format
func (c *ConfigNode) String(level int) string {
	if len(c.Children) == 0 {
		return fmt.Sprintf("{name: %+v, content: %+v}\n", c.Name, c.Content)
	}
	children := ""
	if c.Name != "" {
		children = fmt.Sprintf("{name: %+v, content: %+v, childCount: %+v}\n", c.Name, c.Content, len(c.Children))
	}
	for _, node := range c.Children {
		tabs := ""
		for i := 0; i < level; i++ {
			tabs += "\t"
		}
		children += tabs + node.String(level+1)
	}
	return children
}

// Parse reads a data source and converts the apache config file into a tree-based struct
func Parse(r io.Reader) (*ConfigNode, error) {
	scanner := bufio.NewScanner(r)
	currentNode := newRootNode()
	for scanner.Scan() {
		if commentMatcher.MatchString(scanner.Text()) {
			continue
		} else if sectionOpenMatcher.MatchString(scanner.Text()) {
			groups := sectionOpenMatcher.FindStringSubmatch(scanner.Text())
			name := groups[1]
			content := groups[2]
			sectionNode, err := createChildNode(name, content, currentNode)
			if err != nil {
				return nil, err
			}
			currentNode = sectionNode
		} else if sectionCloseMatcher.MatchString(scanner.Text()) {
			currentNode = currentNode.Parent
		} else if directiveMatcher.MatchString(scanner.Text()) {
			groups := directiveMatcher.FindStringSubmatch(scanner.Text())
			name := groups[1]
			content := groups[2]
			createChildNode(name, content, currentNode)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return currentNode, nil
}
