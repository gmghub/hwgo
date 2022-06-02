package hw04lrucache

import "fmt"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first   *ListItem
	last    *ListItem
	listmap map[*ListItem]*ListItem
}

func NewList() List {
	return &list{
		listmap: make(map[*ListItem]*ListItem),
	}
}

func (l *list) Len() int {
	return len(l.listmap)
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := &ListItem{
		Value: v,
	}
	if l.first == nil && l.last == nil {
		l.listmap[li] = li
		l.first = li
		l.last = li
		return li
	}
	l.first.Prev = li
	li.Next = l.first
	l.listmap[li] = li
	l.first = li
	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := &ListItem{
		Value: v,
	}
	if l.first == nil && l.last == nil {
		l.listmap[li] = li
		l.first = li
		l.last = li
		return li
	}
	l.last.Next = li
	li.Prev = l.last
	l.listmap[li] = li
	l.last = li
	return li
}

func (l *list) Remove(i *ListItem) {
	if _, ok := l.listmap[i]; !ok {
		return
	}
	switch {
	case i.Prev == nil && i.Next == nil:
		delete(l.listmap, i)
		l.first = nil
		l.last = nil
	case i.Prev == nil:
		next := l.listmap[i.Next]
		delete(l.listmap, i)
		next.Prev = nil
		l.first = next
	case i.Next == nil:
		prev := l.listmap[i.Prev]
		delete(l.listmap, i)
		prev.Next = nil
		l.last = prev
	default:
		prev := l.listmap[i.Prev]
		next := l.listmap[i.Next]
		delete(l.listmap, i)
		prev.Next = next
		next.Prev = prev
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if _, ok := l.listmap[i]; !ok {
		return
	}
	if i == l.first {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}

func (l list) String() string {
	s := fmt.Sprintf("first: %v\nlast: %v\n", l.first, l.last)
	for k, v := range l.listmap {
		s += fmt.Sprintf("k: %v v: %v prev: %v next: %v\n", k, v, v.Prev, v.Next)
	}
	return s
}
