package zabbix_test

import (
	"fmt"
	"math/rand"
	"testing"

	zapi "github.com/tpretz/go-zabbix-api"
)

func CreateProxy(t *testing.T) *zapi.Proxy {
	name := fmt.Sprintf("%s-%d", "test-proxy", rand.Int())
	proxy := zapi.Proxy{
		Host:          name,
		OperatingMode: "0", // Active proxy
		Description:   "Test proxy created by go-zabbix-api tests",
	}

	proxies := zapi.Proxies{proxy}
	err := getAPI(t).ProxiesCreate(proxies)
	if err != nil {
		t.Fatal(err)
	}
	return &proxies[0]
}

func DeleteProxy(proxy *zapi.Proxy, t *testing.T) {
	err := getAPI(t).ProxiesDelete([]string{proxy.ProxyID})
	if err != nil {
		t.Fatal(err)
	}
}

func TestProxies(t *testing.T) {
	api := getAPI(t)

	// Test creating a proxy
	proxy := CreateProxy(t)
	defer DeleteProxy(proxy, t)

	// Test getting proxy by ID
	params := zapi.Params{
		"proxyids": []string{proxy.ProxyID},
	}
	proxies, err := api.ProxiesGet(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(proxies) != 1 {
		t.Errorf("Expected 1 proxy, got %d", len(proxies))
	}
	if proxies[0].Host != proxy.Host {
		t.Errorf("Expected host %s, got %s", proxy.Host, proxies[0].Host)
	}

	// Test updating proxy
	proxy.Description = "Updated description"
	proxies = zapi.Proxies{*proxy}
	err = api.ProxiesUpdate(proxies)
	if err != nil {
		t.Fatal(err)
	}

	// Verify update
	proxies, err = api.ProxiesGet(params)
	if err != nil {
		t.Fatal(err)
	}
	if proxies[0].Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got '%s'", proxies[0].Description)
	}
}
