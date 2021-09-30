package app

import (
	// "go/printer"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/mozillazg/go-pinyin"
	"github.com/tealeg/xlsx"
	"gitlab.chenxk.com/test/es"
	"gitlab.chenxk.com/test/excel"
	"gitlab.chenxk.com/test/model"
)

var a model.MigrateConfig

func init() {
	a.Read("qksh_config", "./configs")
}

const times = 10

var names = []string{"腾讯",
	"华为",
	"菜鸟网络",
	"华大智造",
	"大疆创新",
	"比亚迪",
	"微众银行",
	"优必选",
	"柔宇科技",
	"跨越速运",
	"货拉拉",
	"喜茶",
	"土巴兔",
	"丰巢",
	"六度人和",
	"百果园",
	"行云创新",
	"全棉时代",
	"KK集团客",
	"路旅行",
	"随手记",
	"联易融",
	"掌门1对1",
	"奥比中光",
	"酷开网络",
	"碳云智能",
	"辣妈帮",
	"速腾聚创微",
	"众税银",
	"云天励飞",
	"大地影院",
	"一加手机",
	"中国平安",
	"绿米联创",
	"智布互联",
	"开思时代",
	"编程猫",
	"Ucar",
	"富途证券",
	"房多多",
	"光鉴科技",
	"惠牛科技",
	"深慧视",
	"融安网络",
	"视见科技",
	"刷脸呗",
	"倍电科技",
	"迅雷",
	"元戎启行",
	"极光大数据",
	"代来代去",
	"小鹅通",
	"金蝶软件",
	"未知君",
	"奈雪の茶",
	"昂纳信息",
	"海普洛斯",
	"普渡科技",
	"思谋科技",
	"帧观德芯",
	"装速配科技",
	"Nreal",
	"竹云科技",
	"大心电子",
	"Airwallex空中云汇",
	"地上铁",
	"舒可士科技",
	"捷顺科技",
	"芯能半导体",
	"英联利农",
	"优地科技",
	"唯酷光电",
	"云工网络科技",
	"职问",
	"小熊U租",
	"行云集团",
	"美联英文",
	"云英谷科技",
	"e换电",
	"晓微科技",
	"深圳高灯",
	"铁汉生态",
	"唯童教育",
	"杉岩",
	"蘑菇云",
	"极视角",
	"瑞云科技",
	"贝特瑞集团",
	"有方科技",
	"赢合科技",
	"岚锋创视",
	"唯童文化产业",
	"易库易供应链",
	"麦克韦尔科技",
	"皇庭国际",
	"金斧子",
	"华鹏飞",
	"欢创科技",
	"翰宇药业",
	"科陆电子"}

