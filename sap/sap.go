package sap

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/andyzhuzh/gotools/common/excel"
	"github.com/andyzhuzh/gotools/common/xml"
	"github.com/sap/gorfc/gorfc"
)

type ServerParameters struct {
	User         string
	Password     string
	HostName     string
	SystemNumber string
	Client       string
	Lang         string
	SapRouter    string
}

type SapServer struct {
	ServerName     string
	ServerInstance *gorfc.Connection
	ConnParams     ServerParameters
}

type FunctionParameters struct {
	Name          string
	ParameterType string
	Direction     string
	NucLength     uint
	UcLength      uint
	Decimals      uint
	DefaultValue  string
	ParameterText string
	Optional      bool
	TypeName      string
	// TypeDirection string
}

type FunctionParameterType struct {
	Name     string
	TypeName string
	fields   []FunctionParameterTypeField
}
type FunctionParameterTypeField struct {
	Name      string
	FieldType string
	NucLength uint
	NucOffset uint
	UcLength  uint
	UcOffset  uint
	Decimals  uint
	TypeName  string
}

type TableStructure struct {
	TableName  string
	FieldName  string
	Position   int
	OffSet     int
	RollName   string
	Leng       int
	Decimals   int
	DataType   string
	FieldText  string
	FieldTitle string
}

// 获取RFC结果所有参数名称及类别
func GetObjectListFromResult(rfcResultMap map[string]interface{}) map[string]string {
	var resultList = make(map[string]string)
	for key, content := range rfcResultMap {
		resultList[key] = fmt.Sprint(reflect.TypeOf(content))
	}
	return resultList
}

// 获取RFC结果中指定Table内容
func GetTableFromResult(tableName string, rfcResultMap map[string]interface{}) ([]map[string]interface{}, error) {
	var resultTable []map[string]interface{}
	var err error
	resultM, OK := rfcResultMap[tableName]
	if !OK {
		// err = fmt.Errorf("没有返回参数：%s", tableName)
		return resultTable, fmt.Errorf("没有返回参数：%s", tableName)
	}
	switch tabCnt := resultM.(type) {

	case []interface{}: //表  工作区：map[string]interface{}:  其他单值
		for _, value := range tabCnt {
			switch rowMap := value.(type) {
			case map[string]interface{}:
				resultTable = append(resultTable, rowMap)
			default:
				err = fmt.Errorf(tableName, "中行类型是：%s", fmt.Sprint(reflect.TypeOf(value)))
			}
		}
	default:
		err = fmt.Errorf(tableName, "类型是：%s", fmt.Sprint(reflect.TypeOf(resultM)))
	}
	return resultTable, err
}

// 获取RFC结果中指定工作区
func GetWorkAreaFromResult(paraName string, rfcResultMap map[string]interface{}) (map[string]interface{}, error) {
	// var resultMap map[string]interface{}
	var err error
	var resultMap = make(map[string]interface{})

	resultM, OK := rfcResultMap[paraName]
	if !OK {
		return nil, fmt.Errorf("没有返回参数：%s", paraName)
	}
	switch tabCnt := resultM.(type) {
	case map[string]interface{}:
		resultMap = tabCnt
	default:
		err = fmt.Errorf(paraName, "类型是：%s", fmt.Sprint(reflect.TypeOf(resultM)))
	}
	return resultMap, err

}

// 获取RFC结果中指定单值
func GetSingleValueFromResult(paraName string, rfcResultMap map[string]interface{}) (string, error) {
	var err error
	var resultString string
	resultM, OK := rfcResultMap[paraName]

	if !OK {
		return "", fmt.Errorf("没有返回参数：%s", paraName)
	}
	switch tabCnt := resultM.(type) {
	case map[string]interface{}:
		err = fmt.Errorf(paraName, "类型是：%s", fmt.Sprint(reflect.TypeOf(resultM)))
		resultString = ""
	case []interface{}:
		err = fmt.Errorf(paraName, "类型是：%s", fmt.Sprint(reflect.TypeOf(resultM)))
		resultString = ""
	default:
		resultString = fmt.Sprint(tabCnt)
	}
	return resultString, err
}

