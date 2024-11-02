package btree

import (
	"fmt"
	"strconv"
)

func (bt *BTree) DeleteItem(key string) (value string, err error) {

	keyVal, err := strconv.Atoi(key)
	if err != nil {
		return "", err
	}
	var found bool

	value, _, found = bt.Root.deleteItem(keyVal)
	if !found {
		return "", fmt.Errorf("key not found in btree")
	}

	if len(bt.Root.itemArray) == 0 {
		bt.Root = bt.Root.childNodes[0]
	}
	return value, nil
}

func (nd *Node) RetrievePredecessorNode(position int) (predecessorNode *Node) {

	predecessorNode = nd.childNodes[position]

	for !predecessorNode.isLeafNode() {

		predecessorNode = predecessorNode.childNodes[len(predecessorNode.childNodes)-1]
	}

	return predecessorNode
}

func (nd *Node) RetrieveSuccessorNode(position int) (successorNode *Node) {

	successorNode = nd.childNodes[position+1]

	for !successorNode.isLeafNode() {

		successorNode = successorNode.childNodes[0]
	}

	return successorNode
}

func (nd *Node) deleteItem(key int) (value string, adequate bool, found bool) {

	/*

		here position is used to indicate =>

		1) if the item was found : index of the item in the itemArray of the current node.
		2) if the item was not found : index of the child node to recursively search in.

	*/
	position, foundInArray := nd.SearchItemInItemArray(key)

	if nd.isLeafNode() {

		if foundInArray {

			/*
				if item to be deleted is found in a leaf node, follow these steps =>

				1) Delete item from leaf node.
				2) if leaf node has fewer than minimum number of items allowed, return an adequate = false value.

			*/

			// 1) Delete item from leaf node.

			var deletedItem *Item

			deletedItem, nd.itemArray = DeleteItemAtIndex(position, nd.itemArray)
			value = deletedItem.value

			// 2) if leaf node has fewer than minimum number of items allowed, return an adequate = false value.

			if len(nd.itemArray) < nd.tree.maxItemsPerNode/2 {
				return value, false, true

			} else {
				return value, true, true
			}

		} else {

			// if leaf node does not contain the item to be deleted, then this item was not found in the tree, so return found = false value.

			return "", true, false
		}

	} else {

		if foundInArray {

			/*
				if item to be deleted found in an internal node, follow these steps =>

				1) retrieve predecessor node and successor node,
				2) choose the node which has more items, to prevent unwanted underflow condition.
				3) retrieve the replacement item from the chosen node, and replace the item to be deleted with this replacement item.
				4) call this same delete function for the replacement key. as this replacement key is chosen from a leaf node,
				   any and all underflow conditions will be handled by the else condition below.

			*/

			// 1) retrieve predecessor node and successor node.

			predecessorNode := nd.RetrievePredecessorNode(position)
			successorNode := nd.RetrieveSuccessorNode(position)

			// 2) choose the node which has more items, to prevent unwanted underflow condition

			var lenderNode *Node
			if len(predecessorNode.itemArray) > len(successorNode.itemArray) {
				lenderNode = predecessorNode
				logger.Println("predecessor node chosen : ")
				lenderNode.DisplayItemsInNode()
			} else {
				lenderNode = successorNode
				logger.Println("successor node chosen : ")
				lenderNode.DisplayItemsInNode()
			}

			// 3) retrieve the replacement item from the chosen node, and replace the item to be deleted with this replacement item.

			var lentItem *Item

			if lenderNode == predecessorNode {

				// if we're borrowing from the predecessor node, then we must choose the largest value, which is smaller than the item to be deleted.
				lentItem = lenderNode.itemArray[len(lenderNode.itemArray)-1]
			} else {

				// if we're borrowing from the successor node, then we must choose the smallest value, which is larger than the item to be deleted.
				lentItem = lenderNode.itemArray[0]
			}

			logger.Printf("item chosen to replace item to be deleted => key : %s value %s", lentItem.key, lentItem.value)

			var deletedItem *Item

			deletedItem, nd.itemArray = DeleteItemAtIndex(position, nd.itemArray)

			logger.Println("item array after deleting item from node : ")
			nd.DisplayItemsInNode()

			value = deletedItem.value

			nd.itemArray = appendItemToIndex(lentItem, position, nd.itemArray)

			logger.Println("item array after replacing with lent item in node : ")
			nd.DisplayItemsInNode()

			keyVal, _ := strconv.Atoi(lentItem.key)
			logger.Printf("item to be deleted from lender node => key %d value %s", keyVal, lentItem.value)

			// 4) call this same delete function for the replacement key. as this replacement key is chosen from a leaf node,
			// 	  any and all underflow conditions will be handled by the else condition below.

			if lenderNode == predecessorNode {
				nd.childNodes[position].deleteItem(keyVal)
			} else {
				nd.childNodes[position+1].deleteItem(keyVal)
			}

		} else {

			/*

				if item to be deleted was not found in the internal node, follow these steps =>

				1) recursively search the child node.
				2) if this search is unsuccessful, return the result.
				3) if the search was successful, but deleting the item caused underflow, follow steps 4 to 6

				4) choose a child, from whom an item must be borrowed, to satisfy the underflow condition of the child node that was searched.
				5) if either of the child nodes can lend an item, without themselves experiencing underflow, then the parent/internal node
				   borrows an item from the child node, and lends the separator to the underflow child node.
				   any children belonging to the item lent to the parent node, is now shifted to the underflow child node.

				6) if neither of the child nodes can lend an item, we must merge the separator, the underflow child node, and an adjacent child node.


			*/

			// 1) recursively search the child node.
			value, adequate, found = nd.childNodes[position].deleteItem(key)

			// 2) if this search is unsuccessful, return the result.

			if !found {
				return value, adequate, found
			}

			// 3) if the search was successful, but deleting the item caused underflow, follow steps 4 to 6

			if !adequate {

				// 4) choose a child, from whom an item must be borrowed, to satisfy the underflow condition of the child node that was searched.

				logger.Printf("child at position %d is inadequate having %d items", position, len(nd.childNodes[position].itemArray))
				var previousChildNode *Node
				var nextChildNode *Node
				var separator *Item
				var separatorPosition int

				if position-1 >= 0 {
					previousChildNode = nd.childNodes[position-1]

				}
				if position+1 < len(nd.childNodes) {
					nextChildNode = nd.childNodes[position+1]
				}

				var lenderChildNode *Node

				if previousChildNode != nil {
					lenderChildNode = previousChildNode
					separator = nd.itemArray[position-1]
					separatorPosition = position - 1
				}
				if nextChildNode != nil {

					if lenderChildNode == nil || (len(lenderChildNode.itemArray) < len(nextChildNode.itemArray)) {
						lenderChildNode = nextChildNode
						separator = nd.itemArray[position]
						separatorPosition = position
					}
				}

				if lenderChildNode.willBeAdequate() {

					/*

						5) if either of the child nodes can lend an item, without themselves experiencing underflow,
							then the parent/internal node borrows an item from the child node, and lends the separator to the underflow child node.
							Any children belonging to the item lent to the parent node, is now shifted to the underflow child node.

					*/

					// borrow from lender.

					logger.Printf("lender child node is adequate, having %d items", len(lenderChildNode.itemArray))
					lenderChildNode.DisplayItemsInNode()

					var lentItem *Item
					var lentChildNode *Node
					if lenderChildNode == nextChildNode {
						lentItem, lenderChildNode.itemArray = DeleteItemAtIndex(0, lenderChildNode.itemArray)
						lentChildNode, lenderChildNode.childNodes = DeleteChildNodeAtIndex(0, lenderChildNode.childNodes)

					} else {
						lentItem, lenderChildNode.itemArray = DeleteItemAtIndex(len(lenderChildNode.itemArray)-1, lenderChildNode.itemArray)
						lentChildNode, lenderChildNode.childNodes = DeleteChildNodeAtIndex(len(lenderChildNode.itemArray)-1, lenderChildNode.childNodes)
					}
					logger.Printf("item lent by lender child node => key : %s value %s ", lentItem.key, lentItem.value)

					_, nd.itemArray = DeleteItemAtIndex(separatorPosition, nd.itemArray)

					logger.Println("parent node after deleting separator item : ")
					nd.DisplayItemsInNode()

					nd.itemArray = appendItemToIndex(lentItem, separatorPosition, nd.itemArray)

					logger.Println("parent node after replacing with lent item : ")
					nd.DisplayItemsInNode()

					if lenderChildNode == nextChildNode {
						nd.childNodes[position].itemArray = appendItemToIndex(separator, len(nd.childNodes[position].itemArray), nd.childNodes[position].itemArray)
						if lentChildNode != nil {
							nd.childNodes[position].childNodes = appendChildNodeToIndex(lentChildNode, len(nd.childNodes[position].childNodes), nd.childNodes[position].childNodes)
						}

					} else {
						nd.childNodes[position].itemArray = appendItemToIndex(separator, 0, nd.childNodes[position].itemArray)
						if lentChildNode != nil {
							nd.childNodes[position].childNodes = appendChildNodeToIndex(lentChildNode, 0, nd.childNodes[position].childNodes)
						}
					}

					logger.Println("child node after borrowing from parent node : ")
					nd.childNodes[position].DisplayItemsInNode()

				} else {

					// 6) if neither of the child nodes can lend an item, we must merge the separator, the underflow child node, and an adjacent child node.

					if lenderChildNode == previousChildNode {

						logger.Println("previous child node selected")

						newChildNode := MergeChildNodesAndSeparator(previousChildNode, separator, nd.childNodes[position])

						logger.Println("new child node after merging")
						newChildNode.DisplayItemsInNode()

						_, nd.childNodes = DeleteChildNodeAtIndex(position-1, nd.childNodes)

						logger.Println("child nodes in parent node after 1 deletion")
						nd.DisplayChildNodes()

						_, nd.childNodes = DeleteChildNodeAtIndex(position-1, nd.childNodes)

						logger.Println("child nodes in parent node after 2 deletion")
						nd.DisplayChildNodes()

						nd.childNodes = appendChildNodeToIndex(newChildNode, position-1, nd.childNodes)

					} else {
						logger.Println("next child node selected")

						newChildNode := MergeChildNodesAndSeparator(nd.childNodes[position], separator, nextChildNode)

						logger.Println("new child node after merging")
						newChildNode.DisplayItemsInNode()

						_, nd.childNodes = DeleteChildNodeAtIndex(position, nd.childNodes)

						logger.Println("child nodes in parent node after 1 deletion")
						nd.DisplayChildNodes()

						_, nd.childNodes = DeleteChildNodeAtIndex(position, nd.childNodes)

						logger.Println("child nodes in parent node after 2 deletion")
						nd.DisplayChildNodes()

						nd.childNodes = appendChildNodeToIndex(newChildNode, position, nd.childNodes)
					}
					// deleting separator item
					_, nd.itemArray = DeleteItemAtIndex(separatorPosition, nd.itemArray)

				}
			}
		}
	}

	if nd.isAdequate() {
		return value, true, true
	}
	return value, false, true

}

