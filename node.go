package gonode

import (
	"encoding/json"
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
)

// A Node of data, of any kind
//
// Useful for nested data
type Node struct {
	tags     []string
	data     any
	parent   *Node
	children []*Node
}

// Secret util for making a nested map of the children
func (n *Node) tomap() map[string]any {
	pay := map[string]any{}
	if n.data != nil {
		pay["Data"] = n.data
	}
	if len(n.tags) != 0 {
		pay["Tags"] = n.tags
	}
	if n.Len() != 0 {
		kids := []map[string]any{}
		for _, k := range n.children {
			kids = append(kids, k.tomap())
		}
		pay["Children"] = kids
	}
	return pay
}

// Secret util for making a nested Node structure from
func (n *Node) tonode(lvl map[string]any) error {
	o := n.NewChild()
	if lvl["Data"] != nil {
		err := o.SetData(lvl["Data"])
		if err != nil {
			return err
		}
	}
	if lvl["Tags"] != nil {
		tgs := lvl["Tags"].([]any)
		for _, t := range tgs {
			o.AddTag(t.(string))
		}
	}
	if lvl["Children"] != nil {
		nxt := lvl["Children"].([]any)
		for _, l := range nxt {
			err := o.tonode(l.(map[string]any))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Custom Marshaler for json
func (n *Node) MarshalJSON() ([]byte, error) {
	pay := n.tomap()
	return json.Marshal(pay)
}

// Custom Unmarshaler for json
func (n *Node) UnmarshalJSON(data []byte) error {
	pay := map[string]any{}
	err := json.Unmarshal(data, &pay)
	if err != nil {
		return err
	}
	err = n.tonode(pay)
	if err != nil {
		return err
	}
	if n.Child(0).Parent() != nil {
		c := n.Child(0)
		n.data = c.data
		n.tags = c.tags
		n.children = c.children
		c.Destroy()
	}
	return nil
}

// Makes a new "root" Node
func NewNode() *Node {
	return &Node{
		tags: []string{"root"},
	}
}

// Makes a new "root" Node with given data
//
// Can return nil when data is Node or *Node (which are better as children rather than data)
func NewNodeWithData(data any) *Node {
	n := &Node{
		tags: []string{"root"},
	}
	err := n.SetData(data)
	if err != nil {
		return nil
	}
	return n
}

// Makes a new "root" Node with given tag(s)
func NewNodeWithTags(tags ...string) *Node {
	n := NewNode()
	n.AddTag(tags...)
	return n
}

// Makes a new "root" Node with given data and given tag(s)
//
// Can return nil when data is Node or *Node (which are better as children rather than data)
func NewNodeWithDataAndTags(data any, tags ...string) *Node {
	n := NewNodeWithData(data)
	if n == nil {
		return n
	}
	n.AddTag(tags...)
	return n
}

// Obtains the data for this Node
func (n *Node) Data() any {
	return n.data
}

// Assigns the given data for this Node
//
// data types not allowed: Node, *Node (These are better suited for a "root" Node with children)
func (n *Node) SetData(d any) error {
	// Don't allow Node and *Node
	if reflect.TypeOf(d) == reflect.TypeOf(&Node{}) || reflect.TypeOf(d) == reflect.TypeOf(Node{}) {
		return fmt.Errorf("data type of %s not allowed", reflect.TypeOf(d))
	}
	n.data = d
	return nil
}

// Obtain the parent of the Node
//
// Can return nil for the "root" Node
func (n *Node) Parent() *Node {
	return n.parent
}

// Obtain's the number of children below this Node
func (n *Node) Len() int {
	return len(n.children)
}

// Obtains the index of the current node from it's parent
//
// Will return -1 if the current node has no parent (thus error occurred)
func (n *Node) Index() int {
	if n.Parent() == nil {
		return -1
	}
	idx := 0
	for kid := range n.Parent().Iter() {
		if kid == n {
			break
		}
		idx += 1
	}
	return idx
}

// Iterator - Iterates over the Node's children
//
// http://www.golangpatterns.info/object-oriented/iterators#TOC-Implementation
//
// for kid := range Node
func (n *Node) Iter() <-chan *Node {
	ch := make(chan *Node, n.Len()/2)
	go func() {
		defer close(ch)
		for idx := range make([]byte, n.Len()) {
			ch <- n.children[idx]
		}
	}()
	return ch
}

// Adds a given Node (as pointer) below this Node
func (n *Node) AddChild(o *Node) {
	o.parent = n
	n.children = append(n.children, o)
}

// Creates a new Node below this Node
//
// Returns a pointer to the new Node created
func (n *Node) NewChild() *Node {
	o := &Node{
		parent: n,
	}
	n.AddChild(o)
	return o
}

// Creates a new Node below this Node with the given tag(s)
//
// Returns a pointer to the new Node created
func (n *Node) NewChildWithTags(tags ...string) *Node {
	o := &Node{
		parent: n,
		tags:   tags,
	}
	n.AddChild(o)
	return o
}

// Creates a new Node below this Node with the given data
//
// data types not allowed: Node, *Node (These are better suited for a "root" Node with children)
//
// Returns a pointer to the new Node created (can be nil for invalid data type)
func (n *Node) NewChildWithData(data any) *Node {
	o := &Node{
		parent: n,
	}
	n.AddChild(o)
	err := o.SetData(data)
	if err != nil {
		return nil
	}
	return o
}

// Creates a new Node below this Node with the given data and given tag(s)
//
// data types not allowed: Node, *Node (These are better suited for a "root" Node with children)
//
// Returns a pointer to the new Node created (can be nil for invalid data type)
func (n *Node) NewChildWithDataAndTags(data any, tags ...string) *Node {
	o := &Node{
		parent: n,
		tags:   tags,
	}
	n.AddChild(o)
	err := o.SetData(data)
	if err != nil {
		return nil
	}
	return o
}

// Creates a new child after a particular index (Use -1 to place at the beginning/top)
func (n *Node) IndexNewChild(idx int) *Node {
	if n.Len() == 0 || idx > n.Len() || idx < -1 {
		return nil
	}
	var o *Node = nil
	at := 0
	new_kids := make([]*Node, n.Len())
	for kid := range n.Iter() {
		if at == 0 && idx == -1 {
			o = &Node{
				parent: n,
			}
			new_kids = append(new_kids, o)
		}
		if at != idx {
			new_kids = append(new_kids, kid)
			at += 1
			continue
		}
		new_kids = append(new_kids, kid)
		o = &Node{
			parent: n,
		}
		new_kids = append(new_kids, o)
		at += 1
	}
	n.children = new_kids
	return o
}

// Creates a new child after a particular index (Use -1 to place at the beginning/top)
//
// This version includes given tags to assign to the new node
func (n *Node) IndexNewChildWithTags(idx int, tags ...string) *Node {
	if n.Len() == 0 || idx > n.Len() || idx < -1 {
		return nil
	}
	var o *Node = nil
	at := 0
	new_kids := make([]*Node, n.Len())
	for kid := range n.Iter() {
		if at == 0 && idx == -1 {
			o = &Node{
				parent: n,
				tags:   tags,
			}
			new_kids = append(new_kids, o)
		}
		if at != idx {
			new_kids = append(new_kids, kid)
			at += 1
			continue
		}
		new_kids = append(new_kids, kid)
		o = &Node{
			parent: n,
			tags:   tags,
		}
		new_kids = append(new_kids, o)
		at += 1
	}
	n.children = new_kids
	return o
}

// Creates a new child after a particular index (Use -1 to place at the beginning/top)
//
// This version allows setting the new node's data
func (n *Node) IndexNewChildWithData(idx int, data any) *Node {
	if n.Len() == 0 || idx > n.Len() || idx < -1 {
		return nil
	}
	var o *Node = nil
	at := 0
	new_kids := make([]*Node, n.Len())
	for kid := range n.Iter() {
		if at == 0 && idx == -1 {
			o = &Node{
				parent: n,
			}
			err := o.SetData(data)
			if err != nil {
				return nil
			}
			new_kids = append(new_kids, o)
		}
		if at != idx {
			new_kids = append(new_kids, kid)
			at += 1
			continue
		}
		new_kids = append(new_kids, kid)
		o = &Node{
			parent: n,
		}
		err := o.SetData(data)
		if err != nil {
			return nil
		}
		new_kids = append(new_kids, o)
		at += 1
	}
	n.children = new_kids
	return o
}

// Creates a new child after a particular index (Use -1 to place at the beginning/top)
//
// This version allows setting the new node's data, and given tags
func (n *Node) IndexNewChildWithDataAndTags(idx int, data any, tags ...string) *Node {
	if n.Len() == 0 || idx > n.Len() || idx < -1 {
		return nil
	}
	var o *Node = nil
	at := 0
	new_kids := make([]*Node, n.Len())
	for kid := range n.Iter() {
		if at == 0 && idx == -1 {
			o = &Node{
				parent: n,
				tags:   tags,
			}
			err := o.SetData(data)
			if err != nil {
				return nil
			}
			new_kids = append(new_kids, o)
		}
		if at != idx {
			new_kids = append(new_kids, kid)
			at += 1
			continue
		}
		new_kids = append(new_kids, kid)
		o = &Node{
			parent: n,
			tags:   tags,
		}
		err := o.SetData(data)
		if err != nil {
			return nil
		}
		new_kids = append(new_kids, o)
		at += 1
	}
	n.children = new_kids
	return o
}

// Sets all values to empty
func (n *Node) Destroy() {
	n.data = nil
	n.tags = []string{}
	n.RmAllChildren()
	n.parent = nil
}

// Returns how far deep from the parent/"root" Node this Node is
//
// Use "root" to tag a Node as the "root" Node (this reduces the returned depth value to it's correct value)
func (n *Node) Depth() int {
	if n.Parent() != nil {
		depth := 0
		at := n
		for at.Parent() != nil {
			at = at.Parent()
			depth += 1
		}
		if at.HasTag("root") { // Exclude "root" Nodes
			depth -= 1
		}
		return depth
	}
	return 0
}

// Obtains a Node below this Node
//
// Can return nil for invalid Node (which could be because invalid index)
func (n *Node) Child(index int) *Node {
	if index >= n.Len() || index < 0 {
		return nil
	}
	return n.children[index]
}

// Returns the first child which satisfies the given tag(s)
//
// Returns nil if no children match the given tag(s)
func (n *Node) ChildByTag(tags ...string) *Node {
	for _, kid := range n.children {
		if kid.HasTag(tags...) {
			return kid
		} // Diff of ChildByTag and ChildByTagDeep
		/* else if kid.Len() != 0 {
			r := kid.ChildByTag(tags...) // Recursive call!
			if r != nil {
				return r
			}
		}*/
	}
	return nil
}

// Returns the first child which satisfies the given tag(s)
//
// # This Deep version will call it recursively on children too
//
// Returns nil if no children match the given tag(s)
func (n *Node) ChildByTagDeep(tags ...string) *Node {
	for _, kid := range n.children {
		if kid.HasTag(tags...) {
			return kid
		} else if kid.Len() != 0 {
			r := kid.ChildByTagDeep(tags...)
			if r != nil {
				return r
			}
		}
	}
	return nil
}

// Returns the index of the first child which satisfies the given tag(s)
//
// Because this is index (number), there won't be a Deep recursive version.
//
// Returns -1 if no children match the given tag(s)
func (n *Node) ChildIndexByTag(tags ...string) int {
	for idx, kid := range n.children {
		if kid.HasTag(tags...) {
			return idx
		}
	}
	return -1
}

// Replaces the child at index with the given Node
func (n *Node) ReplaceChild(index int, o *Node) {
	if index >= n.Len() || index < 0 {
		return
	}
	// De-couple the old child
	n.children[index].parent = nil
	n.children[index] = o // Replace with new child
	o.parent = n          // Update new child (re-couple)
}

// Removes multiple (or single) children by index(s)
func (n *Node) RmChild(indexs ...int) {
	if len(indexs) == 0 {
		return
	}
	kids := []*Node{}
	for i, kid := range n.children {
		if !slices.Contains(indexs, i) {
			kids = append(kids, kid)
		} else {
			// De-couple the child from us
			kid.parent = nil
		}
	}
	n.children = kids
}

// Removes all children below this Node
func (n *Node) RmAllChildren() {
	for _, kid := range n.children {
		// De-couple the child from us
		kid.parent = nil
	}
	n.children = []*Node{}
}

// Detaches the current node from the parent
func (n *Node) Detach() bool {
	if n.Parent() != nil {
		idx := -1
		for i, kid := range n.Parent().children {
			if kid == n {
				idx = i
				break
			}
		}
		if idx != -1 {
			n.Parent().RmChild(idx)
			return true
		}
	}
	return false
}

// Checks if this Node has the given tag(s)
func (n *Node) HasTag(tags ...string) bool {
	for _, tag := range tags {
		if !slices.Contains(n.tags, tag) {
			return false
		}
	}
	return true
}

// Adds the given tag(s) to this Node
func (n *Node) AddTag(tags ...string) {
	need := []string{}
	for _, tag := range tags {
		if !slices.Contains(n.tags, tag) {
			need = append(need, tag)
		}
	}
	n.tags = append(n.tags, need...)
}

// Removes the given tag(s) from this Node
func (n *Node) RmTag(tags ...string) {
	new_tags := []string{}
	for _, tag := range n.tags {
		if !slices.Contains(tags, tag) {
			new_tags = append(new_tags, tag)
		}
	}
	n.tags = new_tags
}

// Removes all tags from this Node
func (n *Node) RmAllTags() {
	n.tags = []string{}
}

func (n *Node) Tags() []string {
	return n.tags
}