func PadLeft(str string, width int) string {
	if len(str) >= width {
		return str
	}
	return strings.Repeat(" ", width-len(str)) + str
}

func PadRight(str string, width int) string {
	if len(str) >= width {
		return str
	}
	return str + strings.Repeat(" ", width-len(str))
}

func (server *SapServer) ConnectSystem(serverName, clientID, systemNo, userName, password, language, sapRouter string) {
	server.ConnParams.User = userName
	//  server.ConnParams.Password,
	server.ConnParams.HostName = serverName
	server.ConnParams.SystemNumber = systemNo
	server.ConnParams.Client = clientID
	server.ConnParams.Lang = language
	server.ConnParams.SapRouter = sapRouter

	connPara := gorfc.ConnectionParameters{
		"user":      server.ConnParams.User,         //"USERNAME"
		"passwd":    password,                       //"PASSWORD"
		"ashost":    server.ConnParams.HostName,     //"192.168.0.1"
		"sysnr":     server.ConnParams.SystemNumber, //"00"
		"client":    server.ConnParams.Client,       // "100"
		"lang":      server.ConnParams.Lang,         //"ZH"
		"sapRouter": server.ConnParams.SapRouter,    //"/H/192.XX.XX.XX/H/XX.XX.XX.XX/H/"
	}
	server.ServerInstance, _ = gorfc.ConnectionFromParams(connPara)
}

func (server *SapServer) CloseSystem() bool {
	if server.ServerInstance != nil {
		if err := server.ServerInstance.Close(); err != nil {
			return false
		}
		return true
	}
	return true
}

func (server *SapServer) Clear() bool {
	if server.ServerInstance != nil {
		server.ServerInstance.Close()
	}
	server.ConnParams = ServerParameters{}
	server.ServerName = ""
	return true
}

func (server *SapServer) CallRFC(FuncName string, params interface{}) (map[string]interface{}, error) {
	if server.ServerInstance != nil {
		// 调用函数 GetXXXFromResult 处理这个结果
		return server.ServerInstance.Call(FuncName, params) //(result map[string]interface{}, err error)
		// if err != nil {
		// 	return nil, err
		// }

	} else {
		return nil, fmt.Errorf("服务器没链接")
	}
	// resultJson, _ = json.Marshal(returnResult)

}

func (server *SapServer) ReadTable(tableName string, conds, fields []string) ([]map[string]interface{}, error) {
	// 调用READ_TABLE函数获取表内容
	// 如：ReadTable("SPFLI", []string{{ "COUNTRYFR EQ 'US' "}, []string{"CARRID", "CITYFROM"})
	//     ReadTable("DD02T", []string{"TABNAME EQ 'MARA'"," AND DDLANGUAGE EQ '1' "}, []string{})
	var resultTable []map[string]interface{}
	var parameters map[string]interface{}
	DELIMITER := "|"
	parameters = map[string]interface{}{"QUERY_TABLE": tableName,
		"DELIMITER": DELIMITER,
		"OPTIONS":   conds,
		"FIELDS":    fields,
	}
	functionName := "RFC_READ_TABLE"
	if server.ServerInstance != nil {
		// 调用函数 GetXXXFromResult 处理这个结果
		rfcResult, err := server.ServerInstance.Call(functionName, parameters) //(result map[string]interface{}, err error)
		if err != nil {
			return nil, err
		}
		resultFields, err3 := GetTableFromResult("FIELDS", rfcResult)
		if err3 != nil {
			return nil, err3
		}
		resultDATA, err2 := GetTableFromResult("DATA", rfcResult)
		if err2 != nil {
			return nil, err2
		}
		if len(resultDATA) == 0 {
			return nil, fmt.Errorf("没有数据")
		}
		for _, data := range resultDATA {
			var row map[string]interface{}
			row = make(map[string]interface{}, len(resultDATA))
			waList := data["WA"].(string)
			fieldsvalue := strings.Split(waList, DELIMITER)
			for idx, value := range fieldsvalue {
				row[resultFields[idx]["FIELDNAME"].(string)] = value
			}
			resultTable = append(resultTable, row)
			// 		wa1 = [wa.strip() for wa in waList]  # 删除文本空格
		}
		return resultTable, nil

	} else {
		return nil, fmt.Errorf("服务器没链接")
	}
}

