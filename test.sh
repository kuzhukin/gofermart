curl localhost:8080/api/user/register -X POST -d '{"login":"user1","password":"1234"}' -i
curl localhost:8080/api/user/register -X POST -d '{"login":"user3","password":"4321"}' -i

curl localhost:8080/api/user/login -X POST -d '{"login":"user1","password":"1234"}' -i


curl localhost:8080/api/user/orders -X POST -d '12345678903' -i --cookie "Authorization=K7Cask9SYldqCOzkOJAOFLLpXARQLsqVoVY="

curl localhost:8080/api/user/orders -X POST -d '4561261212345467' -i --cookie "Authorization=K7Cask9SYldqCOzkOJAOFLLpXARQLsqVoVY="

curl localhost:8080/api/user/orders -X POST -d '4561261212345467' -i --cookie "Authorization=K7Cask1SZ1ZrDTsMhpnBkiMaN9iVfpXTPxM="

curl localhost:8080/api/user/orders -X GET -i --cookie "Authorization=K7Cask9SYldqCOzkOJAOFLLpXARQLsqVoVY="


curl localhost:8080/api/user/balance -X GET -i --cookie "Authorization=K7Cask9SYldqCOzkOJAOFLLpXARQLsqVoVY="

curl localhost:8080/api/user/balance/withdraw -X POST -i --cookie "Authorization=K7Cask9SYldqCOzkOJAOFLLpXARQLsqVoVY=" -d '{"order":"4561261212345467","sum":750}'


## acural
curl localhost:33555/api/orders -X POST -d '{"order": "12345678903", "goods": [{"description": "111", "price": 7000}]}' -H 'Content-type: application/json' -i

curl localhost:33555/api/orders -X POST -d '{"order": "4561261212345467", "goods": [{"description": "111", "price": 7000}]}' -H 'Content-type: application/json' -i

curl localhost:33555/api/orders/4561261212345467 -X GET
