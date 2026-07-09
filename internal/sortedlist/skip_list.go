package sortedlist

import (
	"encoding/json"
	"fmt"
	"iter"
	"slices"

	"github.com/stochastic-parrots/gollections/internal/list"
	"github.com/stochastic-parrots/gollections/internal/shared/collection"
)

const (
	skipListNil              = 0
	skipListHead             = 1
	skipListMaxLevel         = 32
	skipListInitialBlockSize = 64
)

// SkipList is a probabilistic indexed-arena sorted list with custom comparator order.
//
// Elements are kept in comparator order after every mutation. Each arena link
// stores the next node index and span, allowing index-based operations to skip
// over lower levels instead of walking the bottom list one node at a time.
// Equivalent values are allowed, but their relative order is not part of the
// list contract.
type SkipList[T any] struct {
	nodes         []skipListNode[T]
	links         []skipListLink
	tail          int
	free          int
	freeLength    int
	level         int
	length        int
	seed          uint64
	nextBlockSize int
	compare       func(a, b T) int
}

type skipListNode[T any] struct {
	value        T
	linkStart    int
	linkCapacity int
	backward     int
	freeNext     int
	height       int
}

type skipListLink struct {
	next int
	span int
}

// NewSkipList creates an empty SkipList.
func NewSkipList[T any](compare func(a, b T) int) *SkipList[T] {
	list := &SkipList[T]{
		nodes:         make([]skipListNode[T], skipListHead+1),
		tail:          skipListNil,
		free:          skipListNil,
		level:         1,
		seed:          0x9e3779b97f4a7c15,
		nextBlockSize: skipListInitialBlockSize,
		compare:       compare,
	}
	list.prepareNode(skipListHead, skipListMaxLevel)
	return list
}

// NewSkipListFromSlice creates a SkipList from the provided slice.
//
// The input slice is sorted in place before values are inserted into the list.
func NewSkipListFromSlice[T any](data []T, compare func(a, b T) int) *SkipList[T] {
	slices.SortFunc(data, compare)
	list := NewSkipList(compare)
	list.addSorted(data)
	return list
}

// NewSkipListCloneSlice creates a SkipList from a sorted clone of the provided slice.
func NewSkipListCloneSlice[T any](data []T, compare func(a, b T) int) *SkipList[T] {
	return NewSkipListFromSlice(slices.Clone(data), compare)
}

// NewSkipListFromSeq creates a SkipList from an iterator.
func NewSkipListFromSeq[T any](seq iter.Seq[T], compare func(a, b T) int) *SkipList[T] {
	return NewSkipListFromSlice(slices.Collect(seq), compare)
}

func (list *SkipList[T]) allocateLinks(level int) int {
	start := len(list.links)
	list.links = append(list.links, make([]skipListLink, level)...)
	return start
}

func (list *SkipList[T]) prepareNode(node, level int) {
	entry := &list.nodes[node]
	if entry.linkCapacity < level {
		entry.linkStart = list.allocateLinks(level)
		entry.linkCapacity = level
	} else {
		clear(list.links[entry.linkStart : entry.linkStart+entry.linkCapacity])
	}
	entry.height = level
}

func (list *SkipList[T]) link(node, level int) *skipListLink {
	return &list.links[list.nodes[node].linkStart+level]
}

func (list *SkipList[T]) next(node, level int) int {
	return list.links[list.nodes[node].linkStart+level].next
}

func (list *SkipList[T]) newValueNode(level int, value T) int {
	if list.free == skipListNil {
		list.reserveNodes(list.nextBlockSize)
		if list.nextBlockSize < 1<<20 {
			list.nextBlockSize *= 2
		}
	}

	node := list.free
	list.free = list.nodes[node].freeNext
	list.freeLength--
	list.nodes[node].freeNext = skipListNil
	list.prepareNode(node, level)
	list.nodes[node].value = value
	return node
}

func (list *SkipList[T]) ensureFreeNodes(count int) {
	if count <= list.freeLength {
		return
	}

	list.reserveNodes(count - list.freeLength)
}

func (list *SkipList[T]) reserveNodes(count int) {
	if count <= 0 {
		return
	}

	start := len(list.nodes)
	list.nodes = append(list.nodes, make([]skipListNode[T], count)...)
	for idx := count; idx > 0; idx-- {
		node := start + idx - 1
		list.nodes[node].freeNext = list.free
		list.free = node
	}
	list.freeLength += count
}

