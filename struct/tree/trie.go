package tree

import "fmt"

// Trie(前缀树)

type TrieNode struct {
	isWord bool
	next   map[string]*TrieNode
}

func (t *TrieNode) String() string {
	return fmt.Sprintf("[isWord:%v,next:%v]", t.isWord, t.next)
}

type Trie struct {
	root *TrieNode
	size int
}

func (trie *Trie) String() string {
	return fmt.Sprintf("[size:%v,root:%v]", trie.size, trie.root)
}

func NewTrie() *Trie {
	return &Trie{
		root: NewNode(),
	}
}

// 获得Trie中存储的单词数量
func (trie *Trie) Size() int {
	return trie.size
}

func NewNode() *TrieNode {
	return &TrieNode{
		next: make(map[string]*TrieNode),
	}
}

// 向Trie中添加一个新的单词word
func (trie *Trie) Add(word string) {
	cur := trie.root

	for _, w := range []rune(word) {
		c := string(w)

		if cur.next[c] == nil {
			cur.next[c] = NewNode()
		}
		cur = cur.next[c]
	}

	if cur.isWord == false {
		cur.isWord = true
		trie.size++
	}
}

/** Adds a word into the data structure. */
func (trie *Trie) AddWord(word string) {
	cur := trie.root

	for _, w := range []rune(word) {
		c := string(w)
		if cur.next[c] == nil {
			cur.next[c] = NewNode()
		}
		cur = cur.next[c]
	}

	cur.isWord = true
}

/** Returns if the word is in the data structure. A word could contain the dot character '.' to represent any one letter. */
func (trie *Trie) Search(word string) bool {
	return trie.match(trie.root, word, 0)
}

func (trie *Trie) match(n *TrieNode, word string, index int) bool {
	if index == len(word) {
		return n.isWord
	}

	c := string([]rune(word)[index])

	if c != "." {
		if n.next[c] == nil {
			return false
		}
		return trie.match(n.next[c], word, index+1)
	} else {
		for w := range n.next {
			if trie.match(n.next[w], word, index+1) {
				return true
			}
		}
		return false
	}
}
