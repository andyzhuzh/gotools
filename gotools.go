package gotools

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RegexpFindDataInString(`33333-user=ABCD -pwf23232-dd 33333-user=ABCD -pwf2d32d2a32-dd daa`, `-pwf23232-dd`, "23232")
func RegexpFindDataInString(allString, strTemplate, data string) (result []string) {
	firstIndex := strings.Index(strTemplate, data)
	lastIndex := firstIndex + len(data)
	if firstIndex == -1 || lastIndex > len(strTemplate) {
		return
	}
	regString := strTemplate[:firstIndex] + "(.*?)" + strTemplate[lastIndex:]
	// fmt.Println(regString)
	regExpObj, err := regexp.Compile(regString)
	if err == nil {
		allMatch := regExpObj.FindAllStringSubmatch(allString, -1)
		for _, match := range allMatch {
			if len(match) > 1 {
				result = append(result, match[1])
			}
		}
	}
	return
}

func StructToMap(content interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	marshalContent, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(bytes.NewReader(marshalContent))
	// d.UseNumber()  //float64 转换number
	if err1 := d.Decode(&result); err1 != nil {
		return nil, err1
	}
	// for k, v := range result {
	// 	// fmt.Println(k, v)
	// 	result[k] = v
	// }
	return result, nil
}

func FileReadToSlices(filename string) (context []string, err error) {

	// 文件存在
	// if _, err1 := os.Stat(filefullName); err1 == nil {
	// }
	//读取文件
	file, err := os.Open(filename) // 打开文件
	if err != nil {
		return
	}
	defer file.Close() // 确保文件在函数结束时关闭
	scanner := bufio.NewScanner(file)
	// lines := make([]string, 0) // 初始化切片
	for scanner.Scan() { // 逐行扫描
		line_text := scanner.Text()
		// line_texts := strings.Split(line_text, ",")
		context = append(context, line_text) // 将行添加到切片
	}
	return
}

func FileSaveText(filefullName, content string) error {
	//写文件,打开文件，如果不存在则创建
	var rwMutext sync.RWMutex
	rwMutext.Lock()
	defer rwMutext.Unlock()
	dirPath := filepath.Dir(filefullName)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
		// log.Fatalf("无法创建目录: %v", err)
	}
	f, err := os.OpenFile(filefullName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	byteData := []byte(content)
	_, err = f.Write(byteData)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func FileSaveCSV(filename string, content []map[string]interface{}) error {
	var rwMutext sync.RWMutex
	rwMutext.Lock()
	defer rwMutext.Unlock()
	// 创建CSV文件
	// file, err := os.Create(filename)
	dirPath := filepath.Dir(filename)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
		// log.Fatalf("无法创建目录: %v", err)
	}
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return err
	}
	defer file.Close()
	// 创建CSV写入器
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入头数据
	var headerString []string
	if len(content) > 0 {
		firstMap := content[0]
		mapKeys := reflect.ValueOf(firstMap).MapKeys()
		keyStrings := make([]string, len(mapKeys))
		for i, key := range mapKeys {
			keyStrings[i] = key.String()
			headerString = append(headerString, key.String())
		}
	}
	err = writer.Write(headerString)
	if err != nil {
		fmt.Println("Error writing to file", err)
		return err
	}

	// 写入更多数据行,重复
	for _, row := range content {
		var rowstring []string
		for _, field := range row {
			switch value := field.(type) {
			case string:
				rowstring = append(rowstring, value)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				// str := string(value)
				str := fmt.Sprintf("%d", value)
				rowstring = append(rowstring, str)
				// str := strconv.Itoa(value)
				// strconv.FormatFloat()
			case float64, float32:
				str := fmt.Sprintf("%.5fd", value)
				rowstring = append(rowstring, str)
			case bool:
				str := strconv.FormatBool(value)
				rowstring = append(rowstring, str)
			case time.Duration:
				str := value.String()
				rowstring = append(rowstring, str)
			case time.Time:
				str := value.String()
				rowstring = append(rowstring, str)
			default:
				typeName := reflect.TypeOf(field).String()
				str := "不支持类型：" + typeName
				rowstring = append(rowstring, str)
			}

		}
		err = writer.Write(rowstring)
		if err != nil {
			// fmt.Println("Error writing to file", err)
			return err
		}
	}
	// fmt.Println("CSV file created!")
	return nil
}
