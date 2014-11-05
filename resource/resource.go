package resource

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

type Path []string

func (p *Path) String() string {
	return strings.Join(*p, "/")
}

func NewPath(nodes ...interface{}) *Path {
	path := Path{}
	for _, n := range nodes {
		switch n.(type) {
		case *uint64:
			path = append(path, strconv.FormatUint(*n.(*uint64), 10))
		case uint64:
			path = append(path, strconv.FormatUint(n.(uint64), 10))
		case string:
			path = append(path, n.(string))
		default:
			log.Fatal("Resource paths can only by Ids and string")
		}
	}
	return &path
}

type Edge interface {
	Get() (Node, error)
	Path() *Path
	Link()
}

type Node interface {
	Path() *Path
}

type Collection struct {
	path  *Path
	items []*Path
}

func (c *Collection) Path() *Path {
	return c.path
}

func (c *Collection) Append(paths ...*Path) {
	c.items = append(c.items, paths...)
}

func (c *Collection) Add(nodes ...Node) {
	paths := make([]*Path, len(nodes))
	for i, n := range nodes {
		paths[i] = n.Path()
	}
	c.Append(paths...)
}

func (c *Collection) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.items)
}

func NewCollection(path *Path, items ...*Path) *Collection {
	return &Collection{path, items}
}

type Resources map[*Path]Node

func (r Resources) Serializable() map[string]interface{} {
	serializable := map[string]interface{}{}
	for k, v := range r {
		serializable[k.String()] = v
	}
	return serializable
}

func (r Resources) Add(nodes ...Node) {
	for _, n := range nodes {
		switch node := n.(type) {
		case Edge:
			newNode, err := node.Get()
			if err == nil {
				r.Add(newNode)
				node.Link()
			}
		case Node:
			if node != nil && node.Path() != nil {
				r[node.Path()] = node
			}
		}
	}
}

func NewResources(nodes ...Node) Resources {
	r := Resources{}
	r.Add(nodes...)
	return r
}
