package testing

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/runtime"
)

func CreateRawExtension(data map[string]interface{}) *runtime.RawExtension {
	rawData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	object := runtime.RawExtension{Raw: rawData}
	return &object
}
