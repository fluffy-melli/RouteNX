package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
)

const (
	RouteNXJSON = "RouteNX.json"
)

type Firewall struct {
	Name  string   `json:"name"`
	CIDR  []string `json:"cidr"`
	Block bool     `json:"block"`
}

type Route struct {
	Host     []string `json:"host"`
	Firewall []string `json:"firewall"`
	Endpoint string   `json:"endpoint"`
}

type RouteNX struct {
	Port     uint16     `json:"port"`
	Routes   []Route    `json:"routes"`
	Firewall []Firewall `json:"firewall"`
}

func NewRouteNX() *RouteNX {
	return &RouteNX{
		Port:     8080,
		Routes:   make([]Route, 0),
		Firewall: make([]Firewall, 0),
	}
}

func (c *RouteNX) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal RouteNX: %v", err)
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	return nil
}

func LoadFromFile(filename string) (*RouteNX, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	var routeNX RouteNX
	err = json.Unmarshal(data, &routeNX)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return &routeNX, nil
}

func (c *RouteNX) GetRoute(host string) *Route {
	for _, route := range c.Routes {
		for _, hostL := range route.Host {
			if hostL == host {
				return &route
			}
			if matched, _ := path.Match(hostL, host); matched {
				return &route
			}
		}
	}
	return nil
}

func (c *RouteNX) IsBlock(route *Route, ip string) bool {
	for _, firewall := range c.Firewall {
		for _, RoutefirewallL := range route.Firewall {
			if firewall.Name == RoutefirewallL {
				for _, blockIP := range firewall.CIDR {
					_, cidr, _ := net.ParseCIDR(blockIP)
					ipr := net.ParseIP(ip)
					if firewall.Block && cidr.Contains(ipr) {
						return true
					}
					if !firewall.Block && cidr.Contains(ipr) {
						return false
					}
				}
				if !firewall.Block {
					return true
				}
			}
		}
	}
	return false
}
