package lru_cach_service

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
	ExternalID Key
	Value      interface{}
	Next       *ListItem
	Prev       *ListItem
}

type list struct {
	length int
	head   *ListItem
	tail   *ListItem
}

func NewList() List {
	return new(list)
}

func (list *list) Len() int {
	return list.length
}

func (list *list) Front() *ListItem {
	return list.head
}

func (list *list) Back() *ListItem {
	return list.tail
}

func (list *list) PushFront(value interface{}) *ListItem {
	if list.head == nil {
		insertFirstItem(list, value)
	} else {
		firstItem := list.head
		item := &ListItem{
			Value: value,
			Prev:  nil,
			Next:  firstItem,
		}

		firstItem.Prev = item
		list.head = item
	}

	list.length++

	return list.head
}

func (list *list) PushBack(value interface{}) *ListItem {
	if list.tail == nil {
		insertFirstItem(list, value)
	} else {
		item := &ListItem{
			Value: value,
			Prev:  list.tail,
			Next:  nil,
		}

		list.tail.Next = item
		list.tail = item
	}

	list.length++

	return list.tail
}

func (list *list) Remove(item *ListItem) {
	if item.Prev == nil {
		list.head = nil
		if item.Next != nil {
			list.head = item.Next
			item.Next.Prev = nil
		}
	}

	if item.Next == nil {
		list.tail = nil
		if item.Prev != nil {
			list.tail = item.Prev
			item.Prev.Next = nil
		}
	}

	if item.Prev != nil {
		if item.Next != nil {
			item.Next.Prev = item.Prev
		}
	}

	if item.Next != nil {
		if item.Prev != nil {
			item.Prev.Next = item.Next
		}
	}

	list.length--
}

func (list *list) MoveToFront(item *ListItem) {
	if list.length > 1 {
		list.Remove(item)
		item.Prev = nil
		item.Next = list.head
		list.head.Prev = item
		list.head = item
	}
}

func insertFirstItem(list *list, value interface{}) {
	item := &ListItem{
		Value: value,
		Prev:  nil,
		Next:  nil,
	}

	list.head = item
	list.tail = item
}
