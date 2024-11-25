package excel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelHeader struct {
	ColumnName  string
	ColumnLable string
	ColumnPos   int
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

func ExportExcelByMap(data []map[string]interface{}, fileName, sheetName string, header []ExcelHeader) error {
	var columns = make(map[string]int)
	// var columnsLable = make(map[string]int)
	var sheetLable string
	sheetLable = "Sheet1"
	if sheetName != "" {
		sheetLable = sheetName
	}
	f := excelize.NewFile()
	defer f.Close()
	f.SetSheetName("Sheet1", sheetName)
	sheet, _ := f.NewSheet(sheetLable)
	if len(data) < 1 {
		return fmt.Errorf("没有数据")
	}

	// 抬头
	colInt := 0
	rowID := 1
	if len(header) == 0 {
		for key := range data[0] {
			colInt += 1
			cell := Int2Alpha(colInt) + fmt.Sprint(rowID)
			// println(cell, ":", key)
			f.SetCellValue(sheetLable, cell, key)
			HeaderStyle(f, sheetLable, cell)
			columns[key] = colInt
		}
	} else {
		colInt = len(header)
		for _, colValue := range header {
			if colValue.ColumnPos < 1 {
				colInt += 1
				colValue.ColumnPos = colInt
			}
			cell := Int2Alpha(colValue.ColumnPos) + fmt.Sprint(rowID)
			// println(cell, ":", key)
			f.SetCellValue(sheetLable, cell, colValue.ColumnLable)
			HeaderStyle(f, sheetLable, cell)
			columns[colValue.ColumnName] = colValue.ColumnPos
		}
	}

	// 数据
	for keyID, rowContent := range data {
		rowID = keyID + 2
		colInt = 0
		for field, cellCont := range rowContent {
			colInt = columns[field]
			cellID := Int2Alpha(colInt) + fmt.Sprint(rowID)
			WriteExcelCell(f, sheetLable, cellID, cellCont)
		}
	}
	f.SetActiveSheet(sheet)
	if err := f.SaveAs(fileName); err != nil {
		// println(err.Error())
		return err
	}
	return nil
}

func ExportExcel(imData interface{}, fileName, sheetName string, header []ExcelHeader) error {
	switch data1 := imData.(type) {
	case []interface{}:
		return ExportExcelByInterface(data1, fileName, sheetName, header)
	case []map[string]interface{}:
		return ExportExcelByMap(data1, fileName, sheetName, header)
	default:
		return fmt.Errorf("不支持类型:" + reflect.TypeOf(imData).String())
	}
}

func ExportExcelByInterface(data []interface{}, fileName, sheetName string, header []ExcelHeader) error {
	var columns = make(map[string]int)
	var sheetLable string
	sheetLable = "Sheet1"
	if sheetName != "" {
		sheetLable = sheetName
	}

	f := excelize.NewFile()
	defer f.Close()
	f.SetSheetName("Sheet1", sheetName)
	sheet, _ := f.NewSheet(sheetLable)
	if len(data) < 1 {
		return fmt.Errorf("没有数据")
	}

	// 抬头
	colInt := 0
	rowID := 1
	if len(header) == 0 {
		hddata := data[0]
		hdtype := reflect.TypeOf(hddata).String()
		if hdtype != "map[string]interface {}" {
			return fmt.Errorf("无法处理类型:%v", hdtype)
		}
		switch hddata1 := hddata.(type) {
		case map[string]interface{}:
			hdMap := hddata1 //hddata.(map[string]interface{})
			for key := range hdMap {
				colInt += 1
				cell := Int2Alpha(colInt) + fmt.Sprint(rowID)
				f.SetCellValue(sheetLable, cell, key)
				HeaderStyle(f, sheetLable, cell)
				columns[key] = colInt
			}
		default:
			return fmt.Errorf("无法处理数量类型:%v", hdtype)

		}

	} else {
		colInt = len(header)
		for _, colValue := range header {
			if colValue.ColumnPos < 1 {
				colInt += 1
				colValue.ColumnPos = colInt
			}
			cell := Int2Alpha(colValue.ColumnPos) + fmt.Sprint(rowID)
			// println(cell, ":", key)
			f.SetCellValue(sheetLable, cell, colValue.ColumnLable)
			HeaderStyle(f, sheetLable, cell)
			columns[colValue.ColumnName] = colValue.ColumnPos
		}
	}

	// 数据
	for keyID, rowContent := range data {
		rowID = keyID + 2
		colInt = 0
		switch rowValue := rowContent.(type) {
		case map[string]interface{}:
			rowMap := rowValue // rowContent.(map[string]interface{})
			for field, cellCont := range rowMap {
				// colInt += 1
				colInt = columns[field]
				cellID := Int2Alpha(colInt) + fmt.Sprint(rowID)
				WriteExcelCell(f, sheetLable, cellID, cellCont)
			}
		default:
			// ltype := fmt.Sprintf("%T", fval)
			cellID := Int2Alpha(colInt) + fmt.Sprint(rowID)
			f.SetCellValue(sheetLable, cellID, "不支持类型:"+reflect.TypeOf(rowContent).String())
		}

	}
	f.SetActiveSheet(sheet)
	if err := f.SaveAs(fileName); err != nil {
		// println(err.Error())
		return err
	}
	return nil
}

func HeaderStyle(f *excelize.File, sheetName string, cellID string) error {
	style := &excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"D9E1F3", "BFBFBF"}, // 红色背景
			// Color:   []string{"BFBFBF", "BFBFBF"}, // 红色背景
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold:  true,
			Color: "003296",
			// Size: 12,
			// Name:   "Arial",
			// Family: excelize.FontFamilyProportional,
		},
	}
	styleID, _ := f.NewStyle(style)
	f.SetCellStyle(sheetName, cellID, cellID, styleID)
	return nil
}

