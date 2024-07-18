package tree23

import (
	"fmt"
	"slices"
	"sort"
)

type Node struct {
	keys  []int
	nodes []*Node
}

// Печатаем дерево
func (n *Node) PrintTree() {
	if nil == n {
		return
	}
	lvl := func(node *Node) int {
		var res int
		for node.parent() != nil {
			node = node.parent()
			res++
		}
		return res
	}
	var queue []*Node
	currentLvl := 0
	queue = append(queue, n)
	for len(queue) != 0 {
		currentNode := queue[0]
		if l := lvl(currentNode); currentLvl != l {
			currentLvl = l
			fmt.Println()
		}
		fmt.Printf("\t%v", currentNode.keys)
		for _, child := range currentNode.childs() {
			queue = append(queue, child)
		}
		queue = queue[1:]
	}
}

// Создаем новый узел и вносим ключ
func NewNode(key int) *Node {
	return &Node{
		keys:  []int{key},
		nodes: make([]*Node, 5),
	}
}

// Вставка ключей
func (n *Node) Insert(keys ...int) {
	for _, key := range keys {
		if n == nil {
			n = NewNode(key)
			continue
		}
		p := n
		if p.isLeaf() {
			p.insertToNode(key)
		} else {
			switch {
			case key <= p.keys[0]:
				p.first().Insert(key)
			case len(p.keys) == 1 || (len(p.keys) == 2 && key <= p.keys[1]):
				p.second().Insert(key)
			default:
				p.third().Insert(key)
			}
		}
		n = p.split()
	}
}

// Поиск узла с заданным ключом
func (n *Node) Search(key int) *Node {
	if nil == n {
		return nil
	}
	switch {
	case n.find(key):
		return n
	case key < n.keys[0]:
		return n.first().Search(key)
	case len(n.keys) == 2 && key < n.keys[1] || len(n.keys) == 1:
		return n.second().Search(key)
	default:
		return n.third().Search(key)
	}
}

