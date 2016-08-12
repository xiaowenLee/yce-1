define([
	'utils'
	], function(utils){
	'use strict';

	var getApis = function($http){
		var apis = {};

		apis.getAppList = function(param, success, error){
			return utils.http($http, 'get', '/api/appManage/appList', param, success, error);
		};
		return apis;
	};	

	var services = {
		module: 'appManage',
		name: 'appManageService',
		getApis: getApis
	};

	return services;
});