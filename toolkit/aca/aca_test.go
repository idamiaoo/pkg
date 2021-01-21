package aca

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

var rawContents []string
var content string
var words []string

func loadData() {
	file, err := os.Open("pinyin2phone_20191128.bk")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	rawContents = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineWord := strings.Split(line, "|")[0]
		if len(lineWord) > 1 {
			rawContents = append(rawContents, lineWord)
		}
	}
}

func compose(prefixWords []string, count *int) {
	pwl := len(prefixWords)
	for j := 0; j < pwl; j++ {
		suffixWords := words[pwl:]
		if len(suffixWords) > 0 {
			tmp := make([]string, 0)
			for i := 0; i < len(suffixWords); i++ {
				subWords := prefixWords[j] + suffixWords[i]
				if content == subWords {
					*count++
				} else if strings.HasPrefix(content, subWords) {
					tmp = append(tmp, subWords)
				}
			}
			if len(tmp) > 0 {
				compose(tmp, count)
			}
		}
	}
}

func TestFind(t *testing.T) {
	loadData()
	a := New()
	for _, c := range rawContents {
		a.Add(c)
	}
	a.Build()

	content = "lixian"
	words = a.Find(content)
	fmt.Println("主串: ", content, " 模式串: ", words)

	prefixWords := make([]string, 0)
	for j := 0; j < len(words); j++ {
		if strings.HasPrefix(content, words[j]) {
			prefixWords = append(prefixWords, words[j])
		}
	}
	var count int
	compose(prefixWords, &count)
	require.Equal(t, 2, count)
}

func TestPinYin(t *testing.T) {
	p, _ := NewPinYin().LoadData()
	res := p.Guess("lixian")
	fmt.Println(res)

}