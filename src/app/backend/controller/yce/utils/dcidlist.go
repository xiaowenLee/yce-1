package utils

import (
	mylog "app/backend/common/util/log"
	myerror "app/backend/common/yce/error"
	"encoding/json"
)

type DcIdListType struct {
	DcIdList []int32 `json:"dcIdList"`
}

// Encode DcIdList []int32 to dcIdList Json []string form
func StringDcIdList(dcIdList []int32) (string, *myerror.YceError) {

	if CheckValidate(dcIdList) {
		dcIdListByte, err := EncodeDcIdList(dcIdList)
		if err != nil {
			mylog.Log.Errorf("StringDcIdList Error: error=%s", err)
			return "", err
		} else {
			dcIdListString := string(dcIdListByte)
			return dcIdListString, nil
		}
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		mylog.Log.Errorf("StringDcIdList Error: error=%s", ye.Message)

		return "", ye
	}

}

// Encode dcIdList []int32 to dcIdList Json []byte form
func EncodeDcIdList(dcIdList []int32) (string, *myerror.YceError) {

	if CheckValidate(dcIdList) {
		dcList := &DcIdListType{
			DcIdList: dcIdList,
		}
		dcIdListByte, err := json.Marshal(dcList)

		if err != nil {
			ye := myerror.NewYceError(myerror.EJSON, "")
			mylog.Log.Errorf("EncodeDcIdList Error: error=%s", ye.Message)
			return "", ye
		} else {
			return string(dcIdListByte), nil
		}
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		mylog.Log.Errorf("EncodeDcIdList Error: error=%s", ye.Message)

		return "", ye
	}

}

// Decode dcIdList from Json string form to []int32
func DecodeDcIdList(dcIdListString string) ([]int32, *myerror.YceError) {
	if CheckValidate(dcIdListString) {
		dcIdListByte := []byte(dcIdListString)
		dcIdList := new(DcIdListType)

		err := json.Unmarshal(dcIdListByte, dcIdList)

		if err != nil {
			ye := myerror.NewYceError(myerror.EJSON, "")
			mylog.Log.Errorf("DecodeDcIdList Error: error=%s", ye.Message)
			return nil, ye
		} else {
			return dcIdList.DcIdList, nil
		}
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		mylog.Log.Errorf("DecodeDcIdList Error: error=%s", ye.Message)
		return nil, ye
	}

}
