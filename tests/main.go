package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type data struct {
	Text string `json:"text,omitempty"`
}

const s = `^singleWord(|(\nsingleWord)+)$`

func main() {
	s := "2,000\n" + "0.56 5.36\n" + "6,000"
	singleWord := `((\d+)|(\-\d+)|(\/\d+))(|\,\d+)(|\.\d+)(|\%)`
	// reg := `^` + singleWord + `(|\s+(` + singleWord + `))(|(\n` + singleWord + `(|\s+(` + singleWord + `)))+)$`
	reg := fmt.Sprintf(`^%s(|\s+(%s))(|(\n%s(|\s+(%s)))+)$`, singleWord, singleWord, singleWord, singleWord)
	flag, err := regexp.MatchString(reg, strings.TrimSpace(s))
	log.Println(err)
	log.Println(flag)
}
