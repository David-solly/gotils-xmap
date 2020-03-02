package xmap

//Xmap : Creates an empty xmap Structure
func Xmap() XMap {
	return XMap{
		xmap:     make(map[string]**interface{}),
		xIDKy:    make(map[int]string),
		elements: 0,
	}

}

//XMap : Fast iteration type map struct - at least 10x faster
type XMap struct {
	xmap     map[string]**interface{}
	xslice   []*interface{}
	xIDKy    map[int]string
	elements int
}

//Add : add to the structure
func (p *XMap) Add(k string, v interface{}) {
	p.xslice = append(p.xslice, &v)
	p.xmap[k] = &p.xslice[p.elements]
	p.xIDKy[p.elements] = k
	p.elements = len(p.xmap)

}

//Delete :remove an entry using a key to the entry
func (p *XMap) Delete(key string) {
	if _, ok := p.xmap[key]; ok {
		delete(p.xmap, key)
		p.elements = len(p.xmap)
	}

}

//DeleteElementAt : Delete an entry using a value index
//Only sets the underlying slice location to nil
func (p *XMap) DeleteElementAt(idx int) {
	p.xslice[idx] = nil
}

//RebuildIndex : Use this periodically to garbage collect nil values from underlying structure. WARNING: calling this function resets all the indexes, so if you had 5 items and removed #2, running this will shift the remaining indexes left to fill the empty slot
func (p *XMap) RebuildIndex() {
	newm := Xmap()
	for id, element := range p.xslice {
		if element != nil {
			v := p.xIDKy[id]
			newm.Add(v, element)

		}
	}
	*p = newm
}

// Slice : Reteurns the Slice of the values in the xmap structure, retains insertion order
func (p *XMap) Slice() []*interface{} {
	return p.xslice
}

// Map : Returns the map contained within the structure - map values are pointers to actual values
func (p *XMap) Map() map[string]**interface{} {
	return p.xmap
}

// Indexes : Returns the indexing map used to store indexes
func (p *XMap) Indexes() map[int]string {

	return p.xIDKy
}

// Count : the amount of values in the xmap structure
func (p *XMap) Count() int {
	return p.elements
}

// GetByKey : Get the value by key as a normal map function
func (p *XMap) GetByKey(key string) (interface{}, bool) {
	k, b := p.xmap[key]
	return **k, b
}

//GetByIndex : Retrives a value by its index inside the xmap structure
func (p *XMap) GetByIndex(index int) interface{} {
	k := p.xslice[index]
	return k
}

// GetIndexOf : Return the index of the item if it exists in the xmap structure, check for the boolean first as the default is to return nil on a non hit
func (p *XMap) GetIndexOf(value interface{}) (index interface{}, exists bool) {
	for i, e := range p.xslice {
		if e == value {
			return i, true
		}
	}
	return nil, false
}
