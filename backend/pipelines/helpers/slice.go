package helpers

func ChunkSlice(sliceLen int, chunkSize int) [][]int {
	min := func(a, b int) int {
		if a <= b {
			return a
		}
		return b
	}

	makeRange := func(min, max int) []int {
		a := make([]int, max-min+1)
		for i := range a {
			a[i] = min + i
		}
		return a
	}

	batches := make([][]int, 0)
	for i := 0; i < sliceLen; i += chunkSize {
		batches = append(batches, makeRange(i, min(i+chunkSize, sliceLen)-1))
	}
	return batches
}
