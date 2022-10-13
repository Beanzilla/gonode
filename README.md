# gonode

Node-like/Tree-like data structures for Go

## Nodebase?

Nodes in gonode are all pointers `*Node`, thus you don't have to live in fear of making massive copies of many depths of Nodes.

Nodes in gonode are designed to cover most applications/uses with a Data() and SetData() using go 1.19's `any` feature.

Nodes can be tagged for easier view of your data and easier accessing/retrieving a particular Node.

Nodes have children Nodes, and a parent Node (except for the "root" Node, which has a parent Node of `nil`).

## Examples

So you are parsing a file and want to split up the data into more managable "chunks" or perhaps you just want to represent nested data with a little bit of ease:

```go
package main

import (
  "fmt"
  "github.com/beanzilla/gonode"
)

func main() {
  root := gonode.NewNode() // Initalizes a "root" Node with everything empty (Usually you'll use NewChild() for making a new node from here on)
  kid1 := root.NewChild()
  kid1.AddTag("kid", "1") // Adds the tags 'kid' and '1' to kid1 (root will only have 'root', and kid1 won't have 'root')
  kid1.SetData(42) // Sets a value (most values accepted without error... but)
  other := gonode.NewNodeWithTags("other") // Makes a new Node and adds the tag 'other' to it
  err := kid1.SetData(other) // Example of an invalid value for a Node's data (Node's don't allow Node and *Node as they are better as children than data)
  if err != nil {
    fmt.Println("kid1.SetData()", err)
  }
  // Maybe you want to see your data in json?
  payload, err := json.MarshalIndent(root, "", "    ")
  if err != nil {
    fmt.Println("json.MarshalIndent", err)
    return // This is a fatal thing, so let's stop here
  }
  err := os.WriteFile("_debug.json", payload, 0660)
  if err != nil {
    fmt.Println("os.WriteFile", err)
    // Normally I'd stop here but we are done
  }
}
```
