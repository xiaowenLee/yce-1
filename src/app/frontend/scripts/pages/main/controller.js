/**
 * Created by Jora on 2016/7/29.
 */
define([
        'base64'
    ], function(Base64){
        'use strict';

        var ctrl = ['$scope','$state','mainService', '$sessionStorage', '$stateParams', function($scope, $state,mainService, $sessionStorage, $stateParams){
            // login
            $scope.login = function () {
                mainService.login({
                    'username': $scope.username,
                    'password': $scope.pwd
                }, function (data) {
                    if (data.code == 0) {
                        $scope.loginData = JSON.parse(data.data);
                        $sessionStorage.username =  $scope.loginData.userName;
                        $sessionStorage.sessionId =  $scope.loginData.sessionId;
                        $sessionStorage.orgId = $scope.loginData.orgId;
                        $scope.jump();
                    }else{
                        alert("用户名密码错误");
                    }
                });
            };

            // logout
            $scope.logout = function(){
                mainService.logout({
                    'username':  $sessionStorage.username,
                    'sessionId': $sessionStorage.sessionId
                }, function(data) {
                    if (data.code == 0) {
                        $sessionStorage.$reset();
                        // alert("退出成功~");
                        $state.go('login');
                    }
                });
            }

            $scope.jump = function(){
                $state.go('main.dashboard');
                $scope.data = {
                    username : $sessionStorage.username,
                    showSubnav: [],
                    toggleNav : false
                };
                mainService.getNavlist(null, function (data) {
                    $scope.navList = data.list;
                });

                $scope.showSubnav = function (index) {
                    $scope.data.showSubnav[index] = !$scope.data.showSubnav[index];
                };
            };

            if(!$sessionStorage.username) {
                $state.go('login');
            }else{
                $scope.jump();
            }
        }];



        var controllers = [
            {module: 'yce-manage', name: 'mainController', ctrl: ctrl}
        ];

        return controllers;
    }
);





