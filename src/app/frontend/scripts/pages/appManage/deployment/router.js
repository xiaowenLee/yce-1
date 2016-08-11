define([], function(){
	'use strict';

	var router = {
		deployment: {
			url: '/deployment',
			templateUrl: 'views/appManage/deployment.html',
			controller: 'deploymentController'
		}
	};

	return router;
});