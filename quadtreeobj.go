package quadtree

import (
	"fmt"
	"image"
)

type Object struct {
	bound  image.Rectangle
	curpos image.Point
	node   *QuadTree
}

func (o Object) String() string {
	return fmt.Sprintf("bound:%v, curpos:%v, node in: %v", o.bound, o.curpos, o.node)
}

func NewObject(bound image.Rectangle) *Object {
	return &Object{bound: bound}
}

func (o Object) Bound() image.Rectangle {
	return o.bound
}

func (o Object) CurrentPos() image.Point {
	return o.curpos
}
