package trie

import (
	"net/http"
	"reflect"
	"sort"

	"strings"
	"sync"
)

type router struct {
	tree       []*trie //排除options,//get,post放在最前
	paramsPool sync.Pool
	maxParams  uint16

	//前后调用
	middleware             HandlerFuncs
	SaveMatchedRoutePath   bool
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	HandleOPTIONS          bool
	GlobalOPTIONS          http.Handler
	globalAllowed          string
	NotFound               http.Handler
	MethodNotAllowed       http.Handler
	PanicHandler           func(http.ResponseWriter, *http.Request, interface{})
}

type trie struct {
	maxWith uint8 //节点采用切片的最大长度
}

type nodeType uint8

const (
	static   nodeType = 0b0000_0001 // default
	root     nodeType = 0b0000_0010
	param    nodeType = 0b0000_0100
	catchAll nodeType = 0b0000_1000
)

type node struct {
	path string

	indices   []byte //顺序索引
	schildren []*node

	nType    nodeType
	cType    nodeType //if>3 wildChild,代替原来的wildChild
	priority uint8

	//mchildren map[string]*node
	middleware []http.HandlerFunc
	handle     []*methodHandle
}

type methodHandle struct {
	method string
	handle reflect.Value
}

func (n *node) addRoute(path string, handle *methodHandle) {
	if path == "" || path[0] != '/' {
		path += "/"
	}

	fullPath := path
	n.priority++

	// Empty tree
	if len(n.path) == 0 && len(n.schildren) == 0 {
		n.insertChild(path, fullPath, handle)
		n.nType = root
		return
	}
walk:
	for {
		// Find the longest common prefix.
		// This also implies that the common prefix contains no ':' or '*'
		// since the existing key can't contain those chars.
		i := longestCommonPrefix(path, n.path)

		// Split edge
		if i < len(n.path) {
			child := &node{
				path:      n.path[i:],
				nType:     static,
				indices:   n.indices,
				schildren: n.schildren,
				handle:    n.handle,
				priority:  n.priority - 1,
			}
			n.path = path[:i]
			n.indices = []byte{child.path[0]}
			n.schildren = []*node{child}
			n.cType = static
			n.handle = nil
		}
		/// /:def
		// Make new node a child of this node
		if i < len(path) {
			path = path[i:]

			if n.cType >= param && (path[0] == ':' || path[0] == '*') {
				n = n.schildren[0]
				n.priority++
				// /:name/:names
				if (n.nType == catchAll && path != n.path) ||
					(n.nType == param && (len(path) < len(n.path) || n.path != path[:len(n.path)])) ||
					(n.nType == param && len(path) > len(n.path) && path[len(n.path)] != '/') {

					pathSeg := path
					if n.nType != catchAll {
						pathSeg = strings.SplitN(pathSeg, "/", 2)[0]
					}
					//这里有问题 /test/id/path/path/path
					prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
					panic("'" + pathSeg +
						"' in new path '" + fullPath +
						"' conflicts with existing wildcard '" + n.path +
						"' in existing prefix '" + prefix +
						"'")
				}
				continue walk
			}

			idx := insertGetIndex(n.indices, path[0])
			if idx > -1 { // / path
				n = n.schildren[idx]
				n.priority++

				continue walk
			}

			// []byte for proper unicode char conversion, see #65
			n.indices = append(n.indices, path[0])
			child := &node{}
			n.schildren = append(n.schildren, child)

			n.incrementChildPrio(len(n.indices) - 1)

			child.insertChild(path, fullPath, handle)
			n.cType = n.cType | child.nType
			n.sortIndices()
			return
		}

		// Otherwise add handle to current node
		if n.handle != nil {
			for _, h := range n.handle {
				if h.method == handle.method {
					panic("a handle is already registered for path '" + fullPath + "'")
				}
			}
		}
		if handle != nil {
			n.handle = append(n.handle, handle)
		}
		return
	}
}

func (n *node) insertChild(path, fullPath string, handle *methodHandle) {
	for {
		// Find prefix until first wildcard
		wildcard, i, valid := findWildcard(path)
		if i < 0 { // No wilcard found
			break
		}

		// The wildcard name must not contain ':' and '*'
		if !valid {
			panic("only one wildcard per path segment is allowed, has: '" +
				wildcard + "' in path '" + fullPath + "'")
		}

		// Check if the wildcard has a name
		if len(wildcard) < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		if i == 0 {
			n.path = path
			if wildcard[0] == '*' {
				n.nType = n.nType | catchAll
			} else {
				n.nType += n.nType | param
			}
			if len(wildcard) < len(path) {
				path = path[len(wildcard):]
				child := &node{
					priority: 1,
				}
				n.schildren = append(n.schildren, child)
				n.sortIndices()
				n = child
				continue
			}
			return
		}

		if wildcard[0] == '*' { // catchAll
			if i+len(wildcard) != len(path) {
				panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
			}

			if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
				panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
			}

			// Currently fixed width 1 for '/'
			if path[i-1] != '/' {
				panic("no / before catch-all in path '" + fullPath + "'")
			}
			n.path = path[:i]
			n.cType = n.cType | catchAll
			// First node: catchAll node with empty path
			child := &node{
				path:     wildcard,
				nType:    catchAll,
				priority: 1,
			}
			if handle != nil {
				child.handle = []*methodHandle{handle}
			}
			n.indices = append(n.indices, wildcard[0])
			n.schildren = append(n.schildren, child)
			n.sortIndices()
			// Second node: node holding the variable
		}

		if wildcard[0] == ':' { // param
			n.path = path[:i]
			path = path[i:]
			n.cType = n.cType | param

			child := &node{
				nType: param,
				path:  wildcard,
			}
			n.indices = append(n.indices, wildcard[0])
			n.schildren = append(n.schildren, child)
			n.sortIndices()
			n = child
			n.priority++
			// If the path doesn't end with the wildcard, then there
			// will be another non-wildcard subpath starting with '/'
			if len(wildcard) < len(path) {
				n.cType = n.cType | static
				path = path[len(wildcard):]
				child := &node{
					priority: 1,
				}
				n.indices = append(n.indices, '/')
				n.schildren = append(n.schildren, child)
				n.sortIndices()
				n = child
				continue
			}

			// Otherwise we're done. Insert the handle in the new leaf
			if handle != nil {
				n.handle = []*methodHandle{handle}
			}
		}
		return
	}
	// If no wildcard was found, simply insert the path and handle
	n.path = path
	n.nType = n.nType | static
	if handle != nil {
		n.handle = []*methodHandle{handle}
	}
}

