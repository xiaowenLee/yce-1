function config($stateProvider,$urlRouterProvider){
    $urlRouterProvider.otherwise('/home/index');
    $stateProvider
    .state('home',{
        url:'/home',
        templateUrl:'view/content.html'
    })
    .state('home.index',{
        url:'/index',
        templateUrl:'view/index.html',
        data:{
            'title':''
        },
        controller:'index'
    })
    .state('user',{
        url:'/user',
        templateUrl:'view/content2.html'
    })
    .state('user.info',{
        url:'/info',
        templateUrl:'view/info.html',
        controller:'info'
    })

}

    angular.module('app')
    .config(config)