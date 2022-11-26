# Parf
A concurrent parsing experiment in golang.  The idea is to explore and
understand some light improvements and optimizations we're able to do
in Go when parsing a whole bunch of file.

## Setup
```fish
# Generating 9999 sample yaml files
for i in (seq 0 9999); cp ./samples/sample.yml ./samples/"$i"_sample.yml; end
```

## Benchmarks
```bash
$ go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: parf
cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
BenchmarkNaive-8                1000000000               0.3140 ns/op          0 B/op          0 allocs/op
BenchmarkWithGoKeyword-8        1000000000               0.1857 ns/op          0 B/op          0 allocs/op
BenchmarkWithGoroutines-8       1000000000               0.1675 ns/op          0 B/op          0 allocs/op
BenchmarkWithChannels-8         1000000000               0.1433 ns/op          0 B/op          0 allocs/op
PASS
ok      parf    11.379s
```
