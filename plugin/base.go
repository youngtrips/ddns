package plugin

type DNSPlugin interface {
	Init(params map[string]string) error
	QueryRR(domain string, rr string) (string, error)
	UpdateRR(domain string, rr string, ip string) error
}

var (
	_plugins map[string]DNSPlugin
)

func init() {
	_plugins = make(map[string]DNSPlugin)
}

func Register(name string, plugin DNSPlugin) {
	_plugins[name] = plugin
}

func Get(name string) DNSPlugin {
	p, present := _plugins[name]
	if !present {
		return nil
	}
	return p
}
