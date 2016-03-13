package quadtree

import (
	"image"
	"testing"
)

func TestFailToInsertOverBound(t *testing.T) {
	qt := New(image.Rect(0, 0, 10, 10), 20)

	// objbount over quadtree bound
	if qt.Insert(NewObject(image.Rect(0, 0, 11, 10)), image.Pt(0, 0)) {
		t.Errorf("should be failed to insert over bound")
	}

	// pos over quadtree bound
	if qt.Insert(NewObject(image.Rect(0, 0, 1, 1)), image.Pt(10, 10)) {
		t.Errorf("should be failed to insert over bound")
	}

	// pos over quadtree bound
	if qt.Insert(NewObject(image.Rect(0, 0, 1, 1)), image.Pt(-1, 10)) {
		t.Errorf("should be failed to insert over bound")
	}
}

func TestSetCurretPosAfterInsert(t *testing.T) {
	qt := New(image.Rect(0, 0, 10, 10), 20)

	pos := image.Pt(4, 6)
	newobj := NewObject(image.Rect(0, 0, 2, 2))
	qt.Insert(newobj, pos)
	if curpos := newobj.CurrentPos(); pos != curpos {
		t.Errorf("current pos invalid %v vs %v", pos, curpos)
	}
	if newobj.node != qt {
		t.Errorf("unexpected quadtree for newobj: %v %v", newobj, qt)
	}
}

func TestSplitIntoChildrenWhenOverCapacity(t *testing.T) {
	qt := New(image.Rect(0, 0, 10, 10), 2)
	obj1 := NewObject(image.Rect(0, 0, 1, 1))
	obj2 := NewObject(image.Rect(0, 0, 8, 8))
	if !qt.Insert(obj1, image.Pt(1, 1)) {
		t.Errorf("failed to insert: %v", obj1)
	}
	if !qt.Insert(obj2, image.Pt(1, 1)) {
		t.Errorf("failed to insert: %v", obj2)
	}
	if qt.children[topLeft] != nil {
		t.Errorf("quadtree should not be split yet: %v", qt.children)
	}

	obj3 := NewObject(image.Rect(0, 0, 1, 1))
	if !qt.Insert(obj3, image.Pt(1, 1)) {
		t.Errorf("failed to insert: %v", obj3)
	}
	if qt.children[topLeft] == nil {
		t.Errorf("quadtree should be split: %v", qt)
	}
	// obj2 cannot be split
	if len(qt.objects) != 1 {
		t.Errorf("unexpected len: %v", len(qt.objects))
	}
	if child := qt.children[topLeft]; len(child.objects) != 2 {
		t.Errorf("unexpected len: %v", len(child.objects))
	}
	if obj1.node != qt.children[topLeft] {
		t.Errorf("unexpected quadtree for obj1: %v %v", obj1, qt.children[topLeft])
	}
	if obj2.node != qt {
		t.Errorf("unexpected quadtree for obj2: %v %v", obj2, qt)
	}
	if obj3.node != qt.children[topLeft] {
		t.Errorf("unexpected quadtree for obj3: %v %v", obj3, qt.children[topLeft])
	}
}

func TestInsertIntoParentWhenCannotSplitMore(t *testing.T) {
	qt := New(image.Rect(0, 0, 1, 4), 1)

	qt.Insert(NewObject(image.Rect(0, 0, 1, 1)), image.Pt(0, 0))
	qt.Insert(NewObject(image.Rect(0, 0, 1, 1)), image.Pt(0, 0))
	qt.Insert(NewObject(image.Rect(0, 0, 1, 1)), image.Pt(0, 0))
	qt.Insert(NewObject(image.Rect(0, 0, 1, 1)), image.Pt(0, 0))
	qt.Insert(NewObject(image.Rect(0, 0, 1, 1)), image.Pt(0, 0))

	if len(qt.objects) != 5 {
		t.Errorf("unexpected len: %v", len(qt.objects))
	}
}

func TestRemove(t *testing.T) {
	qt := New(image.Rect(0, 0, 10, 10), 2)
	obj1 := NewObject(image.Rect(0, 0, 1, 1))
	obj2 := NewObject(image.Rect(0, 0, 8, 8))
	obj3 := NewObject(image.Rect(0, 0, 1, 1))

	qt.Insert(obj1, image.Pt(0, 0))
	qt.Insert(obj2, image.Pt(0, 0))
	qt.Insert(obj3, image.Pt(0, 0))

	// top left child len should be 2
	if l := len(qt.children[topLeft].objects); l != 2 {
		t.Errorf("unexpected len: %v", l)
	}
	// root len should be 1
	if l := len(qt.objects); l != 1 {
		t.Errorf("unexpected len: %v", l)
	}

	if !qt.Remove(obj1) {
		t.Errorf("failed to remove: %v", obj1)
	}
	// obj1 was in top left child
	if l := len(qt.children[topLeft].objects); l != 1 {
		t.Errorf("unexpected len: %v", l)
	}
	// root not changed
	if l := len(qt.objects); l != 1 {
		t.Errorf("unexpected len: %v", l)
	}

	if !qt.Remove(obj2) {
		t.Errorf("failed to remove: %v", obj2)
	}
	// obj2 was in root
	if l := len(qt.objects); l != 0 {
		t.Errorf("unexpected len: %v", l)
	}
	// top left not changed
	if l := len(qt.children[topLeft].objects); l != 1 {
		t.Errorf("unexpected len: %v", l)
	}

	obj4 := NewObject(image.Rect(0, 0, 1, 1))
	if qt.Remove(obj4) {
		t.Errorf("not inserted object cannot be removed: %v", obj4)
	}
}

