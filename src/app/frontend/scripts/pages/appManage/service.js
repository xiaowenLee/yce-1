define([
	'utils'
	], function(utils){
	'use strict';

	var getApis = function($http){
		var apis = {};

		apis.getAppList = function(param, success, error, $sessionStorage){
		    var orgId = param.orgId
		    var userId = param.userId
            // console.log("appManage service: userId=" + userId + ", orgId=" + orgId)
			// return utils.http($http, 'get', '/api/v1/organizations/' + orgId + '/users/' + userId + '/deployments', param, success, error);
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