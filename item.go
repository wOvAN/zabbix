package zabbix

import (
	"fmt"
	"strconv"

	"github.com/wOvAN/reflector"
)

type (
	ItemType  interface{}
	ValueType interface{}
	DataType  interface{}
	DeltaType interface{}
)

var (
	ZabbixAgent       ItemType = 0
	SNMPv1Agent       ItemType = 1
	ZabbixTrapper     ItemType = 2
	SimpleCheck       ItemType = 3
	SNMPv2Agent       ItemType = 4
	ZabbixInternal    ItemType = 5
	SNMPv3Agent       ItemType = 6
	ZabbixAgentActive ItemType = 7
	ZabbixAggregate   ItemType = 8
	WebItem           ItemType = 9
	ExternalCheck     ItemType = 10
	DatabaseMonitor   ItemType = 11
	IPMIAgent         ItemType = 12
	SSHAgent          ItemType = 13
	TELNETAgent       ItemType = 14
	Calculated        ItemType = 15
	JMXAgent          ItemType = 16
	SNMPTrap          ItemType = 17
	DependentItem     ItemType = 18

	Float     ValueType = 0
	Character ValueType = 1
	Log       ValueType = 2
	Unsigned  ValueType = 3
	Text      ValueType = 4

	Decimal     DataType = 0
	Octal       DataType = 1
	Hexadecimal DataType = 2
	Boolean     DataType = 3

	AsIs  DeltaType = 0
	Speed DeltaType = 1
	Delta DeltaType = 2

	ItemEnabled  = 0
	ItemDisabled = 1
)

// https://www.zabbix.com/documentation/4.0/manual/api/reference/item/object
type Item struct {
	ItemId      string      `json:"itemid,omitempty"`
	Delay       string      `json:"delay"`
	HostId      string      `json:"hostid"`
	InterfaceId string      `json:"interfaceid,omitempty"`
	Key         string      `json:"key_"`
	Name        string      `json:"name"`
	Type        ItemType    `json:"type"`
	ValueType   ValueType   `json:"value_type"`
	DataType    DataType    `json:"data_type"`
	Delta       DeltaType   `json:"delta"`
	Description string      `json:"description"`
	Error       string      `json:"error"`
	History     string      `json:"history,omitempty"`
	Trends      string      `json:"trends,omitempty"`
	Status      interface{} `json:"status,omitempty"`

	// Fields below used only when creating applications
	ApplicationIds []string `json:"applications,omitempty"`
}

func ItemTypeToText(aItemType ItemType) string {

	v, e := strconv.Atoi(string(aItemType.(string)))
	if e != nil {
		return e.Error()
	}
	switch v {
	case ZabbixAgent:
		return "Zabbix agent"
	case SNMPv1Agent:
		return "SNMPv1 agent"
	case ZabbixTrapper:
		return "Zabbix trapper"
	case SimpleCheck:
		return "simple check"
	case SNMPv2Agent:
		return "SNMPv2 agent"
	case ZabbixInternal:
		return "Zabbix internal"
	case SNMPv3Agent:
		return "SNMPv3 agent"
	case ZabbixAgentActive:
		return "Zabbix agent (active)"
	case ZabbixAggregate:
		return "Zabbix aggregate"
	case WebItem:
		return "Web item"
	case ExternalCheck:
		return "External check"
	case DatabaseMonitor:
		return "Database monitor"
	case IPMIAgent:
		return "IPMI agent"
	case SSHAgent:
		return "SSH agent"
	case TELNETAgent:
		return "TELNET agent"
	case Calculated:
		return "Calculated"
	case JMXAgent:
		return "JMX agent"
	case SNMPTrap:
		return "SNMP trap"
	case DependentItem:
		return "Dependent item"
	default:
		return "Unknown (" + fmt.Sprintf("%s", aItemType) + ")"
	}
}

type Items []Item

// Converts slice to map by key. Panics if there are duplicate keys.
func (items Items) ByKey() (res map[string]Item) {
	res = make(map[string]Item, len(items))
	for _, i := range items {
		_, present := res[i.Key]
		if present {
			panic(fmt.Errorf("Duplicate key %s", i.Key))
		}
		res[i.Key] = i
	}
	return
}

// Wrapper for item.get https://www.zabbix.com/documentation/2.2/manual/appendix/api/item/get
func (api *API) ItemsGet(params Params) (res Items, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("item.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

// Gets items by application Id.
func (api *API) ItemsGetByApplicationId(id string) (res Items, err error) {
	return api.ItemsGet(Params{"applicationids": id})
}

// Wrapper for item.create: https://www.zabbix.com/documentation/2.2/manual/appendix/api/item/create
func (api *API) ItemsCreate(items Items) (err error) {
	response, err := api.CallWithError("item.create", items)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	itemids := result["itemids"].([]interface{})
	for i, id := range itemids {
		items[i].ItemId = id.(string)
	}
	return
}

// Wrapper for item.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/item/delete
// Cleans ItemId in all items elements if call succeed.
func (api *API) ItemsDelete(items Items) (err error) {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ItemId
	}

	err = api.ItemsDeleteByIds(ids)
	if err == nil {
		for i := range items {
			items[i].ItemId = ""
		}
	}
	return
}

// Wrapper for item.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/item/delete
func (api *API) ItemsDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("item.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	itemids1, ok := result["itemids"].([]interface{})
	l := len(itemids1)
	if !ok {
		// some versions actually return map there
		itemids2 := result["itemids"].(map[string]interface{})
		l = len(itemids2)
	}
	if len(ids) != l {
		err = &ExpectedMore{len(ids), l}
	}
	return
}
