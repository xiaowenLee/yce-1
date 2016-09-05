/*
Mock.mock('http://10.69.40.36.com/api/v1/users/login', {
    "code": 0,
    "message": "",
    "data": {
        "userId": "12",
        "userName": "lidawei",
        "orgId": "2",
        "sessionId": "sfssfd-afds-asdf-af32s"
    }
});
*/

Mock.mock('http://10.69.40.36.com/api/main/navlist', {
    list: [
        {id: 1, name: 'Dashboard', state: 'main.dashboard',includeState: 'main.dashboard',className:'fa-dashboard'},
        {id: 2, name: '应用管理', state: 'main.appManage',includeState: 'main.appManage',className:'fa-adn',
            item: [
                {id: 21, name: '发布', state: 'main.appManageDeployment', includeState: 'main.appManageDeployment'},
                {id: 22, name: '回滚', state: 'main.appManageRollback', includeState: 'main.appManageRollback'},
                {id: 22, name: '滚动升级', state: 'main.appManageRollup', includeState: 'main.appManageRollup'},
                {id: 22, name: '撤销', state: 'main.appManageCancel', includeState: 'main.appManageCancel'},
                {id: 22, name: '历史发布', state: 'main.appManageHistory', includeState: 'main.appManageHistory'}
            ]
        },
        {id: 3, name: '镜像管理', state: 'main.imageManage', includeState: 'main.imageManage',className:'fa-file-archive-o',
            item: [
                {id: 31, name: '查找镜像', state: 'main.imageManageSearch', includeState: 'main.imageManageSearch'},
                {id: 32, name: '删除镜像', state: 'main.imageManageDelete', includeState: 'main.imageManageDelete'}
            ]
        },
        {id: 4, name: '云盘管理', state: 'main.rbdManage', includeState: 'main.rbdManage',className:'fa-cloud'},
        {id: 5, name: '扩展功能', state: 'main.extensions', includeState: 'main.extensions',className:'fa-arrows',
            item: [
                {id: 51, name: '创建服务', state: 'main.extensionsService', includeState: 'main.extensionsService'},
                {id: 52, name: '创建访问点', state: 'main.extensionsEndpoint', includeState: 'main.extensionsEndpoint'}
            ]
        },
        {id: 6, name: '计费&充值', state: 'main.costManage', includeState: 'main.costManage',className:'fa-credit-card'}
    ]
});

Mock.mock('http://192.168.0.102:8080/api/appManage/appList', {
    list: [
        {id: 1, name: 'Dashboard1',labels:'标签组合',dataCenter:'数据中心',replicas:'20',time:'20h'},
        {id: 2, name: 'Dashboard2',labels:'标签组合',dataCenter:'数据中心',replicas:'20',time:'20h'},
        {id: 3, name: 'Dashboard3',labels:'标签组合',dataCenter:'数据中心',replicas:'20',time:'20h'},
        {id: 4, name: 'Dashboard4',labels:'标签组合',dataCenter:'数据中心',replicas:'20',time:'20h'},
        {id: 5, name: 'Dashboard',labels:'标签组合',dataCenter:'数据中心',replicas:'20',time:'20h'},
        {id: 6, name: 'Dashboard',labels:'标签组合',dataCenter:'数据中心',replicas:'20',time:'20h'},
        {id: 7, name: 'Dashboard',labels:'标签组合',dataCenter:'数据中心',replicas:'20',time:'20h'},
        {id: 8, name: 'Dashboard',labels:'标签组合',dataCenter:'数据中心',replicas:'20',time:'20h'},
        {id: 9, name: 'Dashboard',labels:'标签组合',dataCenter:'数据中心',replicas:'20',time:'20h'}
    ]
});

Mock.mock('http://192.168.0.102:8080/api/v1/organizations/1/users/2/deployments/new', {
    "code": 0,
    "message": "",
    "data": {
        "orgId":  "1",
        "orgName": "Ops",
        "dataCenter": [{
            "dcId": "1",
            "name": "世纪互联",
            "budget": 10000000,
            "balance": 10000000
        },{
            "dcId": "2",
            "name": "电信机房",
            "budget": 10000000,
            "balance": 10000000
        },{
            "dcId": "3",
            "name": "电子城机房",
            "budget": 10000000,
            "balance": 10000000
        }],
        resourceSpec : [{
            quotaId : '1',
            name : '2X',
            cpu : '2 Core',
            memory : '2G'
        },{
            quotaId : '2',
            name : '4X',
            cpu : '4 Core',
            memory : '8G'
        }],
        "dcQuotas": {
            "dcId": "1",
            "PodMax": 1
            // 第一版用不到...
        }
    }
});

