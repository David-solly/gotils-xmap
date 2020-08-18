package xmap

import "sync"

//Xmap : Creates an empty xmap Structure
func Xmap() XMap {
	return XMap{
		xmap:      make(map[string]**interface{}),
		xIDKy:     make(map[int]string),
		xKyID:     make(map[string]int),
		elements:  0,
		freeSpace: 0,
	}

}

//XMap : Fast iteration type map struct - at least 10x faster
type XMap struct {
	xmap      map[string]**interface{}
	xslice    []*interface{}
	xempty    []int
	xIDKy     map[int]string
	xKyID     map[string]int
	elements  int
	freeSpace int
	mutex     sync.Mutex
}

//Add : add to the structure by appending to structure
func (p *XMap) Add(k string, v interface{}) {
	p.xslice = append(p.xslice, &v)
	p.xmap[k] = &p.xslice[p.elements]
	p.xIDKy[p.elements] = k
	p.xKyID[k] = p.elements
	p.elements = len(p.xmap)

}

//AddEco : add to the structure by filling vacant spaces
func (p *XMap) AddEco(k string, v interface{}) int {
	if p.xempty != nil {
		i := p.xempty[0]
		p.xslice[i] = &v
		p.xKyID[k] = i
		if len(p.xempty) > 1 {
			p.xempty = p.xempty[1:]
			p.freeSpace = len(p.xempty)
		} else {
			p.xempty = nil
			p.freeSpace = 0
		}
		return i
	}

	// Add a new value the usual way
	p.Add(k, v)
	return p.elements

}

//Update : update a value by its index
func (p *XMap) Update(i int, v interface{}) (previous, newValue interface{}) {
	newVP := interface{}(v)
	previous = *p.GetByIndex(i).(*interface{})
	p.Slice()[i] = &newVP
	newValue = *p.GetByIndex(i).(*interface{})
	return
}

//Delete :remove an entry using a key to the entry
func (p *XMap) Delete(key string) {
	if _, ok := p.xmap[key]; ok {
		p.mutex.Lock()
		delete(p.xmap, key)
		p.DeleteElementAt(p.xKyID[key])
		p.elements = len(p.xmap)
		p.mutex.Unlock()
	}
}

//DeleteElementAt : Delete an entry using a value index
//Only sets the underlying slice location to nil
func (p *XMap) DeleteElementAt(idx int) {
	p.xslice[idx] = nil
	p.xempty = append(p.xempty, idx)
	p.freeSpace = len(p.xempty)
}

// IndexFreeSpace : Index the free space freed by Delete function
func (p *XMap) IndexFreeSpace() {
	go p.indexFreeSpace()
}

// FreeSpace : Available free space to re-populate
func (p *XMap) FreeSpace() int {
	return p.freeSpace
}

func (p *XMap) indexFreeSpace() {
	ind := []int{}
	for id, element := range p.xslice {
		if element == nil {
			ind = append(ind, id)
		}
	}

	if len(ind) > 0 {
		p.mutex.Lock()
		p.xempty = ind
		p.freeSpace = len(ind)
		p.mutex.Unlock()
	}

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
	p.mutex.Lock()
	p.elements = newm.elements
	p.xIDKy = newm.xIDKy
	p.xmap = newm.xmap
	p.xslice = newm.xslice
	p.freeSpace = 0
	p.xempty = nil
	p.mutex.Unlock()

}

// Slice : Returns the Slice of the values in the xmap structure, retains insertion order
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
	if !b {
		return nil, b
	}
	return **k, b
}

//GetByIndex : Retrives a value by its index inside the xmap structure
func (p *XMap) GetByIndex(index int) interface{} {
	k := p.xslice[index]
	return k
}

// GetIndexOf : Return the index of the item if it exists in the xmap structure, check for the boolean first as the default is to return nil on a non hit
func (p *XMap) GetIndexOf(value interface{}) (index int, exists bool) {
	for i, e := range p.xslice {
		if *e == value {
			return i, true
		}
	}
	return -1, false
}
