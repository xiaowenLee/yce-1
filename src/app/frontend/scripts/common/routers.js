define([
	'../pages/main/router',
	'../pages/dashboard/router',
	'../pages/appManage/router',
	'../pages/rbdManage/router',
	'../pages/costManage/router',
	'../pages/extensions/router',
	'../pages/imageManage/router',
	'../pages/appManage/deployment/router'
		], function(mainRouter, dashboardRouter, appManageRouter, rbdManageRouter, costManageRouter, extensionsRouter, imageManageRouter,deploymentRouter){

		'use strict';

		var init = function(app){
			app.config(['$stateProvider', '$urlRouterProvider', '$httpProvider', function($stateProvider, $urlRouterProvider, $httpProvider){
				//$urlRouterProvider.otherwise("/login");

				$stateProvider
					.state('login', mainRouter.login)
					.state('main', mainRouter.main)
              //      .state('cancel', mainRouter.cancel)
					.state('main.dashboard', dashboardRouter.dashboard)
					.state('main.appManage', appManageRouter.appManage)
					.state('main.appManageDeployment', deploymentRouter.deployment)
					.state('main.appManageRollback', appManageRouter.rollback)
					.state('main.appManageRollup', appManageRouter.rollup)
					.state('main.appManageCancel', appManageRouter.cancel)
					.state('main.appManageHistory', appManageRouter.history)
					.state('main.rbdManage', rbdManageRouter.rbdManage)
					.state('main.costManage', costManageRouter.costManage)
					.state('main.extensions', extensionsRouter.extensions)
					.state('main.extensionsService', extensionsRouter.service)
					.state('main.extensionsEndpoint', extensionsRouter.endpoint)
					.state('main.imageManage', imageManageRouter.imageManage)
					.state('main.imageManageSearch', imageManageRouter.search)
					.state('main.imageManageDelete', imageManageRouter.delete)
			}]);
		};

		return {
			init: init
		};	
	}
);