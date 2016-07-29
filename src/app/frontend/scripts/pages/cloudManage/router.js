define([], function(){
	'use strict';

	var router = {
		cloudManage: {
			url: '/cloudManage',
			templateUrl: 'views/cloudManage/cloudManage.html',
			controller: 'cloudManageController'
		}
	};

	return router;
});