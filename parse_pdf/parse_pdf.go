package parse_pdf

import (
	// "bytes"
	"context"
	"io"
	"log"
	"os"
	"strings"

	// "time"

	// "strings"
	// "time"

	// "os"
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/minio/minio-go/v7"

	// "github.com/minio/minio-go/v7/pkg/credentials"

	// "encoding/json"
	// "net/http"

	// "github.com/PuerkitoBio/goquery"
	dcu "github.com/dcu/pdf"
	"github.com/google/go-tika/tika"
	ledongthuc "github.com/ledongthuc/pdf"
)

// func main() {
// 	// ctx := context.Background()
// 	endpoint := "localhost:9000"
// 	accessKeyID := "root"
// 	secretAccessKey := "password"
// 	useSSL := false

// 	// Initialize minio client object.
// 	minioClient, err := minio.New(endpoint, &minio.Options{
// 		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
// 		Secure: useSSL,
// 	})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// upload(minioClient, context.Background())
// 	// log.Println(object)
// 	objectNames := listObject(minioClient)
// 	objects, err := getObject(minioClient, objectNames)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// b := bytes.NewReader(objectInfos[0].TransLateToBytes())

// 	// var buffer bytes.Buffer
// 	// if _, err = io.Copy(&buffer, object); err != nil {
// 	// 	panic(err)
// 	// }
// 	// f := bytes.NewReader(buffer.Bytes())

// 	startTime := time.Now()
// 	bodies := readPdfByTika(objects)
// 	endTime := time.Now()
// 	tikaTime := endTime.Sub(startTime).Seconds()
// 	log.Println("tika time:", tikaTime)

// 	requests := make(map[string][]string)
// 	for name, body := range bodies {
// 		strs := strings.Split(body, "。")
// 		requests[name] = strs
// 	}

// 	// log.Println("=====================")

// 	// startTime = time.Now()
// 	// readAllPdfByLedongthuc()
// 	// endTime = time.Now()
// 	// ledongthucTime := endTime.Sub(startTime).Seconds()
// 	// log.Println("ledongthuc time:", ledongthucTime)

// 	// // log.Println("======================")
// 	// var buf bytes.Buffer
// 	// buf.ReadFrom(object)
// 	// log.Println(buf.String())

// 	// readPdfByLedongthuc()
// }

func GetObject(minioClient *minio.Client, names []string) (map[string]*minio.Object, error) {
	objects := make(map[string]*minio.Object)
	for _, name := range names {
		object, err := minioClient.GetObject(context.Background(), "anxin", name, minio.GetObjectOptions{})
		if err != nil {
			// TODO:
			log.Println(err)
			continue
		}
		objects[name] = object
	}

	return objects, nil
}

func ListObject(minioClient *minio.Client) []string {
	names := make([]string, 0)
	chanObjects := minioClient.ListObjects(context.Background(), "anxin", minio.ListObjectsOptions{})
	for object := range chanObjects {
		if object.Err != nil {
			log.Println(object.Err)
			return names
		}
		names = append(names, object.Key)
	}
	// a := <-chanObjects
	// info := model.MyObjectInfo{
	// 	ObjectInfo: a,
	// }
	// log.Println("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	// log.Println("original info:", a.Owner)
	// bytes := info.TransLateToBytes()
	// log.Println("info bytes:", bytes)
	// byteToStruct := *(*model.MyObjectInfo)(unsafe.Pointer(&bytes[0]))
	// log.Println("final info:", byteToStruct.Owner)
	// log.Println("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	return names
}

func ReadPdfByTika(f map[string]*minio.Object) map[string]string {
	client := tika.NewClient(nil, "http://localhost:9998")
	results := make(map[string]string)
	for key, object := range f {
		body, err := client.Parse(context.Background(), object)
		if err != nil {
			panic(err)
		}

		result := takeText(body)

		results[key] = result
	}
	return results
}

