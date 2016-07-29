Mock.mock('http://10.69.40.36.com/api/main/navlist', {
    list: [
        {id: 1, name: 'Dashboard', state: 'main.dashboard'},
        {id: 3, name: '游戏运营', state: 'main.operate', includeState: 'main.operate'},
        {id: 2, name: '资源管理', state: 'main.resource', includeState: 'main.resource',
            subNavList: [
                {id: 21, name: '云主机', state: 'main.resource.cloud', includeState: 'main.resource.cloud'},
                {id: 22, name: '程序库', state: 'main.resource.library', includeState: 'main.resource.library'}
            ]
        }
    ]
});

