package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	// "fmt"
	"log"
	"net/http"

	// "os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/minio/minio-go/v7"
	"golang.org/x/net/html"

	"github.com/google/go-tika/tika"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/yob/pdfreader/pdfread"
	pdf "gitlab.chenxk.com/test/parse_pdf"
)

const first_reg =`(^\d+(\.\d+)?\d$)|((\%$)|(\%\s$))|` + 
`(^\-\d+(\.\d+)?\d$)|(^\d+\,\d+)|` + 
`(^\-\d+(\.\d+)?\s+)|` +
`(^\d+(\.\d+)?\s+)`

type data struct {
	Text string `json:"text,omitempty"`
}

type Response struct {
}

var text string = "本报告版权属于安信证券股份有限公司。各项声明请参见报告尾页。（a+h）】。新能源运营板块推荐坐拥福建台海风资源，有望成长为国内海上风电龙头的【福能股份】，同时建议关注a股新能源运营龙头【三峡能源】。水电板块建议关注分红有承诺、集团机组陆续投产的【长江电力】以及分红比例大幅提升的【川投能源】。燃气板块推荐天然气一体化产业链稀缺标的【新奥股份】■风险提示：全社会用电量增长不及预期、煤价持续高位运行、电价调整不及预期、来水不及预期、新能源装机增速不及预期、天然气消费增速不及预期"

func main() {
	f, err := os.Open("resource/xinda.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	client := tika.NewClient(nil, "http://localhost:9998")

	body, err := client.MetaRecursiveType(context.Background(), f, "html")
	if err != nil {
		panic(err)
	}

	object := strings.NewReader(body[0]["X-TIKA:content"][0])
	doc, err := goquery.NewDocumentFromReader(object)
	if err != nil {
		panic(err)
	}
	// divNodes := doc.Find("div.page").Nodes
	// log.Println(divNodes[0].Data)

	// parseTextP()
	// eachNodes(doc)
	eachP(doc)
	// eachDiv(doc)
}

func eachNodes(doc *goquery.Document) {
	var strs string
	doc.Find("div.page").Each(func(i int, s *goquery.Selection) {
		str := strings.ReplaceAll(s.Find("p").Text(), " ", "")
		// str = strings.ReplaceAll(str, "\n", "")

		// log.Println("p:", str)
		// log.Println("===============================")
		strs += "\n\n" + str
	})

	localFile, err := os.Create("resource/out/xinda.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err = io.Copy(localFile, strings.NewReader(strs)); err != nil {
		fmt.Println(err)
		return
	}
}

func eachP(doc *goquery.Document) {
	// reg := `(^((\d)|(\-\d))$)|(^\d+\d+$)|(^\d+\.\d+$)|(^\-\d+$)|(\d+(\n+\d+)+$)|(^\d+\.\d+(\s+)(\n\d+\.\d+)+$)|` +
	// 	   `(\%$)|` + `(^\d+\.\d+(\s+)(\n\d+\.\d+((\s+)\d+\.\d+$))+$)|` + `(\d+(\n+\/\d+)+$)|(^\/\d+$)|` +
	// 	   `(^\-\d+(\n+((\d+)|(\-\d+)))+$)`
	// reg := `^((\d+)|(\-\d+)|(\/\d+))(|\.\d+)(|\%)(|\s+(((\d+)|(\-\d+)|(\/\d+))(|\.\d+)(|\%)))(|(\n((\d+)|(\-\d+)|(\/\d+))(|\.\d+)(|\%)(|\s+(((\d+)|(\-\d+)|(\/\d+))(|\.\d+)(|\%))))+)$`
	// singleWord := `((\d+)|(\-\d+)|(\/\d+))(|\,\d+)(|\.\d+)(|\%)`

	singleWord := `((\d+)|(\-\d+)|(\/\d+)|(\(\d+(|\.\d+)\)))(|\,\d+)(|\.\d+)(|\d+\%|\%)(|\s+)(|\/+)`
	dateWord := "[\u5e74\u6708\u65e5]"
	reg := fmt.Sprintf(`^%s(|%s|\s+(%s))+(|(\n(|%s|%s(|\s+(%s)))+)+)$`, singleWord, singleWord, singleWord, dateWord, singleWord, singleWord)
	// singleWord := "(\\d+|\\-d+\\/d+)(|\\.)\\d+(|\\%|\u5e74|\u6708|\u65e5)"
	// reg := fmt.Sprintf("^(%s)[%s]+(|%s)$", singleWord, singleWord, singleWord)
	match := regexp.MustCompile(reg)

	regPictureOrForm := regexp.MustCompile("^(\u8868|\u56fe|\u56fe\u8868)(|\\s+)(|\\d+)[\uff1a|\\:].*[\\D+]$")

	chinNumReg := "(\u4e00|\u4e8c|\u4e09|\u56db|\u4e94|\u516d|\u4e03|\u516b|\u4e5d|\u5341)+"
	chinPuctuation := "\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b"
	regOfTitle1Pattern := regexp.MustCompile(fmt.Sprintf("^([[:digit:]]|%s)+[\\.\u3001](|\\s+).*[\\D+]$", chinNumReg))
	regOfTitle2Pattern := regexp.MustCompile(fmt.Sprintf("^([[:digit:]]|%s)+[\\.\u3001](|\\s+).*[^%s]$", chinNumReg, chinPuctuation))

	lengthPattern := regexp.MustCompile("[^\u3002]$")

	linkPattern := regexp.MustCompile(`(http|ftp|https)\:\/\/(www|WWW|[A-Za-z]+)\.[A-Za-z]+\.[A-Za-z]+(|.*)$`)
	tablePattern := regexp.MustCompile(`^(\[Table)\_.*`)
	emailPattern := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)

	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		str := strings.ReplaceAll(s.Text(), " ", "")
		str = strings.ReplaceAll(str, "\n", "")

		if len([]rune(str)) <= 20 && lengthPattern.MatchString(str){
			s.Remove()
		} else if match.MatchString(str) {
			// 去除单行只有数字的段落
			// log.Println(str)
			s.Remove()
		} else if strings.Contains(str, "资料来源") {
			// 去除包含“资料来源”段落
			s.Remove()
		} else if regPictureOrForm.MatchString(str) {
			// 去除图表标签段落
			s.Remove()
		} else if regOfTitle1Pattern.MatchString(str) && regOfTitle2Pattern.MatchString(str) {
			// TODO:会删除少量有语意的句子
			// 去除标题
			// log.Println(str)
			s.Remove()
		} else if linkPattern.MatchString(str) {
			// 去除链接
			s.Remove()
		} else if tablePattern.MatchString(str) {
			// TODO:会删除少量有语意的句子
			// log.Println(str)
			s.Remove()
		} else if emailPattern.MatchString(str) {
			// 去除油箱通讯方式
			log.Println(str)
			s.Remove()
		} else if isItHeaderOrFooter(str) {
			// 去除页眉，页尾
			log.Println(str)
			s.Remove()
		}
	})

	eachNodes(doc)
}

