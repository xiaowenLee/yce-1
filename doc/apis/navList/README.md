
获取用户的功能列表
===============

此项功能是用于控制不同的用户可以拥有不同的权限(功能)
----------------------------


### 请求

请求URL: GET /api/v1/organizations/orgId/users/userId/navList

请求头: Authorization: ${SessionId}

### 返回数据:

* 如果用户是admin:

{
    "code": 0,
    "msg": "OK",
    "data": "{
        "list": [
            {
                "id": 1,
                "name": "Dashboard",
                "state": "main.dashboard",
                "includeState": "main.dashboard",
                "className": "fa-dashboard"
            },
            {
                "id": 4,
                "name": "镜像管理",
                "state": "main.imageManage",
                "includeState": "main.imageManage",
                "className": "fa-file-archive-o",
                "item": [
                    {
                        "id": 41,
                        "name": "基础镜像",
                        "state": "main.imageManageBase",
                        "includeState": "main.imageManageBase"
                    }
                ]
            },
            {
                "id": 5,
                "name": "集群拓扑",
                "state": "main.topology",
                "includeState": "main.topology",
                "className": "fa-share-alt"
            },
            {
                "id": 6,
                "name": "个人中心",
                "state": "main.personalCenter",
                "includeState": "main.personalCenter",
                "className": "fa-adn",
                "item": [
                    {
                        "id": 61,
                        "name": "个人设置",
                        "state": "main.personalSetting",
                        "includeState": "main.personalSetting"
                    },
                    {
                        "id": 62,
                        "name": "修改密码",
                        "state": "main.personalPassword",
                        "includeState": "main.personalPassword"
                    }
                ]
            },
            {
                "id": 7,
                "name": "用户管理",
                "state": "main.userManage",
                "includeState": "main.userManage",
                "className": "fa-adn",
                "item": [
                    {
                        "id": 71,
                        "name": "创建用户",
                        "state": "main.createUser",
                        "includeState": "main.createUser"
                    }
                ]
            },
            {
                "id": 8,
                "name": "数据中心管理",
                "state": "main.dcManage",
                "includeState": "main.dcManage",
                "className": "fa-adn",
                "item": [
                    {
                        "id": 81,
                        "name": "添加数据中心",
                        "state": "main.addDc",
                        "includeState": "main.addDc"
                    }
                ]
            },
            {
                "id": 9,
                "name": "组织管理",
                "state": "main.orgManage",
                "includeState": "main.orgManage",
                "className": "fa-adn",
                "item": [
                    {
                        "id": 91,
                        "name": "添加组织",
                        "state": "main.addOrg",
                        "includeState": "main.addOrg"
                    }
                ]
            }
        ]
    }"
}
    
* 如果是普通用户 

{
    "code": 0,
    "msg": "OK",
    "data": "{
        "list": [
            {
                "id": 1,
                "name": "Dashboard",
                "state": "main.dashboard",
                "includeState": "main.dashboard",
                "className": "fa-dashboard"
            },
            {
                "id": 2,
                "name": "应用管理",
                "state": "main.appManage",
                "includeState": "main.appManage",
                "className": "fa-adn",
                "item": [
                    {
                        "id": 21,
                        "name": "应用发布",
                        "state": "main.appManageDeployment",
                        "includeState": "main.appManageDeployment"
                    },
                    {
                        "id": 22,
                        "name": "历史操作",
                        "state": "main.appManageHistory",
                        "includeState": "main.appManageHistory"
                    }
                ]
            },
            {
                "id": 3,
                "name": "服务管理",
                "state": "main.extensions",
                "includeState": "main.extensions",
                "className": "fa-arrows",
                "item": [
                    {
                        "id": 31,
                        "name": "创建服务",
                        "state": "main.extensionsService",
                        "includeState": "main.extensionsService"
                    },
                    {
                        "id": 32,
                        "name": "创建访问点",
                        "state": "main.extensionsEndpoint",
                        "includeState": "main.extensionsEndpoint"
                    }
                ]
            },
            {
                "id": 4,
                "name": "镜像管理",
                "state": "main.imageManage",
                "includeState": "main.imageManage",
                "className": "fa-file-archive-o",
                "item": [
                    {
                        "id": 41,
                        "name": "基础镜像",
                        "state": "main.imageManageBase",
                        "includeState": "main.imageManageBase"
                    }
                ]
            },
            {
                "id": 5,
                "name": "集群拓扑",
                "state": "main.topology",
                "includeState": "main.topology",
                "className": "fa-share-alt"
            },
            {
                "id": 6,
                "name": "个人中心",
                "state": "main.personalCenter",
                "includeState": "main.personalCenter",
                "className": "fa-adn",
                "item": [
                    {
                        "id": 61,
                        "name": "个人设置",
                        "state": "main.personalSetting",
                        "includeState": "main.personalSetting"
                    },
                    {
                        "id": 62,
                        "name": "修改密码",
                        "state": "main.personalPassword",
                        "includeState": "main.personalPassword"
                    },
                    {
                        "id": 63,
                        "name": "事件提醒",
                        "state": "main.eventAlert",
                        "includeState": "main.eventAlert"
                    },
                    {
                        "id": 64,
                        "name": "计费&充值",
                        "state": "main.recharge",
                        "includeState": "main.recharge"
                    }
                ]
            },
            {
                "id": 10,
                "name": "绿色通道",
                "state": "main.walkthrogh",
                "includeState": "main.costManage",
                "className": "fa-adn"
            }
        ]
    }"
}
