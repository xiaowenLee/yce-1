/**
 * Created by Jora on 2016/7/29.
 */
define([
        'base64'
    ], function(Base64){
        'use strict';

        var ctrl = ['$scope', '$rootScope', 'appManageService', '$sessionStorage', function($scope, $rootScope, appManageService, $sessionStorage){

            $scope.param = {"orgId": $sessionStorage.orgId, "userId": $sessionStorage.userId, "sessionId": $sessionStorage.sessionId}

            appManageService.getAppList($scope.param,function(data){
                 if (data.code == 0) {
                    $scope.appList = JSON.parse(data.data);
                    console.log($scope.appList);
                 }
            });

            $scope.showAppDetail = function(item){
                $scope.appDetailConf = {
                    widgetId : 'widgetAppDetail',
                    widgetTitle : '应用详情',
                    isAppDetail : true,
                    data : item
                };
                $rootScope.widget = {
                    showAppDetail : true
                };
            };

        }];


        var controllers = [
            {module: 'appManage', name: 'appManageController', ctrl: ctrl}
        ];

        return controllers;
    }
);