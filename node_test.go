package gonode_test

import (
	"testing"

	"github.com/beanzilla/gonode"
	"golang.org/x/exp/slices"
)

func TestNewNode(t *testing.T) {
	n := gonode.NewNode()
	if !n.HasTag("root") {
		t.Fail()
		t.Logf("Expected 'root' as a tag, got '%s'", n.Tags())
	}
	if n.Len() != 0 {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected no children, got %d children", n.Len())
	}
	if n.Parent() != nil {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected no parent, got %v", n.Parent())
	}
	if n.Data() != nil {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Didn't assign any data, should be nil, got %#v", n.Data())
	}
}

func TestNewNodeWithData(t *testing.T) {
	n := gonode.NewNodeWithData("Hello World")
	if !n.HasTag("root") {
		t.Fail()
		t.Logf("Expected 'root' as tag, got '%s'", n.Tags())
	}
	if n.Data() != "Hello World" {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 'Hello World' as data, got %#v", n.Data())
	}
	n.SetData(42)
	if n.Data() != 42 {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 42 as data, got %#v", n.Data())
	}
	err := n.SetData(gonode.Node{})
	if err == nil {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected error, as setdata was given Node type")
		t.Logf("Data: %#v", n.Data())
	}
	n.AddTag("Meow", "Test")
	if !n.HasTag("root", "Meow", "Test") {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 'root', 'Meow', and 'Test' as tags, got '%s'", n.Tags())
	}
	if !slices.Equal(n.Tags(), []string{"root", "Meow", "Test"}) {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 'root', 'Meow', and 'Test' as tags, got '%s'", n.Tags())
	}
	n.RmTag("Meow", "root")
	if n.HasTag("root", "Meow") {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 'root' and 'Meow' removed, got '%s'", n.Tags())
	}
	n.RmAllTags()
	if n.HasTag("Test") {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected no tags, got '%s'", n.Tags())
	}
}

func TestNodeHierarchy(t *testing.T) {
	n := gonode.NewNode()
	a := n.NewChild()
	a.AddTag("kid 1")
	if !a.HasTag("kid 1") {
		t.Fail()
		t.Logf("Expected 'kid 1' tag on child")
	}
	if !a.Parent().HasTag("root") {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected child's parent to have 'root' tag")
	}
	if !n.Child(0).HasTag("kid 1") {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected parent's child to have 'kid 1' tag")
	}
	if !a.Detach() {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected successful Detach from parent")
	}
	if a.Parent() != nil {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Child should be independent from parent")
	}
	if n.Len() != 0 {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Parent should have no child, child detached")
	}
	n.AddChild(a)
	if n.Len() == 0 {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected parent excepting child back")
	}
	if a.Depth() != 0 {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected child's depth of 0, not %d", a.Depth())
	}
	n.RmTag("root")
	n.AddTag("main")
	if a.Depth() != 1 {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected child's depth of 1, not %d", a.Depth())
	}
}

func TestNodeChildEdge(t *testing.T) {
	n := gonode.NewNodeWithData(&gonode.Node{})
	if n != nil {
		t.Fail()
		t.Logf("Expected nil from constructor, given invalid data")
	}
	n = gonode.NewNode()
	a := n.NewChild()
	n.NewChild()
	a.NewChild()
	a.NewChild()
	if n.Child(2) != nil {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected nil, invalid child index")
	}
	if n.Depth() != 0 {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected \"root\" Node's depth of 0, not %d", n.Depth())
	}
	n.RmChild()
	if n.Len() != 2 {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 2 children, got %d", n.Len())
	}
	a.RmAllChildren()
	if a.Len() != 0 {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 0 children (after remove all), got %d", a.Len())
	}
	a.NewChild()
	a.NewChild()
	if n.Detach() {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected detach fail on \"root\" Node")
	}
	a.RmChild(1)
	if a.Len() != 1 {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 1 child (after removing the 2nd child), got %d", a.Len())
	}
}

