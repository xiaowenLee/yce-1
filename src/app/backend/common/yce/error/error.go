package error

import (
	"encoding/json"
	mylog "app/backend/common/util/log"
)

var log =  mylog.Log

type YceError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
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
		log.Errorf("YceError DecodeJson Error: err=%s", err)
		return err
	}

	return nil
}

func (ye *YceError) EncodeJson() (string, error) {
	data, err := json.Marshal(ye)

	if err != nil {
		log.Errorf("YceError EncodeJson Error: err=%s", err)
		return "", err
	}
	return string(data), nil
}

func (ye *YceError) EncodeSelf() []byte {
	errJSON, _ := json.Marshal(ye)
	return errJSON
}


func (ye *YceError) SetError() {


}

