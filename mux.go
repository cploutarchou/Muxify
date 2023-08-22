package muxify

import (
	"net/http"
	"strings"
	"sync"
)

type Mux struct {
	root            *pathNode
	notFoundHandler http.Handler
	middlewares     []func(http.Handler) http.Handler
	pool            *sync.Pool
}

type pathNode struct {
	pathSegment string
	children    []*pathNode
	route       *RouteInfo
	isWildcard  bool
}

type RouteInfo struct {
	Method      string
	Handler     http.Handler
	Middlewares []func(http.Handler) http.Handler
}

func NewMux() *Mux {
	return &Mux{
		root: &pathNode{},
		pool: &sync.Pool{
			New: func() interface{} {
				return &Context{}
			},
		},
		notFoundHandler: http.NotFoundHandler(),
	}
}

func (m *Mux) Handle(method, path string, handler http.Handler) {
	segments := strings.Split(path, "/")
	currentNode := m.root
	for _, segment := range segments {
		if segment != "" {
			if segment[0] == '*' {
				// It's a wildcard route, so we don't look any further.
				child := currentNode.addChild(segment)
				child.isWildcard = true
				currentNode = child
				break
			} else {
				currentNode = currentNode.addChild(segment)
			}
		}
	}
	currentNode.route = &RouteInfo{
		Method:  method,
		Handler: handler,
	}
}
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	currentNode := m.root
	for _, segment := range segments {
		if segment != "" {
			nextNode := currentNode.findChild(segment)
			if nextNode == nil {
				if currentNode.isWildcard {
					// We are at a wildcard, so we'll stop the traversal and handle it
					break
				}
				// No matching child node, and we're not at a wildcard, so it's a not found
				m.notFoundHandler.ServeHTTP(w, r)
				return
			}
			currentNode = nextNode
		}
	}

	if currentNode.route != nil && currentNode.route.Method == r.Method {
		finalHandler := currentNode.route.Handler
		for _, middleware := range m.middlewares {
			finalHandler = middleware(finalHandler)
		}

		ctx := m.pool.Get().(*Context)
		ctx.Response = &w
		ctx.Request = r
		defer m.pool.Put(ctx)

		finalHandler.ServeHTTP(w, r)
	} else {
		m.notFoundHandler.ServeHTTP(w, r)
	}
}

func (n *pathNode) addChild(segment string) *pathNode {
	for _, child := range n.children {
		if child.pathSegment == segment {
			return child
		}
	}
	newChild := &pathNode{
		pathSegment: segment,
	}
	n.children = append(n.children, newChild)
	return newChild
}

func (n *pathNode) findChild(segment string) *pathNode {
	for _, child := range n.children {
		if child.pathSegment == segment {
			return child
		}
	}
	return nil
}
