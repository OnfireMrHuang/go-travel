package algorithm


type doublyLinkedNode struct {
	prev,next * doublyLinkedNode
	key,val int
}

type LruCache struct {
	len,cap int
	first,last *doublyLinkedNode
	nodes map[int]*doublyLinkedNode
}

func CreateLruCache(cap int) *LruCache {
	return &LruCache{
		len: 0,
		cap: cap,
		first:nil,
		last: nil,
		nodes: make(map[int]*doublyLinkedNode,cap),
	}
}

func (l *LruCache) Get(key int) int {
	if node,ok := l.nodes[key]; ok {
		l.moveToFirst(node)
		return node.val
	}
	return -1
}

func (l *LruCache) Put(key,val int)  {
	if node,ok := l.nodes[key]; ok {
		node.val = val
		l.moveToFirst(node)
	} else {
		if l.len == l.cap {
			delete(l.nodes,l.last.key)
			l.removeLast()
		} else {
			l.len++
		}
		node := &doublyLinkedNode{
			prev: nil,
			next: nil,
			key: key,
			val: val,
		}
		l.nodes[key] = node
		l.insertToFirst(node)
	}
}

func (l *LruCache) removeLast()  {
	if l.last.prev != nil {
		l.last.prev.next = nil
	} else {
		l.first = nil
	}
	l.last = l.last.prev
}

func (l *LruCache) insertToFirst(node *doublyLinkedNode)  {
	if l.last == nil {
		l.last = node
	} else {
		l.first.prev = node
		node.next = l.first
	}
	l.first = node
}

func (l *LruCache) moveToFirst(node *doublyLinkedNode)  {
	switch node {
	case l.first:
		return
	case l.last:
		l.removeLast()
	default:
		node.prev.next = node.next
		node.next.prev = node.prev
	}
	l.insertToFirst(node)
}
