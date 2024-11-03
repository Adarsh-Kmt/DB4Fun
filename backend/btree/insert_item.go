package btree

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	logger = log.New(os.Stdout, "DB4Fun >> ", 0)
)

type BTree struct {
	Root            *Node
	maxItemsPerNode int
	//minItemsPerNode int
}

type Node struct {
	tree       *BTree
	itemArray  []*Item
	childNodes []*Node
}

type Item struct {
	key   string
	value string
}

func BTreeTraversal(node *Node) {
	//logger.Printf("number of items in child : %d", len(node.itemArray))
	fmt.Print("[")

	for index, item := range node.itemArray {
		fmt.Printf("{ index %d --> key %s : value %s }, ", index, item.key, item.value)
	}
	fmt.Println("]")
	for index, childNode := range node.childNodes {
		fmt.Printf("child number %d => ", index)
		//fmt.Println()
		BTreeTraversal(childNode)
	}

}
func BTreeInit(maxItems int) *BTree {

	bTree := &BTree{maxItemsPerNode: maxItems}

	rootNode := &Node{tree: bTree, itemArray: make([]*Item, 0), childNodes: make([]*Node, 0)}

	bTree.Root = rootNode

	return bTree
}

func (bt *BTree) BtreeSetup() {

	bt.InsertItem("0", "a")
	bt.InsertItem("10", "b")
	bt.InsertItem("20", "c")
	bt.InsertItem("30", "d")
	bt.InsertItem("40", "e")
	bt.InsertItem("50", "f")
	bt.InsertItem("60", "g")
	bt.InsertItem("70", "h")
	bt.InsertItem("80", "i")
	bt.InsertItem("90", "j")
	bt.InsertItem("100", "k")
	bt.InsertItem("120", "l")
	bt.InsertItem("130", "m")
	bt.InsertItem("140", "n")
	bt.InsertItem("150", "o")
	bt.InsertItem("160", "p")
	bt.InsertItem("170", "q")
	bt.InsertItem("180", "r")
	bt.InsertItem("190", "r")
	bt.InsertItem("200", "r")
	bt.InsertItem("210", "r")

}
func (bt *BTree) InsertItem(key string, value string) *BTree {

	// logger.Printf("-------------------------------------------------------")
	// fmt.Println()
	// logger.Printf("before inserting element [ %s - %s ]", key, value)
	// BTreeTraversal(bt.Root)
	// fmt.Println()
	logger.Printf("-------------------------------------------------------")
	fmt.Println()
	logger.Printf("inserting element [ %s - %s ]", key, value)

	newItem := &Item{key: key, value: value}

	extraItem, leftChildNode, rightChildNode := bt.Root.InsertItem(newItem)

	if extraItem != nil {
		newRoot := &Node{tree: bt, itemArray: []*Item{extraItem}, childNodes: []*Node{leftChildNode, rightChildNode}}
		bt.Root = newRoot
	}
	fmt.Println()
	logger.Printf("-------------------------------------------------------")
	fmt.Println()
	logger.Printf("after inserting element [ %s - %s ], bTree Traversal : ", key, value)
	BTreeTraversal(bt.Root)
	fmt.Println()
	logger.Printf("-------------------------------------------------------")

	return bt
}
func (nd *Node) InsertItem(newItem *Item) (extraItem *Item, childNode1 *Node, childNode2 *Node) {

	position := nd.findPosition(newItem.key)
	//logger.Printf("number of items in node : %d max items allowed %d", len(nd.itemArray), nd.tree.maxItemsPerNode)
	if nd.isLeafNode() {
		// logger.Println()
		// logger.Printf("------------------- BTREE TRAVERSAL ALERT ----------------------")
		// BTreeTraversal(nd.tree.Root)
		// logger.Printf("--------------------------- END --------------------------------")
		// logger.Println()
		nd.itemArray = appendItemToIndex(newItem, position, nd.itemArray)
		// logger.Println()
		// logger.Printf("------------------- BTREE TRAVERSAL ALERT ----------------------")
		// BTreeTraversal(nd.tree.Root)
		// logger.Printf("--------------------------- END --------------------------------")
		// logger.Println()

	} else {
		// logger.Println()
		// logger.Printf("------------------- BTREE TRAVERSAL ALERT ----------------------")
		// BTreeTraversal(nd.tree.Root)
		// logger.Printf("--------------------------- END --------------------------------")
		// logger.Println()

		nextChildNode := nd.childNodes[position]
		//logger.Printf("next child node position : %d", position)
		excessItem, ch1, ch2 := nextChildNode.InsertItem(newItem)

		// logger.Println()
		// logger.Printf("------------------- BTREE TRAVERSAL ALERT ----------------------")
		// BTreeTraversal(nd.tree.Root)
		// logger.Printf("--------------------------- END --------------------------------")
		// logger.Println()

		if excessItem != nil {

			// logger.Printf("deleting child node at position : %d", position)
			// fmt.Printf("\nbefore deleting child node at index %d, child nodes are: ", position)
			// for ind, childNode := range nd.childNodes {
			// 	fmt.Printf("\nchild node %d ", ind)
			// 	for ind2, item := range childNode.itemArray {
			// 		fmt.Printf("{ index %d --> key %s : value %s }, ", ind2, item.key, item.value)
			// 	}
			// }
			// fmt.Println()
			nd.childNodes = append(nd.childNodes[:position], nd.childNodes[position+1:]...)
			// fmt.Printf("\nafter deleting child node at index %d, child nodes are: ", position)
			// for ind, childNode := range nd.childNodes {
			// 	fmt.Printf("\nchild node %d ", ind)
			// 	for ind2, item := range childNode.itemArray {
			// 		fmt.Printf("{ index %d --> key %s : value %s }, ", ind2, item.key, item.value)
			// 	}
			// }
			// fmt.Println()

			excessItemPosition := nd.findPosition(excessItem.key)
			//logger.Printf("adding excess item at index : %d", excessItemPosition)
			nd.itemArray = appendItemToIndex(excessItem, excessItemPosition, nd.itemArray)
			//logger.Printf("adding left child at index : %d", excessItemPosition)
			nd.childNodes = appendChildNodeToIndex(ch1, excessItemPosition, nd.childNodes)
			//logger.Printf("adding right child at index : %d", excessItemPosition+1)
			nd.childNodes = appendChildNodeToIndex(ch2, excessItemPosition+1, nd.childNodes)
			// fmt.Print("\nafter adding excess item, child nodes are: ")
			// for ind, childNode := range nd.childNodes {
			// 	fmt.Printf("\nchild node %d ", ind)
			// 	for ind2, item := range childNode.itemArray {
			// 		fmt.Printf("{ index %d --> key %s : value %s }, ", ind2, item.key, item.value)
			// 	}
			// }
			// fmt.Println()
		} else {
			return nil, nil, nil
		}

	}

	if len(nd.itemArray) > nd.tree.maxItemsPerNode {
		return nd.Split()
	}

	return nil, nil, nil
}

