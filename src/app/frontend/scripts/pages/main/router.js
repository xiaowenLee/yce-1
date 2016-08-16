define([], function(){
	'use strict';

	var router = {
		login: {
			url: '/login',
			templateUrl: 'views/main/login.html',
			controller: 'mainController'
		},
		main: {
			url: '/:orgId/:userId/main',
			templateUrl: 'views/main/main.html',
			controller: 'mainController'
		}		
	};

	return router;
});


/*
define([], function(){
	'use strict';

	var router = {
		login: {
			url: '/login',
			templateUrl: 'views/main/login.html',
			controller: 'mainController'
		},
		main: {
			url: '/main',
			templateUrl: 'views/main/main.html',
			controller: 'mainController'
		},
        cancel: {
            url:'/cancel',
            templateUrl: 'views/main/cancel.html',
            controller: 'mainController'
        }  
	};

	return router;
});
*/