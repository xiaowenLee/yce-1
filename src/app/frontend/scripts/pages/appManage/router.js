define([], function(){
	'use strict';

	var router = {
		appManage: {
			url: '/appManage',
			templateUrl: 'views/appManage/appManage.html',
			controller: 'appManageController'
		}
	};

	return router;
});