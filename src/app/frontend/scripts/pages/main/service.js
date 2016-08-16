define([
	'utils'
	], function(utils){
	'use strict';

	var getApis = function($http){
		var apis = {};

		apis.login = function(param, success, error){
			return utils.http($http, 'post', '/api/v1/users/login', param, success, error);
		};

		apis.getNavlist = function(param, success, error){
<<<<<<< HEAD
			return utils.http($http, 'get', '/api/main/navlist', param, success, error);
=======
			return utils.http($http, 'get', '/api/v1/navlist', param, success, error);
>>>>>>> 0dfcf3fca61df983c4b19bff9df1baa9f6ead5e8
		};

		return apis;
	};	

	var services = {
		module: 'yce-manage',
		name: 'mainService',	
		getApis: getApis
	};

	return services;
});