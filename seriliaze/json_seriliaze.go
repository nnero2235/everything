package seriliaze

import (
	"encoding/json"
)

type JsonPersistence struct {
}

func (j *JsonPersistence) Seriliaze(obj interface{}) ([]byte, error) {
	s, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	b := []byte(s)
	return b, nil
}

func (j *JsonPersistence) Deriliaze(bytesArr []byte, vPointer interface{}) error {
	err := json.Unmarshal(bytesArr, vPointer)
	return err
}