var companeNames = []string{"深圳市腾讯计算机系统有限公司",
	"华为技术有限公司",
	"菜鸟网络科技有限公司",
	"深圳华大智造科技股份有限公司",
	"深圳市大疆创新科技有限公司",
	"比亚迪股份有限公司",
	"深圳前海微众银行股份有限公司",
	"深圳市优必选科技股份有限公司",
	"深圳市柔宇科技股份有限公司",
	"跨越速运集团有限公司",
	"深圳货拉拉科技有限公司",
	"深圳美西西餐饮管理有限公司",
	"土巴兔集团股份有限公司",
	"深圳市丰巢科技有限公司",
	"深圳市六度人和科技有限公司",
	"深圳百果园实业（集团）股份有限公司",
	"深圳行云创新科技有限公司",
	"深圳全棉时代科技有限公司",
	"广东快客电子商务有限公司",
	"深圳市客路网络科技有限公司",
	"深圳市随手科技有限公司",
	"联易融数字科技集团有限公司",
	"深圳掌门人教育咨询有限公司",
	"奥比中光科技集团股份有限公司",
	"深圳市酷开网络科技股份有限公司",
	"深圳碳云智能科技有限公司",
	"深圳市辣妈帮科技有限公司",
	"深圳市速腾聚创科技有限公司",
	"深圳微众信用科技股份有限公司",
	"深圳云天励飞技术股份有限公司",
	"广东大地影院建设有限公司",
	"深圳市万普拉斯科技有限公司",
	"中国平安保险（集团）股份有限公司",
	"深圳绿米联创科技有限公司",
	"深圳市智能制造软件开发有限公司",
	"深圳开思时代科技有限公司",
	"深圳点猫科技有限公司",
	"深圳联合能源控股有限公司",
	"深圳市富途网络科技有限公司",
	"深圳市房多多网络科技有限公司",
	"深圳市光鉴科技有限公司",
	"深圳惠牛科技有限公司",
	"深慧视（深圳）科技有限公司",
	"深圳融安网络科技有限公司",
	"深圳视见医疗科技有限公司",
	"汇联易家互联网金融服务（深圳）有限公司",
	"深圳市倍电科技有限公司",
	"深圳市迅雷网络技术有限公司",
	"深圳元戎启行科技有限公司",
	"深圳市和讯华谷信息技术有限公司",
	"深圳市代来代去网络科技有限公司",
	"深圳小鹅网络技术有限公司",
	"金蝶软件(中国）有限公司",
	"深圳未知君生物科技有限公司",
	"深圳市品道餐饮管理有限公司",
	"昂纳信息技术(深圳)有限公司",
	"深圳市海普洛斯生物科技有限公司",
	"深圳市普渡科技有限公司",
	"深圳思谋信息科技有限公司",
	"深圳帧观德芯科技有限公司",
	"深圳装速配科技有限公司",
	"深圳太若科技有限公司",
	"深圳竹云科技有限公司",
	"深圳大心电子科技有限公司",
	"空中云汇（深圳）网络科技有限公司",
	"地上铁租车（深圳）有限公司",
	"深圳素士科技股份有限公司",
	"深圳市捷顺科技实业股份有限公司",
	"深圳芯能半导体技术有限公司",
	"深圳英联利农生物科技股份有限公司",
	"深圳优地科技有限公司",
	"深圳市唯酷光电有限公司",
	"深圳云工网络科技有限公司",
	"深圳市凯为咨询有限公司",
	"深圳市凌雄租赁服务有限公司",
	"深圳市天行云供应链有限公司",
	"深圳市美联国际教育有限公司",
	"深圳云英谷科技有限公司",
	"深圳易马达科技有限公司",
	"深圳市晓微科技有限公司",
	"深圳高灯计算机科技有限公司",
	"中节能铁汉生态环境股份有限公司",
	"深圳唯童教育有限公司",
	"深圳市杉岩数据技术有限公司",
	"深圳市蘑菇财富技术有限公司",
	"深圳极视角科技有限公司",
	"深圳市瑞云科技有限公司",
	"贝特瑞新材料集团股份有限公司",
	"深圳市有方科技股份有限公司",
	"深圳市赢合科技股份有限公司",
	"影石创新科技股份有限公司",
	"深圳唯童文化产业有限公司",
	"前海深蕾科技集团（深圳）有限公司",
	"深圳麦克韦尔科技有限公司",
	"深圳市皇庭国际企业股份有限公司",
	"深圳市金斧子网络科技有限公司",
	"华鹏飞股份有限公司",
	"深圳市欢创科技有限公司",
	"深圳翰宇药业股份有限公司",
	"深圳市科陆电子科技股份有限公司"}
var letters []string

func FillArgs(args []string) {

	for i := range names {
		r := []rune(names[i])
		l := ""
		for index := range r {
			a := string(r[index])
			if strings.Trim(a, " ") == "" || len([]byte(a))<3|| a=="の"{
				l += string(a)
				continue
			}
			println(a)
			str := pinyin.LazyConvert(a, nil)[0]
			l += str
		}
		letters = append(letters, l)
	}


	client, err := es.InitElasticsearch(a.ESConfig)
	if err != nil {
		panic(err)
	}
	SavePath := "./test.xlsx"
	newFile := xlsx.NewFile()
	// testQpsAndSave(client, args, productHunPinYin(names), newFile, "混拼")

	testQpsAndSave(client, args, names, newFile, "汉字")
	// testQpsAndSave(client, args, letters, newFile, "拼音")

	err = excel.DeleteFileIfExist(SavePath)
	err = newFile.Save(SavePath)
	if err != nil {
		panic(err)
	}
	println("+++++++++ End save data +++++++++")

}

func testQpsAndSave(client *elasticsearch.Client, args, keywords []string, file *xlsx.File, sheetName string) {
	datas := []excel.Data{}

	var records []string
	var count int
	for i := range keywords {

		start := time.Now()
		for j := 0; j < times; j++ {
			//xx
			records, count = es.PerformESQuery(a, client, args[0], names[i])
			log.Println(i)
		}
		searchEnd := time.Now()
		duration := searchEnd.Sub(start).Seconds()
		qps := times / duration
		log.Println("qps ===> ", qps)
		data := excel.Data{
			SearchedWord: keywords[i],
			CompanyName:  companeNames[i],
			Count:        count,
			QPS:          qps,
		}

		// TODO:
		items := records[0:]
		if len(records) > 10 {
			items = records[0:10]
		}
		data.TransferSliceToTop10(items)
		datas = append(datas, data)
	}
	excel.SaveDataToExcel(file, sheetName, datas)
}

func productHunPinYin(names []string) []string {
	var res []string
	for i := range names {
		r := []rune(names[i])
		l := ""
		for index := range r {
			a := string(r[index])
			if strings.Trim(a, " ") == "" || len([]byte(a))<3|| a=="の"{
				l += string(a)
				continue
			}
			if index%2 == 0 {
				println(a)
				str := pinyin.LazyConvert(a, nil)[0]
				l += str
				continue
			}
			l += string(a)
		}
		res = append(res, l)
	}
	return res
}
