package tool

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetBody(req *http.Request) map[string]interface{} {
	data := make(map[string]interface{})
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return data
	}

	err = json.Unmarshal(body, &data)

	if err != nil {
		return data
	}

	return data
}
