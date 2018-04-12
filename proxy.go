// proxys

package zabbix

import (
	"github.com/wOvAN/reflector"
)

type ProxyType int

const (
	ActiveProxy  ProxyType = 5
	PassiveProxy ProxyType = 2
)

// https://www.zabbix.com/documentation/3.2/manual/api/reference/proxy/object
type Proxy struct {
	ProxyId     string    `json:"proxyid,omitempty"`
	Host        string    `json:"host"`
	Description string    `json:"description,omitempty"`
	Status      ProxyType `json:"status ,omitempty"`
}
type Proxys []Proxy

type ProxyId struct {
	ProxyId string `json:"proxyid"`
}

type ProxyIds []ProxyId

// Wrapper for proxy.get: https://www.zabbix.com/documentation/3.2/manual/api/reference/proxy/get
func (api *API) ProxyGet(params Params) (res Proxys, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("proxy.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

// Gets host proxy by Id only if there is exactly 1 matching host proxy.
func (api *API) ProxyGetById(id string) (res *Proxy, err error) {
	proxys, err := api.ProxyGet(Params{"proxyids": id})
	if err != nil {
		return
	}

	if len(proxys) == 1 {
		res = &proxys[0]
	} else {
		e := ExpectedOneResult(len(proxys))
		err = &e
	}
	return
}

// Wrapper for proxy.create: https://www.zabbix.com/documentation/2.2/manual/appendix/api/proxy/create
func (api *API) ProxyCreate(proxys Proxys) (err error) {
	response, err := api.CallWithError("proxy.create", proxys)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	proxyids := result["proxyids"].([]interface{})
	for i, id := range proxyids {
		proxys[i].ProxyId = id.(string)
	}
	return
}

// Wrapper for proxy.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/proxy/delete
// Cleans ProxyId in all proxys elements if call succeed.
func (api *API) ProxyDelete(proxys Proxys) (err error) {
	ids := make([]string, len(proxys))
	for i, proxy := range proxys {
		ids[i] = proxy.ProxyId
	}

	err = api.ProxyDeleteByIds(ids)
	if err == nil {
		for i := range proxys {
			proxys[i].ProxyId = ""
		}
	}
	return
}

// Wrapper for proxy.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/proxy/delete
func (api *API) ProxyDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("proxy.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	proxyids := result["proxyids"].([]interface{})
	if len(ids) != len(proxyids) {
		err = &ExpectedMore{len(ids), len(proxyids)}
	}
	return
}
