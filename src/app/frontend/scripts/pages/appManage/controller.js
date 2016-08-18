/**
 * Created by Jora on 2016/7/29.
 */
define([
        'base64'
    ], function(Base64){
        'use strict';

        var ctrl = ['$scope', 'appManageService', '$sessionStorage', function($scope,appManageService, $sessionStorage){

            $scope.param = {"orgId": $sessionStorage.orgId, "userId": $sessionStorage.userId, "sessionId": $sessionStorage.sessionId}

            appManageService.getAppList($scope.param,function(data){
                 if (data.code == 0) {
                    console.log(data.data);
                    $scope.appList = data.data;
                 } else if (data.code == 1 ) {
                    console.log("Data: " + JSON.stringify(data))
                 }
            });

            console.log(appList);
        }];


        var controllers = [
            {module: 'appManage', name: 'appManageController', ctrl: ctrl}
        ];

        return controllers;
    }
);