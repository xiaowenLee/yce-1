define([], function(){
	var yce = {};
	yce.preUrl = 'http://10.69.40.36.com';

	yce.http = function($http, method, url, param, success, error){
		$http[method](yce.preUrl + url, param, {headers: {'Content-Type': 'application/x-www-form-urlencoded'}})
		.success(function(data){
			success(data);
		})
		.error(function(){
			if(error && typeof error == 'function') return error();
			console.log('error');
		});
	};

	return yce;
});