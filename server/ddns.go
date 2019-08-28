package server

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/youngtrips/ddns/internal/config"
	"github.com/youngtrips/ddns/internal/utils"
	"github.com/youngtrips/ddns/plugin"
	_ "github.com/youngtrips/ddns/plugin/alidns"
)

var (
	ERR_INVALID_KEY        = errors.New("invalid access key id")
	ERR_INVALID_KEY_SECRET = errors.New("invalid access key secret")
	ERR_INVALID_DOMAIN     = errors.New("invalid domain")
)

func now() int64 {
	return time.Now().UnixNano() / 1e6
}

func Run(ctx context.Context) error {
	cfg := config.Config()

	dnsPlugin := plugin.Get(cfg.Plugin.Name)
	if dnsPlugin == nil {
		return errors.New("no found dns plugin: " + cfg.Plugin.Name)
	}

	if err := dnsPlugin.Init(cfg.Plugin.Params); err != nil {
		log.Error("init dns plugin failed:", err)
		return err
	}

	domain := cfg.Domain
	records := make(map[string]string)
	for _, rr := range cfg.Records {
		val, _ := dnsPlugin.QueryRR(domain, rr)
		records[rr] = val
	}

	log.Info("check interval: ", cfg.Interval)
	tick := time.Tick(time.Duration(100 * time.Millisecond))
	prev := int64(0)
	interval := int64(cfg.Interval) * 1000
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick:
			curr := now()
			if curr >= prev+interval {
				if err := onTick(dnsPlugin, domain, records); err != nil {
					log.Error("check ip and update domain failed: ", err)
				}
				prev = curr
			}
		}
	}
}

func onTick(dnsPlugin plugin.DNSPlugin, domain string, records map[string]string) error {
	currIp, err := utils.MyIP("")
	if err != nil {
		return err
	}
	log.Info("curr ip: ", currIp)

	for rr, prevIp := range records {
		if prevIp != currIp {
			if err := dnsPlugin.UpdateRR(domain, rr, currIp); err != nil {
				log.Errorf("update domain RR failed: %s.%s, %s", rr, domain, err)
			} else {
				log.Infof("update domain RR: %s.%s %s=>%s", rr, domain, prevIp, currIp)
			}
			records[rr] = currIp
		}
	}
	return nil
}