func skipListPerfectLevel(rank int) int {
	level := 1
	for level < skipListMaxLevel && rank%4 == 0 {
		level++
		rank /= 4
	}
	return level
}

// Length returns the current number of elements in the list.
//
// Complexity: O(1).
func (list *SkipList[T]) Length() int {
	return list.length
}

// IsEmpty returns true if the list contains no elements.
//
// Complexity: O(1).
func (list *SkipList[T]) IsEmpty() bool {
	return list.length == 0
}

// Get retrieves the value at the specified sorted index.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) Get(idx int) (T, error) {
	if idx < 0 || idx >= list.length {
		var zero T
		return zero, listpkgIndexError(idx, list.length)
	}

	rank := idx + 1
	traversed := 0
	current := skipListHead
	nodes := list.nodes
	links := list.links

	for level := list.level - 1; level >= 0; level-- {
		for link := links[nodes[current].linkStart+level]; link.next != skipListNil && traversed+link.span <= rank; link = links[nodes[current].linkStart+level] {
			traversed += link.span
			current = link.next
		}
		if traversed == rank {
			return nodes[current].value, nil
		}
	}

	var zero T
	return zero, listpkgIndexError(idx, list.length)
}

// LowerBound returns the first index whose value does not sort before x.
//
// It returns Length when every value sorts before x.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) LowerBound(x T) int {
	idx := 0
	current := skipListHead
	nodes := list.nodes
	links := list.links
	compare := list.compare

	for level := list.level - 1; level >= 0; level-- {
		for link := links[nodes[current].linkStart+level]; link.next != skipListNil && compare(nodes[link.next].value, x) < 0; link = links[nodes[current].linkStart+level] {
			idx += link.span
			current = link.next
		}
	}

	return idx
}

// UpperBound returns the first index whose value sorts after x.
//
// It returns Length when no value sorts after x.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) UpperBound(x T) int {
	idx := 0
	current := skipListHead
	nodes := list.nodes
	links := list.links
	compare := list.compare

	for level := list.level - 1; level >= 0; level-- {
		for link := links[nodes[current].linkStart+level]; link.next != skipListNil && compare(nodes[link.next].value, x) <= 0; link = links[nodes[current].linkStart+level] {
			idx += link.span
			current = link.next
		}
	}

	return idx
}

// EqualRange returns the half-open index range containing values equivalent to x.
//
// The returned range is [start, end), and is empty when start == end.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) EqualRange(x T) (start, end int) {
	start, node := list.lowerBoundNode(x)
	if node == skipListNil || list.compare(list.nodes[node].value, x) != 0 {
		return start, start
	}

	end, _ = list.upperBoundNode(x)
	return start, end
}

// Count returns the number of values equivalent to x.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) Count(x T) int {
	start, end := list.EqualRange(x)
	return end - start
}

// Find locates the first index equivalent to x according to the comparator.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) Find(x T) (idx int, ok bool) {
	idx, node := list.lowerBoundNode(x)
	if node == skipListNil || list.compare(list.nodes[node].value, x) != 0 {
		return -1, false
	}

	return idx, true
}

// Ceiling returns the first value that does not sort before x.
//
// It returns a zero value, -1, and false when every value sorts before x.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) Ceiling(x T) (value T, idx int, ok bool) {
	idx, node := list.lowerBoundNode(x)
	if node == skipListNil {
		return value, -1, false
	}

	return list.nodes[node].value, idx, true
}

// Floor returns the last value that does not sort after x.
//
// It returns a zero value, -1, and false when every value sorts after x.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) Floor(x T) (value T, idx int, ok bool) {
	rank, node := list.upperBoundPrevious(x)
	if node == skipListHead {
		return value, -1, false
	}

	return list.nodes[node].value, rank - 1, true
}

// Higher returns the first value that sorts after x.
//
// It returns a zero value, -1, and false when no value sorts after x.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) Higher(x T) (value T, idx int, ok bool) {
	idx, node := list.upperBoundNode(x)
	if node == skipListNil {
		return value, -1, false
	}

	return list.nodes[node].value, idx, true
}

// Lower returns the last value that sorts before x.
//
// It returns a zero value, -1, and false when no value sorts before x.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) Lower(x T) (value T, idx int, ok bool) {
	rank, node := list.lowerBoundPrevious(x)
	if node == skipListHead {
		return value, -1, false
	}

	return list.nodes[node].value, rank - 1, true
}

