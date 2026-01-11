package mdx

// RecordBlockRangeTreeNode represents a node in the record block range tree.
// This tree is used to efficiently find the record block corresponding to a given offset.
type RecordBlockRangeTreeNode struct {
	startRange int64
	endRange   int64
	data       *MdictRecordBlockInfoListItem
	left       *RecordBlockRangeTreeNode
	right      *RecordBlockRangeTreeNode
}

// BuildRangeTree constructs a range tree from a list of record block info items.
// This tree allows for efficient querying of record blocks based on an offset.
func BuildRangeTree(list []*MdictRecordBlockInfoListItem, root *RecordBlockRangeTreeNode) {
	if len(list) == 0 {
		return
	}

	if len(list) == 1 {
		root.data = list[0]
		root.startRange = list[0].deCompressAccumulatorOffset
		root.endRange = list[0].deCompressAccumulatorOffset + list[0].deCompressSize
		return
	}

	if len(list) == 2 {
		root.startRange = list[0].deCompressAccumulatorOffset
		root.endRange = list[1].deCompressAccumulatorOffset + list[1].deCompressSize
		root.left = new(RecordBlockRangeTreeNode)
		BuildRangeTree(list[:1], root.left)
		root.right = new(RecordBlockRangeTreeNode)
		BuildRangeTree(list[1:], root.right)
		return
	}

	root.startRange = list[0].deCompressAccumulatorOffset
	root.endRange = list[len(list)-1].deCompressAccumulatorOffset + list[len(list)-1].deCompressSize

	mid := (len(list) - 1) / 2
	if mid > 0 {
		root.left = new(RecordBlockRangeTreeNode)
		BuildRangeTree(list[0:mid], root.left)
	}

	if mid < len(list) {
		root.right = new(RecordBlockRangeTreeNode)
		BuildRangeTree(list[mid:], root.right)
	}
}

// QueryRangeData queries the range tree to find the record block info item
// that contains the given queryRange offset.
func QueryRangeData(root *RecordBlockRangeTreeNode, queryRange int64) *MdictRecordBlockInfoListItem {
	if root == nil {
		return nil
	}

	if root.startRange > queryRange || root.endRange < queryRange {
		return nil
	}

	if root.data != nil {
		return root.data
	}

	if root.left != nil && root.left.endRange > queryRange {
		return QueryRangeData(root.left, queryRange)
	}

	if root.right != nil && root.right.startRange <= queryRange {
		return QueryRangeData(root.right, queryRange)
	}
	return nil
}
