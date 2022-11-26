package main

import (
	"testing"
)

func BenchmarkNaive(b *testing.B) {
	confs := Naive()

	for _, c := range confs {
		if c == nil {
			b.FailNow()
		}
	}
}

func BenchmarkWithGoKeyword(b *testing.B) {
	confs := WithGoKeyword()

	for _, c := range confs {
		if c == nil {
			b.FailNow()
		}
	}
}

func BenchmarkWithGoroutines(b *testing.B) {
	confs := WithGoroutines()

	for _, c := range confs {
		if c == nil {
			b.FailNow()
		}
	}
}

func BenchmarkWithChannels(b *testing.B) {
	confs := WithChannels()

	for _, c := range confs {
		if c == nil {
			b.FailNow()
		}
	}
}
