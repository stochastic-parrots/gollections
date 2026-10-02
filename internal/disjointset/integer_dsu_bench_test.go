package disjointset

const (
	findComponentSize = 256
	findComponents    = 256
	findForestSize    = findComponentSize * findComponents
)

func buildBalancedFindForest(union func(int, int) bool) {
	for base := 0; base < findForestSize; base += findComponentSize {
		for width := 1; width < findComponentSize; width *= 2 {
			for idx := base; idx < base+findComponentSize; idx += 2 * width {
				union(idx, idx+width)
			}
		}
	}
}