func TestSearch(t *testing.T) {
	qt := New(image.Rect(0, 0, 10, 10), 2)
	obj1 := NewObject(image.Rect(0, 0, 2, 2))
	obj2 := NewObject(image.Rect(0, 0, 8, 8))
	obj3 := NewObject(image.Rect(0, 0, 3, 5))
	obj4 := NewObject(image.Rect(0, 0, 4, 8))
	obj5 := NewObject(image.Rect(0, 0, 1, 1))

	qt.Insert(obj1, image.Pt(2, 4)) // max: [4, 6]
	qt.Insert(obj2, image.Pt(2, 0)) // max: [10, 8]
	qt.Insert(obj3, image.Pt(5, 5)) // max: [8, 10]
	qt.Insert(obj4, image.Pt(0, 0)) // max: [4, 8]
	qt.Insert(obj5, image.Pt(9, 9)) // max: [10, 10]

	if objs := qt.Search(qt.bound); len(objs) != 5 {
		t.Errorf("unexpected len: %v", objs)
	}

	// obj3, obj4, obj5
	if objs := qt.Search(image.Rect(6, 6, 10, 10)); len(objs) != 3 {
		t.Errorf("unexpected len: %v", objs)
	}

	// obj5
	if objs := qt.Search(image.Rect(9, 9, 10, 10)); len(objs) != 1 {
		t.Log(obj3.Bound().Add(image.Pt(5, 5)))
		t.Errorf("unexpected len: %v", objs)
	}

	// obj3, obj5
	if objs := qt.Search(image.Rect(7, 9, 10, 10)); len(objs) != 2 {
		t.Log(obj3.Bound().Add(image.Pt(5, 5)))
		t.Errorf("unexpected len: %v", objs)
	}

	// obj1, obj2, obj4
	if objs := qt.Search(image.Rect(0, 0, 5, 5)); len(objs) != 3 {
		t.Errorf("unexpected len: %v", objs)
	}
}

func TestMoveInNotSplitQuadTree(t *testing.T) {
	qt := New(image.Rect(0, 0, 10, 10), 2)
	obj1 := NewObject(image.Rect(0, 0, 2, 2))
	qt.Insert(obj1, image.Pt(2, 2)) // max: [4, 6]

	mvpos := image.Pt(8, 8)

	if !qt.Move(obj1, mvpos) {
		t.Errorf("failed to mv: %v", qt.objects)
	}
	if len(qt.objects) != 1 {
		t.Errorf("unexpected len: %v", qt.objects)
	}
	if curpos := obj1.CurrentPos(); curpos != mvpos {
		t.Errorf("unexpected pos: %v", curpos)
	}
}

func TestMoveIntoAnotherNode(t *testing.T) {
	qt := New(image.Rect(0, 0, 10, 10), 4)
	obj1 := NewObject(image.Rect(0, 0, 2, 2))
	obj2 := NewObject(image.Rect(0, 0, 2, 2))
	obj3 := NewObject(image.Rect(0, 0, 2, 2))
	obj4 := NewObject(image.Rect(0, 0, 2, 2))
	obj5 := NewObject(image.Rect(0, 0, 2, 2))

	qt.Insert(obj1, image.Pt(0, 0))
	qt.Insert(obj2, image.Pt(0, 0))
	qt.Insert(obj3, image.Pt(0, 0))
	qt.Insert(obj4, image.Pt(0, 0))
	qt.Insert(obj5, image.Pt(0, 0))

	beforenode := qt.children[topLeft].children[topLeft]
	if l := len(beforenode.objects); l != 5 {
		t.Errorf("unexpected len: %v", l)
	}

	mvpos := image.Pt(1, 1)
	if !qt.Move(obj1, mvpos) {
		t.Errorf("failed to mv")
	}
	if obj1.CurrentPos() != mvpos {
		t.Errorf("unexpected current pos: %v", obj1.curpos)
	}
	if obj1.node != qt {
		t.Errorf("unexpected node: %v %v", obj1.node, qt)
	}
	if l := len(qt.objects); l != 1 {
		t.Errorf("unexpected len: %v", l)
	}
	if l := len(beforenode.objects); l != 4 {
		t.Errorf("unexpected len: %v", l)
	}
}

func TestMoveIntoSameNode(t *testing.T) {
	qt := New(image.Rect(0, 0, 10, 10), 4)
	obj1 := NewObject(image.Rect(0, 0, 2, 2))
	obj2 := NewObject(image.Rect(0, 0, 2, 2))
	obj3 := NewObject(image.Rect(0, 0, 2, 2))
	obj4 := NewObject(image.Rect(0, 0, 2, 2))
	obj5 := NewObject(image.Rect(0, 0, 2, 2))

	qt.Insert(obj1, image.Pt(0, 0))
	qt.Insert(obj2, image.Pt(0, 0))
	qt.Insert(obj3, image.Pt(4, 4))
	qt.Insert(obj4, image.Pt(4, 4))
	qt.Insert(obj5, image.Pt(4, 4))

	beforenode := qt.children[topLeft]
	if l := len(beforenode.objects); l != 2 {
		t.Errorf("unexpected len: %v", l)
	}

	mvpos := image.Pt(1, 1)
	if !qt.Move(obj1, mvpos) {
		t.Errorf("failed to mv")
	}
	if obj1.CurrentPos() != mvpos {
		t.Errorf("unexpected current pos: %v", obj1.curpos)
	}
	if obj1.node != beforenode {
		t.Errorf("unexpected node: %v %v", obj1.node, qt)
	}
	if l := len(beforenode.objects); l != 2 {
		t.Errorf("unexpected len: %v", l)
	}
}
