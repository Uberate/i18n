package i18n

import "encoding/json"

func ToJSON(i *I18n) (string, error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func FromJson(value string) (*I18n, error) {
	res := &I18n{}
	err := json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
