package gotools

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
)

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

func ReadFileToSlices(filename string) (context []string, err error) {

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

func WriteFileFromText(filefullName, content string) error {
	//写文件,打开文件，如果不存在则创建
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
