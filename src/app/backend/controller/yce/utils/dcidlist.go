package utils

import (
	myerror "app/backend/common/yce/error"
	"encoding/json"
)

type DcIdListType struct {
	DcIdList []int32 `json:"dcIdList"`
}


// Encode DcIdList []int32 to dcIdList Json []string form
func StringDcIdList(dcIdList []int32) (string, myerror.YceError) {
	dcIdListByte, err := EncodeDcIdList(dcIdList)
	if err != nil {
		return "", err
	} else {
		dcIdListString := string(dcIdListByte)
		return dcIdListString, nil
	}
}

// Encode dcIdList []int32 to dcIdList Json []byte form
func EncodeDcIdList(dcIdList []int32) ([]byte, myerror.YceError) {
	dcList := &DcIdListType {
		DcIdList: dcIdList,
	}
	dcIdListByte, err := json.Marshal(dcList)

	if err != nil {
		ye := myerror.NewYceError(myerror.EJSON, "")
		return nil, ye
	} else {
		return dcIdListByte, nil
	}

}


// Decode dcIdList from Json string form to []int32
func DecodeDcIdList(dcIdListString string) ([]int32, myerror.YceError) {
	dcIdListByte := []byte(dcIdListString)
	dcIdList := new(DcIdListType)

	err := json.Unmarshal(dcIdListByte, dcIdList)

	if err != nil {
		ye := myerror.NewYceError(myerror.EJSON, "")
		return nil, nil
	} else {
		return dcIdList.DcIdList, nil
	}
}