func (server *SapServer) ExportTableToEXCEL(resultM map[string]interface{}, imFuncName, imParaName, imFileName, imSheetName string) error {
	// var resultTable []map[string]interface{}
	var excelfilename string
	resultTable, err := GetTableFromResult(imParaName, resultM)
	if err != nil {
		return err
	}
	fileExt := strings.ToUpper(filepath.Ext(imFileName))
	if len(fileExt) < 4 || fileExt[1:4] != "XLS" {
		excelfilename = imFileName + ".XLSX"
	} else {
		excelfilename = imFileName
	}

	var ExcelColumn = make([]excel.ExcelHeader, 0)
	var column excel.ExcelHeader
	if imFileName != "" && imParaName != "" {
		funcPara, err := server.GetFunctionParameterType(imFuncName, imParaName)
		if err != nil {
			return err
		}
		tableStru, err := server.GetTableStructure(funcPara.TypeName)
		if err != nil {
			return err
		}
		for _, fields := range tableStru {
			column.ColumnName = fields.FieldName
			column.ColumnLable = fields.FieldTitle
			column.ColumnPos = fields.Position
			if column.ColumnLable == "" {
				column.ColumnLable = column.ColumnName
			}
			ExcelColumn = append(ExcelColumn, column)
		}
	}
	err = excel.ExportExcelByMap(resultTable, excelfilename, imSheetName, ExcelColumn)
	if err != nil {
		return err
		// println(err.Error())
	}
	return nil
}

func (server *SapServer) GetFunctionDescription(imFuncName string) (goFuncDesc gorfc.FunctionDescription, err error) {
	return server.ServerInstance.GetFunctionDescription(imFuncName)
}

func (server *SapServer) GetFunctionParameters(imFuncName string) (goFuncPara []FunctionParameters, err error) {
	funcDesc, err1 := server.ServerInstance.GetFunctionDescription(imFuncName)

	var para FunctionParameters
	if err1 != nil {
		err = err1
		return
	}
	for _, v := range funcDesc.Parameters {
		para.Name = v.Name
		para.ParameterType = v.ParameterType
		para.Direction = v.Direction
		para.NucLength = v.NucLength
		para.UcLength = v.UcLength
		para.Decimals = v.Decimals
		para.DefaultValue = v.DefaultValue
		para.ParameterText = v.ParameterText
		para.Optional = v.Optional
		para.TypeName = v.TypeDesc.Name
		goFuncPara = append(goFuncPara, para)
	}
	return
}

