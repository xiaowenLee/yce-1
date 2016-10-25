package utils

import (
	mylog "app/backend/common/util/log"
	myerror "app/backend/common/yce/error"
	"encoding/json"
)

type NodePortType struct {
	NodePort []string  `json:"nodePort"`
}

// Encode nodePort []int32 to nodePort Json []byte form
func EncodeNodePort(nodePort []string) (string, *myerror.YceError) {

	if CheckValidate(nodePort) {
		np := &NodePortType{
			NodePort: nodePort,
		}
		nodePortByte, err := json.Marshal(np)

		if err != nil {
			ye := myerror.NewYceError(myerror.EJSON, "")
			mylog.Log.Errorf("EncodeNodePort Error: error=%s", ye.Message)
			return "", ye
		} else {
			return string(nodePortByte), nil
		}
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		mylog.Log.Errorf("EncodeNodePort Error: error=%s", ye.Message)

		return "", ye
	}

}

// Decode nodePort from Json string form to []string
func DecodeNodePort(nodePortString string) ([]string, *myerror.YceError) {
	if CheckValidate(nodePortString) {
		nodePortByte := []byte(nodePortString)
		nodePort := new(NodePortType)

		err := json.Unmarshal(nodePortByte, nodePort)

		if err != nil {
			ye := myerror.NewYceError(myerror.EJSON, "")
			mylog.Log.Errorf("DecodeNodePort Error: error=%s", ye.Message)
			return nil, ye
		} else {
			return nodePort.NodePort, nil
		}
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		mylog.Log.Errorf("DecodeNodePort Error: error=%s", ye.Message)
		return nil, ye
	}

}

