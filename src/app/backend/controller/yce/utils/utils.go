package utils

import (
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/restclient"
	"strconv"
	"github.com/kubernetes/kubernetes/pkg/apis/extensions"
	"github.com/kubernetes/kubernetes/pkg/api"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
)


// Create K8s Client List by ApiServerList
func CreateK8sClientList(apiServerList []string) ([]*client.Client, myerror.YceError) {
	k8sClientList := make([]*client.Client, 0)

	if apiServerList != nil && len(apiServerList) > 0 {
		for _, apiServer := range apiServerList {
			k8sClient, err := CreateK8sClient(apiServer)
			if err != nil {
				ye := err
				return nil, ye
			}
			k8sClientList = append(k8sClientList, k8sClient)
		}

		return k8sClientList, nil
	} else {
		ye := myerror.NewYceError(myerror.EOOM, "")
		return nil, ye
	}
}

// Create K8s Client By ApiServer
func CreateK8sClient(apiServer string) (*client.Client, myerror.YceError) {
	config := &restclient.Config{
		Host: apiServer,
	}

	k8sclient, err := client.New(config)
	if err != nil {
		ye := myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return nil, ye
	} else {
		return k8sclient, nil
	}
}

// Get ApiServer List by Datacenter Id List
func GetApiServerList(dcIdList []int32) ([]string, myerror.YceError) {
	apiServerList := make([]string, 0)

	if dcIdList != nil && len(dcIdList) > 0 {
		for _, dcId := range dcIdList {
			apiServer, err := GetApiServerByDcId(dcId)
			if err != nil {
				ye := err
				return nil, ye
			}
			apiServerList = append(apiServerList, apiServer)
		}

		return apiServerList, nil
	} else {
		ye := myerror.NewYceError(myerror.EOOM, "")
		return nil, ye
	}
}

// Get Single ApiServer by single Datacenter Id
func GetApiServerByDcId(DcId int32) (string, myerror.YceError) {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(DcId)

	if err != nil {
		ye := myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return "", ye
	} else {
		host := dc.Host
		port := strconv.Itoa(int(dc.Port))
		apiServer := host + ":" + port

		return apiServer, nil
	}
}


func GetDeplyomentByNamespace(c *client.Client, namespace string) ([]extensions.Deployment, myerror.YceError) {
	if c != nil && namespace != "" {
		dps, err := c.Extensions().Deployments(namespace).List(api.ListOptions{})
		if err != nil {
			ye := myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
			return nil, ye
		}


		return dps.Items, nil
	} else {
		ye := myerror.NewYceError(myerror.EOOM, "")
		return nil, ye
	}
}

func GetReplicaSetsByDeployment(c *client.Client, d *extensions.Deployment) ([]extensions.ReplicaSet, myerror.YceError) {
	namespace := d.Namespace
	selector, err := unver.LabelSelectorAsSelector(d.Spec.Selector)
	if err != nil {
		ye := myerror.NewYceError(myerror.EKUBE_LABEL_SELECTOR, "")
		return nil, ye
	}

	options := api.ListOptions{LabelSelector: selector}
	if c != nil && namespace != "" {
		rss, err := c.Extensions().ReplicaSets(namespace).List(options)
		if err != nil {
			ye := myerror.NewYceError(myerror.EKUBE_LIST_REPLICASET, "")
			return nil, ye
		}

		return rss.Items, nil
	} else {
		ye := myerror.NewYceError(myerror.EOOM, "")
		return nil, ye
	}
}

func GetPodsByReplicaSets(c *client.Client, rs *extensions.ReplicaSet) ([]api.Pod, myerror.YceError) {
	selector, err := unver.LabelSelectorAsSelector(rs.Spec.Selector)
	if err != nil {
		ye := myerror.NewYceError(myerror.EKUBE_LABEL_SELECTOR, "")
		return nil, ye
	}

	namespace := rs.Namespace
	options := api.ListOptions{LabelSelector:selector}

	if c != nil && namespace != "" {
		pods, err := c.Pods(namespace).List(options)
		if err != nil {
			ye := myerror.NewYceError(myerror.EKUBE_LIST_PODS, "")
			return nil, ye
		}

		return pods.Items, nil
	} else {
		ye := myerror.NewYceError(myerror.EOOM, "")
		return nil, ye
	}
}

func GetNodeByPod(c *client.Client, pod *api.Pod) (api.Node, myerror.YceError) {
	nodeName := pod.Spec.NodeName
	node, err := c.Nodes().Get(nodeName)
	if err != nil {
		ye := myerror.NewYceError(myerror.EKUBE_GET_NODE_BY_POD, "")
		return nil, ye
	}

	return node, nil
}

func GetServicesByNamespace(c *client.Client, namespace string) ([]api.Service, myerror.YceError) {
	if c != nil && namespace != "" {

		svcs, err := c.Services(namespace).List(api.ListOptions{})
		if err != nil {
			ye := myerror.NewYceError(myerror.EKUBE_LIST_SERVICE, "")
			return nil, ye
		}

		return svcs.Items, nil
	} else {
		ye := myerror.NewYceError(myerror.EOOM, "")
		return nil, ye
	}
}

func GetEndpointsByNamespace(c *client.Client, namespace string) ([]api.Endpoints, myerror.YceError) {
	if c != nil && namespace != "" {
		eps, err := c.Endpoints(namespace).List(api.ListOptions{})
		if err != nil {
			ye := myerror.NewYceError(myerror.EKUBE_LIST_ENDPOINTS, "")
			return nil, ye
		}

		return eps.Items, nil
	} else {
		ye := myerror.NewYceError(myerror.EOOM, "")
		return nil, ye
	}
}

func GetPodsByService(c *client.Client, svc *api.Service) ([]api.Pod, myerror.YceError) {
	selector := new(unver.LabelSelector)
	selector.MatchLabels = svc.Spec.Selector

	s, err := unver.LabelSelectorAsSelector(selector)
	if err != nil {
		ye := myerror.NewYceError(myerror.EKUBE_LABEL_SELECTOR, "")
		return nil, ye
	}

	namespace := svc.Namespace
	options := api.ListOptions{LabelSelector: s}

	podList, err := c.Pods(namespace).List(options)
	if err != nil {
		ye := myerror.NewYceError(myerror.EKUBE_GET_PODS_BY_SERVICE, "")
		return nil, err
	}

	return podList.Items, nil
}


//TODO: Get Namespace List By Datacenter Id List
func GetNamespaceListByDcIdList() {

}
