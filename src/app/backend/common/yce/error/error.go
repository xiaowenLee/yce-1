package error

import (
	"encoding/json"
	"log"
)

type YceError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func NewYceError(code int32, message, data string) *YceError {
	return &YceError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func (ye *YceError) DecodeJson(data string) error {
	err := json.Unmarshal([]byte(data), ye)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (ye *YceError) EncodeJson() (string, error) {
	data, err := json.Marshal(ye)

	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(data), nil
}
