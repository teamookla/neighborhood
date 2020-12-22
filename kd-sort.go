package neighborhood

import "math"

func kdSort(ids []int, coords []float64, nodeSize, left, right, axis int) {
	if right-left < nodeSize {
		return
	}
	m := (left + right) >> 1 // middle index

	// sort ids and coords around the middle index so that the halves lie
	// either left/right or top/bottom correspondingly (taking turns)
	selection(ids, coords, m, left, right, axis)

	// recursively kd-sort first half and second half on the opposite axis
	kdSort(ids, coords, nodeSize, left, m-1, 1-axis)
	kdSort(ids, coords, nodeSize, m+1, right, 1-axis)
}

// selection is a custom Floyd-Rivest selection algorithm: sort ids and coords so that
// [left..k-1] items are smaller than k-th item (on either x or y axis)
func selection(ids []int, coords []float64, k, left, right, axis int) {
	for right > left {
		if right-left > 600 {
			n := float64(right - left + 1)
			m := float64(k - left + 1)
			z := math.Log(n)
			s := 0.5 * math.Exp(2*z/3)
			sign := 1.0
			if m-n/2 < 0 {
				sign = -1.0
			}
			sd := 0.5 * math.Sqrt(z*s*(n-s)/n) * sign
			newLeft := int(math.Max(float64(left), math.Floor(float64(k)-m*s/n+sd)))
			newRight := int(math.Min(float64(right), math.Floor(float64(k)+(n-m)*s/n+sd)))
			selection(ids, coords, k, newLeft, newRight, axis)
		}
		t := coords[2*k+axis]
		i := left
		j := right

		swapItem(ids, coords, left, k)
		if coords[2*right+axis] > t {
			swapItem(ids, coords, left, right)
		}

		for i < j {
			swapItem(ids, coords, i, j)
			i++
			j--
			for coords[2*i+axis] < t {
				i++
			}
			for coords[2*j+axis] > t {
				j--
			}
		}

		if coords[2*left+axis] == t {
			swapItem(ids, coords, left, j)
		} else {
			j++
			swapItem(ids, coords, j, right)
		}

		if j <= k {
			left = j + 1
		}
		if k <= j {

			right = j - 1
		}
	}
}

func swapItem(ids []int, coords []float64, i, j int) {
	swapInt(ids, i, j)
	swapFloat(coords, 2*i, 2*j)
	swapFloat(coords, 2*i+1, 2*j+1)
}

func swapInt(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func swapFloat(arr []float64, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}
