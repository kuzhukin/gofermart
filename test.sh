 curl localhost:8080/api/user/register -X POST -d '{"login":"user1","password":"1234"}' -i
 curl localhost:8080/api/user/login -X POST -d '{"login":"user1","password":"1234"}' -i
 curl localhost:8080/api/user/orders -X POST -d '12345678903' -i --cookie "Authorization=K7Cask9SYldqCOzkOJAOFLLpXARQLsqVoVY="

 curl localhost:8080/api/user/orders -X POST -d '4561261212345467' -i --cookie "Authorization=K7Cask9SYldqCOzkOJAOFLLpXARQLsqVoVY="

 curl localhost:8080/api/user/orders -X GET -i --cookie "Authorization=K7Cask9SYldqCOzkOJAOFLLpXARQLsqVoVY="
