package pkg

import (
	"encoding/json"
	"github.com/uberate/i18n/pkg/provider"
)

func ToJSON(i *provider.I18n) (string, error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func FromJson(value string) (*provider.I18n, error) {
	res := &provider.I18n{}
	err := json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