func (nd *Node) Split() (extraItem *Item, leftChildNode *Node, rightChildNode *Node) {

	splitIndex := len(nd.itemArray) / 2
	extraItem = nd.itemArray[splitIndex]
	// logger.Printf("split index %d, split node key %s : value %s ", splitIndex, extraItem.key, extraItem.value)

	// logger.Printf("before splitting, child nodes are: ")
	// for ind, childNode := range nd.childNodes {
	// 	fmt.Printf("\nchild node %d ", ind)
	// 	for ind2, item := range childNode.itemArray {
	// 		fmt.Printf("{ index %d --> key %s : value %s }, ", ind2, item.key, item.value)
	// 	}
	// }
	// fmt.Println()
	//logger.Printf("length of items array : %d", len(nd.itemArray))

	leftChildNodeItemArray := []*Item{}
	leftChildNodeItemArray = append(leftChildNodeItemArray, nd.itemArray[:splitIndex]...)
	rightChildNodeItemArray := []*Item{}
	rightChildNodeItemArray = append(rightChildNodeItemArray, nd.itemArray[splitIndex+1:]...)

	leftChildNode = &Node{tree: nd.tree, itemArray: leftChildNodeItemArray, childNodes: make([]*Node, 0)}
	rightChildNode = &Node{tree: nd.tree, itemArray: rightChildNodeItemArray, childNodes: make([]*Node, 0)}
	// fmt.Println("after splitting, child nodes formed : ")
	// fmt.Println("items in left child node are :")
	// for index2, item := range leftChildNode.itemArray {
	// 	fmt.Printf("{ index %d --> key %s : value %s }, ", index2, item.key, item.value)
	// }
	// fmt.Println("\nitems in right child node are :")
	// for index2, item := range rightChildNode.itemArray {
	// 	fmt.Printf("{ index %d --> key %s : value %s }, ", index2, item.key, item.value)
	// }
	// fmt.Println()
	if !nd.isLeafNode() {
		leftChildNode.childNodes = nd.childNodes[:splitIndex+1]
		rightChildNode.childNodes = nd.childNodes[splitIndex+1:]
		// fmt.Println("after splitting, child nodes formed : ")
		// fmt.Println("left child node :")
		// for index1, childNode := range leftChildNode.childNodes {
		// 	fmt.Printf("\nchild node %d : ", index1)
		// 	for index2, item := range childNode.itemArray {
		// 		fmt.Printf("{ index %d --> key %s : value %s }, ", index2, item.key, item.value)
		// 	}

		// }
		// fmt.Println("\nright child node :")
		// for index1, childNode := range rightChildNode.childNodes {
		// 	fmt.Printf("\nchild node %d : ", index1)
		// 	for index2, item := range childNode.itemArray {
		// 		fmt.Printf("{ index %d --> key %s : value %s }, ", index2, item.key, item.value)
		// 	}

		// }
		// fmt.Println()
	}

	return extraItem, leftChildNode, rightChildNode

}
func (nd *Node) findPosition(key string) (position int) {

	keyVal, err := strconv.Atoi(key)

	if err != nil {
		panic("error")
	}

	for index, item := range nd.itemArray {

		itemKeyVal, err := strconv.Atoi(item.key)
		if err != nil {
			panic("error")
		}
		if keyVal < itemKeyVal {
			return index
		}
	}
	return len(nd.itemArray)
}

