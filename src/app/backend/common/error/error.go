package error
import (
	"encoding/json"
	"log"
)

type YceError struct {
	Code int32 `json:"code"`
	Message string `json:"message"`
}

func NewYceError(code int32, message string) *YceError {
	return &YceError{
		Code: code,
		Message: message,
	}
}

func (ye *YceError) DecodeJson(data string) error {
	err := json.Unmarshal([]byte(data), ye)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (ye *YceError) EncodeJson() (string, error) {
	data, err := json.MarshalIndent(ye, "", " ")

	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(data), nil
}

/*
const (
	ErrorMap = map[int32][string]{

	}
)
*/