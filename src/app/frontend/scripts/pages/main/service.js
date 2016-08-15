define([
	'utils'
	], function(utils){
	'use strict';

	var getApis = function($http){
		var apis = {};

		apis.login = function(param, success, error){
			return utils.http($http, 'post', '/api/v1/users/' + param.username + '/login', param, success, error);
		};

		apis.getNavlist = function(param, success, error){
			return utils.http($http, 'post', '/api/main/navlist', param, success, error);
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