package ziva

import "fmt"

type ProxyIP struct {
	IP    string
	Port  int
	User  string
	Pass  string
	IsTLS bool
}

func NewProxyIP(ip string, port int, user, pass string, isTLS bool) *ProxyIP {
	return &ProxyIP{
		IP:    ip,
		Port:  port,
		User:  user,
		Pass:  pass,
		IsTLS: isTLS,
	}
}

func (p *ProxyIP) String() string {
	scheme := "http"
	if p.IsTLS {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s:%s@%s:%d", scheme, p.User, p.Pass, p.IP, p.Port)
}
