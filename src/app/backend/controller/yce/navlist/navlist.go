package navlist

import (
	"github.com/kataras/iris"
)

type NavListController struct {
	*iris.Context
}

var navList = `
{
    "list": [
        {"id": 1, "name": "Dashboard", "state": "main.dashboard","includeState": "main.dashboard","className":"fa-dashboard"},
        {"id": 2, "name": "应用管理", "state": "main.appManage","includeState": "main.appManage","className":"fa-adn",
            "item": [
                {"id": 21, "name": "应用发布", "state": "main.appManageDeployment", "includeState": "main.appManageDeployment"},
                {"id": 22, "name": "历史操作", "state": "main.appManageHistory", "includeState": "main.appManageHistory"}
            ]
        },
        {"id": 3, "name": "镜像管理", "state": "main.imageManage", "includeState": "main.imageManage","className":"fa-file-archive-o",
            "item": [
                {"id": 31, "name": "基础镜像", "state": "main.imageManageDelete", "includeState": "main.imageManageDelete"}
            ]
        },
        {"id": 4, "name": "服务管理", "state": "main.extensions", "includeState": "main.extensions","className":"fa-arrows",
            "item": [
                {"id": 51, "name": "创建服务", "state": "main.extensionsService", "includeState": "main.extensionsService"},
                {"id": 52, "name": "创建访问点", "state": "main.extensionsEndpoint", "includeState": "main.extensionsEndpoint"}
            ]
        },
        {"id": 5, "name": "云盘管理", "state": "main.rbdManage", "includeState": "main.rbdManage","className":"fa-cloud"},
        {"id": 6, "name": "计费&充值", "state": "main.costManage", "includeState": "main.costManage","className":"fa-credit-card"}
    ]
}
`

func (nlc NavListController) Get() {
	nlc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	nlc.Write(navList)
}