// Contains returns true if an equivalent value exists in the list.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) Contains(x T) bool {
	current := skipListHead
	nodes := list.nodes
	links := list.links
	compare := list.compare

	for level := list.level - 1; level >= 0; level-- {
		for link := links[nodes[current].linkStart+level]; link.next != skipListNil; link = links[nodes[current].linkStart+level] {
			order := compare(nodes[link.next].value, x)
			if order == 0 {
				return true
			}
			if order > 0 {
				break
			}
			current = link.next
		}
	}

	return false
}

// First returns the first value in comparator order.
//
// Complexity: O(1).
func (list *SkipList[T]) First() (T, bool) {
	if list.IsEmpty() {
		var zero T
		return zero, false
	}

	return list.nodes[list.next(skipListHead, 0)].value, true
}

// Last returns the last value in comparator order.
//
// Complexity: O(1).
func (list *SkipList[T]) Last() (T, bool) {
	if list.IsEmpty() {
		var zero T
		return zero, false
	}

	return list.nodes[list.tail].value, true
}

// Add inserts one or more values while preserving the sorted invariant.
//
// Complexity: expected O(K log (N+K)).
func (list *SkipList[T]) Add(xs ...T) {
	if list.length == 0 && len(xs) > 1 {
		values := slices.Clone(xs)
		slices.SortFunc(values, list.compare)
		list.buildFromSorted(values)
		return
	}

	if len(xs) > 1 {
		list.ensureFreeNodes(len(xs))
	}

	for _, x := range xs {
		list.insert(x)
	}
}

func (list *SkipList[T]) addSorted(xs []T) {
	if list.length == 0 {
		list.buildFromSorted(xs)
		return
	}

	list.ensureFreeNodes(len(xs))

	for _, x := range xs {
		list.insert(x)
	}
}

func (list *SkipList[T]) buildFromSorted(xs []T) {
	if len(xs) == 0 {
		return
	}

	list.ensureFreeNodes(len(xs))
	clear(list.links[list.nodes[skipListHead].linkStart : list.nodes[skipListHead].linkStart+list.nodes[skipListHead].linkCapacity])
	list.tail = skipListNil
	list.level = 1
	list.length = 0

	previous := skipListHead
	for idx, x := range xs {
		level := skipListPerfectLevel(idx + 1)
		if level > list.level {
			list.level = level
		}

		node := list.newValueNode(level, x)
		if previous != skipListHead {
			list.nodes[node].backward = previous
		}

		link := list.link(previous, 0)
		link.next = node
		link.span = 1
		previous = node
	}

	list.tail = previous
	list.length = len(xs)

	for level := 1; level < list.level; level++ {
		previous = skipListHead
		previousRank := 0
		rank := 1
		for node := list.next(skipListHead, 0); node != skipListNil; node = list.next(node, 0) {
			if list.nodes[node].height > level {
				link := list.link(previous, level)
				link.next = node
				link.span = rank - previousRank
				previous = node
				previousRank = rank
			}
			rank++
		}
	}
}

// Replace changes the value at idx while preserving the sorted invariant.
//
// It returns an index error when idx is out of bounds and ErrOrderViolation
// when x would sort before the previous value or after the next value.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) Replace(idx int, x T) error {
	if idx < 0 || idx >= list.length {
		return listpkgIndexError(idx, list.length)
	}

	node := list.nodeAt(idx)
	if previous := list.nodes[node].backward; previous != skipListNil && list.compare(list.nodes[previous].value, x) > 0 {
		return ErrOrderViolation
	}
	if next := list.next(node, 0); next != skipListNil && list.compare(x, list.nodes[next].value) > 0 {
		return ErrOrderViolation
	}

	list.nodes[node].value = x
	return nil
}

// Remove removes the first value equivalent to x according to the comparator.
//
// Complexity: expected O(log N).
func (list *SkipList[T]) Remove(x T) bool {
	var update [skipListMaxLevel]int
	current := skipListHead
	nodes := list.nodes
	links := list.links
	compare := list.compare

	for level := list.level - 1; level >= 0; level-- {
		for link := links[nodes[current].linkStart+level]; link.next != skipListNil && compare(nodes[link.next].value, x) < 0; link = links[nodes[current].linkStart+level] {
			current = link.next
		}
		update[level] = current
	}

	node := links[nodes[current].linkStart].next
	if node == skipListNil || compare(nodes[node].value, x) != 0 {
		return false
	}

	list.deleteNode(node, &update)
	return true
}

// All returns a sequence that yields elements in sorted order.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (list *SkipList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for node := list.next(skipListHead, 0); node != skipListNil; node = list.next(node, 0) {
			if !yield(list.nodes[node].value) {
				return
			}
		}
	}
}

