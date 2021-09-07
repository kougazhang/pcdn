package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	LogfileFormatTypeHour = "hour"
	LogfileFormatType5Min = "5min"
)

type PCDN struct {
	Auth
	Host string `json:"host"`
}

type PCDNParams struct {
	Host string
	Auth Auth
}

func NewPCDN(params PCDNParams) PCDN {
	if len(params.Host) == 0 {
		params.Host = "http://api.instafogging.com"
	}
	params.Host = strings.TrimSuffix(params.Host, "/")

	return PCDN{
		Auth: params.Auth,
		Host: params.Host,
	}
}

type AccessTokenResp struct {
	Retcode     string `json:"retcode"`
	Retmsg      string `json:"retmsg"`
	AccessToken string `json:"access_token"`
	Expire      string `json:"expire"`
}

func (p PCDN) GetAccessToken(ctime int64) (*AccessTokenResp, error) {
	sign := p.Sign(ctime)
	reqData, err := json.Marshal(map[string]string{
		"cp":    p.CP,
		"ctime": strconv.FormatInt(ctime, 10),
		"sign":  sign,
	})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/auth/token/get", p.Host)
	response, err := http.Post(url, "Content-Type: application/json; charset=utf-8", bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var resp AccessTokenResp
	err = json.Unmarshal(data, &resp)
	return &resp, err
}

func (p PCDN) Request(method, url, accessToken string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", accessToken)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

type PreheatResp struct {
	Retcode string `json:"retcode"`
	Retmsg  string `json:"retmsg"`
	Data    Data   `json:"data"`
}

type Data struct {
	TaskId string `json:"TaskId"`
}

func (p PCDN) Preheat(urls []string, accessToken string) (*PreheatResp, error) {
	if len(urls) > 100 {
		return nil, fmt.Errorf("the length of urls is beyond of 100")
	}

	url := fmt.Sprintf("%s/pcdn/cdm/prefetch", p.Host)
	data, err := json.Marshal(map[string]interface{}{"Urls": urls})
	if err != nil {
		return nil, err
	}

	response, err := p.Request("POST", url, accessToken, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var resp PreheatResp
	err = json.Unmarshal(response, &resp)
	return &resp, err
}

type PurgeResp struct {
	Retcode string `json:"retcode"`
	Retmsg  string `json:"retmsg"`
	Data    Data   `json:"data"`
}

func (p PCDN) Purge(urls []string, accessToken string) (*PurgeResp, error) {
	if len(urls) > 100 {
		return nil, fmt.Errorf("the length of urls is beyond of 100")
	}

	url := fmt.Sprintf("%s/pcdn/cdm/purge", p.Host)
	data, err := json.Marshal(map[string]interface{}{
		"Urls": urls,
		"Type": "file",
	})
	if err != nil {
		return nil, err
	}

	response, err := p.Request("POST", url, accessToken, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var resp PurgeResp
	err = json.Unmarshal(response, &resp)
	return &resp, err
}

type DomainListResp struct {
	Retcode string             `json:"retcode"`
	Retmsg  string             `json:"retmsg"`
	Data    []ItemOfDomainList `json:"data"`
}

type ItemOfDomainList struct {
	Value      string `json:"value"`
	DomainName string `json:"domain_name"`
}

func (p PCDN) DomainList(accessToken string) (*DomainListResp, error) {
	url := fmt.Sprintf("%s/pcdn/cdm/domainlist", p.Host)
	response, err := p.Request("POST", url, accessToken, nil)
	if err != nil {
		return nil, err
	}

	var resp DomainListResp
	err = json.Unmarshal(response, &resp)
	return &resp, err
}

type logfileRequest struct {
	Domain     string `json:"domain"`
	Begintime  string `json:"begintime"`
	Endtime    string `json:"endtime"`
	FormatType string `json:"format_type"`
}

type LogfileRequest struct {
	Domain     string `json:"domain"`
	Begintime  int64  `json:"begintime"`
	Endtime    int64  `json:"endtime"`
	FormatType string `json:"format_type"`
}

type LogfileResponse struct {
	Retcode string          `json:"retcode"`
	Retmsg  string          `json:"retmsg"`
	Data    []ItemOfLogfile `json:"data"`
}

type ItemOfLogfile struct {
	Host       string `json:"host"`
	LogDate    string `json:"log_date"`
	FormatType string `json:"format_type"`
	LogPath    string `json:"log_path"`
}

func (p PCDN) GetLogfile(requestBody LogfileRequest, accessToken string) (*LogfileResponse, error) {
	format := func(timestamp int64) string {
		return time.Unix(timestamp, 0).Format("200601021504")
	}

	req := logfileRequest{
		Domain:     requestBody.Domain,
		Begintime:  format(requestBody.Begintime),
		Endtime:    format(requestBody.Endtime),
		FormatType: requestBody.FormatType,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/pcdn/cdm/getlogfile", p.Host)
	response, err := p.Request("POST", url, accessToken, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var resp LogfileResponse
	err = json.Unmarshal(response, &resp)
	return &resp, err
}
