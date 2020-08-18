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

func TestXmapFeatures(t *testing.T) {
	xm := Xmap()
	limit := 100
	t.Run("Test xmap Add insert", func(t *testing.T) {
		for i := 0; i < limit; i++ {
			xm.Add(fmt.Sprintf("key%d", i), i)
		}
		assert.Equal(t, xm.Count(), limit)
		assert.Equal(t, xm.elements, limit)
	})

	t.Run("Test xmap Delete", func(t *testing.T) {
		xm.Delete(fmt.Sprintf("key%d", 89))
		assert.Equal(t, xm.Count(), limit-1)
		assert.Equal(t, xm.elements, limit-1)
		assert.Equal(t, len(xm.xslice), limit) // Still keep 100 items

	})

	t.Run("Test xmap DeleteElementAt", func(t *testing.T) {
		eToDelete := 15
		xm.DeleteElementAt(eToDelete)
		assert.Equal(t, xm.Count(), limit-1)
		assert.Equal(t, xm.elements, limit-1)
		assert.Equal(t, len(xm.xslice), limit) // Still keep 100 items

		b := xm.GetByIndex(eToDelete)
		assert.Equal(t, true, b != nil)
		assert.Equal(t, true, b.(*interface{}) == nil)
		c, k := xm.GetByKey(fmt.Sprintf("key%d", eToDelete))
		assert.Equal(t, true, k)
		assert.NotNil(t, c)
		assert.Equal(t, xm.FreeSpace(), 2)

	})

	t.Run("Test xmap index free space", func(t *testing.T) {
		// Deleting elements directly from the underlying struct will
		// require manual indexing
		xm.Slice()[16] = nil
		xm.IndexFreeSpace()
		assert.Equal(t, xm.FreeSpace() > 0, true)
	})

	t.Run("Test xmap Eco Add over free space", func(t *testing.T) {
		eToDelete := 12
		xm.DeleteElementAt(eToDelete)

		diff := []struct {
			key   string
			Value interface{}
			Index int
		}{
			{"hi", 2332, 89}, //first deleted index
			{"key--1234", "121254", 15},
			{"key--123", "121255", 12},
			{"key--12", 125.554, 100}, //list full start expanding
			{"key--1", 125.554, 101},  // expanded structure
		}

		for _, data := range diff {
			newIDX := xm.AddEco(data.key, data.Value)
			assert.Equal(t, newIDX, data.Index)
			assert.Equal(t, xm.elements, len(xm.Map()))
		}
	})

	t.Run("Test xmap RebuildIndex", func(t *testing.T) {
		xm.Delete(fmt.Sprintf("key%d", 0))
		xm.RebuildIndex()
		assert.Equal(t, xm.Count(), limit)
		assert.Equal(t, xm.elements, limit)
		c, k := xm.GetByKey(fmt.Sprintf("key%d", 0))
		assert.Equal(t, k, false)
		assert.Equal(t, c == nil, true)

	})

	t.Run("Test xmap Get underlying Slice", func(t *testing.T) {
		sl := xm.Slice()
		assert.Equal(t, len(sl), limit)

	})

	t.Run("Test xmap ADD different datatypes", func(t *testing.T) {
		xmi := Xmap()
		diff := []struct {
			key  string
			Data interface{}
		}{
			{"hi", 2332},
			{"key--1234", "121254"},
			{"key--123", "121255"},
			{"key--12", 125.554},
			{"key--1", 125.554},
		}

		for i, data := range diff {
			xmi.Add(data.key, data.Data)

			t.Run("RETREIVE different value types BY INDEX and BY KEY", func(t *testing.T) {
				d, ok := xmi.GetByKey(data.key)
				assert.Equal(t, ok, true)
				assert.Equal(t, *xmi.GetByIndex(i).(*interface{}), data.Data)
				assert.Equal(t, *xmi.GetByIndex(i).(*interface{}), d)
				assert.Contains(t, fmt.Sprintf("%v", *xmi.GetByIndex(i).(*interface{})), fmt.Sprintf("%v", data.Data))
				assert.Equal(t, reflect.TypeOf(*xmi.GetByIndex(i).(*interface{})), reflect.TypeOf(data.Data))
			})
		}
	})

}

func TestXmapValueAlteration(t *testing.T) {
	xm := Xmap()
	value := 0
	newVal := "123456789"
	key := "key-0"
	xm.Add(key, value)

	t.Run("Change value of data - long way", func(t *testing.T) {
		newVP := interface{}(newVal)
		a := xm.GetByIndex(0)
		assert.Equal(t, true, a != nil)
		assert.Equal(t, 0, *a.(*interface{}))
		xm.Slice()[0] = &newVP
		b := xm.GetByIndex(0)
		assert.Equal(t, true, b != nil)
		assert.Equal(t, *b.(*interface{}), newVal)
		assert.Equal(t, *b.(*interface{}), newVal)
		c, k := xm.GetByKey(key)
		assert.Equal(t, true, k)
		assert.Equal(t, true, c != nil)
		assert.Equal(t, c, newVal)

	})

	t.Run("Change value of data - recommended way", func(t *testing.T) {
		newValue := "1234567891"
		prev := newVal //previous step value
		newVal = newValue
		a := xm.GetByIndex(0)
		assert.Equal(t, true, a != nil)
		assert.Equal(t, prev, *a.(*interface{}))

		//Overwrite key value
		xm.Add(key, newValue)

		c, k := xm.GetByKey(key)
		assert.Equal(t, true, k)
		assert.Equal(t, true, c != nil)
		assert.Equal(t, c, newValue)

	})

	t.Run("Change value of data - via index recommended way", func(t *testing.T) {
		xm := Xmap()
		value := 0
		key := "key-0"
		index := 0
		xm.Add(key, value)

		newValue := "1234567891"

		a := xm.GetByIndex(0)
		assert.Equal(t, true, a != nil)
		assert.Equal(t, value, *a.(*interface{}))

		//Update value
		prev, current := xm.Update(index, newValue)
		assert.Equal(t, prev, value)
		assert.Equal(t, current, newValue)

		// test new value is stored in underlying index
		c, k := xm.GetByKey(key)
		assert.Equal(t, true, k)
		assert.Equal(t, true, c != nil)
		assert.Equal(t, c, newValue)

	})

	t.Run("Get index of value", func(t *testing.T) {
		xm := Xmap()
		diff := []struct {
			Key   string
			Value interface{}
			Found bool
			Index int
		}{
			{"key-0", "1234567891", true, 0},
			{"", "12lkl34567891", false, -1},
		}

		xm.Add(diff[0].Key, diff[0].Value)

		for _, data := range diff {
			a, k := xm.GetIndexOf(data.Value)
			assert.Equal(t, k, data.Found)
			assert.Equal(t, a, data.Index)

			if k {
				keyV := xm.Indexes()[a]
				assert.Equal(t, keyV, data.Key)

			}
		}
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
