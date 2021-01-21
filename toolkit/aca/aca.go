// Aho-Corasick automation
package aca

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

type node struct {
	next       map[rune]*node
	fail       *node
	wordLength int
}

type ACA struct {
	root      *node
	nodeCount int
}

// New
// 创建一个只有root节点的Trie树
func New() *ACA {
	return &ACA{root: &node{}, nodeCount: 1}
}

// Add
// 添加模式串到Trie树
func (a *ACA) Add(word string) {
	n := a.root
	for _, r := range word {
		if n.next == nil {
			n.next = make(map[rune]*node)
		}
		if n.next[r] == nil {
			n.next[r] = &node{}
			a.nodeCount++
		}
		n = n.next[r]
	}
	n.wordLength = len(word)
}

// Del
// 从Trie删除模式串
func (a *ACA) Del(word string) {
	rs := []rune(word)
	stack := make([]*node, len(rs))
	n := a.root

	for i, r := range rs {
		if n.next[r] == nil {
			return
		}
		stack[i] = n
		n = n.next[r]
	}

	// if it is NOT the leaf node
	if len(n.next) > 0 {
		n.wordLength = 0
		return
	}

	// if it is the leaf node
	for i := len(rs) - 1; i >= 0; i-- {
		stack[i].next[rs[i]].next = nil
		stack[i].next[rs[i]].fail = nil

		delete(stack[i].next, rs[i])
		a.nodeCount--
		if len(stack[i].next) > 0 ||
			stack[i].wordLength > 0 {
			return
		}
	}
}

// Build
// 构建AC
func (a *ACA) Build() {
	// allocate enough memory as a queue
	q := append(make([]*node, 0, a.nodeCount), a.root)

	for len(q) > 0 {
		n := q[0]
		q = q[1:]

		for r, c := range n.next {
			q = append(q, c)

			p := n.fail
			for p != nil {
				if p.next[r] != nil {
					c.fail = p.next[r]
					break
				}
				p = p.fail
			}
			if p == nil {
				c.fail = a.root
			}
		}
	}
}

// find
func (a *ACA) find(s string, cb func(start, end int)) {
	n := a.root
	for i, r := range s {
		for n.next[r] == nil && n != a.root {
			n = n.fail
		}
		n = n.next[r]
		if n == nil {
			n = a.root
			continue
		}

		end := i + utf8.RuneLen(r)
		for t := n; t != a.root; t = t.fail {
			if t.wordLength > 0 {
				cb(end-t.wordLength, end)
			}
		}
	}
}

// Find
// 从AC查找模式串
func (a *ACA) Find(s string) (words []string) {
	a.find(s, func(start, end int) {
		words = append(words, s[start:end])
	})
	return
}


type pinYin struct {
	FileName    string
	Content     string
	Aca 		*ACA
	RawContents []string
	Words       []string
	PreResult   []string
	Trim        map[string]struct{}
}

func NewPinYin() *pinYin {
	p := &pinYin{
		FileName:    "pinyin2phone_20191128.bk",
		Trim:        make(map[string]struct{}),
	}
	return p
}

func (p *pinYin) LoadData() (*pinYin, error) {
	file, err := os.Open(p.FileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineWord := strings.Split(line, "|")[0]
		if len(lineWord) > 1 {
			p.RawContents = append(p.RawContents, lineWord)
		}
	}

	p.Aca = New()
	for _, c := range p.RawContents {
		p.Aca.Add(c)
	}
	p.Aca.Build()
	return p, nil
}

func (p *pinYin) Guess(content string) []string {
	p.Content = strings.ToLower(content)
	p.Words = p.Aca.Find(p.Content)

	prefixWords := make([]string, 0)
	for j := 0; j < len(p.Words); j++ {
		//单字名
		if p.Words[j] == p.Content {
			p.PreResult = append(p.PreResult, p.Words[j])
			continue
		}
		if strings.HasPrefix(p.Content, p.Words[j]) {
			//去重，避免叠字重复匹配
			if _, ok := p.Trim[p.Words[j]]; !ok {
				prefixWords = append(prefixWords, p.Words[j])
				p.Trim[p.Words[j]] = struct{}{}
			}
		}
	}
	p.compose(prefixWords)
	return p.PreResult
}

func (p *pinYin) compose(prefixWords []string) {
	pwl := len(prefixWords)
	sign := "'"
	for j := 0; j < pwl; j++ {
		suffixWords := p.Words[pwl:]
		if len(suffixWords) > 0 {
			tmp := make([]string, 0)
			for i := 0; i < len(suffixWords); i++ {
				showSubWords := fmt.Sprintf("%s%s%s", prefixWords[j], sign, suffixWords[i])
				compPrefixWords := strings.ReplaceAll(prefixWords[j], sign, "")
				compSubWords := compPrefixWords + suffixWords[i]
				if p.Content == compSubWords {
					if _, ok := p.Trim[showSubWords]; !ok {
						p.PreResult = append(p.PreResult, showSubWords)
						p.Trim[showSubWords] = struct{}{}
					}
				} else if strings.HasPrefix(p.Content, compSubWords) {
					tmp = append(tmp, showSubWords)
				}
			}
			if len(tmp) > 0 {
				p.compose(tmp)
			}
		}
	}
}
