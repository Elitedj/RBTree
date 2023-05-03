## Red Black Tree

Red Black Tree implemented by Go.

## Usage
```go
// Make a rbt
rbt := NewRBTree[int]()

// Insert a value into tree
rbt.Insert(1)

// Check if a certain value exists
rbt.Has(1)

// Remove a value from the tree
rbt.Delete(1)

// Release
rbt.Clear()
```

## Benchmark
```
goos: linux
goarch: amd64
pkg: github.com/elitedj/rbtree
cpu: AMD EPYC 7K62 48-Core Processor
BenchmarkInsert          4255707               324.2 ns/op
BenchmarkDelete         10957584               113.8 ns/op
```