func WriteExcelCell(f *excelize.File, sheetName string, cellID string, cellContent interface{}) error {
	switch fval := cellContent.(type) {
	case string:
		// f.SetCellValue(sheetName, cellID, cellContent.(string))
		f.SetCellValue(sheetName, cellID, fval)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		f.SetCellValue(sheetName, cellID, fval)
	case float64, float32:
		f.SetCellValue(sheetName, cellID, fval)
	case bool:
		f.SetCellValue(sheetName, cellID, fval)
	case time.Duration:
		f.SetCellValue(sheetName, cellID, fval)
	case time.Time:
		f.SetCellValue(sheetName, cellID, fval)
	case []byte:
		f.SetCellValue(sheetName, cellID, fval)
	case []map[string]interface{}:
		cellValueString := "表记录数:" + strconv.Itoa(len(fval))
		f.SetCellValue(sheetName, cellID, cellValueString)
	case []interface{}:
		cellValueString := "表记录数:" + strconv.Itoa(len(fval))
		f.SetCellValue(sheetName, cellID, cellValueString)
	default:
		typeName := reflect.TypeOf(cellContent).String()
		// ltype := fmt.Sprintf("%T", fval)
		f.SetCellValue(sheetName, cellID, "不支持类型:"+typeName)
	}
	return nil
}

func FindColum(fieldName string, columns []ExcelHeader) (column ExcelHeader) {
	if len(columns) == 0 {
		return
	}
	for _, value := range columns {
		if value.ColumnName == fieldName {
			column = value
			return
		}
	}
	return
}

func Int2Alpha(Num int) string {
	var (
		Str  string = ""
		k    int
		temp []int //保存转化后每一位数据的值，然后通过索引的方式匹配A-Z
	)
	//用来匹配的字符A-Z
	Slice := []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O",
		"P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	if Num > 26 { //数据大于26需要进行拆分
		for {
			k = Num % 26 //从个位开始拆分，如果求余为0，说明末尾为26，也就是Z，如果是转化为26进制数，则末尾是可以为0的，这里必须为A-Z中的一个
			if k == 0 {
				temp = append(temp, 26)
				k = 26
			} else {
				temp = append(temp, k)
			}
			Num = (Num - k) / 26 //减去Num最后一位数的值，因为已经记录在temp中
			if Num <= 26 {       //小于等于26直接进行匹配，不需要进行数据拆分
				temp = append(temp, Num)
				break
			}
		}
	} else {
		return Slice[Num]
	}
	for _, value := range temp {
		Str = Slice[value] + Str //因为数据切分后存储顺序是反的，所以Str要放在后面
	}
	return Str
}

func ReadExcelXLSX(excelFileName string) [][]interface{} {
	xlsData := make([][]interface{}, 0)
	// 打开一个现有的Excel文件
	f, err := excelize.OpenFile(excelFileName)
	if err != nil {
		// fmt.Println(err)
		return xlsData
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 读取Sheet1的A1单元格的值
	// cell, err := f.GetCellValue("Sheet1", "A1")

	// 获取所有工作表名
	sheets := f.GetSheetList()
	for _, sheetName := range sheets {
		// 读取整个Sheet的数据
		rows, _ := f.GetRows(sheetName)
		for _, row := range rows {
			xlsRows := make([]interface{}, 0)
			for _, colCell := range row {
				xlsRows = append(xlsRows, colCell)
			}
			xlsData = append(xlsData, xlsRows)
		}
	}
	return xlsData
}
