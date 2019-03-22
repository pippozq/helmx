package spec

import (
	"fmt"
	"github.com/go-courier/helmx/constants"
	"net/url"
	"strconv"
	"strings"
)

func ParsePort(s string) (*Port, error) {
	if s == "" {
		return nil, fmt.Errorf("missing port value")
	}

	port := uint16(0)
	targetPort := uint16(0)
	protocol := ""

	parts := strings.Split(s, "/")

	s = parts[0]

	if len(parts) == 2 {
		protocol = strings.ToLower(parts[1])
	}

	ports := strings.Split(s, ":")

	p, err := strconv.ParseUint(ports[0], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid port %v", ports[0])
	}

	port = uint16(p)

	if len(ports) == 2 {
		p, err := strconv.ParseUint(ports[1], 10, 16)
		if err != nil {
			panic(fmt.Errorf("invalid target port %v", ports[1]))
		}
		targetPort = uint16(p)
	} else {
		targetPort = port
	}

	return &Port{
		Port:          port,
		ContainerPort: targetPort,
		Protocol:      constants.Protocol(strings.ToUpper(protocol)),
	}, nil
}

type Port struct {
	Port          uint16
	ContainerPort uint16
	Protocol      constants.Protocol
}

func (s Port) String() string {
	v := ""
	if s.Protocol != "" {
		v = "/" + strings.ToLower(string(s.Protocol))
	}

	if s.ContainerPort != 0 && s.ContainerPort != s.Port {
		v = ":" + strconv.FormatUint(uint64(s.ContainerPort), 10) + v
	}

	if s.Port != 0 {
		v = strconv.FormatUint(uint64(s.Port), 10) + v
	}

	return v
}

func (s Port) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Port) UnmarshalText(data []byte) error {
	servicePort, err := ParsePort(string(data))
	if err != nil {
		return err
	}
	*s = *servicePort
	return nil
}

func ParseIngressRule(s string) (*IngressRule, error) {
	if s == "" {
		return nil, fmt.Errorf("invalid ingress rule")
	}

	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	r := &IngressRule{
		Scheme: u.Scheme,
		Host:   u.Hostname(),
		Path:   u.Path,
	}

	if r.Scheme == "" {
		r.Scheme = "http"
	}

	p := u.Port()
	if p == "" {
		r.Port = 80
	} else {
		port, _ := strconv.ParseUint(p, 10, 16)
		r.Port = uint16(port)
	}

	return r, nil
}

type IngressRule struct {
	Scheme string
	Host   string
	Path   string
	Port   uint16
}

func (r IngressRule) String() string {
	if r.Scheme == "" {
		r.Scheme = "http"
	}
	if r.Port == 0 {
		r.Port = 80
	}

	return (&url.URL{
		Scheme: r.Scheme,
		Host:   r.Host + ":" + strconv.FormatUint(uint64(r.Port), 10),
		Path:   r.Path,
	}).String()
}

func (r IngressRule) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *IngressRule) UnmarshalText(data []byte) error {
	ir, err := ParseIngressRule(string(data))
	if err != nil {
		return err
	}
	*r = *ir
	return nil
}