func (n *node) incrementChildPrio(pos int) int {
	cs := n.schildren
	cs[pos].priority++
	prio := cs[pos].priority

	// Adjust position (move to front)
	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		// Swap node positions
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
		n.indices[newPos-1], n.indices[newPos] = n.indices[newPos], n.indices[newPos-1]
	}

	return newPos
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func longestCommonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}

// Shift bytes in array by n bytes left
func shiftNRuneBytes(rb [4]byte, n int) [4]byte {
	switch n {
	case 0:
		return rb
	case 1:
		return [4]byte{rb[1], rb[2], rb[3], 0}
	case 2:
		return [4]byte{rb[2], rb[3]}
	case 3:
		return [4]byte{rb[3]}
	default:
		return [4]byte{}
	}
}

func findWildcard(path string) (wilcard string, i int, valid bool) {
	// Find start
	for start, c := range []byte(path) {
		// A wildcard starts with ':' (param) or '*' (catch-all)
		if c != ':' && c != '*' {
			continue
		}

		valid = true
		// Find end and check for invalid characters
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '/':
				return path[start : start+1+end], start, valid
			case ':', '*':
				valid = false
			}
		}
		return path[start:], start, valid
	}
	return "", -1, false
}

func countParams(path string) uint16 {
	var n uint
	for i := range []byte(path) {
		switch path[i] {
		case ':', '*':
			n++
		}
	}
	return uint16(n)
}

//根据索引获取位置
func insertGetIndex(indices []byte, b byte) int {
	for i, c := range indices {
		if c == b {
			return i
		}
	}
	return -1
}

//排序
func (n *node) sortIndices() {
	sort.Slice(n.indices, func(i, j int) bool {
		return n.indices[i] < n.indices[j]
	})
	sort.Slice(n.schildren, func(i, j int) bool {
		return n.schildren[i].path[0] < n.schildren[j].path[0]
	})
}

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (n *node) getValue(path string, params func() *Params) (handle *methodHandle, ps *Params, tsr bool) {
walk: // Outer loop for walking the tree
	for {
		prefix := n.path
		if len(path) > len(prefix) {
			if path[:len(prefix)] == prefix {
				path = path[len(prefix):]

				// If this node does not have a wildcard (param or catchAll)
				// child, we can just look up the next child node and continue
				// to walk down the tree
				if n.cType < param {
					for i, c := range n.indices {
						if c == path[0] {
							n = n.schildren[i]
							continue walk
						}
					}

					// Nothing found.
					// We can recommend to redirect to the same URL without a
					// trailing slash if a leaf exists for that path.
					tsr = path == "/" && n.handle != nil
					return

				}

				// Handle wildcard child
				n = n.schildren[0]
				switch n.nType {
				case param:
					// Find param end (either '/' or path end)
					end := 0
					for end < len(path) && path[end] != '/' {
						end++
					}

					// Save param value
					if params != nil {
						if ps == nil {
							ps = params()
						}
						// Expand slice within preallocated capacity
						i := len(*ps)
						*ps = (*ps)[:i+1]
						(*ps)[i] = Param{
							Key:   n.path[1:],
							Value: path[:end],
						}
					}

					// We need to go deeper!
					if end < len(path) {
						if len(n.schildren) > 0 {
							path = path[end:]
							n = n.schildren[0]
							continue walk
						}

						// ... but we can't
						tsr = len(path) == end+1
						return
					}

					if handle = n.handle; handle.IsValid() {
						return
					} else if len(n.children) == 1 {
						// No handle found. Check if a handle for this path + a
						// trailing slash exists for TSR recommendation
						n = n.children[0]
						tsr = n.path == "/" && n.handle.IsValid()
					}

					return

				case catchAll:
					// Save param value
					if params != nil {
						if ps == nil {
							ps = params()
						}
						// Expand slice within preallocated capacity
						i := len(*ps)
						*ps = (*ps)[:i+1]
						(*ps)[i] = Param{
							Key:   n.path[2:],
							Value: path,
						}
					}

					handle = n.handle
					return

				default:
					panic("invalid node type")
				}
			}
		} else if path == prefix {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if handle = n.handle; handle != nil {
				return
			}

			// If there is no handle for this route, but this route has a
			// wildcard child, there must be a handle for this path with an
			// additional trailing slash
			if path == "/" && n.wildChild && n.nType != root {
				tsr = true
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i, c := range n.indices {
				if c == '/' {
					n = n.schildren[i]
					tsr = (len(n.path) == 1 && n.handle != nil) ||
						(n.nType == catchAll && n.schildren[0].handle != nil)
					return
				}
			}
			return
		}

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		tsr = (path == "/") ||
			(len(prefix) == len(path)+1 && prefix[len(path)] == '/' &&
				path == prefix[:len(prefix)-1] && n.handle != nil)
		return
	}
}
