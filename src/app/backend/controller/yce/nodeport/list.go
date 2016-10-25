package nodeport

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	mynodeport "app/backend/model/mysql/nodeport"
	yceutils "app/backend/controller/yce/utils"
	"encoding/json"
)

type ListNodePortController struct {
	yce.Controller

	params *NodePortList
}

type NodePortList struct {
	NodePorts []mynodeport.NodePort `json:"nodePorts"`
	DcList 	     map[int32]string `json:"dcList"`
}


func (lnc *ListNodePortController) getDcList() {
	datacenters, ye := yceutils.QueryAllDatacenters()
	if ye != nil {
		lnc.Ye = ye
		return
	}

	for _, dc := range datacenters {
		lnc.params.DcList[dc.Id] = dc.Name
	}
}

func (lnc *ListNodePortController) getNodePortList() string {
	nodeports, err := mynodeport.QueryAllInvalidNodePort()
	if err != nil {
		lnc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	lnc.getDcList()
	if lnc.Ye != nil {
		return ""
	}

	lnc.params.NodePorts = nodeports

	npListJSON, err := json.Marshal(lnc.params)
	if err != nil {
		lnc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}

	npListString := string(npListJSON)
	return npListString
}


func (lnc ListNodePortController) Get() {
	//TODO: rethink of session authroization. Here it is omitted.
	//SessionIdFromClient := iuc.RequestHeader("Authorization")

	lnc.params = new(NodePortList)
	lnc.params.DcList = make(map[int32]string)

	orgList := lnc.getNodePortList()
	if lnc.CheckError() {
		return
	}

	lnc.WriteOk(orgList)
	log.Infoln("ListNodePortController Get Over!")
}