func (server *SapServer) PrintFunctionParameters(imFuncName string) {
	funcPara, err := server.GetFunctionParameters(imFuncName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// for _, v := range funcPara {
	// 	fmt.Printf("\nDirection:%v, Name:%v, Text:%v, Type:%v, Length:%v, Opt:%v, Default:%v, TypeName:%v", v.Direction, v.Name, v.ParameterText, v.ParameterType, v.NucLength, v.Optional, v.DefaultValue, v.TypeName)
	// }
	fmt.Println(PadRight("Direction", 10), PadRight("Name", 20), PadRight("Type", 20), PadRight("Opt", 5), PadRight("Default", 8), PadRight("TypeName", 20), PadRight("Desc", 30))
	for _, v := range funcPara {
		fmt.Println(PadRight(v.Direction, 10), PadRight(v.Name, 20), PadRight(v.ParameterType, 20), PadRight(strconv.FormatBool(v.Optional), 5), PadRight(v.DefaultValue, 8), PadRight(v.TypeName, 20), PadRight(v.ParameterText, 30))
	}
}

func (server *SapServer) GetFunctionParameterType(imFuncName, imParaName string) (goFuncParaType FunctionParameterType, err error) {
	funcDesc, err1 := server.ServerInstance.GetFunctionDescription(imFuncName)
	var para FunctionParameterTypeField
	if err1 != nil {
		err = err1
		return
	}
	for _, paraValue := range funcDesc.Parameters {
		if paraValue.Name == imParaName && paraValue.TypeDesc.Name != "" {
			goFuncParaType.Name = imParaName
			goFuncParaType.TypeName = paraValue.TypeDesc.Name
			fields := paraValue.TypeDesc.Fields
			for _, v := range fields {
				para.Name = v.Name
				para.FieldType = v.FieldType
				para.NucLength = v.NucLength
				para.NucOffset = v.NucOffset
				para.UcLength = v.UcLength
				para.UcOffset = v.UcOffset
				para.Decimals = v.Decimals
				para.TypeName = v.TypeDesc.Name
				goFuncParaType.fields = append(goFuncParaType.fields, para)
			}
		}
	}
	return
}

func (server *SapServer) PrintFunctionParameterType(imFuncName, imParaName string) {
	funcPara, err := server.GetFunctionParameterType(imFuncName, imParaName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range funcPara.fields {
		fmt.Printf("\nName:%v,FieldType:%v,Length:%v,NucOffset:%v,Decimals:%v,TypeName:%v", v.Name, v.FieldType, v.NucLength, v.NucOffset, v.Decimals, v.TypeName)
	}
}
func (server *SapServer) GetTableDescription(imTableName string) (desc string, err error) {
	talcond := "TABNAME EQ '" + imTableName + "'"
	connAttributes, _ := server.ServerInstance.GetConnectionAttributes()
	language := connAttributes["language"]
	secCond := " AND DDLANGUAGE EQ '" + language + "'"
	tabDesc, err1 := server.ReadTable("DD02T", []string{talcond, secCond}, []string{})
	if err1 != nil {
		err = err1
		desc = ""
		return
	}
	desc = tabDesc[0]["DDTEXT"].(string)
	return
}

func (server *SapServer) GetTableStructure(imTableName string) (goTableStruc []TableStructure, err error) {
	var tabStru TableStructure
	params := map[string]interface{}{
		"TABNAME":   imTableName,
		"ALL_TYPES": "X",
	}
	functionName := "DDIF_FIELDINFO_GET"
	resultMap, err := server.CallRFC(functionName, params)
	if err != nil {
		return
	}
	iTabType, _ := resultMap["DDOBJTYPE"].(string)
	if iTabType == "" {
		err = fmt.Errorf("类型错误")
		return
	}
	if iTabType == "TTYP" {
		resultWa, err1 := GetWorkAreaFromResult("DFIES_WA", resultMap)
		if err1 != nil {
			err = err1
			return
		}
		tableName := resultWa["ROLLNAME"].(string)
		if tableName != "" {
			return server.GetTableStructure(tableName)
		}
	}
	tableFields, err := GetTableFromResult("DFIES_TAB", resultMap)
	for _, field := range tableFields {
		tabStru.TableName = field["TABNAME"].(string)
		tabStru.FieldName = field["FIELDNAME"].(string)
		tabStru.Position, _ = strconv.Atoi(field["POSITION"].(string))
		tabStru.OffSet, _ = strconv.Atoi(field["OFFSET"].(string))
		tabStru.RollName = field["ROLLNAME"].(string)
		tabStru.Leng, _ = strconv.Atoi(field["LENG"].(string))
		tabStru.Decimals, _ = strconv.Atoi(field["DECIMALS"].(string))
		tabStru.DataType = field["DATATYPE"].(string)
		tabStru.FieldText = field["FIELDTEXT"].(string)
		tabStru.FieldTitle = field["SCRTEXT_L"].(string)
		goTableStruc = append(goTableStruc, tabStru)
	}
	return
}

func (server *SapServer) TableStructureToExcel(imTableName, fileName string) {
	var resultMap []map[string]interface{}
	// var headers []excel.ExcelHeader
	// headers= make([]excel.ExcelHeader,0 )
	tabStru, err := server.GetTableStructure(imTableName)
	// resultMap = make([]map[string]interface{}, len(tabStru))

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, field := range tabStru {
		fieldmap, _ := excel.StructToMap(field)
		resultMap = append(resultMap, fieldmap)
	}
	headers := []excel.ExcelHeader{
		{ColumnName: "FieldName", ColumnLable: "FieldName", ColumnPos: 1},
		{ColumnName: "FieldText", ColumnLable: "FieldText", ColumnPos: 2},
		{ColumnName: "Position", ColumnLable: "Position", ColumnPos: 3},
		{ColumnName: "DataType", ColumnLable: "DataType", ColumnPos: 4},
		{ColumnName: "Leng", ColumnLable: "Length", ColumnPos: 5},
		{ColumnName: "Decimals", ColumnLable: "Decimals", ColumnPos: 6},
		{ColumnName: "RollName", ColumnLable: "RollName", ColumnPos: 7},
		{ColumnName: "FieldTitle", ColumnLable: "FieldTitle", ColumnPos: 8},
		{ColumnName: "TableName", ColumnLable: "TableName", ColumnPos: 9},
	}

	excel.ExportExcel(resultMap, fileName, imTableName, headers)

}

func (server *SapServer) PrintTableStructure(imTableName string) {
	tabStru, err := server.GetTableStructure(imTableName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// for _, v := range tabStru {
	// 	fmt.Printf("\nTabName:%v, FieldName:%v, Pos:%v, Desc:%v, Offset:%v, RollName:%v, Length:%v, Decimals:%v, DataType:%v, Title:%v", v.TableName, v.FieldName, v.Position, v.FieldText, v.OffSet, v.RollName, v.Leng, v.Decimals, v.DataType, v.FieldTitle)
	// }
	fmt.Println(PadRight("TabName", 20), PadRight("FieldName", 10), PadRight("Pos", 5), PadRight("Length", 6), PadRight("DataType", 10), PadRight("RollName", 20), PadRight("Title", 30), PadRight("Desc", 30))
	for _, v := range tabStru {
		fmt.Println(PadRight(v.TableName, 20), PadRight(v.FieldName, 10), PadRight(strconv.Itoa(v.Position), 5), PadRight(strconv.Itoa(v.Leng), 6), PadRight(v.DataType, 10), PadRight(v.RollName, 20), PadRight(v.FieldTitle, 30), PadRight(v.FieldText, 30))
	}

}

func (server *SapServer) RestExecuteSql(sqlStatement string) (result []map[string]interface{}, err error) {

	return
}

// 获取对象代码imObjType ： class，PROG，FUNC，TABL，如果INFO，获取对象清单及描述
func (server *SapServer) RestReadSource(imObjName string, imObjType string) (string, error) {
	var url, urltxtsymbol, urltxtselection, retString string
	var errall error
	var objName string
	acceptTypeDefault := "text/plain"
	if imObjName == "" {
		return "", fmt.Errorf("对象名称错误")
	}
	objName = strings.TrimSpace(strings.ToLower(imObjName))
	switch imObjType {
	case "CLASS":
		url = fmt.Sprintf(`/sap/bc/adt/oo/classes/%s/source/main`, objName)
		urltxtsymbol = fmt.Sprintf(`/sap/bc/adt/textelements/classes/%s/source/symbols`, objName)

	case "PROG":
		url = fmt.Sprintf(`/sap/bc/adt/programs/programs/%s/source/main`, objName)
		urltxtsymbol = fmt.Sprintf(`/sap/bc/adt/textelements/programs/%s/source/symbols`, objName)
		urltxtselection = fmt.Sprintf(`/sap/bc/adt/textelements/programs/%s/source/selections`, objName)

	case "FUNC":
		condString := "FUNCNAME EQ '" + strings.TrimSpace(strings.ToUpper(imObjName)) + "'"
		funcs, err := server.ReadTable("TFDIR", []string{condString}, []string{"PNAME"})
		if err != nil {
			return "", nil
		}
		funcGrupName := funcs[0]["PNAME"].(string)[4:]
		funcGrupName = strings.TrimSpace(strings.ToLower(funcGrupName))
		url = fmt.Sprintf(`/sap/bc/adt/functions/groups/%s/fmodules/%s/source/main`, funcGrupName, objName)
		urltxtsymbol = fmt.Sprintf(`/sap/bc/adt/textelements/functiongroups/%s/source/symbols`, funcGrupName)
		urltxtselection = fmt.Sprintf(`/sap/bc/adt/textelements/functiongroups/%s/source/selections`, funcGrupName)
	case "TABL":
		retString, _ = server.RestReadStructure(objName)
		return retString, nil
	case "INFO":
		retString, _ = server.RestReadObjectList(objName)
		return retString, nil
	default:
		return "", fmt.Errorf("类型错误")
	}
	retString, errall = server.RestCallURL(url, acceptTypeDefault)
	if errall != nil {
		return retString, errall
	}

	retString += `\n\n"=====以下是相关文本信息=====\n\n`

	// text symbols
	if urltxtsymbol != "" {
		retTxtSymbol, errt := server.RestCallURL(urltxtsymbol, "application/vnd.sap.adt.textelements.symbols.v1")
		if errt != nil {
			return retString, errt
		}
		txtList := strings.Split(retTxtSymbol, `\r`)
		retString += `\n\n"=====以下是文本元素=====\n`
		for _, txt := range txtList {
			if txt != "" && txt[:1] != "@" {
				retString += `\n` + txt
			}
		}
	}

	// 以下是Selection文本
	if urltxtselection != "" {
		retTxtSymbol, errts := server.RestCallURL(urltxtselection, "application/vnd.sap.adt.textelements.selections.v1")
		if errts != nil {
			return retString, errts
		}
		txtList := strings.Split(retTxtSymbol, `\r`)
		retString += `\n\n"=====以下是Selection文本=====\n`
		for _, txt := range txtList {
			if txt != "" && txt[:1] != "@" {
				retString += `\n` + txt
			}
		}
	}
	return retString, nil

}

// 获取表结构
func (server *SapServer) RestReadStructure(imObjName string) (string, []TableStructure) {
	var fields []TableStructure
	url := fmt.Sprintf(`/sap/bc/adt/ddic/elementinfo?path=%s`, strings.TrimSpace(strings.ToLower(imObjName)))
	resultXmlString, err := server.RestCallURL(url, "application/vnd.sap.adt.elementinfo+xml")
	if err != nil {
		return "", fields
	}

	result, err1 := xml.XmlParse(resultXmlString)

	if err1 != nil {
		log.Fatalf("Error parsing XML: %v", err)
	}

	if len(result.Items) == 0 {
		return "", fields
	}
	tabname := result.Attribute["name"]
	tabdesc := result.GetItem(1).Text
	fieldstring := fmt.Sprintf("Talble/Structure:%s, Description:%s\n", tabname, tabdesc)
	for _, item := range result.Items[2:] {
		var field TableStructure
		field.TableName = tabname
		field.FieldTitle = tabdesc
		field.FieldName = item.Attribute["name"]
		if item.Attribute["type"] != "TABL/DTF" {
			continue
		}
		for _, itm := range item.Items {
			if itm.Tag == "documentation" {
				field.FieldText = itm.Text
			} else {
				for _, fielda := range itm.Items {
					switch fielda.Attribute["key"] {
					case "ddicDataElement":
						field.RollName = fielda.Text
					case "ddicDataType":
						field.DataType = fielda.Text
					case "ddicLength":
						field.Leng, _ = strconv.Atoi(fielda.Text)
					case "ddicDecimals":
						field.Decimals, _ = strconv.Atoi(fielda.Text)
						// case "ddicDataElement":
					}

				}
			}
		}
		fields = append(fields, field)
	}

	for idx, value := range fields {
		if idx == 0 {
			fieldstring += "\n" + "\t" + "fieldname" + "\t" + "fieldType" + "\t" + "fieldText"
			fieldstring += "\n=====================================================\n"
		}

		fieldstring += "\n" + "\t" + value.FieldName + "\t" + value.DataType + "\t" + value.FieldText

	}
	return fieldstring, fields

}

// 获取对象清单
func (server *SapServer) RestReadObjectList(imObjName string) (orgString string, results []map[string]string) {
	url := fmt.Sprintf(`/sap/bc/adt/repository/informationsystem/search?operation=quickSearch&query=%s*&maxResults=51`, strings.TrimSpace(strings.ToLower(imObjName)))
	resultXmlString, err := server.RestCallURL(url, "application/xml")
	if err != nil {
		return
	}

	resultNodes, err1 := xml.XmlParse(resultXmlString)

	if err1 != nil {
		log.Fatalf("Error parsing XML: %v", err)
	}
	for _, item := range resultNodes.Items {
		var result = make(map[string]string)
		result["name"] = item.Attribute["name"]
		result["type"] = item.Attribute["type"]
		result["description"] = item.Attribute["description"]
		result["uri"] = item.Attribute["uri"]
		result["packageName"] = item.Attribute["packageName"]
		results = append(results, result)
	}

	for idx, value := range results {
		if idx == 0 {
			orgString += "\n" + "\t" + "Name" + "\t" + "Type" + "\t" + "Desc" + "\t" + "packageName" + "\t" + "uri"
			orgString += "\n=====================================================\n"
		}

		orgString += "\n" + "\t" + value["name"] + "\t" + value["type"] + "\t" + value["description"] + "\t" + value["packageName"] + "\t" + value["uri"]

	}

	return
}

// 有ADT 开发权限
// call function:SADT_REST_RFC_ENDPOINT
func (server *SapServer) RestCallURL(url, acceptType string) (string, error) {
	// # abap URL list 代码
	// # DATA(lr_inst) = cl_adt_tools_core_factory=>get_instance( ).
	// # DATA(lr_urlmapper) = lr_inst->get_uri_mapper( ).
	// # " CL_ADT_URI_TEMPLATES_SHM_ROOT->BUILD_TEMPLATES
	// # DATA(type_provider) = cl_wb_registry=>get_objtype_provider( ).
	// # type_provider->get_objtypes( IMPORTING p_objtype_data = DATA(wbobjtype_data) ).
	// # LOOP AT wbobjtype_data->mt_wbobjtype ASSIGNING FIELD-SYMBOL(<type>) WHERE pgmid              = 'R3TR'
	// #                                                                     AND   is_main_subtype_wb = abap_true.
	// #   DATA ls_type TYPE wbobjtype.
	// #   ls_type-objtype_tr = <type>-objecttype.
	// #   ls_type-subtype_wb = <type>-subtype_wb.
	// # *<type>-description
	// #   DATA(lv_type) = |{ <type>-objecttype }{ <type>-subtype_wb }|.
	// #   DATA(uri) = cl_adt_tools_core_factory=>get_instance( )->get_uri_mapper( )->get_adt_object_ref_uri(
	// #                   name = 'ZCL_RTR_COMMON'
	// #                   type = ls_type  ).
	// # ENDLOOP.

	// url:
	// class source:/sap/bc/adt/oo/classes/{classname.lower().strip()}/source/main
	// report source: /sap/bc/adt/programs/programs/{programname.lower().strip()}/source/main source
	// report text: /sap/bc/adt/textelements/programs/{programname.lower().strip()/source/selections?version=workingArea
	//              /sap/bc/adt/textelements/programs/{programname.lower().strip()/source/symbols?version=workingArea
	// funct source: /sap/bc/adt/functions/groups/{funcgrp.lower().strip()}/fmodules/{objName.lower().strip()}/source/main
	// /sap/bc/adt/ddic/tables/mara/source/main
	var acceptTypeString string
	acceptTypeString = acceptType
	if acceptType == "" {
		acceptTypeString = "text/plain"
	}

	request_line := map[string]string{"METHOD": "GET",
		"URI":     url,
		"VERSION": "HTTP/1.1",
	}
	headers := []map[string]string{
		{"NAME": "Cache-Control",
			"VALUE": "no-cache",
		},
		{"NAME": "Accept",
			"VALUE": acceptTypeString,
		},
		{"NAME": "User-Agent",
			"VALUE": "Eclipse/4.28.0.v20230605-0440 (win32; x86_64; Java 17.0.7) ADT/3.38.2 (devedition)",
		},
		{"NAME": "X-sap-adt-profiling",
			"VALUE": "server-time",
		},
	}
	request := map[string]interface{}{"REQUEST_LINE": request_line, "HEADER_FIELDS": headers}
	parameters := map[string]interface{}{"REQUEST": request}

	functionName := "SADT_REST_RFC_ENDPOINT"
	returnResult, err := server.CallRFC(functionName, parameters)
	if err != nil {
		return "", err
	}
	response := returnResult["RESPONSE"].(map[string]interface{})
	body := string(response["MESSAGE_BODY"].([]byte))
	bodyString := strings.ReplaceAll(body, "\n", "") // str(body, encoding="utf8").replace("\n","")
	return bodyString, nil
}
