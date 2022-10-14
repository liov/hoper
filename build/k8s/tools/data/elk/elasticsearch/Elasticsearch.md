修改密码 curl --location --request POST 'http://es.d/_security/user/elastic/_password' \
--header 'Authorization: Basic xxx' \
--header 'Content-Type: application/json' \
--data-raw '{
"password": "xxx"
}'