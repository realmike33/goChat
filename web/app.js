angular.module('Chat', ['ui.router'])

.config(function($stateProvider, $urlRouterProvider){
  $urlRouterProvider.otherwise('/');
  $stateProvider
    .state('chat', {
      url: '/',
      templateUrl: 'templates/chat.html',
      controller: 'ChatCtrl'
    })
})

.controller('ChatCtrl', function($scope){

  $scope.messages = [];

  var conn = new WebSocket('ws://localhost:8000/ws');

  conn.onclose = function(e){
    $scope.$apply(function(){
      $scope.connection = "DISCONNECTED";
    })
  };

  conn.onopen = function(e){
    $scope.$apply(function(){
      $scope.connection = "CONNECTED";
    })
  };

  conn.onmessage = function(e){
    var userArray = e.data.split(',')
    var username = $scope.username = userArray[0];
    var message = userArray[1];
    $scope.$apply(function(){
      $scope.messages.push({name: username, message: message});
    })
  }

  $scope.send = function(){
    if($scope.name){
      var arry = [$scope.name, $scope.msg]
      conn.send(arry);
      $scope.msg = '';
    } else{
      alert('Insert username')
    }
  }
})
