package parking

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// GetFeeModels reads an io.Reader and parses it to FeeModels
func GetFeeModels(r io.Reader) (FeeModels, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	models := []FeeModel{}
	err = json.Unmarshal(b, &models)
	if err != nil {
		return nil, err
	}
	var m FeeModels = make(FeeModels)
	for _, model := range models {
		m[model.Model] = model
	}
	return m, nil
}
