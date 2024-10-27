# curl -X POST localhost:19818/v1/admin \
#     -H "Content-Type: application/json" \
#     -H "Authorization: Api-Key YOUR_API_KEY_HERE" \
#     -H "Origin: https://your-frontend-domain.com" \
#     -H "Access-Control-Allow-Origin: *" \
#     -d '{"name": "John", "role": "manager"}'

curl -X POST localhost:19818/v1/auth \
    -H "Content-Type: application/json" \
    -H "X-API-Key: $RE_UP_SECURE_API_KEY" \