func TestNodeChildByTag(t *testing.T) {
	n := gonode.NewNode()
	a := n.NewChild()
	a.AddTag("kid 1")
	b := n.NewChild()
	b.AddTag("kid 2")
	c := a.NewChild()
	c.AddTag("kid 1.1")
	d := b.NewChild()
	d.AddTag("kid 2.1")

	kid := n.ChildByTag("kid 2")
	if kid == nil || !kid.HasTag("kid 2") {
		t.Fail()
		t.Logf("Expected to find child(1), 'kid 2'")
	}
	kid = n.ChildByTag("kid 1.1")
	if kid != nil {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected nil, that child is nested, 'kid 1.1'")
	}
	kid = n.ChildByTagDeep("kid 1.1")
	if kid == nil || !kid.HasTag("kid 1.1") {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected to find child(0).child(0), 'kid 1.1'")
	}
	kid = n.ChildByTagDeep("kid 3")
	if kid != nil {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected nil, that child doesn't exist, 'kid 3'")
	}
	kid = n.ChildByTagDeep("kid 2")
	if kid == nil || !kid.HasTag("kid 2") {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected to find child(1), 'kid 2' again")
	}
}

func TestNodeChildByIndex(t *testing.T) {
	n := gonode.NewNode()
	a := n.NewChild()
	a.AddTag("kid 1")
	b := gonode.NewNode()
	b.AddTag("kid 2")
	c := a.NewChild()
	c.AddTag("kid 1.1")
	d := b.NewChild()
	d.AddTag("kid 2.1")

	if !n.Child(0).HasTag("kid 1") {
		t.Fail()
		t.Logf("Expected to find 'kid 1' as child 0, got '%s'", n.Child(0).Tags())
	}
	n.ReplaceChild(1, b)
	if !n.Child(0).HasTag("kid 1") {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected to find 'kid 1', as child 0, got '%s' (invalid index on replace)", n.Child(0).Tags())
	}
	n.ReplaceChild(0, b)
	if n.Child(0).HasTag("kid 1") {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected to find 'kid 2', as child 0, got '%s'", n.Child(0).Tags())
	}

	n.ReplaceChild(0, a)
	n.AddChild(b)
	id := n.ChildIndexByTag("kid 2")
	if id != 1 {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 'kid 2' at index 1")
	}
	id = n.ChildIndexByTag("kid 3")
	if id != -1 {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 'kid 3' at index -1, it doesn't exist")
	}
}

func TestNewNodes(t *testing.T) {
	n := gonode.NewNodeWithTags("Meow", "Test")
	if !n.HasTag("Meow", "Test", "root") {
		t.Fail()
		t.Logf("Expected 'Meow', 'Test' and 'root' as tags, got '%s'", n.Tags())
	}
	n = gonode.NewNodeWithDataAndTags("Meow? Glug?", "cat", "fish")
	if !n.HasTag("cat", "fish", "root") {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 'cat', 'fish' and 'root' as tags, got '%s'", n.Tags())
	}
	n = gonode.NewNodeWithDataAndTags(gonode.Node{}, "failwhale")
	if n != nil {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected nil due to invalid data assignment")
	}
}

func TestNodeNewChilds(t *testing.T) {
	n := gonode.NewNode()
	a := n.NewChildWithData(42)
	if a.Data() != 42 {
		t.Fail()
		t.Logf("Expected 42, got %#v", a.Data())
	}
	a = n.NewChildWithData(gonode.NewNode())
	if a != nil {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected nil, given invalid data type")
	}
	b := n.NewChildWithTags("kid", "2")
	if !b.HasTag("kid", "2") {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 'kid' and '2' as tags, got '%s'", b.Tags())
	}
	c := n.NewChildWithDataAndTags(9.81, "gravity")
	if c.Data() != 9.81 || !c.HasTag("gravity") {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected 9.81 as data and 'gravity' as tags, got %#v as data and '%s' tags", c.Data(), c.Tags())
	}
	c = n.NewChildWithDataAndTags(gonode.NewNode(), "failwhale")
	if c != nil {
		if !t.Failed() {
			t.Fail()
		}
		t.Logf("Expected nil, given invalid data type")
	}
}
