/**
 * Created by Jora on 2016/7/29.
 */
define([
        'base64'
    ], function(Base64){
        'use strict';

        var ctrl = ['$scope','$state','mainService', '$sessionStorage', function($scope, $state,mainService, $sessionStorage){
            $scope.login = function () {
                mainService.login({
                    'username': $scope.username,
                    'password': $scope.pwd
                }, function (data) {
                    if (data.code == 0) {
                        alert('登录成功！');
                        $sessionStorage.login = true;
                        $scope.jump();
                    }
                });
            };
            $scope.logout = function(){
                delete $sessionStorage.login;
                alert('退出成功！');
                $state.go('login');
            }
            $scope.jump = function(){
                $state.go('main.dashboard');
                $scope.data = {
                    showSubnav: []
                };

                mainService.getNavlist(null, function (data) {
                    $scope.navList = data.list;
                });

                $scope.showSubnav = function (index) {
                    $scope.data.showSubnav[index] = !$scope.data.showSubnav[index];
                };
            };

            if(!$sessionStorage.login) {
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






/*



define([
        'base64'
    ], function(Base64){
        'use strict';

        var ctrl = ['$scope','$state','mainService', '$sessionStorage', function($scope, $state,mainService, $sessionStorage){
            $scope.login = function () {
                mainService.login({
                    'username': $scope.username,
                    'password': $scope.pwd
                }, function (data) {
                    if (data.code == 0) {
                        alert('登录成功！');
                        $sessionStorage.login = true;
                        $scope.jump();
                    }
                });
            };
            $scope.logout = function(){
                delete $sessionStorage.login;
                alert('退出成功！');
                alert('ddddd');
                $state.go('cancel');
            }
            $scope.jump = function(){
                $state.go('main.dashboard');
                $scope.data = {
                    showSubnav: []
                };

                mainService.getNavlist(null, function (data) {
                    $scope.navList = data.list;
                });

                $scope.showSubnav = function (index) {
                    $scope.data.showSubnav[index] = !$scope.data.showSubnav[index];
                };
            };

//            if(!$sessionStorage.login) {
//                $state.go('login');
//            }else{
//                $scope.jump();
//            }
        }];



        var controllers = [
            {module: 'yce-manage', name: 'mainController', ctrl: ctrl}
        ];

        return controllers;
    }
);
*/