func MergeChildNodesAndSeparator(leftChildNode *Node, separator *Item, rightChildNode *Node) (newChildNode *Node) {

	newItemArray := make([]*Item, 0)
	newItemArray = append(newItemArray, leftChildNode.itemArray...)
	newItemArray = append(newItemArray, separator)
	newItemArray = append(newItemArray, rightChildNode.itemArray...)

	newChildNodes := make([]*Node, 0)
	newChildNodes = append(newChildNodes, leftChildNode.childNodes...)
	newChildNodes = append(newChildNodes, rightChildNode.childNodes...)

	newChildNode = &Node{tree: rightChildNode.tree, itemArray: newItemArray, childNodes: newChildNodes}

	return newChildNode

}

func MergeChildNodes(leftChildNode *Node, rightChildNode *Node) (newChildNode *Node) {

	newItemArray := make([]*Item, 0)
	newItemArray = append(newItemArray, leftChildNode.itemArray...)
	newItemArray = append(newItemArray, rightChildNode.itemArray...)

	newChildNodes := make([]*Node, 0)
	newChildNodes = append(newChildNodes, leftChildNode.childNodes...)
	newChildNodes = append(newChildNodes, rightChildNode.childNodes...)

	newChildNode = &Node{tree: rightChildNode.tree, itemArray: newItemArray, childNodes: newChildNodes}

	return newChildNode
}
func (nd *Node) SearchItemInItemArray(key int) (position int, found bool) {

	position = 0
	var item *Item
	var keyVal int
	var err error
	for _, item = range nd.itemArray {

		keyVal, err = strconv.Atoi(item.key)
		if err != nil {
			return -1, false
		}
		if key <= keyVal {
			break
		}
		position++
	}

	if key == keyVal {
		return position, true
	}
	return position, false
}

func DeleteChildNodeAtIndex(position int, a []*Node) (deletedChildNode *Node, b []*Node) {

	if len(a) == 0 {
		return nil, a
	}
	deletedChildNode = a[position]
	a = append(a[:position], a[position+1:]...)
	return deletedChildNode, a
}

func DeleteItemAtIndex(position int, a []*Item) (deletedItem *Item, b []*Item) {

	if len(a) == 0 {
		return nil, a
	}
	deletedItem = a[position]
	a = append(a[:position], a[position+1:]...)
	return deletedItem, a
}

func (nd *Node) willBeAdequate() (ok bool) {
	return len(nd.itemArray) > nd.tree.maxItemsPerNode/2
}

func (nd *Node) isAdequate() (ok bool) {
	return len(nd.itemArray) >= nd.tree.maxItemsPerNode/2
}
