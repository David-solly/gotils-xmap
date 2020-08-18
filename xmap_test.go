package xmap

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
)

func TestXmap(t *testing.T) {
	t.Run("Test gotils xmap", func(t *testing.T) {
		assert.Equal(t, reflect.TypeOf(Xmap()), reflect.TypeOf(XMap{}))
	})
}

func BenchmarkInsert(b *testing.B) {
	var im map[interface{}]interface{}
	var jm map[string]int
	var xm XMap
	var keys []string

	b.Run("Initiate keys", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			keys = append(keys, fmt.Sprintf("String%d", i))
		}

	})

	// interface typed map
	b.Run("Run inserts interface map", func(b *testing.B) {
		im = make(map[interface{}]interface{})
		for i := 0; i < b.N && i < len(keys); i++ {
			im[keys[i]] = i
		}

	})

	b.Run("Iterate all elements interface map", func(b *testing.B) {
		for _, val := range im {
			_ = val
		}
	})

	b.Run("Run deletes interface map", func(b *testing.B) {
		for id := 0; id < b.N && id < len(im); id++ {
			delete(im, keys[id])
		}
	})

	// Typed map
	b.Run("Run inserts typed map", func(b *testing.B) {
		jm = make(map[string]int)
		for i := 0; i < b.N && i < len(keys); i++ {
			jm[keys[i]] = i
		}
	})

	b.Run("Iterate all elements typed map", func(b *testing.B) {
		for _, val := range jm {
			_ = val
		}
	})

	b.Run("Run deletes typed map", func(b *testing.B) {
		for id := 0; id < b.N && id < len(jm); id++ {
			delete(jm, keys[id])
		}
	})

	// Xmap structure
	b.Run("Run inserts xmap", func(b *testing.B) {
		xm = Xmap()
		for i := 0; i < b.N && i < len(keys); i++ {
			xm.Add(keys[i], i)
		}
	})

	b.Run("Iterate all elements xmap", func(b *testing.B) {
		for _, val := range xm.Slice() {
			_ = val
		}
	})

	b.Run("Run deletes xmap", func(b *testing.B) {
		for id := 0; id < b.N && id < xm.Count(); id++ {
			xm.DeleteElementAt(id)
		}
	})

	b.Run("Resize xmap indexes", func(b *testing.B) {
		xm.RebuildIndex()
	})

}
