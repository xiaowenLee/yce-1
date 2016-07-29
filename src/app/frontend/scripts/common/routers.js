define([
	'../pages/main/router',
	'../pages/dashboard/router',
	'../pages/appManage/router',
	'../pages/cloudManage/router',
	'../pages/costManage/router',
	'../pages/extensions/router',
	'../pages/imageManage/router'
		], function(mainRouter, dashboardRouter, appManageRouter, cloudManageRouter, costManageRouter, extensionsRouter, imageManageRouter){

		'use strict';

		var init = function(app){
			app.config(['$stateProvider', '$urlRouterProvider', function($stateProvider, $urlRouterProvider){
				$urlRouterProvider.otherwise("/index");

				$stateProvider
					.state('index', mainRouter.index)
					.state('main', mainRouter.main)
					.state('main.dashboard', dashboardRouter.dashboard)
					.state('main.appManage', appManageRouter.appManage)
					.state('main.cloudManage', cloudManageRouter.cloudManage)
					.state('main.costManage', costManageRouter.costManage)
					.state('main.extensions', extensionsRouter.extensions)
					.state('main.imageManage', imageManageRouter.imageManage)
			}]);
		};

		return {
			init: init
		};	
	}
);