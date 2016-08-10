function index($scope,$http,$location){
    
    $scope.formData={}
    $scope.postData=function(){
        // console.log($scope.formData)
        // console.log($scope.formData.name)

        $http({
            method  : 'POST',
            url     : '/api/v1/users/'+ $scope.formData.name +'/login',
            data    : 'password=' + $scope.formData.password,
            headers : { 'Content-Type': 'application/x-www-form-urlencoded' }
        })
        .success(function(data){
            console.log(JSON.stringify(data) + "-------console的data");
            if(data.code == 0){
                alert(123);
                $scope.getInfo=function(){
                    $location.path('/user/info')
                }
            }
        });
    }   
}

function info($scope,$http, $cookies){
    
    $scope.out=function(){
       // alert(1234);
       alert($cookies.get("sessionId"))
       $http({
           method  : 'GET',
           url     : '/api/v1/users/email/logout?sessionld=20acdcd4-2e9e-475c-805c-6e943b7442ef',
           headers : { 'Content-Type': 'application/x-www-form-urlencoded' }
       })
       .success(function(data){
           alert(退出成功);
       });
    }
    
}

angular.module('app')
       .controller('index',index)
       .controller('info',info)