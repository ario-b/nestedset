// Copyright 2018 Ara Israelyan. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package nestedset

import (
	"fmt"
	"testing"
)

var (
	ns    *NestedSet
	nodes []NodeInterface
)

func createTestNestedSet(t *testing.T) {

	ns = NewNestedSet()

	nodes = make([]NodeInterface, 0)

	nodes = append(nodes, ns.rootNode)

	for i := 1; i <= 6; i++ {
		n := NewNode()
		n.NodeName = fmt.Sprintf("node %d", i)
		nodes = append(nodes, n)
	}

	err := ns.Add(nodes[1], nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ns.Add(nodes[2], nodes[1])
	if err != nil {
		t.Fatal(err)
	}
	err = ns.Add(nodes[3], nodes[0])
	if err != nil {
		t.Fatal(err)
	}
	err = ns.Add(nodes[4], nodes[3])
	if err != nil {
		t.Fatal(err)
	}
	err = ns.Add(nodes[5], nodes[1])
	if err != nil {
		t.Fatal(err)
	}
	err = ns.Add(nodes[6], nodes[4])
	if err != nil {
		t.Fatal(err)
	}

}

func checkNode(t *testing.T, node NodeInterface, level, left, right int64) {
	if node.Level() != level {
		t.Errorf("Invalid level for node '%s', expected %d, get %d", node.Name(), right, node.Level())
	}
	if node.Left() != left {
		t.Errorf("Invalid left for node '%s', expected %d, get %d", node.Name(), left, node.Left())
	}
	if node.Right() != right {
		t.Errorf("Invalid right for node '%s', expected %d, get %d", node.Name(), right, node.Right())
	}
}

func TestNestedSet_Add(t *testing.T) {

	createTestNestedSet(t)

	checkNode(t, nodes[0], 0, 0, 13)
	checkNode(t, nodes[1], 1, 1, 6)
	checkNode(t, nodes[2], 2, 2, 3)
	checkNode(t, nodes[3], 1, 7, 12)
	checkNode(t, nodes[4], 2, 8, 11)
	checkNode(t, nodes[5], 2, 4, 5)
	checkNode(t, nodes[6], 3, 9, 10)

}

func TestNestedSet_Delete(t *testing.T) {

	createTestNestedSet(t)

	err := ns.Delete(nodes[1])
	if err != nil {
		t.Fatal(err)
	}

	if ns.exists(nodes[1]) {
		t.Fatalf("Error deleting node '%s'", nodes[0].Name())
	}
	if ns.exists(nodes[2]) {
		t.Fatalf("Error deleting node '%s'", nodes[2].Name())
	}
	if ns.exists(nodes[5]) {
		t.Fatalf("Error deleting node '%s'", nodes[3].Name())
	}

	checkNode(t, nodes[0], 0, 0, 7)
	checkNode(t, nodes[3], 1, 1, 6)
	checkNode(t, nodes[4], 2, 2, 5)
	checkNode(t, nodes[6], 3, 3, 4)

}

func TestNestedSet_Move(t *testing.T) {

	createTestNestedSet(t)

	ns.Move(nodes[4], nodes[2])

	checkNode(t, nodes[0], 0, 0, 13)
	checkNode(t, nodes[1], 1, 1, 10)
	checkNode(t, nodes[2], 2, 2, 7)
	checkNode(t, nodes[3], 1, 11, 12)
	checkNode(t, nodes[4], 3, 3, 6)
	checkNode(t, nodes[5], 2, 8, 9)
	checkNode(t, nodes[6], 4, 4, 5)

}

func TestNestedSet_Branch(t *testing.T) {

	branch := ns.Branch(nodes[1])
	if branch == nil {
		t.Error("Returned nil branch for node 1")
		return
	}
	printBranch(branch)

}

func printBranch(branch []NodeInterface) {

	for _, n := range branch {
		var i int64
		for i = 0; i < n.Level(); i++ {
			fmt.Print("..")
		}
		fmt.Printf("%s lvl:%d, left:%d, right:%d\n", n.Name(), n.Level(), n.Left(), n.Right())
	}
}

func TestNewNestedSetWithNodes(t *testing.T) {

	rootNode := createNodeSample(0, 0, 5)
	node1 := createNodeSample(1, 1, 4)
	node2 := createNodeSample(2, 2, 3)

	nodes := make([]NodeInterface, 0)
	nodes = append(nodes, node1, node2)

	ns := NewNestedSetWithNodes(rootNode, nodes)

	node3 := NewNode()
	node4 := NewNode()
	ns.Add(node3, nil)
	ns.Add(node4, node3)

	branch := ns.Branch(nil)

	resultNodes := make([]NodeInterface, 0)
	for _, n := range branch {
		resultNodes = append(resultNodes, n)
	}

	checkNode(t, resultNodes[0], 0, 0, 9)
	checkNode(t, resultNodes[1], 1, 1, 4)
	checkNode(t, resultNodes[2], 2, 2, 3)
	checkNode(t, resultNodes[3], 1, 5, 8)
	checkNode(t, resultNodes[4], 2, 6, 7)

}

func createNodeSample(level, lft, rgt int64) *Node {
	node := NewNode()
	node.SetLeft(lft)
	node.SetRight(rgt)
	node.SetLevel(level)

	return node
}