func takeText(body string) string {
	p := strings.NewReader(body)
	doc, _ := goquery.NewDocumentFromReader(p)

	doc.Find("script").Each(func(i int, el *goquery.Selection) {
		el.Remove()
	})

	result := strings.ReplaceAll(doc.Text(), " ", "")
	result = strings.ReplaceAll(result, "\n", "")

	return result
}

func ReadPdfByLedongthuc() error {
	path := "resource/anxin.pdf"
	f, r, err := ledongthuc.Open(path)
	if err != nil {
		return err
	}

	// remember close file
	defer f.Close()
	// totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= 1; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		// value := p.V.Text()
		// var lastTextStyle pdf.Text
		// texts := p.Content().Text
		// for _, text := range texts {
		// 	log.Printf("Font: %s, Font-size: %f, x: %f, y: %f, content: %s \n", text.Font, text.FontSize, text.X, text.Y, text.S)
		// 	// if isSameSentence(text, lastTextStyle) {
		// 	// 	lastTextStyle.S = lastTextStyle.S + text.S
		// 	// } else {
		// 	// 	log.Printf("Font: %s, Font-size: %f, x: %f, y: %f, content: %s \n", lastTextStyle.Font, lastTextStyle.FontSize, lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
		// 	// 	lastTextStyle = text
		// 	// }
		// 	log.Println("=====================================")
		// }
		// log.Println(p.Fonts())
		// for _, name := range p.Fonts() {
		// 	log.Println(name,  ":" , p.Font(name).BaseFont())
		// }
		fontMap := make(map[string]*ledongthuc.Font)
		font := p.Font("F5")
		fontMap["F1"] = &font
		res, err := p.GetPlainText(fontMap)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(res)
		log.Println("+++++++++++++++++++")
	}
	return nil
}

func ReadAllPdfByLedongthuc() {
	path := "resource/anxin.pdf"
	f, r, err := ledongthuc.Open(path)
	if err != nil {
		panic(err)
	}
	// remember close file
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()

	if err != nil {
		panic(err)
	}
	localFile, err := os.Create("resource/out/anxinPdf.txt")
	if err != nil {
		log.Println(err)
		return
	}
	if _, err = io.Copy(localFile, b); err != nil {
		log.Println(err)
		return
	}

	buf.ReadFrom(b)
	text := buf.String()
	log.Println(text)

	

	// TODO: split text per 400 words

	// data := make(map[string]interface{})
	// data["text"] = text
	// jsonData, _ := json.Marshal(data)
	// body := bytes.NewBuffer([]byte(jsonData))
	// resp, err := http.Post("http://192.168.88.201:8485/simple_task", "application/json;charset=utf-8", body)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(resp)
}

func readPdfByDcu(path string) (string, error) {
	f, r, err := dcu.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		p.GetPlainText(nil)
		for _, row := range rows {
			println(">>>> row: ", row.Position)
			for _, word := range row.Content {
				log.Println(word.S)
			}
		}
	}
	return "", nil
}

func upload(minioClient *minio.Client, ctx context.Context) {
	// Make a new bucket called mymusic.
	bucketName := "anxin"
	location := "us-east-1"
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Upload the zip file
	objectName := "hexin.pdf"
	filePath := "hexin.pdf"
	contentType := "application/pdf"

	// Upload the zip file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}

// objectsCh := make(chan minio.ObjectInfo)

// // Send object names that are needed to be removed to objectsCh
// go func() {
//     defer close(objectsCh)
//     // List all objects from a bucket-name with a matching prefix.
//     for object := range minioClient.ListObjects(context.Background(), "my-bucketname", "my-prefixname", true, nil) {
//         if object.Err != nil {
//             log.Fatalln(object.Err)
//         }
//         objectsCh <- object.Key
//     }
// }()

// opts := minio.RemoveObjectsOptions{
//     GovernanceBypass: true,
// }

// for rErr := range minioClient.RemoveObjects(context.Background(), "my-bucketname", objectsCh, opts) {
//     fmt.Println("Error detected during deletion: ", rErr)
// }
