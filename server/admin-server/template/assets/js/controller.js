angular.module('viqeoJs', ['ngCookies', 'ngRoute'])
    .controller('DeployController', function($rootScope, dataService, $cookies, $scope) {

        $scope.loading = false;
        $rootScope.debug = false;
        $rootScope.tab = 'tag';

        $rootScope.changeTab = function (tab) {
            $rootScope.tab = tab;
        };

        if ($cookies.get('debug') == "1") {
            $rootScope.debug = true;
        }

        $rootScope.debugAction = function () {
            dataService.debugAction().then(function (response) {
                if ($cookies.get('debug') == "1") {
                    $rootScope.debug = true;
                } else {
                    $rootScope.debug = false;
                }
            })
        };

        $rootScope.checkoutCommit = function(id){
            $rootScope.loading = true;
            dataService.checkoutCommit(id).then(function (response) {
                $rootScope.loading = false;
                $rootScope.showAlert = true;
                $rootScope.message = response.data;
                dataService.getCommits();
            })
        };

        $rootScope.checkoutTag = function(tag){
            $rootScope.loading = true;
            dataService.checkoutTag(tag).then(function (response) {
                $rootScope.showAlert = true;
                $rootScope.loading = false;
                $rootScope.message = response.data;
                dataService.getTags();
            })
        };
        dataService.getTags();
        dataService.getCommits();

    })
    .controller('CustomController', function($rootScope, $scope, dataService, $window) {

        dataService.getBranches();
        dataService.getCustoms();

        $scope.custom = {
            "params": []
        };

        $scope.save = function (model) {
            dataService.saveCustom(model).then(function(){
                dataService.getCustoms();
                $('#createCustom').modal('toggle');
            });
        };

        $scope.update = function (model) {
            dataService.patchCustom(model).then(function(){
                dataService.getCustoms();
                $('#createCustom').modal('toggle');
            });
        };

        $scope.edit = function (model) {
            $scope.custom = model;
            $scope.isEdit = true;
            $('#createCustom').modal('toggle');
        };

        $scope.new = function () {
            $scope.isEdit = false;
            $scope.custom = {
                "params": []
            };
            $('#createCustom').modal('toggle');
        };

        $scope.build = function (model) {
            model.loading = true;
            dataService.buildCustom(model).then(function(response) {
                dataService.getCustoms();
            })
        };

        $scope.remove = function (model) {
            model.loading = true;
            dataService.removeCustom(model).then(function(response) {
                dataService.getCustoms();
            })
        };

        $scope.delParam = function (index, type) {
            $scope.custom[type].splice(index, 1);
        };

        $scope.addParam = function (type) {
            if ($scope.custom[type] == undefined) {
                $scope.custom[type] = [];
            }
            $scope.custom[type].push({
                "Key" : "",
                "Value" : ""
            })
        }

    })
    .service('dataService', function($rootScope, $http, $window) {

        delete $http.defaults.headers.common['X-Requested-With'];

        this.getTags = function() {
            $http.get("/api/tags").then(function(response) {
                $rootScope.tags = response.data
            });
        };

        this.getBranches = function() {
            $http.get("/api/branch").then(function(response) {
                $rootScope.branches = response.data;

                // чет позабыл я уже как скоуп биндить на изменения... так что сорян, пока так
                setTimeout(function(){
                    $('.branch-list select').selectpicker('refresh');
                }, 1000)
            });
        };

        this.getCommits = function() {
            $http.get("/api/commits").then(function(response) {
                $rootScope.commits = response.data
            });
        };

        this.checkoutCommit = function(id) {
            return $http.get("/api/checkout/commit/" + id)
        };

        this.debugAction = function() {
            return $http.get("/api/debug?d=" + ($rootScope.debug ? "0" : "1"))
        };

        this.checkoutTag = function(tag) {

            var result = tag.match( /([^tag: ][^,]*)/g );

            if (result !== undefined && result.length > 0) {
                return $http.get("/api/checkout/tag/" + result[0])
            } else {
            }
        };

        this.buildCustom = function (model) {
            return $http.post("/api/custom/build", JSON.stringify(model))
        };

        this.saveCustom = function(model) {
            return $http.post("/api/custom/save", JSON.stringify(model))
        };

        this.patchCustom = function(model) {
            return $http.patch("/api/custom/" + model.Name, JSON.stringify(model))
        };

        this.removeCustom = function(model) {
            return $http.post("/api/custom/remove", JSON.stringify(model))
        };

        this.getCustoms = function() {
            $http.get("/api/custom/list").then(function(response) {

                var items;

                if (response.data.results) {

                    items = response.data.results;

                    for (var i in items) {
                        var query = [];
                        if (items[i].Params.length) {
                            items[i].Params.forEach(function (p) {
                                query.push(p.Key + "=" + p.Value);
                            });
                        }

                        items[i].Query = query.join('&');

                        tBuild = new Date(items[i].Build).getTime();
                        tUpdate = new Date(items[i].Update).getTime();
                        if (tUpdate != tBuild) {
                            items[i].needBuild = true;
                        } else {
                            items[i].needBuild = false;
                        }


                    }
                }
                $rootScope.customs = items
            });
        }

    })
    .config(function($routeProvider) {
        $routeProvider
            .when("/", {
                templateUrl : "assets/main.html"
            })
            .when("/custom", {
                templateUrl : "assets/custom.html"
            });
    });