// scripts

package zabbix

import (
	"github.com/wOvAN/reflector"
)

type (
	ScriptType     interface{}
	ExecuteOnType  interface{}
	HostAccessType interface{}
)

var (
	// ScriptType
	ScriptTypeScript ScriptType = 0
	ScriptTypeIPMI   ScriptType = 1

	// ExecuteOnType
	ExecuteOnAgent  ExecuteOnType = 0
	ExecuteOnServer ExecuteOnType = 1
	ExecuteOnProxy  ExecuteOnType = 2 // (default) run on Zabbix server (proxy).

	// HostAccessType
	HostAccessRO HostAccessType = 2
	HostAccessRW HostAccessType = 3
)

// https://www.zabbix.com/documentation/4.0/manual/api/reference/script/object
type Script struct {
	ScriptId string `json:"scriptid,omitempty"`
	Command  string `json:"command"`
	Name     string `json:"name"`

	// строка 	Текст подтверждения во всплывающем окне.
	// Всплывающее окно появляется при попытке выполнения скрипта из Zabbix веб-интерфейса.
	Confirmation string `json:"confirmation,omitempty"`
	// строка 	Описание скрипта.
	Description string `json:"description,omitempty"`

	/* целое число 	Где выполнять скрипт.
	   Возможные значения:
		0 - выполнение на Zabbix агенте;
		1 - (по умолчанию) выполнение на Zabbix сервере.
	*/
	ExecuteOn ExecuteOnType `json:"execute_on,omitempty"`

	// строка 	ID группы узлов сети для которой можно выполнять скрипт.
	// Если задано значение 0, скрипт можно выполнять по всем группам узлов сети.
	// По умолчанию: 0.
	GroupId string `json:"groupid,omitempty"`

	// целое число 	Требуемые права доступа к узлу сети для выполнения скрипта.
	// Возможные значения:
	// 2 - (по умолчанию) чтение;
	// 3 - запись.

	HostAccess HostAccessType `json:"host_access,omitempty"`

	// type 	целое число 	Тип скрипта.
	// Возможные значения:
	// 0 - (по умолчанию) скрипт;
	// 1 - IPMI.
	Type ScriptType `json:"type,omitempty"`

	// строка 	ID группы пользователей, которой разрешено выполнение скрипта.
	// Если задано значение 0, скрипт доступен всем группам пользователей.
	// По умолчанию: 0.
	UsrGrpId string `json:"usrgrpid,omitempty"`
}
type (
	Scripts []Script

	ScriptId struct {
		ScriptId string `json:"scriptid"`
	}
)

type ScriptIds []ScriptId

// Wrapper for script.get: https://www.zabbix.com/documentation/3.2/manual/api/reference/script/get
func (api *API) ScriptGet(params Params) (res Scripts, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("script.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

// Gets host script by Id only if there is exactly 1 matching host script.
func (api *API) ScriptGetById(id string) (res *Script, err error) {
	scripts, err := api.ScriptGet(Params{"scriptids": id})
	if err != nil {
		return
	}

	if len(scripts) == 1 {
		res = &scripts[0]
	} else {
		e := ExpectedOneResult(len(scripts))
		err = &e
	}
	return
}

// Wrapper for script.create: https://www.zabbix.com/documentation/4.0/manual/api/reference/script/create
func (api *API) ScriptCreate(scripts Scripts) (err error) {
	response, err := api.CallWithError("script.create", scripts)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	scriptids := result["scriptids"].([]interface{})
	for i, id := range scriptids {
		scripts[i].ScriptId = id.(string)
	}
	return
}

// Wrapper for script.execute: https://www.zabbix.com/documentation/4.0/manual/api/reference/script/execute
const (
	ScriptExecRespSuccess = "success"
	ScriptExecRespFailed  = "failed"
)

//
func (api *API) ScriptExecute(aScriptId string, aHostId string) (rResponse string, rOutput string, err error) {
	params := Params{}
	params["hostid"] = aHostId
	params["scriptid"] = aScriptId

	response, err := api.CallWithError("script.execute", params)
	if err != nil {
		return
	}
	result := response.Result.(map[string]interface{})
	rResponse = result["response"].(string)
	rOutput = result["value"].(string)
	return
}

// Wrapper for script.delete: https://www.zabbix.com/documentation/4.0/manual/api/reference/script/delete
// Cleans ScriptId in all scripts elements if call succeed.
func (api *API) ScriptDelete(scripts Scripts) (err error) {
	ids := make([]string, len(scripts))
	for i, script := range scripts {
		ids[i] = script.ScriptId
	}

	err = api.ScriptDeleteByIds(ids)
	if err == nil {
		for i := range scripts {
			scripts[i].ScriptId = ""
		}
	}
	return
}

// Wrapper for script.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/script/delete
func (api *API) ScriptDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("script.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	scriptids := result["scriptids"].([]interface{})
	if len(ids) != len(scriptids) {
		err = &ExpectedMore{len(ids), len(scriptids)}
	}
	return
}