// Удаление ключей
func (n *Node) Remove(keys ...int) {
	for _, key := range keys {
		node := n.Search(key)
		if nil == node {
			continue
		}
		var minNode *Node
		if node.keys[0] == key {
			minNode = node.second().findMinNode()
		} else {
			minNode = node.third().findMinNode()
		}
		if minNode != nil {
			if key == node.keys[0] {
				node.keys[0], minNode.keys[0] = minNode.keys[0], node.keys[0]
			} else {
				node.keys[1], minNode.keys[0] = minNode.keys[0], node.keys[1]
			}
		}
		node.removeFromNode(key)
		*n = *fix(node)

	}
}
func fix(n *Node) *Node {
	if len(n.keys) == 0 && n.parent() == nil {
		n = nil
		return n
	}
	if len(n.keys) != 0 {
		if n.parent() != nil {
			return fix(n.parent())
		}
		return n
	}
	parent := n.parent()
	if len(parent.first().keys) == 2 || len(parent.keys) == 2 || len(parent.second().keys) == 2 {
		n.redistribute()
	} else if len(parent.keys) == 2 && len(parent.third().keys) == 2 {
		n.redistribute()
	} else {
		n = merge(n)
	}
	return fix(n)

}
func (n *Node) redistribute() {
	parent := n.parent()
	first := parent.first()
	second := parent.second()
	third := parent.third()
	if len(parent.keys) == 2 && len(first.keys) < 2 && len(second.keys) < 2 {
		switch n {
		case first:
			parent.addFirst(second)
			parent.addSecond(third)
			parent.addThird(nil)
			parent.first().insertToNode(parent.keys[0])
			parent.first().addThird(parent.first().second())
			parent.first().addSecond(parent.first().first())

			if n.first() != nil {
				parent.first().addFirst(n.first())
			} else {
				parent.first().addFirst(n.second())
			}
			if parent.first().first() != nil {
				parent.first().first().addParent(parent.first())
			}
			parent.removeFromNode(parent.keys[0])
		case second:
			first.insertToNode(parent.keys[0])
			parent.removeFromNode(parent.keys[0])
			if n.first() != nil {
				first.addThird(n.first())
			} else if n.second() != nil {
				first.addThird(n.second())
			}

			if first.third() != nil {
				first.third().addParent(first)
			}
			parent.addSecond(parent.third())
			parent.addThird(nil)
		case third:
			second.insertToNode(parent.keys[1])
			parent.addThird(nil)
			parent.removeFromNode(parent.keys[1])
			if n.first() != nil {
				second.addThird(n.first())
			} else if n.second() != nil {
				second.addThird(n.second())
			}
			if second.third() != nil {
				second.third().addParent(second)
			}
		}
	} else if len(parent.keys) == 2 && (len(first.keys) == 2 || len(second.keys) == 2 || len(third.keys) == 2) {
		switch {
		case third == n:
			if n.first() != nil {
				n.addSecond(n.first())
				n.addFirst(nil)
			}
			n.insertToNode(parent.keys[1])
			if len(second.keys) == 2 {
				parent.keys[1] = second.keys[1]
				second.removeFromNode(second.keys[1])
				n.addFirst(second.third())
				second.addThird(nil)
				if n.first() != nil {
					n.first().addParent(n)
				}
			} else if len(first.keys) == 2 {
				parent.keys[1] = second.keys[0]
				n.addFirst(second.second())
				second.addSecond(second.first())
				if n.first() != nil {
					n.first().addParent(n)
				}

				second.keys[0] = parent.keys[0]
				parent.keys[0] = first.keys[1]
				first.removeFromNode(first.keys[1])
				second.addFirst(first.third())
				if second.first() != nil {
					second.first().addParent(second)
				}
				first.addThird(nil)
			}
		case second == n:
			if len(third.keys) == 2 {
				if n.first() == nil {
					n.addFirst(n.second())
					n.addSecond(nil)
				}
				second.insertToNode(parent.keys[1])
				parent.keys[1] = third.keys[0]
				third.removeFromNode(third.keys[0])
				second.addSecond(third.first())
				if second.second() != nil {
					second.second().addParent(second)
				}
				third.addFirst(third.second())
				third.addSecond(third.third())
				third.addThird(nil)
			} else if len(first.keys) == 2 {
				if n.second() == nil {
					n.addSecond(n.first())
					n.addFirst(nil)
				}
				second.insertToNode(parent.keys[0])
				parent.keys[0] = first.keys[1]
				first.removeFromNode(first.keys[1])
				second.addFirst(first.third())
				if second.first() != nil {
					second.first().addParent(second)
				}
				first.addThird(nil)
			}
		case first == n:
			if n.first() == nil {
				n.addFirst(n.second())
				n.addSecond(nil)
			}
			first.insertToNode(parent.keys[0])
			if len(second.keys) == 2 {
				parent.keys[0] = second.keys[0]
				second.removeFromNode(second.keys[0])
				first.addSecond(second.first())
				if first.second() != nil {
					first.second().addParent(first)
				}
				second.addFirst(second.second())
				second.addSecond(second.third())
				second.addThird(nil)
			} else if len(third.keys) == 2 {
				parent.keys[0] = second.keys[0]
				second.keys[0] = parent.keys[1]
				parent.keys[1] = third.keys[0]
				third.removeFromNode(third.keys[0])
				first.addSecond(second.first())
				if first.second() != nil {
					first.second().addParent(first)
				}
				second.addFirst(second.second())
				second.addSecond(third.first())
				if second.second() != nil {
					second.second().addParent(second)
				}
				third.addFirst(third.second())
				third.addSecond(third.third())
				third.addThird(nil)
			}

		}
	} else if len(parent.keys) == 1 {
		n.insertToNode(parent.keys[0])

		if first == n && len(second.keys) == 2 {
			parent.keys[0] = second.keys[0]
			second.removeFromNode(second.keys[0])

			if n.first() == nil {
				n.addFirst(n.second())
			}

			n.addSecond(second.first())
			second.addFirst(second.second())
			second.addSecond(second.third())
			second.addThird(nil)
			if n.second() != nil {
				n.second().addParent(n)
			}
		} else if second == n && len(first.keys) == 2 {
			parent.keys[0] = first.keys[1]
			first.removeFromNode(first.keys[1])

			if n.second() == nil {
				n.addSecond(n.first())
			}

			n.addFirst(first.third())
			first.addThird(nil)
			if n.first() != nil {
				n.first().addParent(n)
			}
		}
	}
	*n = *parent
}
func merge(n *Node) *Node {
	parent := n.parent()

	if parent.first() == n {
		parent.second().insertToNode(parent.keys[0])
		parent.second().addThird(parent.second().second())
		parent.second().addSecond(parent.second().first())

		if n.first() != nil {
			parent.second().addFirst(n.first())
		} else if n.second() != nil {
			parent.second().addFirst(n.second())
		}
		if parent.second().first() != nil {
			parent.second().first().addParent(parent.second())
		}

		parent.removeFromNode(parent.keys[0])
		parent.addFirst(nil)
	} else if parent.second() == n {
		parent.first().insertToNode(parent.keys[0])

		if n.first() != nil {
			parent.first().addThird(n.first())
		} else if n.second() != nil {
			parent.first().addThird(n.second())
		}

		if parent.first().third() != nil {
			parent.first().third().addParent(parent.first())
		}

		parent.removeFromNode(parent.keys[0])
		parent.addSecond(nil)
	}

	if parent.parent() == nil {
		var tmp *Node
		if parent.first() != nil {
			tmp = parent.first()
		} else {
			tmp = parent.second()
		}
		tmp.addParent(nil)
		return tmp
	}
	return parent
}
func (n *Node) findMinNode() *Node {
	if nil == n {
		return n
	}
	if nil == n.first() {
		return n
	}
	return n.first().findMinNode()
}

