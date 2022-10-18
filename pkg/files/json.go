package files

import (
	"encoding/json"
	"github.com/uberate/i18n/pkg/provider"
	"os"
)

func WriteToJSONFile(file string, instance *provider.I18n) error {
	value, err := ToJSON(instance)
	if err != nil {
		return err
	}

	return os.WriteFile(file, []byte(value), os.ModePerm)
}

func ReadFromJSONFile(file string) (*provider.I18n, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return FromJson(string(fileBytes))
}

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
