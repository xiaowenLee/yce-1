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
        }];



        var controllers = [
            {module: 'appManage', name: 'deploymentController', ctrl: ctrl}
        ];

        return controllers;
    }
);