func isItHeaderOrFooter(str string) bool {
	words := []string{"敬请参阅", "务必阅读", "各项声明", "特别声明", "免责声明", "免责条款", "执业证书编号", "联系邮箱"}
	for _, v := range words {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}

func eachDiv(doc *goquery.Document) {
	size := len(doc.Find("div.page").Nodes)
	doc.Find("div.page").Each(func(i int, s *goquery.Selection) {
		
		s.Find("p").Each(func(i int, s *goquery.Selection) {
			if i == 4 {
				log.Println(s.Text())
				log.Println("================================")
			}
		})
	})
	log.Println("size:", size)
}

func parseTextP() {
	f, err := os.Open("resource/xinda.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	client := tika.NewClient(nil, "http://localhost:9998")

	body, err := client.MetaRecursiveType(context.Background(), f, "html")
	if err != nil {
		panic(err)
	}

	object := strings.NewReader(body[0]["X-TIKA:content"][0])

	doc, err := goquery.NewDocumentFromReader(object)
	if err != nil {
		panic(err)
	}
	pNodes := doc.Find("p").Nodes
	str := ""
	strs := make([]string, 0)
	for _, v := range pNodes {
		if v.FirstChild == nil {
			str = ""
			str = strings.ReplaceAll(str, " ", "")
			strs = append(strs, str)
			continue
		}
		text := v.FirstChild.Data
		textWithoutSpace := strings.TrimSpace(text)
		if textWithoutSpace == "" {
			// log.Println("==========================")
			str = strings.ReplaceAll(str, " ", "")
			strs = append(strs, str)
			str = ""
		}
		str += text
		// log.Println(text)
	}

	log.Println(len(strs))

	// TODO:
	reg := regexp.MustCompile(`[0-9]\d*\.?\d*%?$`)
	reslut := make([]string, 0)
	for _, v := range strs {
		if strings.TrimSpace(v) == "" || reg.MatchString(v) {
			continue
		}
		reslut = append(reslut, v)
	}

	resultStr := ""
	for _, v := range reslut {
		// log.Println("len:", len(v))
		// log.Println("value:", v)
		resultStr += v + "\n\n"
		// log.Println("==================================")
	}

	localFile, err := os.Create("resource/out/merge.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err = io.Copy(localFile, strings.NewReader(resultStr)); err != nil {
		fmt.Println(err)
		return
	}
}

func clear_dom(divNopde *html.Node) error {
	// var err error
	for nd := divNopde.NextSibling; nd != nil; {
		switch nd.Type {
		case html.TextNode:
			log.Println(nd.LastChild.Data)

			if nd.FirstChild != nil && len([]rune(nd.FirstChild.Data)) <= 14 {
				// delete the element
				log.Println("fuck")
				tmp := nd
				nd = tmp.NextSibling
				divNopde.RemoveChild(tmp)
			}
		default:
			nd = nd.NextSibling
		}
	}

	return nil
}

func tikaMeta() {
	f, err := os.Open("resource/anxin.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	client := tika.NewClient(nil, "http://localhost:9998")

	body, err := client.MetaRecursiveType(context.Background(), f, "xml")
	if err != nil {
		panic(err)
	}
	// log.Println(body[0]["X-TIKA:content"][0])
	object := strings.NewReader(body[0]["X-TIKA:content"][0])
	localFile, err := os.Create("resource/out/anxin.xml")
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err = io.Copy(localFile, object); err != nil {
		fmt.Println(err)
		return
	}
}

func pdfRead() {
	pd := pdfread.Load(os.Args[1])
	if pd != nil {
		pg := pd.Pages()
		for k := range pg {
			fmt.Printf("Page %d - MediaBox: %s\n",
				k+1, pd.Att("/MediaBox", pg[k]))
			fonts := pd.PageFonts(pg[k])
			for l := range fonts {
				fontname := pd.Dic(fonts[l])["/BaseFont"]
				fmt.Printf("  %s = \"%s\"\n",
					l, fontname[1:])
			}
		}
	}
}

func allTextChoice() {
	endpoint := "localhost:9000"
	accessKeyID := "root"
	secretAccessKey := "password"
	useSSL := false
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	objectNames := pdf.ListObject(minioClient)
	objects, err := pdf.GetObject(minioClient, objectNames)
	if err != nil {
		panic(err)
	}
	startTime := time.Now()
	bodies := pdf.ReadPdfByTika(objects)
	endTime := time.Now()
	tikaTime := endTime.Sub(startTime).Seconds()
	log.Println("tika time:", tikaTime)

	requests := make(map[string][]string)
	for name, body := range bodies {
		strs := strings.Split(body, "。")
		requests[name] = strs
	}

	log.Println(len(requests["anxin.pdf"]))

	log.Println("================================================================")

	startTime = time.Now()
	keywordMap := make(map[string][]string)
	for _, value := range requests["anxin.pdf"] {
		s := data{
			Text: value,
		}
		jsonData, _ := json.Marshal(s)
		// data := make(map[string]interface{})
		// data["text"] = "年上半年公共事业板块整体保持稳健增长态势，虽然第二季度增幅较第一季度有所回落，但总体仍处于发展阶段"
		// jsonData, _ := json.Marshal(data)
		body := bytes.NewBuffer([]byte(jsonData))
		resp, err := http.Post("http://192.168.88.201:8485/simple_task", "application/json;charset=utf-8", body)
		if err != nil {
			panic(err)
		}
		log.Println("=======================")
		log.Println("resp:", resp)

		r := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&r)
		if err != nil {
			panic(err)
		}
		log.Println("body:", r["body"])

		var result map[string]interface{}
		if r["body"] != nil {
			json.Unmarshal([]byte(r["body"].(string)), &result)
		}
		// if err != nil {
		// 	panic(err)
		// }

		if result["labels"] == nil {
			continue
		}
		labels := result["labels"].([]interface{})
		text := result["text"].(string)
		e := []rune(text)
		keywords := make([]string, 0)
		log.Println("text:", text)
		for _, v := range labels {
			arr := v.([]interface{})
			startIndex := arr[0].(float64)
			endIndex := arr[1].(float64)
			item := e[int(startIndex):int(endIndex)]
			keywords = append(keywords, string(item))
		}

		log.Println("keywords:", keywords)
		log.Println("=========================")
		_, ok := keywordMap["anxin.pdf"]
		if ok {
			keywordMap["anxin.pdf"] = append(keywordMap["anxin.pdf"], keywords...)
			continue
		}
		keywordMap["anxin.pdf"] = keywords
	}

	duplicateRemovalMap := make(map[string]string)
	for _, v := range keywordMap["anxin.pdf"] {
		duplicateRemovalMap[v] = "anxin.pdf"
	}
	list := make([]string, 0)
	for k := range duplicateRemovalMap {
		list = append(list, k)
	}
	endTime = time.Now()
	keywordsTime := endTime.Sub(startTime).Seconds()

	log.Println("list:", list)
	log.Println("size:", len(list))
	log.Println("keywords time:", keywordsTime)
}
