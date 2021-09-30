package es

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"gitlab.chenxk.com/test/model"

	// "gitlab.mvalley.com/adam/common/pkg/errs"
	cfg "gitlab.mvalley.com/datapack/cain/pkg/config"
)

func ConstructCommonESQuery(configuration model.MigrateConfig, keyword string, from, size int) interface{} {
	boost := configuration.QuickSearchBoost
	esQuery := make(map[string]interface{}, 0)
	esQuery["from"] = from
	esQuery["size"] = size
	esQuery["query"] = map[string]interface{}{
		"function_score": map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						{
							"match": map[string]interface{}{
								"ik.primary_name": map[string]interface{}{
									"operator": "or",
									"query":    keyword,
									"boost":    boost.PrimaryName,
								},
							},
						},
						{
							"match": map[string]interface{}{
								"ik.keywords": map[string]interface{}{
									"operator":             "and",
									"query":                keyword,
									"boost":                boost.KeyWord,
									"minimum_should_match": "75%",
								},
							},
						},
						{
							"match_phrase": map[string]interface{}{
								"ik.description": map[string]interface{}{
									"query": keyword,
									"boost": boost.Description,
								},
							},
						},
					},
					"filter": map[string]interface{}{
						"match": map[string]interface{}{
							"ik.keywords": map[string]interface{}{
								"operator": "and",
								"query":    keyword,
							},
						},
					},
				}},
			"script_score": map[string]interface{}{
				"script": map[string]interface{}{
					"source": fmt.Sprintf("_score + doc['ranking_score'].value*%.2f", 0.0),
				},
			},
			"boost_mode": "sum",
		},
	}
	return esQuery
}

func constructWhxESQuery(configuration model.MigrateConfig, keyword string, from, size int) interface{} {
	// boost := configuration.QuickSearchBoost
	esQuery := make(map[string]interface{}, 0)
	esQuery["from"] = from
	esQuery["size"] = size
	esQuery["query"] = map[string]interface{}{
		"function_score": map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						{
							"match": map[string]interface{}{
								"ik.primary_name": map[string]interface{}{
									"operator": "and",
									"query":    keyword,
									"boost":    4,
								},
							},
						},
						{
							"match": map[string]interface{}{
								"ik.keywords": map[string]interface{}{
									"operator":             "and",
									"query":                keyword,
									"boost":                1,
								},
							},
						},
						{
							"match_phrase": map[string]interface{}{
								"ik.primary_name": map[string]interface{}{
									"query":    keyword,
									"analyzer": "ik_search_analyzer",
									"slop":     1,
									"boost":    8,
								},
							},
						},
					},
					"filter": map[string]interface{}{
						"match": map[string]interface{}{
							"ik.keywords": map[string]interface{}{
								"operator": "and",
								"query":    keyword,
							},
						},
					},
				}},
			"script_score": map[string]interface{}{
				"script": map[string]interface{}{
					"source": fmt.Sprintf("_score + doc['ranking_score'].value*%.2f", 0.0),
				},
			},
			"boost_mode": "sum",
		},
	}
	return esQuery
}

func GetQuery(from, size int) interface{} {
	// boost := configuration.QuickSearchBoost
	esQuery := make(map[string]interface{}, 0)
	esQuery["from"] = from
	esQuery["size"] = size
	esQuery["query"] = map[string]interface{}{
		"function_score": map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						{
							"match": map[string]interface{}{
								"ik.primary_name": map[string]interface{}{
									"operator": "and",
									"query":    "keyword",
									"boost":    4,
								},
							},
						},
						{
							"match": map[string]interface{}{
								"ik.keywords": map[string]interface{}{
									"operator":             "and",
									"query":                "keyword",
									"boost":                1,
								},
							},
						},
						{
							"match_phrase": map[string]interface{}{
								"ik.primary_name": map[string]interface{}{
									"query":    "keyword",
									"analyzer": "ik_search_analyzer",
									"slop":     1,
									"boost":    8,
								},
							},
						},
					},
					"filter": map[string]interface{}{
						"match": map[string]interface{}{
							"ik.keywords": map[string]interface{}{
								"operator": "and",
								"query":    "keyword",
							},
						},
					},
				}},
			"script_score": map[string]interface{}{
				"script": map[string]interface{}{
					"source": fmt.Sprintf("_score + doc['ranking_score'].value*%.2f", 0.0),
				},
			},
			"boost_mode": "sum",
		},
	}
	return esQuery
}

