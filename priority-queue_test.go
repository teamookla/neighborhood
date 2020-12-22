package neighborhood


import (
"testing"
)

func TestPriorityQueue_Peek(t *testing.T) {
	q := newPriorityQueue(10)

	seattle := namedPoint("seattle")
	memphis := namedPoint("memphis")
	woodinville := namedPoint("woodinville")

	q.PushPoint(woodinville, 234)
	q.PushPoint(seattle, 123)
	q.PushPoint(memphis, 2000)

	assertEqual(t, 3, q.Len())
	peeked := q.Peek().point
	assertEqual(t, 3, q.Len()) // didn't remove
	assertEqual(t, seattle.Name, peeked.(*NamedPoint).Name)
	popped := q.PopItem().point
	assertEqual(t, peeked.(*NamedPoint).Name, popped.(*NamedPoint).Name)
	assertEqual(t, 2, q.Len()) // did remove
}

func TestPriorityQueue_Empty(t *testing.T) {
	q := newPriorityQueue(10)
	assertNil(t, q.Peek())
	assertNil(t, q.PopItem())
}
