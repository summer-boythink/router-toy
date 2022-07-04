package router_toy

import (
	"net/http"
	"strings"
)

type PathNode struct {
	url      string
	part     string
	children []*PathNode
	isWild   bool
	Handle   http.HandlerFunc
}

func (n *PathNode) insert(url string, fn http.HandlerFunc) {
	pathArr := strings.Split(url, "/")
	if len(pathArr[0]) == 0 {
		pathArr = pathArr[1:]
	}
	head := n
	for k, path := range pathArr {
		for _, node := range n.children {
			if node.part == path {
				node.insert(url, fn)
			}
		}

		var newNode *PathNode
		if path[0] == '*' || path[0] == ':' {
			newNode = &PathNode{
				url:      url,
				part:     path,
				children: []*PathNode{},
				isWild:   true,
			}
		} else {
			newNode = &PathNode{
				url:      url,
				part:     path,
				children: []*PathNode{},
				isWild:   false,
			}
		}
		if k == len(pathArr)-1 || path[0] == '*' {
			newNode.Handle = fn
		}
		head.children = append(head.children, newNode)
		head = newNode
	}
}

func search(url string, root *PathNode) http.HandlerFunc {
	urls := strings.Split(url, "/")
	if fn := match(root, urls, 0); fn != nil {
		return fn
	}
	return nil
}

func match(p *PathNode, urls []string, height int) http.HandlerFunc {
	if urls[height] == p.part || p.isWild {
		if height == len(urls)-1 || (len(p.part) > 0 && p.part[0] == '*') {
			return p.Handle
		}
		height++
		for _, v := range p.children {
			if fn := match(v, urls, height); fn != nil {
				return fn
			}
		}
	}
	return nil
}

type Router struct {
	store map[string]*PathNode
}

func (router *Router) Get(url string, handle http.HandlerFunc) {
	router.addRoute("get", url, handle)
}

func (router *Router) Post(url string, handle http.HandlerFunc) {
	router.addRoute("post", url, handle)
}

func (router *Router) Put(url string, handle http.HandlerFunc) {
	router.addRoute("put", url, handle)
}

func (router *Router) Delete(url string, handle http.HandlerFunc) {
	router.addRoute("delete", url, handle)
}

func (router *Router) Option(url string, handle http.HandlerFunc) {
	router.addRoute("option", url, handle)
}

func (router *Router) addRoute(method string, url string, handle http.HandlerFunc) {
	root := router.store[method]
	if root != nil {
		root.insert(url, handle)
		return
	}
	root = &PathNode{
		url:      url,
		part:     "",
		children: []*PathNode{},
		isWild:   false,
	}
	root.insert(url, handle)
	router.store[method] = root
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	method := r.Method
	method = strings.ToLower(method)
	if root := router.store[method]; root != nil {
		if handle := search(url, root); handle != nil {
			handle(w, r)
		}
	}
}

func New() *Router {
	return &Router{
		store: map[string]*PathNode{},
	}
}
