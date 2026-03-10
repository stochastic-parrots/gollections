# Performance: Gollections vs Standard Library

Last Updated: March 2026
Environment: Linux, AMD Ryzen 9 7900, Go 1.23+

Gollections Binary Heap significantly outperforms Go's `container/heap` package in every tested scenario,
delivering massive gains in speed, memory efficiency, and allocation reduction.

## 🛠️ Reproduce These Results

```bash
go test -bench=. -benchmem -count=10 -benchtime=1000x ./internal/heaps/binary/... > result.txt
grep "Library=Gollections" result.txt | sed 's/Library=Gollections//' > gollections.txt
grep "Library=Container/Heap" result.txt | sed 's|Library=Container/Heap||' > stdlib.txt
benchstat stdlib.txt gollections.txt
```

## 📊 Results Summary

| Metric | Gollections | StdLib | Improvement |
|--------|------------|--------|-------------|
| **Push Operations** | 3.3µs | 13.0µs | **74% faster** |
| **Pop Operations** | 6.0ms | 8.9ms | **32% faster** |
| **Memory per Op** | 0 B | 4.1 MB | **100% reduction** |
| **Allocations per Op** | 0 | 199k | **100% fewer** |

## 🏃 Execution Time Comparison

| Operation | StdLib | Gollections | Improvement |
|-----------|--------|-------------|-------------|
| Push (100 items) | 13.01µs | 3.34µs | **74.36%** |
| Push (300 items) | 25.01µs | 9.53µs | **61.91%** |
| Push (400 items) | 34.54µs | 12.75µs | **63.09%** |
| Push (500 items) | 40.08µs | 15.98µs | **60.12%** |
| Push (1000 items) | 93.19µs | 31.99µs | **65.67%** |
| Push (10000 items) | 1384.9µs | 225.4µs | **83.72%** |
| Pop (100k items) | 8.87ms | 6.05ms | **31.79%** |
| K-Largest Selection | 5.62ms | 2.64ms | **52.98%** |
| **Geometric Mean** | **203.1µs** | **72.30µs** | **64.41%** |

## 💾 Memory Allocations

| Operation | StdLib | Gollections | Improvement |
|-----------|--------|-------------|-------------|
| Push (100 items) | 1 alloc | 0 allocs | **100%** |
| Push (1000 items) | 3 allocs | 0 allocs | **100%** |
| Push (10000 items) | 10 allocs | 0 allocs | **100%** |
| Pop (100k items) | 99.7k allocs | 0 allocs | **100%** |
| K-Largest Selection | 199.9k allocs | 1 alloc | **99.94%** |

## 📦 Bytes Allocated

| Operation | StdLib | Gollections | Improvement |
|-----------|--------|-------------|-------------|
| Push (100 items) | 106.5 KB | 1 B | **100%** |
| Push (1000 items) | 425.9 KB | 2 B | **100%** |
| Push (10000 items) | 4.1 MB | 16 B | **100%** |
| Pop (100k items) | 779.3 KB | 0 B | **100%** |
| K-Largest Selection | 1.6 MB | 896 B | **99.94%** |

## 📈 Feature Comparison

| Feature | Gollections | container/heap | Notes |
|---------|------------|-----------------|-------|
| Type-Safe Generics | ✅ | ❌ | Uses interface{} |
| Zero Allocations | ✅ | ❌ | Reduces GC pressure |
| CPU Cache Friendly | ✅ | ❌ | Data locality matters |
| Compiler Inlining | ✅ | ❌ | Better optimization |
| Custom Comparators | ✅ | ✅ | Both supported |
| Standard Library | ❌ | ✅ | Maturity & adoption |

## 🎯 Key Advantages

### Zero Allocations After Initialization

Gollections performs **zero heap allocations** during operations, while `container/heap` triggers between
1 and 199,910 allocations per operation.

**Impact**: Reduced garbage collection pressure, fewer stop-the-world pauses, and more predictable latency under high load.

### Data Locality & Cache Efficiency

- **Gollections**: Elements stored contiguously in a typed slice (`[]T`)
- **StdLib**: Elements stored in a slice of interfaces (`[]interface{}`), requiring pointer indirection

Contiguous memory dramatically improves CPU cache utilization, resulting in faster traversals and comparisons.

### Monomorphization

The Go compiler generates specific, optimized code for each concrete type (e.g., `BinaryHeap[int]`, `BinaryHeap[string]`).
This eliminates the cost of interface{} indirection and enables aggressive inlining.