func (nd *Node) isLeafNode() bool {
	return len(nd.childNodes) == 0
}

func appendChildNodeToIndex(childNode *Node, index int, a []*Node) []*Node {

	if len(a) == 0 || index == len(a) {
		a = append(a, childNode)
		return a
	}
	// logger.Printf("child node append function called, before appending at index %d :  ", index)
	// for ind1, ch := range a {
	// 	fmt.Printf("\nchild node %d : ", ind1)
	// 	for ind, it := range ch.itemArray {
	// 		fmt.Printf("[ index : %d key : %s, value : %s ], ", ind, it.key, it.value)
	// 	}
	// }
	// fmt.Println()
	a = append(a[:index+1], a[index:]...)
	a[index] = childNode
	// logger.Printf("child node append function called, after appending %d :  ", index)
	// for ind1, ch := range a {
	// 	fmt.Printf("\nchild node %d : ", ind1)
	// 	for ind, it := range ch.itemArray {
	// 		fmt.Printf("[ index : %d key : %s, value : %s ], ", ind, it.key, it.value)
	// 	}
	// }
	return a
}
func appendItemToIndex(item *Item, index int, a []*Item) []*Item {

	// logger.Println()
	// logger.Printf("------------------- BTREE TRAVERSAL ALERT ----------------------")
	// BTreeTraversal(root)
	// logger.Printf("--------------------------- END --------------------------------")
	// logger.Println()
	if len(a) == 0 || index == len(a) {
		a = append(a, item)
		return a
	}
	// logger.Printf("item append function called, before appending at index %d :  ", index)
	// for ind, it := range a {
	// 	fmt.Printf("[ index : %d key : %s, value : %s ], ", ind, it.key, it.value)
	// }
	a = append(a[:index+1], a[index:]...)
	//fmt.Println()
	// logger.Println()
	// logger.Printf("------------------- BTREE TRAVERSAL ALERT ----------------------")
	// BTreeTraversal(root)
	// logger.Printf("--------------------------- END --------------------------------")
	// logger.Println()

	a[index] = item

	//fmt.Println()
	// fmt.Println()
	// logger.Println()
	// logger.Printf("------------------- BTREE TRAVERSAL ALERT ----------------------")
	// BTreeTraversal(root)
	// logger.Printf("--------------------------- END --------------------------------")
	// logger.Println()
	// logger.Printf("item append function called, after appending at index %d:  ", index)
	// for ind, it := range a {
	// 	fmt.Printf("[ index : %d key : %s, value : %s ], ", ind, it.key, it.value)
	// }
	// fmt.Println()
	// logger.Println()
	// logger.Printf("------------------- BTREE TRAVERSAL ALERT ----------------------")
	// BTreeTraversal(root)
	// logger.Printf("--------------------------- END --------------------------------")
	// logger.Println()
	return a
}

func (nd *Node) DisplayChildNodes() {

	for index1, childNode := range nd.childNodes {
		fmt.Printf("\nchild node %d : ", index1)
		for index2, item := range childNode.itemArray {
			fmt.Printf("{ index %d --> key %s : value %s }, ", index2, item.key, item.value)
		}
		fmt.Printf("\n")
	}

}

func (nd *Node) DisplayItemsInNode() {

	for ind, it := range nd.itemArray {
		fmt.Printf("[ index : %d key : %s, value : %s ], ", ind, it.key, it.value)
	}
	fmt.Printf("\n")
}
