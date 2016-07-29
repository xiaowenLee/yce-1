Mock.mock('http://10.69.40.36.com/api/main/navlist', {
    list: [
        {id: 1, name: 'Dashboard', state: 'main.dashboard',includeState: 'main.dashboard'},
        {id: 2, name: '应用管理', state: 'main.appManage',includeState: 'main.appManage',
            item: [
                {id: 21, name: '发布', state: 'main.appManage.deployment', includeState: 'main.resource.deployment'},
                {id: 22, name: '回滚', state: 'main.appManage.rollback', includeState: 'main.resource.library'},
                {id: 22, name: '滚动升级', state: 'main.appManage.rollup', includeState: 'main.resource.library'},
                {id: 22, name: '撤销', state: 'main.appManage.cancel', includeState: 'main.resource.library'},
                {id: 22, name: '历史发布', state: 'main.appManage.history', includeState: 'main.resource.library'}
            ]
        },
        {id: 3, name: '镜像管理', state: 'main.imageManage', includeState: 'main.imageManage',
            item: [
                {id: 31, name: '查找镜像', state: 'main.resource.search', includeState: 'main.resource.search'},
                {id: 32, name: '删除镜像', state: 'main.resource.delete', includeState: 'main.resource.delete'}
            ]
        },
        {id: 4, name: '云盘管理', state: 'main.cloudManage', includeState: 'main.cloudManage'},
        {id: 5, name: '扩展功能', state: 'main.extensions', includeState: 'main.extensions',
            item: [
                {id: 51, name: '创建服务', state: 'main.extensions.service', includeState: 'main.extensions.service'},
                {id: 52, name: '创建访问点', state: 'main.extensions.endpoint', includeState: 'main.extensions.endpoint'}
            ]
        },
        {id: 6, name: '计费&充值', state: 'main.costManage', includeState: 'main.costManage'}
    ]
});

