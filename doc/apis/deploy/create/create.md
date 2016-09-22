应用发布
-----------

### 发布准备

点击发布应用时用户先请求: GET /api/v1/organizations/{orgId}/users/{userId}/deployments/init

得到关于用户所有数据中心的信息、预算、配额信息、用户所属组织名称等
 

返回信息的data对应的结构体如下:

type InitDeployment struct {
	OrgId string `json:"orgId"`
	OrgName string `json:"orgName"`
	DataCenters []mydatacenter.DataCenter `json:"dataCenters"`
	Quotas []myqouta.Quota `json:"quotas"`
}

相应的InitDeploymentController控制器结构体为:

type InitDeployController struct {
	*iris.Context
	org  *myorganization.Organization
	Init deploy.InitDeployment
	Ye *myerror.YceError
}

其中的Init即为InitDeployment的对象

可以改为:


type InitDeploymentController struct {
    // basic elements
    *iris.Context
    Ye *myerror.YceError
    
    // Params
    org *myorganization.Organization
    
    // response 
    Init deploy.InitDeployment 
    
}






### 提交发布

然后用户填好发布信息后点击发布请求: POST /api/v1/organization/{orgId}/users/{userId}/deployments

请求时提交的结构体为:

type CreateDeployment struct {
	AppName  string `json:"appName"`
	OrgName  string `json:"orgName"`
	DcIdList []int32 `json:"dcIdList"`
	//DcIdList DcIdListType `json:"dcIdList"`
	Deployment extensions.Deployment `json:"deployment"`
}

相应的CreateDeploymentController控制器结构体为:

type CreateDeployController struct {
	*iris.Context
	k8sClients []*client.Client
	apiServers []string
	Ye         *myerror.YceError
}

可以改为:

type CreateDeploymentController struct {
    // basic elements
    *iris.Context
    Ye *myerror.YceError
    
    // Params
    ...
    
    // DAO
    deploy mydeployment.deployment
}

