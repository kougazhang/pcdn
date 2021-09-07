package api

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	auth = Auth{
		CP:     os.Getenv("pcdn_cp"),
		Seckey: os.Getenv("pcdn_seckey"),
	}
	pcdn = NewPCDN(PCDNParams{
		Auth: auth,
	})
	token = ""
)

func TestPCDN_GetAccessToken(t *testing.T) {
	resp, err := pcdn.GetAccessToken(time.Now().Unix())
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func TestPCDN_Preheat(t *testing.T) {
	resp, err := pcdn.Preheat([]string{"http://www.test.com/1"}, token)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func TestPCDN_Purge(t *testing.T) {
	resp, err := pcdn.Purge([]string{"http://<bindDomain>/test.now"}, token)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func TestPCDN_DomainList(t *testing.T) {
	resp, err := pcdn.DomainList(token)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func TestPCDN_GetLogfile(t *testing.T) {
	reqBody := LogfileRequest{
		Domain:     "",
		Begintime:  1630944000,
		Endtime:    1630944000 + 3600*2,
		FormatType: LogfileFormatTypeHour,
	}

	resp, err := pcdn.GetLogfile(reqBody, token)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
