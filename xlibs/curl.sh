# curl -X POST localhost:19818/v1/admin \
#     -H "Content-Type: application/json" \
#     -H "Authorization: Api-Key YOUR_API_KEY_HERE" \
#     -H "Origin: https://your-frontend-domain.com" \
#     -H "Access-Control-Allow-Origin: *" \
#     -d '{"name": "John", "role": "manager"}'

curl -X POST localhost:19818/v1/auth \
    -H "Content-Type: application/json" \
    -H "X-API-Key: $RE_UP_SECURE_API_KEY"


# ACCOUNT_ID
# 39edf942-75e9-4bca-b71a-29c161be9b28
# GROUP_ID
# 4b224c94-dd5f-46c9-bf5a-178d3de50986
# USER_ID
#
