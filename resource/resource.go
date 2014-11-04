package resource

import (
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
	Get() Node
	Path() *Path
	Link()
}

type Node interface {
	Path() *Path
}

type Resources map[*Path]Node

func (r Resources) Serializable() map[string]interface{} {
	serializable := map[string]interface{}{}
	for k, v := range r {
		serializable[k.String()] = v
	}
	return serializable
}
