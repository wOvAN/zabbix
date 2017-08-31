package zabbix

import (
	_ "encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/AlekSi/reflector"
)

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	stamp := fmt.Sprint(ts)
	return []byte(stamp), nil
}
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*t = Timestamp(time.Unix(int64(ts), 0))
	return nil
}
func (t *Timestamp) String() string {
	return time.Time(*t).String()
}

type HistoryItem struct {
	ItemId string    `json:"itemid"`
	Clock  Timestamp `json:"clock"`
	Value  float32   `json:"value"`
	Ns     int       `json:"ns"`
}

type HistoryItems []HistoryItem

func (api *API) HistoryGet(params Params) (res HistoryItems, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	if _, presentl := params["limit"]; !presentl {
		params["limit"] = "100"
	}
	if _, presenth := params["history"]; !presenth {
		params["history"] = "0"
	}
	response, err := api.CallWithError("history.get", params)
	if err != nil {
		return
	}
	//fmt.Println(response.Result)

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}
