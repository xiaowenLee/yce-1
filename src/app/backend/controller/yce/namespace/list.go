package namespace

import (
	//myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	myorganization "app/backend/model/mysql/organization"
)

type ListNamespaceController struct {
	yce.Controller
}

type NamespaceList struct {

}

func (lnc *ListNamespaceController) getNamespaceList() {
	// get Namespace
	myorganization.QueryAllOrganizations()
}

func (lnc ListNamespaceController) Get() {
	//TODO: rethink of session authroization. Here it is omitted.
	//SessionIdFromClient := iuc.RequestHeader("Authorization")

	// getNamespaceList()
}
