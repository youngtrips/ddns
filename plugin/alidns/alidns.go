package alidns

import (
	"errors"
	"fmt"

	"github.com/youngtrips/ddns/plugin"
)

type AliDNSPlugin struct {
	client *Client
}

const (
	NAME = "alidns"
)

func init() {
	plugin.Register(NAME, &AliDNSPlugin{})
}

func (p *AliDNSPlugin) Init(params map[string]string) error {
	accessKeyId := params["access_key_id"]
	accessKeySecret := params["access_key_secret"]

	p.client = NewClient(accessKeyId, accessKeySecret)
	return nil
}

func (p *AliDNSPlugin) QueryRR(domain string, rr string) (string, error) {
	if p.client == nil {
		return "", errors.New("invalid client")
	}

	records, err := p.client.DescribeSubDomainRecords(rr+"."+domain, DNS_RECORD_TYPE_A)
	if err != nil {
		return "", err
	}

	if len(records) == 0 {
		return "", nil
	}

	if len(records) != 1 {
		return "", errors.New(fmt.Sprintf("do not support domain(%s) has multi records.", domain))
	}
	return records[0].Value, nil
}

func (p *AliDNSPlugin) UpdateRR(domain string, rr string, ip string) error {
	if p.client == nil {
		return errors.New("invalid client")
	}
	return p.client.UpdateDomainRecord(domain, rr, ip)
}
