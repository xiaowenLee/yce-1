package namespace

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	mydatacenter "app/backend/model/mysql/datacenter"
	yceutils "app/backend/controller/yce/utils"
	"encoding/json"
)

type InitNamespaceController struct {
	yce.Controller

	params *InitNamespaceParams
}


type InitNamespaceParams struct {
	DcList []mydatacenter.DataCenter `json:"dcList"`
	Account AccountType `json:"account"`
	QuotaPkg []QuotaPkgType `json:"quotaPkg"`
}

func (inc InitNamespaceController) getDatacenters() {
	inc.params.DcList, inc.Ye = yceutils.QueryAllDatacenters()
	if inc.Ye != nil {
		return
	}
}

func (inc InitNamespaceController) getAccount() {
	acc := &AccountType{
		Budget: 1000.00,
		Balance: 1000.00,
	}

	inc.params.Account = *acc
}

func (inc InitNamespaceController) getQuotaPkg() {
	quotaPkgList := make([]QuotaPkgType, 0)
	quotaPkg1 := &QuotaPkgType{
		Cpu: 200,
		Mem: 400,
		Cost: 200.00,
	}

	quotaPkg2 := &QuotaPkgType{
		Cpu: 500,
		Mem: 1000,
		Cost: 500.00,
	}

	quotaPkg3 := &QuotaPkgType{
		Cpu: 1000,
		Mem: 2000,
		Cost: 1000.00,
	}

	quotaPkgList = append(quotaPkgList, *quotaPkg1)
	quotaPkgList = append(quotaPkgList, *quotaPkg2)
	quotaPkgList = append(quotaPkgList, *quotaPkg3)

	inc.params.QuotaPkg = quotaPkgList
}


func (inc *InitNamespaceController) prepare() string {
	inc.getDatacenters()
	inc.getAccount()
	inc.getQuotaPkg()

	resultJSON, err := json.Marshal(inc.params)
	if err != nil {
		inc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}

	resultString := string(resultJSON)
	return resultString
}

func (inc InitNamespaceController) Get() {

	//TODO: admin authorication

	inc.params = new(InitNamespaceParams)
	if inc.CheckError() {
		return
	}

	result := inc.prepare()
	if inc.CheckError() {
		return
	}

	inc.WriteOk(result)
	return
}
