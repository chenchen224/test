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

const singleWord = `((\d+)|(\-\d+)|(\/\d+)|(\(\d+\)))(|\,\d+)(|\.\d+)(|\%)(|\s+)(|\/+)`
const regPictureOrForm = "^(\u8868|\u56fe|\u56fe\u8868)(|\\s+)\\d+\uff1a.*[^\\d+]$"

func preservation() {
	reg := fmt.Sprintf(`^%s(|\s+(%s))+(|(\n%s(|\s+(%s))+)+)$`, singleWord, singleWord, singleWord, singleWord)
	log.Println(reg)
}

func main() {

	s1 := "一、纺织服装板块"
	regTitle := "^([[:digit:]]|[\u4e00-\u5341]{0,})(|\\.|\u3001)(|\\s+).*[^\\d+]$"
	flag, err := regexp.MatchString(regTitle, strings.TrimSpace(s1))
	log.Println(err)
	log.Println(flag)

	regNum := "\u3001"
	s := "、"
	flag, err = regexp.MatchString(regNum, s)
	log.Println(err)
	log.Println(flag)
}
