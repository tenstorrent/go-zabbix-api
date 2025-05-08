package zabbix

import (
	"fmt"
	"strconv"
)

// ProxyGroup represents a Zabbix proxy group object.
// See: https://www.zabbix.com/documentation/current/en/manual/api/reference/proxygroup/object
type ProxyGroup struct {
	ProxyGroupID  string `json:"proxy_groupid,omitempty"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`    // Proxy group description
	FailoverDelay string `json:"failover_delay,omitempty"` // Time unit with suffix (e.g., 1m, 30s)
	MinOnline     string `json:"min_online,omitempty"`     // Minimum number of online proxies
	State         string `json:"state,omitempty"`          // State of the proxy group
	Proxies       []struct {
		ProxyID string `json:"proxyid"`
		Name    string `json:"name"`
	} `json:"proxies,omitempty"` // Proxies in this group (when selectProxies is used)
}

// GetMinOnline returns the min_online value as an integer.
// Returns 0 if the value is empty or invalid.
func (p *ProxyGroup) GetMinOnline() int {
	if p.MinOnline == "" {
		return 0
	}
	val, err := strconv.Atoi(p.MinOnline)
	if err != nil {
		return 0
	}
	return val
}

// SetMinOnline sets the min_online value from an integer.
func (p *ProxyGroup) SetMinOnline(val int) {
	p.MinOnline = strconv.Itoa(val)
}

// ProxyGroups is an array of ProxyGroup structs.
type ProxyGroups []ProxyGroup

// ProxyGroupsGet retrieves proxy groups based on the provided parameters.
// See: https://www.zabbix.com/documentation/current/en/manual/api/reference/proxygroup/get
func (api *API) ProxyGroupsGet(params Params) (ProxyGroups, error) {
	// Convert proxy_groupids to numbers if present
	if ids, ok := params["proxy_groupids"]; ok {
		switch v := ids.(type) {
		case string:
			idNum, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("invalid proxy group ID '%s': %w", v, err)
			}
			params["proxy_groupids"] = []int{idNum}
		case []string:
			nums := make([]int, len(v))
			for i, id := range v {
				num, err := strconv.Atoi(id)
				if err != nil {
					return nil, fmt.Errorf("invalid proxy group ID '%s': %w", id, err)
				}
				nums[i] = num
			}
			params["proxy_groupids"] = nums
		case int:
			params["proxy_groupids"] = []int{v}
		case []int:
			// Already in correct format
		default:
			return nil, fmt.Errorf("unsupported proxy_groupids type: %T", ids)
		}
	}

	// Convert proxyids to numbers if present
	if ids, ok := params["proxyids"]; ok {
		switch v := ids.(type) {
		case string:
			idNum, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("invalid proxy ID '%s': %w", v, err)
			}
			params["proxyids"] = []int{idNum}
		case []string:
			nums := make([]int, len(v))
			for i, id := range v {
				num, err := strconv.Atoi(id)
				if err != nil {
					return nil, fmt.Errorf("invalid proxy ID '%s': %w", id, err)
				}
				nums[i] = num
			}
			params["proxyids"] = nums
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

	var result ProxyGroups
	err := api.CallWithErrorParse("proxygroup.get", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ProxyGroupGetByID retrieves a specific proxy group by its ID.
func (api *API) ProxyGroupGetByID(id string) (*ProxyGroup, error) {
	// Convert string ID to number
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy group ID '%s': %w", id, err)
	}

	result, err := api.ProxyGroupsGet(Params{
		"proxy_groupids": []int{idNum}, // Wrap in array as per API docs
		"selectProxies":  "extend",     // Optionally fetch related proxies if needed later
		"output":         "extend",     // Ensure all fields are returned
	})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no proxy group found with id %s", id)
	}
	return &result[0], nil
}

// ProxyGroupGetByName retrieves a specific proxy group by its name.
func (api *API) ProxyGroupGetByName(name string) (*ProxyGroup, error) {
	result, err := api.ProxyGroupsGet(Params{
		"filter": map[string]string{
			"name": name,
		},
		"selectProxies": "extend", // Optionally fetch related proxies
		"output":        "extend", // Ensure all fields are returned
	})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no proxy group found with name %s", name)
	}
	if len(result) > 1 {
		return nil, fmt.Errorf("multiple proxy groups found with name %s", name)
	}
	return &result[0], nil
}

// ProxyGroupsCreate creates new proxy groups.
// See: https://www.zabbix.com/documentation/current/en/manual/api/reference/proxygroup/create
func (api *API) ProxyGroupsCreate(groups ProxyGroups) error {
	response, err := api.Call("proxygroup.create", groups)
	if err != nil {
		return err
	}

	// Update the IDs of the created groups
	result := response.Result.(map[string]interface{})
	proxygroupids := result["proxygroupids"].([]interface{})
	for i, id := range proxygroupids {
		groups[i].ProxyGroupID = id.(string)
	}
	return nil
}

// ProxyGroupsUpdate updates existing proxy groups.
// See: https://www.zabbix.com/documentation/current/en/manual/api/reference/proxygroup/update
func (api *API) ProxyGroupsUpdate(groups ProxyGroups) error {
	response, err := api.CallWithError("proxygroup.update", groups)
	if err != nil {
		return err
	}

	// Check response.Result for expected IDs if necessary
	result := response.Result.(map[string]interface{})
	proxygroupids := result["proxygroupids"].([]interface{})
	if len(proxygroupids) != len(groups) {
		return fmt.Errorf("expected %d proxy groups to be updated, but API reported %d", len(groups), len(proxygroupids))
	}

	return nil
}

// ProxyGroupsDelete deletes proxy groups by their IDs.
// See: https://www.zabbix.com/documentation/current/en/manual/api/reference/proxygroup/delete
func (api *API) ProxyGroupsDelete(ids []string) error {
	// Zabbix API expects an array of IDs for deletion
	params := make([]string, len(ids))
	copy(params, ids)

	response, err := api.Call("proxygroup.delete", params)
	if err != nil {
		// Handle potential API error structure if needed
		// Check if the error indicates that the group doesn't exist and handle gracefully?
		// For now, just return the raw error.
		return err
	}

	// Verify the result if necessary, e.g., check affected group IDs
	result := response.Result.(map[string]interface{})
	proxygroupids := result["proxygroupids"].([]interface{})
	if len(proxygroupids) != len(ids) {
		// Log or return a specific error if not all requested IDs were deleted
		return fmt.Errorf("expected %d proxy groups to be deleted, but API reported %d", len(ids), len(proxygroupids))
	}

	return nil
}
