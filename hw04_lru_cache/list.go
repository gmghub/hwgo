package hw04lrucache

import (
	"fmt"
)

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
	first *ListItem
	last  *ListItem
	len   int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
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
		l.last = li
	} else {
		l.first.Prev = li
		li.Next = l.first
	}
	l.first = li
	l.len++
	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := &ListItem{
		Value: v,
	}
	if l.first == nil && l.last == nil {
		l.first = li
	} else {
		l.last.Next = li
		li.Prev = l.last
	}
	l.last = li
	l.len++
	return li
}

func (l *list) Remove(i *ListItem) {
	switch {
	case i.Prev == nil && i.Next == nil:
		if l.first != i && l.last != i {
			return
		}
		l.first = nil
		l.last = nil
	case i.Prev == nil:
		if i.Next.Prev != i {
			return
		}
		i.Next.Prev = nil
		l.first = i.Next
	case i.Next == nil:
		if i.Prev.Next != i {
			return
		}
		i.Prev.Next = nil
		l.last = i.Prev
	default:
		if i.Prev.Next != i || i.Next.Prev != i {
			return
		}
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.first {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}

func (l list) String() string {
	s := fmt.Sprintf("first: %v last: %v len: %v items:", l.first, l.last, l.len)
	for i := l.first; i != nil; i = i.Next {
		s += fmt.Sprintf(" %v", i)
	}
	return s
}
