package zabbix

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Proxy represent Zabbix proxy object
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/object
type Proxy struct {
	ProxyID              string `json:"proxyid,omitempty"`
	Host                 string `json:"host"`
	OperatingMode        string `json:"status,omitempty"` // Operating mode: 0 = active, 1 = passive
	Description          string `json:"description,omitempty"`
	LastAccess           string `json:"lastaccess,omitempty"`
	TLSConnect           string `json:"tls_connect,omitempty"`
	TLSAccept            string `json:"tls_accept,omitempty"`
	TLSIssuer            string `json:"tls_issuer,omitempty"`
	TLSSubject           string `json:"tls_subject,omitempty"`
	TLSPSKIdentity       string `json:"tls_psk_identity,omitempty"`
	TLSPSK               string `json:"tls_psk,omitempty"`
	CustomTimeouts       string `json:"custom_timeouts,omitempty"`
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
	Version              string `json:"version,omitempty"`
	Compatibility        string `json:"compatibility,omitempty"`
	State                string `json:"state,omitempty"`
	ProxyGroupID         string `json:"proxy_groupid,omitempty"`
	LocalAddress         string `json:"local_address,omitempty"` // IP address or hostname for active proxy
	LocalPort            string `json:"local_port,omitempty"`    // Port for active proxy
	AllowedAddresses     string `json:"allowed_addresses,omitempty"`
	Address              string `json:"address,omitempty"`
	Port                 string `json:"port,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling to handle string-to-int conversion
func (p *Proxy) UnmarshalJSON(data []byte) error {
	// First unmarshal into a map to get all fields as strings
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Create a temporary struct to hold string values
	type Alias Proxy
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	// Unmarshal into the temporary struct
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Use reflection to convert string fields to int
	val := reflect.ValueOf(p).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Int {
			// Get the JSON tag name
			jsonTag := typ.Field(i).Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				continue
			}
			jsonTag = strings.Split(jsonTag, ",")[0] // Remove omitempty if present

			// If we have a string value in the raw map, convert it to int
			if strVal, ok := raw[jsonTag].(string); ok {
				if num, err := strconv.Atoi(strVal); err == nil {
					field.SetInt(int64(num))
				}
			}
		}
	}

	return nil
}

// Proxies is an array of Proxy
type Proxies []Proxy

// Int returns the integer value of a string field, or 0 if invalid
func (p *Proxy) Int(field string) int {
	val := reflect.ValueOf(p).Elem().FieldByName(field).String()
	if val == "" {
		return 0
	}
	num, _ := strconv.Atoi(val)
	return num
}

// SetInt sets a string field from an integer value
func (p *Proxy) SetInt(field string, value int) {
	reflect.ValueOf(p).Elem().FieldByName(field).SetString(strconv.Itoa(value))
}

// ProxiesGet Wrapper for proxy.get
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/get
func (api *API) ProxiesGet(params Params) (res Proxies, err error) {
	// Convert proxyids to numbers if present
	if ids, ok := params["proxyids"]; ok {
		switch v := ids.(type) {
		case string:
			// Try to convert to number first
			if idNum, err := strconv.Atoi(v); err == nil {
				params["proxyids"] = []int{idNum}
			} else {
				// If not a number, use as name filter
				delete(params, "proxyids")
				if _, hasFilter := params["filter"]; !hasFilter {
					params["filter"] = make(map[string]interface{})
				}
				filter := params["filter"].(map[string]interface{})
				filter["name"] = v
			}
		case []string:
			// Check if all strings are numbers
			allNumeric := true
			nums := make([]int, 0, len(v))

			for _, id := range v {
				num, err := strconv.Atoi(id)
				if err != nil {
					allNumeric = false
					break
				}
				nums = append(nums, num)
			}

			if allNumeric {
				params["proxyids"] = nums
			} else {
				// If any are not numeric, use as name filter
				delete(params, "proxyids")
				if _, hasFilter := params["filter"]; !hasFilter {
					params["filter"] = make(map[string]interface{})
				}
				filter := params["filter"].(map[string]interface{})
				filter["name"] = v
			}
		case int:
			params["proxyids"] = []int{v}
		case []int:
			// Already in correct format
		default:
			return nil, fmt.Errorf("unsupported proxyids type: %T", ids)
		}
	}

	// Set output to extend if not specified
	if _, ok := params["output"]; !ok {
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