func (n *Node) addParent(p *Node) *Node {
	n.nodes[0] = p
	return n
}
func (n *Node) addFirst(p *Node) *Node {
	n.nodes[1] = p
	return n
}
func (n *Node) addSecond(p *Node) *Node {
	n.nodes[2] = p
	return n
}
func (n *Node) addThird(p *Node) *Node {
	n.nodes[3] = p
	return n
}
func (n *Node) addFourth(p *Node) *Node {
	n.nodes[4] = p
	return n
}
func (n *Node) split() *Node {
	if len(n.keys) < 3 {
		return n
	}
	p := n
	x := NewNode(p.keys[0]).addParent(p.parent()).addFirst(p.first()).addSecond(p.second())
	y := NewNode(p.keys[2]).addParent(p.parent()).addFirst(p.third()).addSecond(p.fourth())

	if x.first() != nil {
		x.first().addParent(x)
	}
	if x.second() != nil {
		x.second().addParent(x)
	}
	if y.first() != nil {
		y.first().addParent(y)
	}
	if y.second() != nil {
		y.second().addParent(y)
	}

	if nil == p.parent() {
		x.addParent(p)
		y.addParent(p)
		p.becomeNode2(p.keys[1], x, y)
		n = p
		return n
	}

	p.parent().insertToNode(p.keys[1])

	switch p {
	case p.parent().first():
		p.parent().addFirst(nil)
	case p.parent().second():
		p.parent().addSecond(nil)
	case p.parent().third():
		p.parent().addThird(nil)
	}

	switch {
	case nil == p.parent().first():
		p.parent().addFourth(p.parent().third()).addThird(p.parent().second()).addSecond(y).addFirst(x)
	case nil == p.parent().second():
		p.parent().addFourth(p.parent().third()).addThird(y).addSecond(x)
	default:
		p.parent().addFourth(y).addThird(x)
	}
	n = p.parent()
	return n

}

func (n *Node) find(k int) bool {
	return slices.Contains(n.keys, k)
}
func (n *Node) sort() {
	sort.Ints(n.keys)
}

// Вставляем ключ в вершину
func (n *Node) insertToNode(k int) {
	n.keys = append(n.keys, k)
	sort.Ints(n.keys)
}

// Удаляем ключ из вершины
func (n *Node) removeFromNode(key int) {
	for i, k := range n.keys {
		if k == key {
			n.keys = append(n.keys[:i], n.keys[i+1:]...)
		}
	}
}

// Преобразовать в 2-вершину
func (n *Node) becomeNode2(k int, first, second *Node) {
	*n = *NewNode(k).addFirst(first).addSecond(second)
}

// Проверка является ли узел листом
func (n *Node) isLeaf() bool {
	return n.first() == nil && n.second() == nil && n.third() == nil
}
func (n *Node) parent() *Node {
	return n.nodes[0]
}
func (n *Node) first() *Node {
	return n.nodes[1]
}
func (n *Node) second() *Node {
	return n.nodes[2]
}
func (n *Node) third() *Node {
	return n.nodes[3]
}
func (n *Node) fourth() *Node {
	return n.nodes[4]
}

func (n *Node) childs() []*Node {
	res := make([]*Node, 0, 4)
	for _, node := range n.nodes[1:] {
		if node != nil {
			res = append(res, node)
		}
	}
	return res
}
