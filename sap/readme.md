
# 配置Windows
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


