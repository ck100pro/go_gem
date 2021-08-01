<!-- get get_jwt_token termnal command -->
curl -v -X GET \
http://localhost:8080/get_jwt_token \
-H 'content-type: application/json' \
-d '{ "account": "hello_account", "password": "12345678" }'

<!-- POST valid_jwt_token termnal command -->
curl -v -X POST \
http://localhost:8080/valid_jwt_token \
-H 'content-type: application/json' \
-H 'Authorization: Bearer <insert_your_JWT_here>' \