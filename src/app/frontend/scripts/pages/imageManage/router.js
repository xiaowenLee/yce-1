define([], function(){
	'use strict';

	var router = {
		imageManage: {
			url: '/imageManage',
			templateUrl: 'views/imageManage/imageManage.html',
			controller: 'imageManageController'
		}
	};

	return router;
});