// Backward returns a sequence that yields elements in reverse sorted order.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (list *SkipList[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		for node := list.tail; node != skipListNil; node = list.nodes[node].backward {
			if !yield(list.nodes[node].value) {
				return
			}
		}
	}
}

// Range returns a sequence that yields values in the half-open range [from, to).
//
// Values are selected according to comparator order: yielded values do not sort
// before from and do sort before to.
//
// Complexity: expected O(log N + K), where K is the number of yielded values.
func (list *SkipList[T]) Range(from, to T) iter.Seq[T] {
	return func(yield func(T) bool) {
		start, node := list.lowerBoundNode(from)
		end, _ := list.lowerBoundNode(to)
		if end < start {
			return
		}

		for idx := start; node != skipListNil && idx < end; idx++ {
			if !yield(list.nodes[node].value) {
				return
			}
			node = list.next(node, 0)
		}
	}
}

// Enumerate returns a sequence that yields the sorted index and value of each element.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (list *SkipList[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		idx := 0
		for node := list.next(skipListHead, 0); node != skipListNil; node = list.next(node, 0) {
			if !yield(idx, list.nodes[node].value) {
				return
			}
			idx++
		}
	}
}

// ToSlice exports the sorted elements into a native Go slice.
//
// Complexity: O(N).
func (list *SkipList[T]) ToSlice() []T {
	if list.length == 0 {
		return nil
	}

	values := make([]T, 0, list.length)
	for value := range list.All() {
		values = append(values, value)
	}
	return values
}

// Clear removes all elements from the list.
//
// Complexity: O(N).
func (list *SkipList[T]) Clear() {
	for node := list.next(skipListHead, 0); node != skipListNil; {
		next := list.next(node, 0)
		list.releaseNode(node)
		node = next
	}

	clear(list.links[list.nodes[skipListHead].linkStart : list.nodes[skipListHead].linkStart+list.nodes[skipListHead].linkCapacity])
	list.tail = skipListNil
	list.level = 1
	list.length = 0
}

// MarshalJSON converts the list into a JSON array in sorted order.
//
// Complexity: O(N).
func (list *SkipList[T]) MarshalJSON() ([]byte, error) {
	return collection.Marshal(list)
}

// UnmarshalJSON populates the list from a JSON array and restores sorted order.
//
// The input order is not preserved; values are sorted according to the list
// comparator before replacing the current contents.
//
// Complexity: expected O(M + N log N), where M is the current length and N is
// the number of decoded values.
func (list *SkipList[T]) UnmarshalJSON(data []byte) error {
	var values []T
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}

	slices.SortFunc(values, list.compare)
	list.Clear()
	list.addSorted(values)
	return nil
}

// Format implements fmt.Formatter.
//
// Complexity: O(1) as it respects a fixed display limit.
func (list *SkipList[T]) Format(s fmt.State, verb rune) {
	collection.Format(s, verb, list, list.Length())
}

// String returns a string representation of the list.
//
// Complexity: O(1) as it respects a fixed display limit.
func (list *SkipList[T]) String() string {
	return fmt.Sprint(list)
}

func (list *SkipList[T]) lowerBoundNode(x T) (idx int, node int) {
	current := skipListHead
	nodes := list.nodes
	links := list.links
	compare := list.compare

	for level := list.level - 1; level >= 0; level-- {
		for link := links[nodes[current].linkStart+level]; link.next != skipListNil && compare(nodes[link.next].value, x) < 0; link = links[nodes[current].linkStart+level] {
			idx += link.span
			current = link.next
		}
	}
	return idx, links[nodes[current].linkStart].next
}

func (list *SkipList[T]) upperBoundNode(x T) (idx int, node int) {
	current := skipListHead
	nodes := list.nodes
	links := list.links
	compare := list.compare

	for level := list.level - 1; level >= 0; level-- {
		for link := links[nodes[current].linkStart+level]; link.next != skipListNil && compare(nodes[link.next].value, x) <= 0; link = links[nodes[current].linkStart+level] {
			idx += link.span
			current = link.next
		}
	}
	return idx, links[nodes[current].linkStart].next
}

func (list *SkipList[T]) lowerBoundPrevious(x T) (rank int, node int) {
	current := skipListHead
	nodes := list.nodes
	links := list.links
	compare := list.compare

	for level := list.level - 1; level >= 0; level-- {
		for link := links[nodes[current].linkStart+level]; link.next != skipListNil && compare(nodes[link.next].value, x) < 0; link = links[nodes[current].linkStart+level] {
			rank += link.span
			current = link.next
		}
	}
	return rank, current
}

