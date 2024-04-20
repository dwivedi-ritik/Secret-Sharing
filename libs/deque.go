package lib

type Node struct {
	value      interface{}
	prev, next *Node
}

type Deque struct {
	front, rear *Node
}

func (d *Deque) IsEmpty() bool {
	return d.front == nil
}

func (d *Deque) AddFront(value interface{}) {
	newNode := &Node{value: value}

	if d.front == nil {
		d.front, d.rear = newNode, newNode
		return
	}

	newNode.next = d.front
	d.front.prev = newNode
	d.front = newNode
}

func (d *Deque) AddRear(value interface{}) {
	newNode := &Node{value: value}

	if d.rear == nil {
		d.front, d.rear = newNode, newNode
		return
	}

	newNode.prev = d.rear
	d.rear.next = newNode
	d.rear = newNode
}

func (d *Deque) RemoveFront() interface{} {
	if d.IsEmpty() {
		return nil
	}

	value := d.front.value
	d.front = d.front.next
	if d.front == nil {
		d.rear = nil
	} else {
		d.front.prev = nil
	}
	return value
}

func (d *Deque) RemoveRear() interface{} {
	if d.IsEmpty() {
		return nil
	}

	value := d.rear.value
	d.rear = d.rear.prev
	if d.rear == nil {
		d.front = nil
	} else {
		d.rear.next = nil
	}
	return value
}

func (d *Deque) PeekFront() interface{} {
	if d.IsEmpty() {
		return nil
	}
	return d.front.value
}

func (d *Deque) PeekRear() interface{} {
	if d.IsEmpty() {
		return nil
	}
	return d.rear.value
}
