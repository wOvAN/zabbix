// templates

package zabbix

import (
	"github.com/wOvAN/reflector"
)

// https://www.zabbix.com/documentation/3.2/manual/api/reference/template/object
type Template struct {
	TemplateId  string `json:"templateid,omitempty"`
	Host        string `json:"host"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
}
type Templates []Template

type TemplateId struct {
	TemplateId string `json:"templateid"`
}

type TemplateIds []TemplateId

// Wrapper for template.get: https://www.zabbix.com/documentation/3.2/manual/api/reference/template/get
func (api *API) TemplatesGet(params Params) (res Templates, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("template.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

// Gets host template by Id only if there is exactly 1 matching host template.
func (api *API) TemplateGetById(id string) (res *Template, err error) {
	templates, err := api.TemplatesGet(Params{"templateids": id})
	if err != nil {
		return
	}

	if len(templates) == 1 {
		res = &templates[0]
	} else {
		e := ExpectedOneResult(len(templates))
		err = &e
	}
	return
}

// Wrapper for template.create: https://www.zabbix.com/documentation/2.2/manual/appendix/api/template/create
func (api *API) TemplatesCreate(templates Templates) (err error) {
	response, err := api.CallWithError("template.create", templates)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	templateids := result["templateids"].([]interface{})
	for i, id := range templateids {
		templates[i].TemplateId = id.(string)
	}
	return
}

// Wrapper for template.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/template/delete
// Cleans TemplateId in all templates elements if call succeed.
func (api *API) TemplatesDelete(templates Templates) (err error) {
	ids := make([]string, len(templates))
	for i, template := range templates {
		ids[i] = template.TemplateId
	}

	err = api.TemplatesDeleteByIds(ids)
	if err == nil {
		for i := range templates {
			templates[i].TemplateId = ""
		}
	}
	return
}

// Wrapper for template.delete: https://www.zabbix.com/documentation/2.2/manual/appendix/api/template/delete
func (api *API) TemplatesDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("template.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	templateids := result["templateids"].([]interface{})
	if len(ids) != len(templateids) {
		err = &ExpectedMore{len(ids), len(templateids)}
	}
	return
}
