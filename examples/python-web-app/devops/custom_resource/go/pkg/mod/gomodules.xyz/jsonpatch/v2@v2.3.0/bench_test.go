package jsonpatch_test

import (
	"encoding/json"
	"testing"

	"gomodules.xyz/jsonpatch/v2"
)

func BenchmarkCreatePatch(b *testing.B) {
	cases := []struct {
		name string
		a, b string
	}{
		{
			"complex",
			superComplexBase,
			superComplexA,
		},
		{
			"large array",
			largeArray(1000),
			largeArray(1000),
		},
		{
			"simple",
			simpleA,
			simpleB,
		},
	}

	for _, tt := range cases {
		b.Run(tt.name, func(b *testing.B) {
			at := []byte(tt.a)
			bt := []byte(tt.b)
			for n := 0; n < b.N; n++ {
				_, _ = jsonpatch.CreatePatch(at, bt)
			}
		})
	}
}

func largeArray(size int) string {
	type nested struct {
		A, B string
	}
	type example struct {
		Objects []nested
	}
	a := example{}
	for i := 0; i < size; i++ {
		a.Objects = append(a.Objects, nested{A: "a", B: "b"})
	}
	res, _ := json.Marshal(a)
	return string(res)
}
