package alidns

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	DNS_RECORD_TYPE_A            = "A"
	DNS_RECORD_TYPE_MX           = "MX"
	DNS_RECORD_TYPE_CNAME        = "CNAME"
	DNS_RECORD_TYPE_TXT          = "TXT"
	DNS_RECORD_TYPE_REDIRECT_URL = "REDIRECT_URL"
	DNS_RECORD_TYPE_FORWORD_URL  = "FORWORD_URL"
	DNS_RECORD_TYPE_NS           = "NS"
	DNS_RECORD_TYPE_AAAA         = "AAAA"
	DNS_RECORD_TYPE_SRV          = "SRV"

	DNS_API_HOST     = "http://alidns.aliyuncs.com"
	HTTP_METHOD_GET  = "GET"
	HTTP_METHOD_POST = "POST"
)

type RecordInfo struct {
	DomainName string `json:"DomainName"`
	RecordId   string `json:"RecordId"`
	RR         string `json:"RR"`
	Type       string `json:"Type"`
	Value      string `json:"Value"`
	TTL        int32  `json:"TTL"`
	Priority   int32  `json:"Priority"`
	Line       string `json:"Line"`
	Status     string `json:"Status"`
	Locked     bool   `json:"Locked"`
	Weight     int32  `json:"Weight"`
}

type RecordType struct {
	Record []RecordInfo `json:"Record"`
}

type DescribeSubDomainRecordsRes struct {
	RequestId     string     `json:"RequestId"`
	TotalCount    int32      `json:"TotalCount"`
	PageNumber    int32      `json:"PageNumber"`
	PageSize      int32      `json:"PageSize"`
	DomainRecords RecordType `json:"DomainRecords"`
}

type Client struct {
	accessKeyId     string
	accessKeySecret string
}

func NewClient(accessKeyId string, accessKeySecret string) *Client {
	return &Client{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
	}
}

/*
 * subDomain	名名称，如www.abc.com，如果输入的是abc.com，则认为是@.abc.com；
 * dnsRecordType 如果不填写，则返回子域名对应的全部解析记录类型。解析类型包括(不区分大小写)：A、MX、CNAME、TXT、REDIRECT_URL、FORWORD_URL、NS、AAAA、SRV
 */
func (c *Client) DescribeSubDomainRecords(subDomain string, dnsRecordType string) ([]RecordInfo, error) {
	res := &DescribeSubDomainRecordsRes{}
	err := c.doRequest(
		DNS_API_HOST,
		HTTP_METHOD_GET,
		res,
		ParamInfo{"Action", "DescribeSubDomainRecords"},
		ParamInfo{"SubDomain", subDomain},
		ParamInfo{"PageNumber", "1"},
		ParamInfo{"PageSize", "500"},
		ParamInfo{"Type", dnsRecordType},
	)
	if err != nil {
		return nil, err
	}
	return res.DomainRecords.Record, nil
}

func (c *Client) doRequest(host string, method string, res interface{}, params ...ParamInfo) error {

	encoded := BuildParam(c.accessKeyId, c.accessKeySecret, method, params)
	if method == "GET" {
		url := host + "/?" + encoded
		resp, err := http.Get(url)
		if err != nil {
			return err
		} else {
			defer resp.Body.Close()
			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				return err
			} else {
				return json.Unmarshal(body, res)
			}
		}
	}
	return errors.New("unsupported method: " + method)
}

func (c *Client) SetDomainRecorad(domain string, rr string, dnsRecordType string, recordId string, value string) error {
	res := &DescribeSubDomainRecordsRes{}
	err := c.doRequest(
		DNS_API_HOST,
		HTTP_METHOD_GET,
		res,
		ParamInfo{"Action", "UpdateDomainRecord"},
		ParamInfo{"RecordId", recordId},
		ParamInfo{"RR", rr},
		ParamInfo{"Type", dnsRecordType},
		ParamInfo{"Value", value},
	)
	return err
}

func (c *Client) UpdateDomainRecord(domain string, rr string, ip string) error {
	records, err := c.DescribeSubDomainRecords(rr+"."+domain, DNS_RECORD_TYPE_A)
	if err != nil {
		return err
	}

	if len(records) != 1 {
		return errors.New(fmt.Sprintf("do not support domain(%s) has multi records.", domain))
	}

	return c.SetDomainRecorad(domain, rr, DNS_RECORD_TYPE_A, records[0].RecordId, ip)
}
