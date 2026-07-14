// Package list provides mutable generic linear collections with indexed access.
//
// Lists support positional reads, writes, insertion, removal, forward and
// backward iteration, and reversal. Empty or invalid indexed operations return
// package-owned errors rather than panicking.
//
// # Implementations
//
// [ArrayList] stores values in a contiguous slice. Choose it for O(1) indexed
// access and writes, amortized O(1) append, and memory-local traversal.
// Insertion and removal at arbitrary positions are O(N).
//
// [LinkedList] stores values in a doubly linked sequence. Choose it when O(1)
// reversal and boundary insertion or removal matter more than indexed access.
// Indexed operations require O(N) traversal.
//
// Select an implementation with [Array] or [Linked]. Each selector returns a
// reusable concrete factory whose construction methods return pointers to the
// corresponding concrete list type.
//
// # Construction And Ownership
//
// Array factories provide New, From, Clone, and FromSeq. [ArrayFactory.From]
// takes ownership of the provided slice and may reuse its backing array during
// later mutations. [ArrayFactory.Clone] preserves the source by making a shallow
// copy. FromSeq collects an iterator into new storage.
//
// Linked factories provide New, From, and FromSeq. Linked construction always
// copies values into newly allocated nodes and never retains a source slice.
//
// # Views And Iteration
//
// [List] exposes mutation, while [Readonly] contains observation and traversal
// operations. [AsReadonly] prevents callers from recovering the mutable list
// through a type assertion. Iterators follow the standard iter package and stop
// without further work when the yield function returns false.
//
// All traverses the current logical list order from first to last, including
// after Reverse. Enumerate uses the same order and assigns consecutive indexes
// starting at zero. Backward traverses that logical order in reverse.
package list