func (list *SkipList[T]) upperBoundPrevious(x T) (rank int, node int) {
	current := skipListHead
	nodes := list.nodes
	links := list.links
	compare := list.compare

	for level := list.level - 1; level >= 0; level-- {
		for link := links[nodes[current].linkStart+level]; link.next != skipListNil && compare(nodes[link.next].value, x) <= 0; link = links[nodes[current].linkStart+level] {
			rank += link.span
			current = link.next
		}
	}
	return rank, current
}

func (list *SkipList[T]) nodeAt(idx int) int {
	rank := idx + 1
	traversed := 0
	current := skipListHead
	nodes := list.nodes
	links := list.links

	for level := list.level - 1; level >= 0; level-- {
		for link := links[nodes[current].linkStart+level]; link.next != skipListNil && traversed+link.span <= rank; link = links[nodes[current].linkStart+level] {
			traversed += link.span
			current = link.next
		}
		if traversed == rank {
			return current
		}
	}

	return skipListNil
}

func (list *SkipList[T]) insert(x T) {
	var update [skipListMaxLevel]int
	var rank [skipListMaxLevel]int
	current := skipListHead
	nodes := list.nodes
	links := list.links
	compare := list.compare

	for level := list.level - 1; level >= 0; level-- {
		if level == list.level-1 {
			rank[level] = 0
		} else {
			rank[level] = rank[level+1]
		}

		for link := links[nodes[current].linkStart+level]; link.next != skipListNil && compare(nodes[link.next].value, x) <= 0; link = links[nodes[current].linkStart+level] {
			rank[level] += link.span
			current = link.next
		}
		update[level] = current
	}

	height := list.randomLevel()
	if height > list.level {
		for idx := list.level; idx < height; idx++ {
			rank[idx] = 0
			update[idx] = skipListHead
			links[nodes[skipListHead].linkStart+idx].span = list.length
		}
		list.level = height
	}

	node := list.newValueNode(height, x)
	nodes = list.nodes
	links = list.links

	for idx := range height {
		nodeLink := &links[nodes[node].linkStart+idx]
		updateLink := &links[nodes[update[idx]].linkStart+idx]
		nodeLink.next = updateLink.next
		updateLink.next = node

		skipped := rank[0] - rank[idx]
		nodeLink.span = updateLink.span - skipped
		updateLink.span = skipped + 1
	}

	for idx := height; idx < list.level; idx++ {
		links[nodes[update[idx]].linkStart+idx].span++
	}

	if update[0] != skipListHead {
		nodes[node].backward = update[0]
	}
	if next := links[nodes[node].linkStart].next; next != skipListNil {
		nodes[next].backward = node
	} else {
		list.tail = node
	}
	list.length++
}

func (list *SkipList[T]) deleteNode(node int, update *[skipListMaxLevel]int) {
	nodes := list.nodes
	links := list.links

	for level := 0; level < list.level; level++ {
		updateLink := &links[nodes[update[level]].linkStart+level]
		if updateLink.next == node {
			nodeLink := links[nodes[node].linkStart+level]
			updateLink.span += nodeLink.span - 1
			updateLink.next = nodeLink.next
		} else {
			updateLink.span--
		}
	}

	if next := links[nodes[node].linkStart].next; next != skipListNil {
		nodes[next].backward = nodes[node].backward
	} else {
		list.tail = nodes[node].backward
	}

	for list.level > 1 && links[nodes[skipListHead].linkStart+list.level-1].next == skipListNil {
		list.level--
	}
	list.length--
	list.releaseNode(node)
}

func (list *SkipList[T]) randomLevel() int {
	level := 1
	for level < skipListMaxLevel && list.nextRand()&3 == 0 {
		level++
	}
	return level
}

func (list *SkipList[T]) nextRand() uint64 {
	x := list.seed
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	list.seed = x
	return x
}

func (list *SkipList[T]) releaseNode(node int) {
	var zero T
	entry := &list.nodes[node]
	entry.value = zero
	clear(list.links[entry.linkStart : entry.linkStart+entry.linkCapacity])
	entry.backward = skipListNil
	entry.freeNext = list.free
	entry.height = 0
	list.free = node
	list.freeLength++
}

func listpkgIndexError(idx, length int) error {
	return list.NewIndexOutOfBoundError(idx, length-1)
}
