/**
 * Created by Jora on 2016/7/29.
 */
define([
        'base64'
    ], function(Base64){
        'use strict';

        var ctrl = ['$scope', 'deploymentService','$stateParams', function($scope, deploymentService, $stateParams){
            $scope.param = {
                dataCenter : []
            };
            $scope.stepNum = 1;

            deploymentService.getDeploymentIint({
                orgId : $stateParams.orgId,
                userId :$stateParams.userId
            },function(data){
                if(data.code == 0){
                    $scope.initData = data.data;
                }
                $scope.nextStep = function(stepNum){
                    $scope.stepNum = stepNum;
                };
            });

            // Image
            $scope.shows=false;
            //  模拟
            $scope.names = [
                {"name":"name"},
                {"name":"number"},
                {"name":"sex"},
                {"name":"del"},
                {"name":"name1"},
                {"name":"number1"},
                {"name":"sex1"},
                {"name":"del1"},
                {"name":"name2"},
                {"name":"number2"},
                {"name":"sex2"},
                {"name":"del2"},
                {"name":"name3"},
                {"name":"number3"},
                {"name":"sex3"},
                {"name":"del3"},
                {"name":"name4"},
                {"name":"number4"},
                {"name":"sex4"},
                {"name":"del4"}
            ];
        }];



        var controllers = [
            {module: 'appManage', name: 'deploymentController', ctrl: ctrl}
        ];

        return controllers;
    }
);