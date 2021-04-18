'use strict';
var app = angular.module('application', []);
app.controller('AppCtrl', function($scope, appFactory){
        $("#success_getcertificate").hide();
        $("#success_setJobPosting").hide();
        $("#success_setRating").hide();
        $("#success_getwallet").hide();
        $("#success_setApply").hide();
        $("#success_getJobPosting").hide();
        
        $scope.getJobPosting = function(){
                alert("성공적으로 조회하였습니다.");
                appFactory.getJobPosting(function(data){
                        var array = [];
                        for (var i = 0; i < data.length; i++){
                                parseInt(data[i].Key);
                                data[i].Record.Key = data[i].Key;
                                array.push(data[i].Record);
                                $("#success_getJobPosting").hide();
                        }
                        array.sort(function(a, b) {
                            return parseFloat(a.Key) - parseFloat(b.Key);
                        });
                        $scope.getJobPosting = array;
                });
        }
        $scope.setJobPosting = function(){
                alert("구인공고 게시 완료");
                appFactory.setJobPosting($scope.offer, function(data){
                        $scope.create_offer = data;
                        $("#success_setjobposting").show();
                });
        }
        
        $scope.setRating = function(){
                alert("평점 설정 완료");
                appFactory.setRating($scope.Offer, function(data){
                        $scope.set_rating = data;
                        $("#success_setRating").show();
                });
        }
        $scope.setApply = function(dockey,psword, key){
                alert("지원 완료");
                appFactory.setApply(dockey,psword ,key, function(data){
                        var array = [];
                        for (var i = 0; i < data.length; i++){
                                parseInt(data[i].Key);
                                data[i].Record.Key = data[i].Key;
                                array.push(data[i].Record);
                                $("#success_getJobPosting").show();
                        }
                        array.sort(function(a, b) {
                            return parseFloat(a.Key) - parseFloat(b.Key);
                        });
                        $scope.getJobPosting = array;
                });
        }
        $scope.getWallet = function(){
                appFactory.getWallet($scope.walletid, function(data){
                        $scope.search_wallet = data;
                        $("#success_getwallet").show();
                });
        }

});
 app.factory('appFactory', function($http){
        var factory = {};
        factory.getWallet = function(key, callback){
                $http.get('/api/getWallet?walletid='+key).success(function(output){
                            callback(output)
                });
        }
        factory.getCertificate = function(key, callback){
                $http.get('/api/getCertificate?certkey='+key).success(function(output){
                        callback(output)
                });
        }
        factory.getJobPosting = function(callback){
            $http.get('/api/getJobPosting/').success(function(output){
                        callback(output)
                });
        }
        factory.setJobPosting = function(data, callback){
            $http.get('/api/setJobPosting?field='+data.field+'&name='+data.name+'&psword='+data.psword+'&conditions='+data.conditions+'&pay='+data.pay+'&endDate='+data.endDate).success(function(output){
                        callback(output)
                });
        }
        factory.setRating = function(data, callback){
                $http.get('/api/setRating?companykey='+data.companykey+'&psword='+data.psword+'&rating='+data.rating+'&jobpostingkey='+data.jobpostingkey+'&freelancerkey='+data.freelancerkey).success(function(output){
                callback(output)
                });
        }
        factory.setApply = function(dockey,psword, key, callback){
                $http.get('/api/setApply?dockey='+dockey+'&psword='+psword+'&freelancerkey='+key).success(function(output){
                         $http.get('/api/getJobPosting/').success(function(output){
                                            callback(output)
                                    });
                });
            }
        
        return factory;
});