define([], function(){
	'use strict';

	var router = {
		extensions: {
			url: '/extensions',
			templateUrl: 'views/extensions/extensions.html',
			controller: 'extensionsController'
		}
	};

	return router;
});