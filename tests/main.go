package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type data struct {
	Text string `json:"text,omitempty"`
}

func main() {
	s := "<p>行业动态分析/电力及公用事业 " +
		"</p>" +
		"<p>3 本报告版权属于安信证券股份有限公司。" +
		"</p>" +
		"<p>各项声明请参见报告尾页。 " +
		"</p>"
	str := strings.NewReader(s)
	doc, err := goquery.NewDocumentFromReader(str)
	if err != nil {
		panic(err)
	}
	nodes := doc.Find("p").Nodes
	for _, node := range nodes {
		log.Println(node.Data)
	}
}
