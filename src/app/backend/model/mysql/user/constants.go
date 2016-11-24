package user

const (
	USER_PASSWORD = "SELECT id, name, password, orgId, createdAt, modifiedAt, modifiedOp, navList FROM user WHERE name=? and password=?"

	USER_SELECT          = "SELECT id, name, password, orgId, createdAt, modifiedAt, modifiedOp, navList FROM user WHERE id=? "
	USER_CHECK_DUPLICATE = "SELECT id, name, password, orgId, createdAt, modifiedAt, modifiedOp, navList FROM user WHERE name=? AND orgId=?"

	USER_INSERT = "INSERT INTO " +
		"user(name, password, orgId, status, createdAt, modifiedAt, modifiedOp, comment, navList) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	USER_INSERT_ON_DUPLICATE_KEY_UPDATE = "INSERT INTO " +
		"user(name, password, orgId, status, createdAt, modifiedAt, modifiedOp, comment, navList) " +
		"values(?,?,?,?,?,?,?,?,?) " +
		"ON DUPLICATE KEY UPDATE " +
		"orgId=?, status=?"

	USER_UPDATE = "UPDATE user SET password=?, orgId=?, modifiedAt=?, modifiedOp=?, navList=? WHERE id=?"

	USER_DELETE = "UPDATE user SET status=?, modifiedAt=?, modifiedOp=?, navList=? WHERE id=?"
	USER_SELECT_ALL = "SELECT id, name, password, orgId, createdAt, modifiedAt, modifiedOp, comment FROM user where status=?"
	USER_CHECK_DUPLICATE_NAME = "SELECT id, name, password, orgId, createdAt, modifiedAt, modifiedOp FROM user WHERE name=? AND status=?"
	// USER_NAVLIST = "SELECT navList FROM user WHERE id=?"
	USER_SELECT_BY_ORGID = "SELECT id, name, password, orgId, status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM user " +
		"WHERE status=? AND orgId=?"

	VALID   = 1
	INVALID = 0

	// All of navList
	NAVLIST = `
{
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
        },
        {
            "id": 10,
            "name": "绿色通道",
            "state": "main.walkthrogh",
            "includeState": "main.costManage",
            "className": "fa-adn"
        },
	{
            "id": 11,
            "name": "模板管理",
            "state": "main.tpManage",
            "includeState": "main.tpManage",
            "className": "fa-adn",
            "item": [
                {
                    "id": 111,
                    "name": "添加模板",
                    "state": "main.addTp",
                    "includeState": "main.addTp"
                }
            ]
        }
    ]
}
`

	// NavList for yce user
	USER_NAVLIST = `
{
    "list": [
        {
            "id": 1,
            "name": "Dashboard",
            "state": "main.dashboard",
            "includeState": "main.dashboard",
            "className": "fa-dashboard",
            "href": "dashBoard"
        },
        {
            "id": 2,
            "name": "应用管理",
            "state": "main.appManage",
            "includeState": "main.appManage",
            "className": "fa-adn",
            "href": "appManage",
            "item": [
                {
                    "id": 21,
                    "name": "应用发布",
                    "state": "main.appManageDeployment",
                    "includeState": "main.appManageDeployment",
                    "href": "appManageDeployment"
                },
                {
                    "id": 22,
                    "name": "历史操作",
                    "state": "main.appManageHistory",
                    "includeState": "main.appManageHistory",
                    "href": "appManageHistory"
                }
            ]
        },
        {
            "id": 3,
            "name": "服务管理",
            "state": "main.extensions",
            "includeState": "main.extensions",
            "className": "fa-arrows",
            "href": "extensions",
            "item": [
                {
                    "id": 31,
                    "name": "创建服务",
                    "state": "main.extensionsService",
                    "includeState": "main.extensionsService",
                    "href": "extensionsService"
                },
                {
                    "id": 32,
                    "name": "创建访问点",
                    "state": "main.extensionsEndpoint",
                    "includeState": "main.extensionsEndpoint",
                    "href": "extensionsEndpoint"
                }
            ]
        },
        {
            "id": 4,
            "name": "镜像管理",
            "state": "main.imageManage",
            "includeState": "main.imageManage",
            "className": "fa-file-archive-o",
            "href": "imageManage",
            "item": [
                {
                    "id": 41,
                    "name": "基础镜像",
                    "state": "main.imageManageBase",
                    "includeState": "main.imageManageBase",
                    "href": "imageBase"
                }
            ]
        },
        {
            "id": 5,
            "name": "集群拓扑",
            "state": "main.topology",
            "includeState": "main.topology",
            "className": "fa-share-alt",
            "href": "topology"
        },
        {
            "id": 6,
            "name": "个人中心",
            "state": "main.personalCenter",
            "includeState": "main.personalCenter",
            "className": "fa-adn",
            "href": "personalCenter",
            "item": [
                {
                    "id": 61,
                    "name": "个人设置",
                    "state": "main.personalSetting",
                    "includeState": "main.personalSetting",
                    "href": "personalSetting"
                },
                {
                    "id": 62,
                    "name": "修改密码",
                    "state": "main.personalPassword",
                    "includeState": "main.personalPassword",
                    "href": "personalPassword"
                },
                {
                    "id": 63,
                    "name": "事件提醒",
                    "state": "main.eventAlert",
                    "includeState": "main.eventAlert",
                    "href": "eventAlert"
                },
                {
                    "id": 64,
                    "name": "计费&充值",
                    "state": "main.recharge",
                    "includeState": "main.recharge",
		    "href": "recharge"
                }
            ]
        },
        {
            "id": 10,
            "name": "绿色通道",
            "state": "main.walkthrogh",
            "includeState": "main.costManage",
            "className": "fa-adn",
            "href": "walkthrogh"
        },
        {
            "id": 11,
            "name": "模板管理",
            "state": "main.tpManage",
            "includeState": "main.tpManage",
            "className": "fa-adn",
            "href": "tpManage",
            "item": [
                {
                    "id": 111,
                    "name": "添加模板",
                    "state": "main.addTp",
                    "includeState": "main.addTp",
            	    "href": "tpAdd"
                }
            ]
        }
    ]
}
`

	// NavList for yce amdin
	ADMIN_NAVLIST = `
{
    "list": [
    	{
            "id": 2,
            "name": "应用管理",
            "state": "main.appManage",
            "includeState": "main.appManage",
            "className": "fa-adn",
            "href": "appManage",
            "item": [
                {
                    "id": 22,
                    "name": "历史操作",
                    "state": "main.appManageHistory",
                    "href": "appManageHistory",
                    "includeState": "main.appManageHistory"
                }
            ]
        },
        {
            "id": 4,
            "name": "镜像管理",
            "state": "main.imageManage",
            "includeState": "main.imageManage",
            "className": "fa-file-archive-o",
            "href": "imageManage",
            "item": [
                {
                    "id": 41,
                    "name": "基础镜像",
                    "state": "main.imageManageBase",
                    "href":"imageManageBase",
                    "includeState": "main.imageManageBase"
                }
            ]
        },
        {
            "id": 5,
            "name": "集群拓扑",
            "state": "main.topology",
            "includeState": "main.topology",
            "href":"topology",
            "className": "fa-share-alt"
        },
        {
            "id": 6,
            "name": "个人中心",
            "state": "main.personalCenter",
            "includeState": "main.personalCenter",
            "href":"personalCenter",
            "className": "fa-adn",
            "item": [
                {
                    "id": 61,
                    "name": "个人设置",
                    "state": "main.personalSetting",
                    "href":"personalSetting",
                    "includeState": "main.personalSetting"
                },
                {
                    "id": 62,
                    "name": "修改密码",
                    "state": "main.personalPassword",
                    "href":"personalPassword",
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
            "href":"userManage",
            "item": [
                {
                    "id": 71,
                    "name": "创建用户",
                    "state": "main.createUser",
                    "href":"createUser",
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
            "href":"dcManage",
            "item": [
                {
                    "id": 81,
                    "name": "添加数据中心",
                    "state": "main.addDc",
                    "href":"addDc",
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
            "href":"orgManage",
            "item": [
                {
                    "id": 91,
                    "name": "添加组织",
                    "state": "main.addOrg",
                    "href":"addOrg",
                    "includeState": "main.addOrg"
                }
            ]
        }
    ]
}
`
)
