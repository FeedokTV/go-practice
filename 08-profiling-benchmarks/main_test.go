package main_test

import (
	profilingbenchmarks "go-practice/08-profiling-benchmarks"
	"testing"
)

func TestBothBuildersReturnSame(t *testing.T) {
	strs := []string{"This", "is", "such", "a", "nice", "benchmark"}
	modelValue := "Thisissuchanicebenchmark"

	v1 := profilingbenchmarks.BuildStringV1(strs)
	v2 := profilingbenchmarks.BuildStringV2(strs)

	if v1 != modelValue {
		t.Errorf("expected from v1: %s, got: %s", modelValue, v1)
	}

	if v2 != modelValue {
		t.Errorf("expected from v2: %s, got: %s", modelValue, v2)
	}

}

func BenchmarkBuildStringV1(b *testing.B) {
	strs := []string{"This", "is", "such", "a", "nice", "benchmark"}

	for b.Loop() {
		profilingbenchmarks.BuildStringV1(strs)
	}
}

func BenchmarkBuildStringV2(b *testing.B) {
	strs := []string{"This", "is", "such", "a", "nice", "benchmark"}

	for b.Loop() {
		profilingbenchmarks.BuildStringV2(strs)
	}
}
