// triggers

package zabbix

import (
	"fmt"

	"strconv"

	"github.com/wOvAN/reflector"
)

const (
	// Trigger Selectors
	/*
		selectGroups 	query 	Return the host groups that the trigger belongs to in the groups property.
		selectHosts 	query 	Return the hosts that the trigger belongs to in the hosts property.
		selectItems 	query 	Return items contained by the trigger in the items property.
		selectFunctions 	query 	Return functions used in the trigger in the functions property.
			The function objects represents the functions used in the trigger expression and has the following properties:
			functionid - (string) ID of the function;
			itemid - (string) ID of the item used in the function;
			function - (string) name of the function;
			parameter - (string) parameter passed to the function.
		selectDependencies 	query 	Return triggers that the trigger depends on in the dependencies property.
		selectDiscoveryRule 	query 	Return the low-level discovery rule that created the trigger.
		selectLastEvent 	query 	Return the last significant trigger event in the lastEvent property.
		selectTags 	query 	Return the trigger tags in tags property.
	*/

	// Trigget expandors
	expandComment     = "expandComment"
	expandDescription = "expandDescription"
	expandExpression  = "expandExpression"
)

// https://www.zabbix.com/documentation/4.0/manual/api/reference/trigger/object
type (
	TriggerPriority interface{}
	TriggerStatus   interface{}

	Trigger struct {
		//string `json:",omitempty"`
		TriggerId   string          `json:"triggerid,omitempty"`
		Description string          `json:"description"`
		Expression  string          `json:"expression"`
		Comments    string          `json:"comments,omitempty"`
		Priority    TriggerPriority `json:"priority,omitempty"`
		Status      TriggerStatus   `json:"status,omitempty"`

		//	Recovery_mode       recovery_mode
		Recovery_expression string `json:"recovery_expression,omitempty"`
		//	Correlation_mode    correlation_mode
		Correlation_tag string `json:"correlation_tag,omitempty"`
		// Extended fields

	}

	Triggers []Trigger

	TriggerId struct {
		TriggerId string `json:"triggerid"`
	}

	TriggerIds []TriggerId
)

var (
	// Priorities
	TriggerPriorityDefault     TriggerPriority = 0
	TriggerPriorityInformation TriggerPriority = 1
	TriggerPriorityWarning     TriggerPriority = 2
	TriggerPriorityAverage     TriggerPriority = 3
	TriggerPriorityHigh        TriggerPriority = 4
	TriggerPriorityDisaster    TriggerPriority = 5
	// Status
	TriggerStatusEnabled  TriggerStatus = 0
	TriggerStatusDisabled TriggerStatus = 0
)

func TriggerPriorityToText(aTriggerPriority TriggerPriority) string {

	v, e := strconv.Atoi(string(aTriggerPriority.(string)))
	if e != nil {
		return e.Error()
	}

	switch v {
	case TriggerPriorityDefault:
		return "Default"
	case TriggerPriorityInformation:
		return "Information"
	case TriggerPriorityWarning:
		return "Warning"
	case TriggerPriorityAverage:
		return "Average"
	case TriggerPriorityHigh:
		return "High"
	case TriggerPriorityDisaster:
		return "Disaster"
	default:
		return "Unknown (" + fmt.Sprintf("%s", aTriggerPriority) + ")"
	}
}

func TriggerStatusToText(aTriggerStatus TriggerStatus) string {
	v, e := strconv.Atoi(string(aTriggerStatus.(string)))
	if e != nil {
		return e.Error()
	}
	switch v {
	case TriggerStatusEnabled:
		return "Enabled"
	case TriggerStatusDisabled:
		return "Disabled"
	default:
		return "Unknown (" + fmt.Sprintf("%s", aTriggerStatus) + ")"
	}
}

// Wrapper for trigger.get: https://www.zabbix.com/documentation/4.0/manual/api/reference/trigger/get
func (api *API) TriggersGet(params Params) (res Triggers, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("trigger.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

// Gets host trigger by Id only if there is exactly 1 matching host trigger.
func (api *API) TriggerGetById(id string) (res *Trigger, err error) {
	triggers, err := api.TriggersGet(Params{"triggerids": id})
	if err != nil {
		return
	}

	if len(triggers) == 1 {
		res = &triggers[0]
	} else {
		e := ExpectedOneResult(len(triggers))
		err = &e
	}
	return
}

// Wrapper for trigger.create: https://www.zabbix.com/documentation/2.2/manual/appendix/api/trigger/create
func (api *API) TriggersCreate(triggers Triggers) (err error) {
	response, err := api.CallWithError("trigger.create", triggers)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	triggerids := result["triggerids"].([]interface{})
	for i, id := range triggerids {
		triggers[i].TriggerId = id.(string)
	}
	return
}

// Wrapper for trigger.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/trigger/delete
// Cleans TriggerId in all triggers elements if call succeed.
func (api *API) TriggersDelete(triggers Triggers) (err error) {
	ids := make([]string, len(triggers))
	for i, trigger := range triggers {
		ids[i] = trigger.TriggerId
	}

	err = api.TriggersDeleteByIds(ids)
	if err == nil {
		for i := range triggers {
			triggers[i].TriggerId = ""
		}
	}
	return
}

// Wrapper for trigger.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/trigger/delete
func (api *API) TriggersDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("trigger.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	triggerids := result["triggerids"].([]interface{})
	if len(ids) != len(triggerids) {
		err = &ExpectedMore{len(ids), len(triggerids)}
	}
	return
}
