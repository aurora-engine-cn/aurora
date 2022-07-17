package aurora

type element struct {
	value *node
	next  *element
}

type queue struct {
	head *element
	end  *element
}

func (qe *queue) next() *element {
	if qe.head == nil {
		return nil
	}
	if qe.head.next != nil {
		t := qe.head
		qe.head = qe.head.next
		return t
	} else {
		t := qe.head
		qe.head = nil
		return t
	}
}

func (qe *queue) en(value *node) {
	if qe.head == nil {
		e := &element{value: value, next: nil}
		qe.end = e
		qe.head = e
		return
	}
	e := &element{value: value, next: nil}
	qe.end.next = e
	qe.end = e
}
