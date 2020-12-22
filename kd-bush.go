package neighborhood

import "math"

// KDBush implements the Index interface with a flat kd-tree index. This is the default Index implementation.
type KDBush struct {
	nodeSize int
	points   []Point
	ids      []int
	coords   []float64
}

// KDBushOptions defines configurable options for the KDBush index
type KDBushOptions struct {
	NodeSize int
}

// DefaultKDBushOptions gets the default KDBush options, which you can use directly or modify before creating an Index
func DefaultKDBushOptions() KDBushOptions {
	return KDBushOptions{
		NodeSize: 64,
	}
}

// NewKDBushIndex creates a new KDBush Index implementation with given KDBushOptions
func NewKDBushIndex(opts KDBushOptions, points ...Point) Index {
	// store indices to the input array and coordinates in separate typed arrays
	ids := make([]int, len(points))
	coords := make([]float64, 2*len(points))

	for i := 0; i < len(points); i++ {
		ids[i] = i
		coords[2*i] = points[i].Lon()
		coords[2*i+1] = points[i].Lat()
	}

	// kd-sort both arrays for efficient search (see comments in sort.go)
	kdSort(ids, coords, opts.NodeSize, 0, len(ids)-1, 0)

	return &KDBush{
		nodeSize: opts.NodeSize,
		points:   points,
		ids:      ids,
		coords:   coords,
	}
}

// Nearby finds the k nearest Points to the origin that meet the Accepter criteria.
// If there are multiple Points that are the same distance from the origin and the Points implement the Ranker
// interface, the higher ranking Points will be preferred. Nearby may return less than k results if it cannot
// find k Points in the Index that meet the Accepter criteria.
func (idx *KDBush) Nearby(origin Point, k int, accept Accepter) []Point {
	result := make([]Point, 0, k)

	// a distance-sorted rank queue that will contain both points and kd-tree nodes
	q := newPriorityQueue(k)

	// an object that represents the top kd-tree node (the whole Earth)
	node := &kdTreeNode{
		Left:   0,
		Right:  len(idx.ids) - 1,
		Axis:   0,
		MinLon: -180,
		MinLat: -90,
		MaxLon: 180,
		MaxLat: 90,
	}

	cosLat := math.Cos(origin.Lat() * rad)

	for node != nil {
		left := node.Left
		right := node.Right

		if right-left <= idx.nodeSize { // leaf node
			// add all points of the leaf node to the queue
			for i := left; i <= right; i++ {
				pt := idx.points[idx.ids[i]]
				if accept(pt) {
					dist := haverSinDist(origin.Lon(), origin.Lat(), idx.coords[2*i], idx.coords[2*i+1], cosLat)
					q.PushPoint(pt, dist)
				}
			}
		} else { // not a leaf node (has child nodes)
			m := (left + right) >> 1 // middle index
			midLng := idx.coords[2*m]
			midLat := idx.coords[2*m+1]

			// add middle point to the queue
			pt := idx.points[idx.ids[m]]
			if accept(pt) {
				dist := haverSinDist(origin.Lon(), origin.Lat(), midLng, midLat, cosLat)
				q.PushPoint(pt, dist)
			}

			nextAxis := (node.Axis + 1) % 2

			// first half of the node
			leftNode := &kdTreeNode{
				Left:   left,
				Right:  m - 1,
				Axis:   nextAxis,
				MinLon: node.MinLon,
				MinLat: node.MinLat,
			}
			if node.Axis == 0 {
				leftNode.MaxLon = midLng
				leftNode.MaxLat = node.MaxLat
			} else {
				leftNode.MaxLon = node.MaxLon
				leftNode.MaxLat = midLat
			}

			// second half of the node
			rightNode := &kdTreeNode{
				Left:   m + 1,
				Right:  right,
				Axis:   nextAxis,
				MaxLon: node.MaxLon,
				MaxLat: node.MaxLat,
			}
			if node.Axis == 0 {
				rightNode.MinLon = midLng
				rightNode.MinLat = node.MinLat
			} else {
				rightNode.MinLon = node.MinLon
				rightNode.MinLat = midLat
			}

			leftNode.Dist = boxDist(origin.Lon(), origin.Lat(), cosLat, leftNode)
			rightNode.Dist = boxDist(origin.Lon(), origin.Lat(), cosLat, rightNode)

			// add child nodes to the queue
			q.PushNode(leftNode)
			q.PushNode(rightNode)
		}

		// fetch closest points from the queue; they're guaranteed to be closer
		// than all remaining points (both individual and those in kd-tree nodes),
		// since each node's distance is a lower bound of distances to its children
		for q.Len() > 0 && q.Peek().point != nil {
			itm := q.PopItem()
			result = append(result, itm.point)
			if len(result) == k {
				return result
			}
		}

		// the next closest kd-tree node
		itm := q.PopItem()
		if itm != nil && itm.node != nil {
			node = itm.node
		} else {
			node = nil
		}
	}

	return result
}

// kdTreeNode defines a box of points in the kd-tree
type kdTreeNode struct {
	Left  int     // left index in the kd-tree array
	Right int     // right index
	Axis  int     // 0 for longitude axis and 1 for latitude axis
	Dist  float64 // will hold the lower bound of children's distances to the query point

	// bounding box of the node
	MinLon float64
	MinLat float64
	MaxLon float64
	MaxLat float64
}