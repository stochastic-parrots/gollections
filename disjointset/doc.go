// Package disjointset provides disjoint-set unions over fixed, inclusive
// integer ranges. Every integer between min and max belongs to the structure
// from construction onward, initially as a singleton set. Union merges sets;
// Reset restores the initial partition. Individual values cannot be removed.
//
// [IntsRangeByRank] uses union by rank. Choose [NewIntsRangeByRank] when callers
// need the rank heuristic but not set sizes. [IntsRangeBySize] uses union by
// size. Choose [NewIntsRangeBySize] when callers need the size of a set. Both
// implementations use path compression, and their zero values are invalid.
//
// Len counts values in the range and remains constant; Disjoints counts the
// current sets. IsEmpty is always false for a constructed instance. All and
// Enumerate traverse the range in increasing order, regardless of unions.
// Find returns a current representative, which may change after Union.
// Connected reports whether two values share a set and distinguishes values
// outside the range with its second result.
// Path compression combined with union by rank or size gives Find, Connected,
// Union, Rank, and Size O(α(N)) amortized time across operations between resets,
// where N is the range length and α is the inverse Ackermann function. One
// operation can still take O(log N) time; Reset takes O(N).
//
// [AsReadonly] restricts access to observation operations but is not a
// snapshot. Find and Connected may compress internal paths. These structures
// are not safe for concurrent use; callers must synchronize shared access
// when any goroutine may call an operation that changes internal paths or sets.
//
// The partition has no JSON representation in this package. An array of range
// values alone would not identify which values belong to the same set.
package disjointset
