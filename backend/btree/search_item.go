package btree

import "strconv"

func (bt *BTree) SearchForItem(key string) (value string, found bool) {

	keyint, err := strconv.Atoi(key)
	if err != nil {
		return "", false
	}
	return bt.Root.searchForItem(keyint)
}

func (nd *Node) searchForItem(key int) (value string, found bool) {

	position := 0
	var item *Item

	var keyInt int
	var err error

	for _, item = range nd.itemArray {
		logger.Println("item : " + item.key)
		keyInt, err = strconv.Atoi(item.key)
		if err != nil {
			continue
		}
		if key == keyInt || key < keyInt {
			break
		}
		position++
	}

	if key == keyInt {
		return item.value, true
	}

	if nd.isLeafNode() {
		return "", false
	}

	childNode := nd.childNodes[position]

	return childNode.searchForItem(key)

}