func constructXcESQuery(configuration model.MigrateConfig, keyword string, from, size int) interface{} {
	// boost := configuration.QuickSearchBoost
	esQuery := make(map[string]interface{}, 0)
	esQuery["from"] = from
	esQuery["size"] = size
	esQuery["query"] = map[string]interface{}{
		"function_score": map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": []map[string]interface{}{
						{
							"multi_match": map[string]interface{}{
								"query": keyword,
								"fields": []string{
									"ik.keywords.ik",
									"ik.keywords.py",
								},
								"operator":             "and",
								"minimum_should_match": "75%",
								"boost":                3,
							},
						},
						{
							"match_phrase": map[string]interface{}{
								"ik.primary_name.ik": map[string]interface{}{
									"query":    keyword,
									"boost":    3,
									"slop":     2,
									"analyzer": "ik_smart",
								},
							},
						},
						{
							"match_phrase": map[string]interface{}{
								"ik.primary_name.py": map[string]interface{}{
									"query":    keyword,
									"boost":    3,
									"slop":     2,
									"analyzer": "pinyin",
								},
							},
						},
						{
							"multi_match": map[string]interface{}{
								"query": keyword,
								"fields": []string{
									"ik.primary_name.ik",
									"ik.primary_name.py",
								},
								"operator": "and",
								"boost":    10,
							},
						},
					},
				}},
			"script_score": map[string]interface{}{
				"script": map[string]interface{}{
					"source": fmt.Sprintf("_score + doc['ranking_score'].value*%.2f", 0.0),
				},
			},
			"boost_mode": "sum",
		},
	}
	// esQuery["hightlight"] = map[string]interface{}{
	// 	"fields": map[string]interface{}{
	// 		"ik.keywords.ik": map[string]interface{}{},
	// 		"ik.keywords.py": map[string]interface{}{},
	// 		"ik.primary_name.ik": map[string]interface{}{},
	// 		"ik.primary_name.py": map[string]interface{}{},
	// 	  },
	// }
	return esQuery
}

func getQuery(confiuration model.MigrateConfig, index string, keyword string) interface{} {
	var query interface{}
	switch index {
	case "qksh_saic_prod_test3":
		query = constructXcESQuery(confiuration, keyword, 0, 10)
		return query
	case "qksh_saic_prod_test2":
		query = ConstructCommonESQuery(confiuration, keyword, 0, 10)
		return query
	case "qksh_saic_prod_test1":
		query = constructWhxESQuery(confiuration, keyword, 0, 10)
		return query
	default:
		return nil
	}
}

func PerformESQuery(confiuration model.MigrateConfig, client *elasticsearch.Client, index string, keyword string) ([]string, int) {
	reponses := make([]string, 0)

	query := getQuery(confiuration, index, keyword)
	jsonBody, _ := json.Marshal(query)
	body := bytes.NewReader(jsonBody)

	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(index),
		client.Search.WithBody(body),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
		client.Search.WithTimeout(10*time.Second),
		client.Search.WithSize(10000),
	)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	r := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		panic(err)
	}

	hits := r["hits"].(map[string]interface{})
	total := hits["total"].(map[string]interface{})

	for _, v := range hits["hits"].([]interface{}) {
		vm := v.(map[string]interface{}) //highlight and _source are both in hits
		// score := vm["_score"].(float64)
		_source := vm["_source"].(map[string]interface{})
		legal_name := _source["kw.legal_name"].(string)
		// person_name := _source["kw.legal_person_surface_name"].(string)
		// item := model.Response{
		// 	Score:       score,
		// 	PramaryName: legal_name,
		// 	PersonName:  person_name,
		// }
		reponses = append(reponses, legal_name)
	}

	count := int(total["value"].(float64))
	return reponses, count
}

func InitElasticsearch(config cfg.ESConfiguration) (*elasticsearch.Client, error) {
	var cfg = elasticsearch.Config{
		Addresses: []string{
			config.Host,
		},
		Username: config.User,
		Password: config.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Duration(config.ResponseHeaderTimeoutSeconds) * time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS11,
				InsecureSkipVerify: true,
			},
		},
	}

	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	_, err = esClient.Ping()
	if err != nil {
		panic(err)
	}

	return esClient, nil
}

func GetEs() (es *elasticsearch.Client) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://127.0.0.1:9200/",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic("Error init es")
	}
	return
}
