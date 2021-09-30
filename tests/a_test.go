package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gitlab.chenxk.com/test/es"
	"gitlab.chenxk.com/test/model"
)

var a model.MigrateConfig

func init() {
	a.Read("config", "../configs")
	log.Println(a)
}

func BenchmarkSum(b *testing.B) {
	// client, err := es.InitElasticsearch(a.ESConfig)
	// if err != nil {
	// 	panic(err)
	// }
	client := es.GetEs()
	for i := 0; i < b.N; i++ {
		es.PerformESQuery(a, client, "fund_raising_events", "")
	}
}

func TestFileOpen(t *testing.T) {
	f, err := os.Open("file")
	if err != nil {
		fmt.Println("f: ", f)
		f.Close()
		log.Println(err)
	}
}
