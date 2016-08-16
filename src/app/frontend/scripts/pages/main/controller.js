/**
 * Created by Jora on 2016/7/29.
 */
define([
        'base64'
    ], function(Base64){
        'use strict';

        var ctrl = ['$scope','$state','mainService', '$sessionStorage', function($scope, $state,mainService, $sessionStorage){
            // login
            $scope.login = function () {
                mainService.login({
                    'username': $scope.username,
                    'password': $scope.pwd
                }, function (data) {
                    if (data.code == 0) {
                        // alert('登录成功！');
                        $sessionStorage.username = JSON.parse(data.data).userName;
                        $sessionStorage.sessionId = JSON.parse(data.data).sessionId;
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





