/**
 * Created by Jora on 2016/7/29.
 */
define([
        'base64'
    ], function(Base64){
        'use strict';

        var ctrl = ['$scope','$http','deploymentService','$sessionStorage', function($scope,$http, deploymentService, $sessionStorage){

            $scope.param = {"orgId": $sessionStorage.orgId, "userId": $sessionStorage.userId, "sessionId": $sessionStorage.sessionId}

            deploymentService.getDeploymentIint($scope.param, function(data){
                if(data.code == 0){
                    $scope.initData = JSON.parse(data.data);
                    console.log(JSON.stringify($scope.initData))
                }
            });

            // Image
            $scope.shows=false;
            //  模拟
            /*
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
            */
            //$scope.getImages = function() {
                console.log("getImages test")
                $http({
                    method: 'GET',
                    url: '/api/v1/registry/images'
                })
                .success(function(data) {
                    var dataObject = JSON.parse(data.data);
                    console.log("getImages success");


                    // cycle print images name and tags
                    /*
                    $scope.imageList=dataObject;
                    for (var i in dataObject) {
                        console.log("images: " + dataObject[i].name);
                        var list = dataObject[i].tags;
                        $scope.tagList=list;
                        for (var j in list) {
                            console.log("tags: " + list[j]);
                        }
                    }
                    */

                    // make new images:tags
                    var imageArr = new Array();
                    var k = 0
                    for (var i in dataObject) {
                        var list = dataObject[i].tags;
                        for (var j in list) {
                            imageArr[k] = dataObject[i].name + ":" + list[j]
                            k=k+1
                        }
                    }

                    $scope.imageList=imageArr;
                    $scope.getImages = function(x) {
                        $scope.param.image=x;
                    }
                    /*
                    $scope.getImages = function($index) {
                        $scope.param.image=imageArr[$index];
                    }
                    */

                })
                .error(function() {
                    console.log("getImages error")
                })
            //}
        }];



        var controllers = [
            {module: 'appManage', name: 'deploymentController', ctrl: ctrl}
        ];

        return controllers;
    }
);