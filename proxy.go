package zabbix

// Proxy represent Zabbix proxy object
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/object
type Proxy struct {
	ProxyID              string `json:"proxyid,omitempty"`
	Host                 string `json:"host"`
	Status               int    `json:"status,omitempty"`
	Description          string `json:"description,omitempty"`
	LastAccess           int    `json:"lastaccess,omitempty"`
	TLSConnect           int    `json:"tls_connect,omitempty"`
	TLSAccept            int    `json:"tls_accept,omitempty"`
	TLSIssuer            string `json:"tls_issuer,omitempty"`
	TLSSubject           string `json:"tls_subject,omitempty"`
	TLSPSKIdentity       string `json:"tls_psk_identity,omitempty"`
	TLSPSK               string `json:"tls_psk,omitempty"`
	CustomTimeouts       int    `json:"custom_timeouts,omitempty"`
	TimeoutZabbixAgent   string `json:"timeout_zabbix_agent,omitempty"`
	TimeoutSimpleCheck   string `json:"timeout_simple_check,omitempty"`
	TimeoutSNMPAgent     string `json:"timeout_snmp_agent,omitempty"`
	TimeoutExternalCheck string `json:"timeout_external_check,omitempty"`
	TimeoutDBMonitor     string `json:"timeout_db_monitor,omitempty"`
	TimeoutHTTPAgent     string `json:"timeout_http_agent,omitempty"`
	TimeoutSSHAgent      string `json:"timeout_ssh_agent,omitempty"`
	TimeoutTelnetAgent   string `json:"timeout_telnet_agent,omitempty"`
	TimeoutScript        string `json:"timeout_script,omitempty"`
	TimeoutBrowser       string `json:"timeout_browser,omitempty"`
	Version              int    `json:"version,omitempty"`
	Compatibility        int    `json:"compatibility,omitempty"`
	State                int    `json:"state,omitempty"`
	ProxyGroupID         string `json:"proxy_groupid,omitempty"`
}

// Proxies is an array of Proxy
type Proxies []Proxy

// ProxiesGet Wrapper for proxy.get
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/get
func (api *API) ProxiesGet(params Params) (res Proxies, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("proxy.get", params, &res)
	return
}

// ProxiesCreate Wrapper for proxy.create
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/create
func (api *API) ProxiesCreate(proxies Proxies) (err error) {
	response, err := api.CallWithError("proxy.create", proxies)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	proxyids := result["proxyids"].([]interface{})
	for i, id := range proxyids {
		proxies[i].ProxyID = id.(string)
	}
	return
}

// ProxiesUpdate Wrapper for proxy.update
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/update
func (api *API) ProxiesUpdate(proxies Proxies) (err error) {
	_, err = api.CallWithError("proxy.update", proxies)
	return
}

// ProxiesDelete Wrapper for proxy.delete
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/delete
func (api *API) ProxiesDelete(proxyids []string) (err error) {
	_, err = api.CallWithError("proxy.delete", proxyids)
	return
}
