package main

import (
	"fmt"
	"os"
	"strings"
)

/*
 * infTree structure constructor
 */
func newInfTree() *infTree {
	var t infTree
	t.fact = nil
	t.head = nil
	t.left = nil
	t.right = nil
	t.precedence = 10
	return &t
}

/*
 * fact structure constructor
 */
func newFact() *fact {
	var f fact
	f.op = ""
	f.isKnown = false
	f.isTrue = false
	return &f
}

/*
 * Build the inference tree with all facts and statements
 * https://www.rhyscitlema.com/algorithms/expression-parsing-algorithm/
 */
func buildTree() {
	var root *infTree

	for _, rule := range env.rules {
		root = newInfTree()
		root.precedence = 1
		//root.fact.op = openBra
		var current = root
		for i := 0; i < len(rule); i++ {
			if rule[i] != ' ' && rule[i] != '\t' {
				if i+3 < len(rule) && rule[i:i+3] == ioi {
					current = buildLeaf(root, current, ioi)
					i += 2
				} else if i+2 < len(rule) && rule[i:i+2] == imp {
					current = buildLeaf(root, current, imp)
					i++
				} else {
					current = buildLeaf(root, current, string(rule[i]))
				}
			}
		}
		if root.right != nil {
			root.right.head = nil
		}
		root = root.right
		env.trees = append(env.trees, *root)
	}
}

func buildLeaf(root *infTree, current *infTree, c string) *infTree {
	var node = newInfTree()
	node.fact = newFact()
	var info = noInfo

	if c == openBra {
		node.precedence = openBraPre
		node.fact.op = openBra
		info = skipClimbUp
	} else if c == closeBra {
		node.precedence = closeBraPre
		node.fact.op = closeBra
		info = rightAssociative
	} else if c == ioi {
		node.precedence = ioiPre
		node.fact.op = ioi
		info = rightAssociative
	} else if c == imp {
		node.precedence = impPre
		node.fact.op = imp
		info = rightAssociative
	} else if c == not {
		node.precedence = notPre
		node.fact.op = not
	} else if c == and {
		node.precedence = andPre
		node.fact.op = and
	} else if c == or {
		node.precedence = orPre
		node.fact.op = or
	} else if c == xor {
		node.precedence = xorPre
		node.fact.op = xor
	} else if strings.Contains(factSymbol, c) {
		node.precedence = factPre
		node.fact = env.factList[c]
	} else {
		fmt.Printf("bug parse : '%s'\n", c)
		os.Exit(1)
	}
	current = insertNodeItem(current, *node, info)
	return (current)
}

func insertNodeItem(current *infTree, item infTree, info nodeInfo) *infTree {
	var node *infTree

	if info != skipClimbUp {
		/* step 4: climb up */
		if info != rightAssociative {
			/* for left-associative */
			for current.precedence >= item.precedence {
				current = current.head
			}
		} else {
			/* for right-associative */
			for current.precedence > item.precedence {
				current = current.head
			}
		}
	}
	if item.fact.op == closeBra {
		/* step 5.1: remove the '(' node */
		node = current.head
		node.right = current.right
		if current.right != nil {
			current.right.head = node
		}
		/* step 6: Set the 'current node' to be the parent node */
		current = node
	} else {
		/* step 5.1: create the new node */
		node = newInfTree()
		*node = item
		node.right = nil
		/* step 5.2: add the new node */
		node.left = current.right
		if current.right != nil {
			current.right.head = node
		}
		current.right = node
		node.head = current
		/* step 6: Set the 'current node' to be the new node */
		current = node
	}
	return current
}

/*
 * Print infTree with variable indetation
 */
func printNode(node *infTree, indent int) {
	if node == nil {
		return
	}
	printNode(node.right, indent+4)

	for i := 0; i < indent; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("%v\n", node.fact)
	printNode(node.left, indent+4)
}
