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
// Choose an implementation through [NewArray], [ArrayFrom], [ArrayClone],
// [ArrayFromSeq], [NewLinked], [LinkedFrom], or [LinkedFromSeq]. Each
// constructor returns a pointer to the corresponding concrete list type.
// The zero values of [ArrayList] and [LinkedList] are ready for use.
//
// # Construction And Ownership
//
// [ArrayFrom] takes ownership of the provided slice and may reuse its backing
// array during later mutations. The caller must stop using the source and all
// aliases afterward. [ArrayClone] makes an independent shallow copy of its
// storage. [ArrayFromSeq] collects an iterator into new storage once.
//
// [LinkedFrom] copies values into newly allocated nodes and never retains a
// source slice. [LinkedFromSeq] consumes an iterator once into new nodes.
//
// # JSON
//
// Lists marshal and unmarshal as arrays in index order. Unmarshal replaces the
// current contents only after the complete JSON array is decoded successfully.
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
//
// A readonly view observes the same underlying list; it is not a snapshot and
// does not make concurrent mutation safe.
package list
