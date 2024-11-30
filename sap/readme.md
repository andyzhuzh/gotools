
# 配置Windows
+ gorfc v0.1.1 对应版本 nwrfc750P_11
+ 系统变量设置 ：SAPNWRFC_HOME=C:\nwrfcsdk
+ PATH 增加 ：C:\nwrfcsdk\lib;C:\nwrfcsdk\bin
+ GO 环境变量设置：
  
  go env -w GOARCH=amd64
  
  go env -w CGO_ENABLED=1

# 安装minGW64

# 修改gorfc.go 代码：
+ Windows 编译,修改gorfc.go文件：
```go
//go:build windows
#cgo windows CFLAGS: -IC:/nwrfcsdk/include/
#cgo windows LDFLAGS: -LC:/nwrfcsdk/lib/ -lsapnwrfc -llibsapucum
```
  
  如果有类似下面警告信息，想屏蔽( -Wall 显示警告，-w屏蔽警告): 
    warning: 'align'
    warning: ignoring '#pragma warning ' 
    warning: 'unicodeId'
```go
#cgo windows CFLAGS: -Wall -Wno-uninitialized -Wno-long-long
```
改成：
```go
#cgo windows CFLAGS: -w -Wno-uninitialized -Wno-long-long
```

# gorfc 报错 exit status 0xc0000374
+ 如果 返回有日期或时间，程序可能报错：exit status 0xc0000374
  
  解决方法：
  
  修改 gorfc.go 函数 func wrapVariable
```go
大概832行：case C.RFCTYPE_DATE:
// dateValue = (*C.RFC_CHAR)(C.malloc(8)) 
dateValue = (*C.RFC_CHAR)(C.GoMallocU(8))
大概852行： case C.RFCTYPE_TIME:
// timeValue = (*C.RFC_CHAR)(C.malloc(6))
timeValue = (*C.RFC_CHAR)(C.GoMallocU(6))
```

# 调用说明
+ 创建连接
```go
func NewServer(serverName, clientID, systemNo, userName, password, language, sapRouter string) (server SapServer) {
	server.ConnectSystem(serverName, clientID, systemNo, userName, password, language, sapRouter)
	return
}
```
+ call Function 
```go
func testfunction(imSysName, imFileName string) {
	sapinstance := NewServer(...)
	defer sapinstance.Clear()
	// datef, _ := time.Parse(time.DateOnly, "20210108")
	datet, _ := time.Parse(time.DateOnly, "2024-09-09")
	params := map[string]interface{}{
		"HOLIDAY_CALENDAR": "ZT",
		"DATE_FROM":        func() time.Time { datef, _ := time.Parse(time.DateOnly, "2021-01-08"); return datef }(), //"20210108",
		"DATE_TO":          datet,
	}
	functionName := "DAY_ATTRIBUTES_GET"

	fmt.Println("开始查询：", time.Now().Format("2006-01-02 15:04:05"))
	resultMap, err := sapinstance.CallRFC(functionName, params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("查询完：", time.Now().Format("2006-01-02 15:04:05"))

	resultM := resultMap
	resultParaList := sap.GetObjectListFromResult(resultM)
	fmt.Println("返回参数列表：")
	for k, v := range resultParaList {
		fmt.Println(PadRight(k, 20), " : ", v)
	}
	// fmt.Println(sap.GetObjectListFromResult(resultM))

	// 保存excel文件
	excelfilename := imFileName + ".XLSX"
	err1 := sapinstance.ExportTableToEXCEL(resultM, functionName, "DAY_ATTRIBUTES", excelfilename, "DAY_ATTRIBUTES")
	if err1 != nil {
		println(err1.Error())
	}

	fmt.Println("文件保存完成：", time.Now().Format("2006-01-02 15:04:05"))
}

```

+ 显示函数参数列表
```go
func PrintFunctionParameters(imSysName, imFuncName string) {
	sapinstance := NewServer(...)
	defer sapinstance.Clear()
	sapinstance.PrintFunctionParameters(imFuncName)
}
```

+ 显示函数参数类型
```go
func PrintFunctionParameterType(imSysName, imFuncName, imParaName string) {
	sapinstance := NewServer(...)
	defer sapinstance.Clear()
	sapinstance.PrintFunctionParameterType(imFuncName, imParaName)
}
```

+ 显示表结构
```go
func PrintTableStructure(imSysName, imTabName string) {
	sapinstance := NewServer(...)
	defer sapinstance.Clear()
	sapinstance.PrintTableStructure(imTabName)

}
```
+ 表结构保存excel
```go
func TableStructureToExcel(imSysName, imTabName, fileName string) {
	sapinstance := NewServer(...)
	defer sapinstance.Clear()
	sapinstance.TableStructureToExcel(imTabName, fileName)

}
```


