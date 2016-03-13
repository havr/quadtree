package quadtree

import "image"

const (
	topLeft int = iota
	topRight
	bottomLeft
	bottomRight
)

type QuadTree struct {
	bound    image.Rectangle
	objects  []*Object
	children [4]*QuadTree
	capacity int
}

func New(bound image.Rectangle, capacity int) *QuadTree {
	return &QuadTree{
		capacity: capacity,
		bound:    bound,
		objects:  make([]*Object, 0, capacity),
	}
}

// It inserts an object into quadtree recursively.
// Objects are moved when tree has been split.
// An object can't be inserted into children will be inserted into parent again.
// node field in Object struct will be set after being inserted for fast move and remove.
// CurrentPos funcion of an object returns point given when this function called.
func (qt *QuadTree) Insert(o *Object, p image.Point) bool {
	return qt.insert(o, p)
}

// It removes an object from quadtree node refered by itself.
// Refered node is set when inserted.
func (qt *QuadTree) Remove(o *Object) bool {
	return o.node != nil && o.node.remove(o)
}

// It searchs all objects overlapping given bound.
func (qt *QuadTree) Search(bound image.Rectangle) []*Object {
	return qt.search(bound)
}

// It moves an object into given point.
func (qt *QuadTree) Move(o *Object, p image.Point) bool {
	if o.node == nil {
		return false
	}
	mvbound := o.Bound().Add(p)
	if mvbound.In(o.node.bound) {
		o.curpos = p
		return true
	}
	return o.node.remove(o) && qt.insert(o, p)
}

func (qt *QuadTree) search(bound image.Rectangle) []*Object {
	if !bound.Overlaps(qt.bound) {
		return nil
	}
	var rtn []*Object
	if qt.children[topLeft] != nil {
		for i, _ := range qt.children {
			rtn = append(rtn, qt.children[i].search(bound)...)
		}
	}
	for _, o := range qt.objects {
		if bound.Overlaps(o.Bound().Add(o.CurrentPos())) {
			rtn = append(rtn, o)
		}
	}
	return rtn
}

func (qt *QuadTree) insert(o *Object, p image.Point) bool {
	if !o.Bound().Add(p).In(qt.bound) {
		return false
	}

	qt.objects = append(qt.objects, o)
	o.curpos = p
	o.node = qt

	// objects over capacity, has no children, can split more,
	if len(qt.objects) > qt.capacity && qt.children[topLeft] == nil && qt.bound.Dx()/2 > 0 && qt.bound.Dy()/2 > 0 {
		minX, minY := qt.bound.Min.X, qt.bound.Min.Y
		maxX, maxY := qt.bound.Max.X, qt.bound.Max.Y
		halfX, halfY := minX+qt.bound.Dx()/2, minY+qt.bound.Dy()/2

		qt.children[topLeft] = New(image.Rect(minX, minY, halfX, halfY), qt.capacity)
		qt.children[topRight] = New(image.Rect(halfX, minY, maxX, halfY), qt.capacity)
		qt.children[bottomLeft] = New(image.Rect(minX, halfY, halfX, maxY), qt.capacity)
		qt.children[bottomRight] = New(image.Rect(halfX, halfY, maxX, maxY), qt.capacity)

		var newobjs []*Object
		for _, mvobj := range qt.objects {
			if qt.insertIntoChildren(mvobj, mvobj.CurrentPos()) {
				continue
			}
			newobjs = append(newobjs, mvobj)
		}
		qt.objects = newobjs
	}
	return true
}

func (qt *QuadTree) insertIntoChildren(o *Object, p image.Point) bool {
	for i, _ := range qt.children {
		if qt.children[i].insert(o, p) {
			return true
		}
	}
	return false
}

// It removes in own objects only, not recursively.
func (qt *QuadTree) remove(o *Object) bool {
	for i, myobj := range qt.objects {
		if myobj == o {
			qt.objects = append(qt.objects[:i], qt.objects[i+1:]...)
			o.node = nil
			return true
		}
	}
	return false
}
