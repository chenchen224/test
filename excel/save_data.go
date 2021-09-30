package excel

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/tealeg/xlsx"
)

type Data struct {
	SearchedWord string  `json:"搜索词"`
	QPS          float64 `json:"qps"`
	CompanyName  string  `json:"公司名称"`
	Count        int     `json:"count"`
	Top1         string  `json:"top_1"`
	Top2         string  `json:"top_2"`
	Top3         string  `json:"top_3"`
	Top4         string  `json:"top_4"`
	Top5         string  `json:"top_5"`
	Top6         string  `json:"top_6"`
	Top7         string  `json:"top_7"`
	Top8         string  `json:"top_8"`
	Top9         string  `json:"top_9"`
	Top10        string  `json:"top_10"`
}

func (data *Data) TransferSliceToTop10(arr []string) {
	if len(arr) > 0 {
		value := reflect.ValueOf(data)
		for i := range arr {
			elem := value.Elem()
			name := elem.FieldByName("Top" + strconv.Itoa(i+1))
			*(*string)(unsafe.Pointer(name.Addr().Pointer())) = arr[i]
		}
	}
}

func SaveDataToExcel(xlsxFile *xlsx.File, sheetName string, secondDataList []Data) error {

	println("+++++++++ Start save data +++++++++")

	var sheet *xlsx.Sheet
	sheet, err := xlsxFile.AddSheet(sheetName)
	if err != nil {
		return err
	}
	AddColumnTitle(Data{}, sheet)
	if len(secondDataList) == 0 {
		return nil
	}
	for i := range secondDataList {
		data := secondDataList[i]
		row := sheet.AddRow()
		TransferDataToRow(data, row)
	}

	return nil
}

func AddColumnTitle(structModel interface{}, sheet *xlsx.Sheet) {
	objValue := reflect.ValueOf(structModel)
	// objType:=reflect.TypeOf(structModel)
	row := sheet.AddRow()
	for i := 0; i < objValue.NumField(); i++ {
		fieldInfo := objValue.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("json")
		cell := row.AddCell()
		// str := TransferInterfaceToString(objType.Field(i), name)
		cell.SetString(name)
	}
}

func TransferDataToRow(data Data, row *xlsx.Row) {
	objValue := reflect.ValueOf(data)
	objType := reflect.TypeOf(data)
	var originalData = TransferInterfaceToString(objType.Field(2), objValue.Field(2).Interface())
	for i := 0; i < objValue.NumField(); i++ {
		cell := row.AddCell()
		str := TransferInterfaceToString(objType.Field(i), objValue.Field(i).Interface())
		fmt.Println("originalData: ", originalData)
		if originalData == str {
			println(str, " ：is right company in index ", i)
			style := xlsx.NewStyle()
			alignment := xlsx.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			}
			style.Alignment = alignment
			font := xlsx.Font{
				Size:      12,
				Name:      "宋体",
				Family:    0,
				Charset:   0,
				Color:     "#A377D2",
				Bold:      true,
				Italic:    false,
				Underline: false,
			}
			fill := xlsx.NewFill("", "#000000", "#C2C5aa")

			style.Font = font
			style.Fill = *fill
			style.ApplyAlignment = true
			style.ApplyFill = true
			style.ApplyFont = true
			cell.SetStyle(style)
		}
		cell.SetString(str)

	}
}

// PathExists 文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// DeleteFileIfExist 删除文件
func DeleteFileIfExist(path string) error {
	flag, err := PathExists(path)
	if err != nil {
		return err
	}
	if flag {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}
func TransferInterfaceToString(structField reflect.StructField, data interface{}) string {
	var float32Ptr *float32
	var float64Ptr *float64
	var intPtr *int
	var int8Ptr *int8
	var int16Ptr *int16
	var int32Ptr *int32
	var int64Ptr *int64
	var stringPtr *string

	switch structField.Type {

	case reflect.TypeOf(float64Ptr):
		d := data.(*float64)
		if d == nil {
			return ""
		} else {
			return strconv.FormatFloat(*d, 'f', 6, 64)
		}
	case reflect.TypeOf(float32Ptr):
		d := data.(*float32)
		if d == nil {
			return ""
		} else {
			return strconv.FormatFloat(float64(*d), 'f', 6, 32)
		}
	case reflect.TypeOf(float32(1)):
		return strconv.FormatFloat(float64(data.(float32)), 'f', 6, 32)
	case reflect.TypeOf(float64(1)):
		return strconv.FormatFloat(data.(float64), 'f', 6, 64)
	case reflect.TypeOf(intPtr):
		d := data.(*int)
		if d == nil {
			return ""
		} else {
			return strconv.Itoa(*d)
		}
	case reflect.TypeOf(int8Ptr):
		d := data.(*int)
		if d == nil {
			return ""
		} else {
			return strconv.Itoa(*d)
		}
	case reflect.TypeOf(int16Ptr):
		d := data.(*int)
		if d == nil {
			return ""
		} else {
			return strconv.Itoa(*d)
		}
	case reflect.TypeOf(int32Ptr):
		d := data.(*int)
		if d == nil {
			return ""
		} else {
			return strconv.Itoa(*d)
		}
	case reflect.TypeOf(int64Ptr):
		d := data.(*int)
		if d == nil {
			return ""
		} else {
			return strconv.Itoa(*d)
		}
	case reflect.TypeOf(1):
		return strconv.Itoa(data.(int))
	case reflect.TypeOf(int8(1)):
		return strconv.Itoa(int(data.(int8)))
	case reflect.TypeOf(int16(1)):
		return strconv.Itoa(int(data.(int16)))
	case reflect.TypeOf(int32(1)):
		return strconv.Itoa(int(data.(int32)))
	case reflect.TypeOf(int64(1)):
		return strconv.Itoa(int(data.(int64)))
	case reflect.TypeOf(stringPtr):
		d := data.(*string)
		if d == nil {
			return ""
		} else {
			return *d
		}
	case reflect.TypeOf(""):
		return data.(string)
	default:
		return ""
	}
}
