package utils

import (
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	myorganization "app/backend/model/mysql/organization"
	myqouta "app/backend/model/mysql/quota"
	myuser "app/backend/model/mysql/user"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	deploymentutil "k8s.io/kubernetes/pkg/controller/deployment/util"
	"reflect"
	"strconv"
)

// Create K8s Client List by ApiServerList
func CreateK8sClientList(apiServerList []string) ([]*client.Client, *myerror.YceError) {
	k8sClientList := make([]*client.Client, 0)

	if CheckValidate(apiServerList) {
		for _, apiServer := range apiServerList {
			k8sClient, err := CreateK8sClient(apiServer)
			if err != nil {
				ye := err
				log.Errorf("CreateK8sClientList Error: error=%s", myerror.Errors[ye.Code].LogMsg)
				return nil, ye
			}
			k8sClientList = append(k8sClientList, k8sClient)
		}

		log.Infof("CreateK8sClient Success: len(k8sClientList)=%d", len(k8sClientList))

		return k8sClientList, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("CreateK8sCilentList Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}
}

// Create K8s Client By ApiServer
func CreateK8sClient(apiServer string) (*client.Client, *myerror.YceError) {
	if CheckValidate(apiServer) {
		config := &restclient.Config{
			Host: apiServer,
		}

		k8sclient, err := client.New(config)
		if err != nil {
			ye := myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			log.Errorf("CreateK8sClient Error: error=%s", err)
			return nil, ye
		} else {

			log.Infof("CreateK8sClient Success: &k8sClient=%p", k8sclient)
			return k8sclient, nil
		}
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("CreateK8sCilent Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}

}

// Get ApiServer List by Datacenter Id List
func GetApiServerList(dcIdList []int32) ([]string, *myerror.YceError) {
	apiServerList := make([]string, 0)

	if CheckValidate(dcIdList) {
		for _, dcId := range dcIdList {
			apiServer, err := GetApiServerByDcId(dcId)
			if err != nil {
				ye := err
				log.Errorf("GetApiServerList Error: error=%s", myerror.Errors[ye.Code].LogMsg)
				return nil, ye
			}
			apiServerList = append(apiServerList, apiServer)
		}

		log.Infof("GetApiServerList Success: len(apiServerList)=%d", len(apiServerList))

		return apiServerList, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetApiServerList Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}
}

// Get Single ApiServer by single Datacenter Id
func GetApiServerByDcId(DcId int32) (string, *myerror.YceError) {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(DcId)

	if err != nil {
		ye := myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		log.Errorf("GetApiServerByDcId Error: error=%s", err)
		return "", ye
	} else {
		host := dc.Host
		port := strconv.Itoa(int(dc.Port))
		apiServer := host + ":" + port

		log.Infof("GetApiServerByDcId Success: apiServer=%s", apiServer)
		return apiServer, nil
	}
}

// Get Deployment By Namespace
func GetDeploymentByNamespace(c *client.Client, namespace string) ([]extensions.Deployment, *myerror.YceError) {
	if CheckValidate(c) && CheckValidate(namespace) {
		dps, err := c.Extensions().Deployments(namespace).List(api.ListOptions{})
		if err != nil {
			log.Errorf("GetDeploymentByNamespace Error: error=%s", err)
			ye := myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
			return nil, ye
		}

		log.Infof("GetDeploymentByNamespace Success: len(deployment.Items)=%d", len(dps.Items))
		return dps.Items, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetDeploymentByNamespace Errror: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}
}

// Get ReplicaSets By Deployment
func GetReplicaSetsByDeployment(c *client.Client, d *extensions.Deployment) ([]extensions.ReplicaSet, *myerror.YceError) {
	namespace := d.Namespace
	selector, err := unver.LabelSelectorAsSelector(d.Spec.Selector)
	if err != nil {
		log.Errorf("GetReplicaSetsByDeployment Errror: error=%s", err)
		ye := myerror.NewYceError(myerror.EKUBE_LABEL_SELECTOR, "")
		return nil, ye
	}

	options := api.ListOptions{LabelSelector: selector}
	if CheckValidate(c) && CheckValidate(namespace) {
		rss, err := c.Extensions().ReplicaSets(namespace).List(options)
		if err != nil {
			log.Errorf("GetReplicaSetsByDeployment Error: error=%s", err)
			ye := myerror.NewYceError(myerror.EKUBE_LIST_REPLICASET, "")
			return nil, ye
		}

		log.Infof("GetReplicaSetsByDeployment Success: len(replicaset.Items)=%d", len(rss.Items))
		return rss.Items, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetReplicaSetsByDeployment Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}
}

func GetPodListByReplicaSet(c *client.Client, rs *extensions.ReplicaSet) (*api.PodList, *myerror.YceError) {
	selector, err := unver.LabelSelectorAsSelector(rs.Spec.Selector)
	if err != nil {
		ye := myerror.NewYceError(myerror.EKUBE_LABEL_SELECTOR, "")
		log.Errorf("GetPodListByReplicaSets Error: error=%s", err)
		return nil, ye
	}

	namespace := rs.Namespace
	options := api.ListOptions{LabelSelector: selector}

	if CheckValidate(c) && CheckValidate(namespace) {
		podList, err := c.Pods(namespace).List(options)
		if err != nil {
			ye := myerror.NewYceError(myerror.EKUBE_LIST_PODS, "")
			log.Errorf("GetPodListByReplicaSets Error: error=%s", err)
			return nil, ye
		}

		log.Infof("GetPodListByReplicaSets Success: len(PodList.Items)=%d", len(podList.Items))
		return podList, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetPodListByReplicaSets Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}
}

func GetPodsByReplicaSets(c *client.Client, replicaSets []extensions.ReplicaSet) ([]api.Pod, *myerror.YceError) {

	if CheckValidate(c) && CheckValidate(replicaSets) {
		pods := make([]api.Pod, 0)
		for _, rs := range replicaSets {
			podList, ye := GetPodListByReplicaSet(c, &rs)
			if ye != nil {
				log.Errorf("GetPodsByReplicaSets Error: error=%s", myerror.Errors[ye.Code].LogMsg)
				return nil, ye
			}
			for _, pod := range podList.Items {
				pods = append(pods, pod)
			}

			log.Infof("GetPodsByReplicaSets getPods from ReplicaSet Success: len(podList.Items)=%d", len(podList.Items))
		}

		log.Infof("GetPodsByReplicaSets Success: len(pods)=%d", len(pods))
		return pods, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetPodsByReplicaSets Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}

}

func GetNodeByPod(c *client.Client, pod *api.Pod) (*api.Node, *myerror.YceError) {
	nodeName := pod.Spec.NodeName
	node, err := c.Nodes().Get(nodeName)
	if err != nil {
		ye := myerror.NewYceError(myerror.EKUBE_GET_NODE_BY_POD, "")
		log.Errorf("GetNodeByPod Error: error=%s", err)
		return nil, ye
	}

	log.Infof("GetNodeByPod Success: nodeName=%s", node.Name)
	return node, nil
}

func GetServicesByNamespace(c *client.Client, namespace string) ([]api.Service, *myerror.YceError) {
	if CheckValidate(c) && CheckValidate(namespace) {
		svcs, err := c.Services(namespace).List(api.ListOptions{})
		if err != nil {
			ye := myerror.NewYceError(myerror.EKUBE_LIST_SERVICE, "")
			log.Errorf("GetServicesByNamespace Error: error=%s", err)
			return nil, ye
		}

		log.Infof("GetServicesByNamespace Success: len(service.Items)=%d", len(svcs.Items))
		return svcs.Items, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetServicesByNamespace Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}
}

func GetEndpointsByNamespace(c *client.Client, namespace string) ([]api.Endpoints, *myerror.YceError) {
	if CheckValidate(c) && CheckValidate(namespace) {
		eps, err := c.Endpoints(namespace).List(api.ListOptions{})
		if err != nil {
			log.Errorf("GetEndpointsByNamespace Error: error=%s", err)
			ye := myerror.NewYceError(myerror.EKUBE_LIST_ENDPOINTS, "")
			return nil, ye
		}

		log.Infof("GetEndpointsByNamespace Success: len(endpoints.Items)=%d", len(eps.Items))
		return eps.Items, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetEndpointsByNamespace Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}
}

func GetPodsByService(c *client.Client, svc *api.Service) ([]api.Pod, *myerror.YceError) {

	if CheckValidate(c) && CheckValidate(svc) {
		selector := new(unver.LabelSelector)
		selector.MatchLabels = svc.Spec.Selector
		s, err := unver.LabelSelectorAsSelector(selector)
		if err != nil {
			ye := myerror.NewYceError(myerror.EKUBE_LABEL_SELECTOR, "")
			log.Errorf("GetPodsByService Error: error=%s", err)
			return nil, ye
		}

		namespace := svc.Namespace
		options := api.ListOptions{LabelSelector: s}

		podList, err := c.Pods(namespace).List(options)
		if err != nil {
			ye := myerror.NewYceError(myerror.EKUBE_GET_PODS_BY_SERVICE, "")
			log.Errorf("GetPodsByService Error: error=%s", err)
			return nil, ye
		}

		log.Infof("GetPodsByService Success: len(podList.Items)=%d", len(podList.Items))
		return podList.Items, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetPodsByService Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}

}

// deleteReplicaSet
func DeleteReplicaSets(c *client.Client, replicaSets []extensions.ReplicaSet) *myerror.YceError {
	if CheckValidate(c) && CheckValidate(replicaSets) {
		for _, rs := range replicaSets {
			falseVar := false
			deleteOptions := &api.DeleteOptions{OrphanDependents: &falseVar}

			log.Debugf("DeletReplicaSet Name: replicaSetName=%s", rs.Name)
			err := c.Extensions().ReplicaSets(rs.Namespace).Delete(rs.Name, deleteOptions)
			if err != nil {
				log.Errorf("DeleteReplicaSet Error: name=%s, err=%s", rs.Name, err)
				ye := myerror.NewYceError(myerror.EKUBE_DELETE_REPLICASET, "")
				return ye
			}
		}

		log.Infof("DeleteReplicaSet successfully")
		return nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("DeleteReplicaSet Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return ye
	}

}

// delete Pods

func DeletePods(c *client.Client, pods []api.Pod) *myerror.YceError {
	if CheckValidate(c) && CheckValidate(pods) {
		for _, pod := range pods {
			falseVar := false
			deleteOptions := &api.DeleteOptions{OrphanDependents: &falseVar}

			log.Infof("DeletePods: podName=%s", pod.Name)
			err := c.Pods(pod.Namespace).Delete(pod.Name, deleteOptions)

			if err != nil {
				log.Errorf("DeletePods: Error: name=%s, err=%s", pod.Name, err)
				ye := myerror.NewYceError(myerror.EKUBE_DELETE_POD, "")
				return ye
			}

		}

		log.Infof("Delete pods successfully")
		return nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("DeletePods Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return ye
	}
}

// delete Deployment
func DeleteDeployment(c *client.Client, deployment *extensions.Deployment) *myerror.YceError {

	if CheckValidate(c) && CheckValidate(deployment) {
		err := c.Extensions().Deployments(deployment.Namespace).Delete(deployment.Name, nil)
		if err != nil {
			log.Errorf("DeleteDeployment Error: name=%s, err=%s", deployment.Name, err)
			ye := myerror.NewYceError(myerror.EKUBE_DELETE_DEPLOYMENT, "")
			return ye
		}

		log.Infof("DeleteDeployment success: name=%s", deployment.Name)
		return nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("DeleteDeployment Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return ye
	}

}

// Get Datacenters by OrgId
func GetDatacentersByOrgId(orgId string) ([]mydatacenter.DataCenter, *myerror.YceError) {
	if CheckValidate(orgId) {
		org, err := organization.GetOrganizationById(orgId)

		if err != nil {
			log.Errorf("GetDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
			ye := myerror.NewYceError(myerror.EYCE_ORGTODC, "")
			return nil, ye

		}

		dcList, err := organization.GetDataCentersByOrganization(org)
		if err != nil {
			log.Errorf("GetDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
			ye := myerror.NewYceError(myerror.EYCE_ORGTODC, "")
			return nil, ye
		}

		log.Infof("GetDatacentersByOrgId: len(Datacenters)=%d", len(dcList))
		return dcList, nil
	}

	ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
	log.Errorf("GetDatacentersByOrgId Error: error=%s", myerror.Errors[ye.Code].LogMsg)
	return nil, ye
}

// Get All Quotas
func GetAllQuotasOrderByCpu() ([]myqouta.Quota, *myerror.YceError) {
	// Get all quotas
	quotas, err := myqouta.QueryAllQuotasOrderByCpu()
	if err != nil {
		log.Errorf("GetAllQuotasOrderByCpu error: error=%s", err)
		ye := myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return nil, ye
	}
	return quotas, nil
}

type DatacenterList struct {
	DcIdList []int32
	DcName   []string
}

// Get Datacenter List By OrgId
func GetDatacenterListByOrgId(orgId string) (*DatacenterList, *myerror.YceError) {

	if CheckValidate(orgId) {
		org, err := organization.GetOrganizationById(orgId)

		if err != nil {
			log.Errorf("GetDatacenterListByOrgId Error: orgId=%s, error=%s", orgId, err)
			ye := myerror.NewYceError(myerror.EYCE_ORGTODC, "")
			return nil, ye

		}

		dcList, err := organization.GetDataCentersByOrganization(org)
		if err != nil {

			log.Errorf("GetDatacenterListByOrgId Error: orgId=%s, error=%s", orgId, err)
			ye := myerror.NewYceError(myerror.EYCE_ORGTODC, "")
			return nil, ye
		}

		DcIdList := make([]int32, 0)
		DcName := make([]string, 0)

		for _, dc := range dcList {
			DcIdList = append(DcIdList, dc.Id)
			DcName = append(DcName, dc.Name)
		}

		datacenterList := &DatacenterList{
			DcIdList: DcIdList,
			DcName:   DcName,
		}

		log.Infof("GetDatacenterListByOrgId: len(DcIdList)=%d, len(DcName)=%d", len(DcIdList), len(DcName))
		return datacenterList, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetDatacenterListByOrgId Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}

}

func GetDcIdListByOrgId(orgId string) ([]int32, *myerror.YceError) {
	if CheckValidate(orgId) {
		org, err := organization.GetOrganizationById(orgId)

		if err != nil {
			log.Errorf("GetDcIdListByOrgId Error: orgId=%s, error=%s", orgId, err)
			ye := myerror.NewYceError(myerror.EYCE_ORGTODC, "")
			return nil, ye

		}

		dcList, err := organization.GetDataCentersByOrganization(org)
		if err != nil {
			log.Errorf("GetDcIdListByOrgId Error: orgId=%s, error=%s", orgId, err)
			ye := myerror.NewYceError(myerror.EYCE_ORGTODC, "")
			return nil, ye
		}

		DcIdList := make([]int32, 0)

		for _, dc := range dcList {
			DcIdList = append(DcIdList, dc.Id)
		}

		log.Infof("GetDcIdListByOrgId: len(DcIdList)=%d", len(DcIdList))
		return DcIdList, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetDcIdListByOrgId Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}
}

// get Pod By podName
type LogOptionType struct {
	Container    string      `json:"container,omitempty"`    //暂时不做
	Follow       bool        `json:"follow,omitempty"`       //false 暂时不做, 页面开关,默认为关闭
	Previous     bool        `json:"previous,omitempty"`     //暂时不做
	SinceSeconds *int64      `json:"sinceSeconds,omitempty"` //暂时不做
	SinceTime    *unver.Time `json:"sinceTime,omitempty"`    //暂时不做
	Timestamps   bool        `json:"timeStamps,omitempty"`   //true, 时间戳,默认打开
	TailLines    *int64      `json:"tailLines,omitempty"`    //用户设定
	LimitBytes   *int64      `json:"limitBytes,omitempty"`   //暂时不做
}

func GetPodLogsByPodName(c *client.Client, LogOption *LogOptionType, podName, orgId string) (string, *myerror.YceError) {
	if CheckValidate(c) && CheckValidate(LogOption) && CheckValidate(podName) && CheckValidate(orgId) {
		options := &api.PodLogOptions{
			Container:    LogOption.Container,
			Follow:       LogOption.Follow,
			Previous:     LogOption.Previous,
			SinceSeconds: LogOption.SinceSeconds,
			SinceTime:    LogOption.SinceTime,
			Timestamps:   LogOption.Timestamps,
			TailLines:    LogOption.TailLines,
			LimitBytes:   LogOption.LimitBytes,
		}

		namespace, ye := GetOrgNameByOrgId(orgId)
		reader, err := c.Pods(namespace).GetLogs(podName, options).Stream()
		if err != nil {
			log.Errorf("GetPodLogsByPodName Error: podName=%s, error=%s", podName, err)
			ye = myerror.NewYceError(myerror.EKUBE_LOGS_POD, "")
			return "", ye
		}
		defer reader.Close()

		b, err := ioutil.ReadAll(reader)
		if err != nil {

			log.Errorf("GetPodLogsByPodName Error: podName=%s, error=%s", podName, err)
			ye = myerror.NewYceError(myerror.EKUBE_LOGS_POD, "")
			return "", ye
		}

		logs := string(b)

		log.Infof("GetPodLogsByPodName successfully: len(bytes)=%d", len(b))

		return logs, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")

		log.Errorf("GetPodLogsByPodName Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return "", ye
	}

}

// original: queryDcNameByDcId(dcIdList []int32) (dcNameList []string)
func GetDcNameListByDcIdList(dcIdList []int32) ([]string, *myerror.YceError) {
	if CheckValidate(dcIdList) {
		dcNameList := make([]string, 0)

		for _, dcId := range dcIdList {
			dc, ye := GetDatacenterByDcId(dcId)
			if ye != nil {
				log.Errorf("GetDcNameListByDcIdList Error: error=%s", ye)
				return nil, ye
			}
			dcName := dc.Name
			dcNameList = append(dcNameList, dcName)

			log.Infof("GetDcNameListByDcIdList get Name: %s", dcName)
		}
		log.Infof("GetDcNameListByDcIdList len(dcNameList)=%d", len(dcNameList))

		return dcNameList, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetDcNameListByDcIdList Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}

}

// getDatacenter by DcId
func GetDatacenterByDcId(dcId int32) (*mydatacenter.DataCenter, *myerror.YceError) {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("GetDatacenterByDcId QueryDataCenterById Error: dcId=%d, error=%s", dcId, err)
		ye := myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return nil, ye
	}

	log.Infof("GetDatacenterByDcId successfully: name=%s, id=%d", dc.Name, dc.Id)
	return dc, nil
}

// get OrgNameByOrgId
func GetOrgNameByOrgId(OrgId string) (string, *myerror.YceError) {
	if CheckValidate(OrgId) {
		organization := new(myorganization.Organization)

		orgId, _ := strconv.Atoi(OrgId)
		organization.QueryOrganizationById(int32(orgId))
		log.Infof("GetOrgNameByOrgId successfully: orgName=%s, orgId=%d", organization.Name, orgId)
		return organization.Name, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetOrgNameByOrgId Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return "", ye
	}

}

// Get DeployAndPodList Pair by deploymentList

type DeployAndPodList struct {
	UserName string                 `json:"userName"`
	Deploy   *extensions.Deployment `json:"deploy"`
	PodList  *api.PodList           `json:"podList"`
}

func GetDeployAndPodList(userId int32, c *client.Client, deploymentList *extensions.DeploymentList) ([]DeployAndPodList, *myerror.YceError) {

	if CheckValidate(c) && CheckValidate(deploymentList) {
		dap := make([]DeployAndPodList, 0)

		for _, deployment := range deploymentList.Items {

			dp := new(DeployAndPodList)

			dp.UserName = myuser.QueryUserNameByUserId(userId)

			dp.Deploy = new(extensions.Deployment)

			*dp.Deploy = deployment

			rsList, ye := GetReplicaSetsByDeployment(c, dp.Deploy)
			if ye != nil {
				log.Errorf("GetDeployAndPodList Error: error=%s", myerror.Errors[ye.Code].LogMsg)
				return nil, ye
			}

			newRs, err := deploymentutil.FindNewReplicaSet(dp.Deploy, rsList)
			if err != nil {

				log.Errorf("GetDeployAndPodList Error: error=%s", err)
				ye := myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
				return nil, ye
			}

			PodList, ye := GetPodListByReplicaSet(c, newRs)
			if ye != nil {
				log.Errorf("GetDeployAndPodList Error: error=%s", myerror.Errors[ye.Code].LogMsg)
				return nil, ye
			}

			//dp.PodList = new(api.PodList)
			dp.PodList = PodList

			dap = append(dap, *dp)

		}
		log.Infof("GetDeployAndPodList successfully: len(deployAndPodList)=%d", len(dap))
		return dap, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetDeployAndPodList Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}

}

func GetNewReplicaSetByDeployment(c *client.Client, deployment *extensions.Deployment) (*extensions.ReplicaSet, *myerror.YceError) {
	rsList, ye := GetReplicaSetsByDeployment(c, deployment)
	if ye != nil {
		log.Errorf("GetReplicaSetsByDeployment Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}

	rsNew, err := deploymentutil.FindNewReplicaSet(deployment, rsList)
	if err != nil {

		log.Errorf("GetReplicaSetsByDeployment Error: error=%s", err)
		ye := myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
		return nil, ye
	}

	return rsNew, nil
}

// Get Deployment by deployment-name
func GetDeploymentByNameAndNamespace(c *client.Client, deploymentName, namespace string) (*extensions.Deployment, *myerror.YceError) {

	// Get namespace(org.Name) by orgId
	if CheckValidate(c) && CheckValidate(deploymentName) {
		dp, err := c.Extensions().Deployments(namespace).Get(deploymentName)
		if err != nil {
			log.Errorf("GetDeployByNameAndNamespace Error: namespace=%s, deployment-name=%s, err=%s\n", dp.Namespace, dp.Name, err)
			ye := myerror.NewYceError(myerror.EKUBE_GET_DEPLOYMENT, "")
			return nil, ye
		}

		log.Infof("GetDeploymentByNameAndNamespace over: apiServer=%s, namespace=%s, name=%s, deployment=%p\n", dp.Namespace, dp.Name, dp)
		return dp, nil
	} else {
		ye := myerror.NewYceError(myerror.EINVALID_PARAM, "")
		log.Errorf("GetDeploymentByNameAndNamespace Error: error=%s", myerror.Errors[ye.Code].LogMsg)
		return nil, ye
	}

}

// check the value if it is validate

func CheckValidate(value interface{}) bool {

	if reflect.TypeOf(value).Kind() == reflect.String && value != "" {
		return true
	}

	if value != nil {
		return true
	}

	return false

	/*
		if reflect.TypeOf(value).Kind() == reflect.Array && value != nil && len(value) > 0 {
			flag = true
		}

		if reflect.TypeOf(value).Kind() == reflect.Slice && value != nil && len(value) > 0 {
			flag = true
		}

		if reflect.TypeOf(value).Kind() == reflect.Ptr && value != nil {
			flag = true
		}

		return flag
	*/
}

func QueryDuplicatedNameAndOrgId(name string, orgId int32) (bool, *myerror.YceError) {
	u := new(myuser.User)
	err := u.QueryUserByNameAndOrgId(name, orgId)
	// not found
	if err != nil {
		ye := myerror.NewYceError(myerror.EYCE_NOTFOUND, "")
		return false, ye
	}
	// found
	return true, nil
}

//TODO: Get Namespace List By Datacenter Id List
func GetNamespaceListByDcIdList() {

}
