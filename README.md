# the cdn sdk of instafogging

## sign
```go 
func TestAuth_Sign(t *testing.T) {
	a := Auth{
		CP:     "your-account",
		Seckey: "your-Seckey",
	}
	ctime := int64(your-request-timestamp)
	sign := a.Sign(ctime)
	fmt.Println(a.Sign(ctime))
}
```

## NewPCDN
初始化结构体 PCDN

## GetAccessToken
获取 Access Token

## Preheat
提交预热

## Purge
提交刷新

## DomainList
域名列表

## GetLogfile
获取日志下载链接