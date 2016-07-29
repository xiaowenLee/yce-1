define([
	'utils'
	], function(utils){
	'use strict';

	var getApis = function($http){
		var apis = {};

		return apis;
	};	

	var services = {
		module: 'cloudManage',
		name: 'cloudManageService',
		getApis: getApis
	};

	return services;
});