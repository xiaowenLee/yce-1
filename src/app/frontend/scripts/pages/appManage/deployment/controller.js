/**
 * Created by Jora on 2016/7/29.
 */
define([
        'base64'
    ], function(Base64){
        'use strict';

        var ctrl = ['$scope', 'deploymentService', function($scope, deploymentService){
            $scope.param = {
                dataCenter : []
            }

            deploymentService.getDeploymentIint(null,function(data){
                if(data.code == 0){
                    $scope.initData = data.data;
                }
                $scope.nextStep = function(){
                    console.log($scope.param);
                }
            });
        }];



        var controllers = [
            {module: 'appManage', name: 'deploymentController', ctrl: ctrl}
        ];

        return controllers;
    }
);