curl -X POST http://157.230.193.252:19818/one-time \
    -H "Content-Type: application/json" \
    -H "X-API-Key: $RE_UP_SECURE_API_KEY" \
    -H "Origin: http://localhost" \
    -H "Access-Control-Allow-Origin: *" \
    -d '{"name": "APEX", "nickname": "apex","group_code": "APEX", "email": "fastinsuretech@gmail.com", "phone_number": "+639156984277", "photo_url": "https://lh3.googleusercontent.com/a/ACg8ocLQgj_iSgaeF4xiErBpwfFXNsESSVx3JKcDb4O8YJTgWI7OTjo=s96-c", "uid": "N7yCd3kCViMA0jD3eNuv5rqKxgy1", "account_id": "7bcc41a1-4c22-457c-a4de-750ca8f5d419"}'


# curl -X POST localhost:19818/v1/auth/verify-id-token \
#     -H "Content-Type: application/json" \
#     -H "X-API-Key: $RE_UP_SECURE_API_KEY" \
#     -d '{
#       "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjkyODg2OGRjNDRlYTZhOThjODhiMzkzZDM2NDQ1MTM2NWViYjMwZDgiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiUmV2ZXJzZSBFbnRyb3B5IiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FDZzhvY0lZcjQyR1lPTUtVQXpJdDY5OExka1RFbm1ZUU1GbGtpcnFtdkxhTTRHUERkM29ZX0FqZkE9czk2LWMiLCJhZ2VudCI6dHJ1ZSwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL2Zhc3RpbnN1cmUtZjE4MDEiLCJhdWQiOiJmYXN0aW5zdXJlLWYxODAxIiwiYXV0aF90aW1lIjoxNzMyNDM0Mzc2LCJ1c2VyX2lkIjoiS1F2OTRKYzg0Yk5Od2NCMzlYMGN2d1pSQk51MSIsInN1YiI6IktRdjk0SmM4NGJOTndjQjM5WDBjdndaUkJOdTEiLCJpYXQiOjE3MzI0MzQzNzYsImV4cCI6MTczMjQzNzk3NiwiZW1haWwiOiJwaHRuNDU4QGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7Imdvb2dsZS5jb20iOlsiMTA4Mzc5MDAwNzQ0MTU5NDUwNjAzIl0sImVtYWlsIjpbInBodG40NThAZ21haWwuY29tIl19LCJzaWduX2luX3Byb3ZpZGVyIjoiZ29vZ2xlLmNvbSJ9fQ.b8K2promGQbx0yRZY0LXB5_yT6PTMM8HUeVSDLTPAIqw2LOsXSW0H9shVIS5vqKbUFhphOMRfpxfaVESlL_rU9ljcPDlRSJ-CJeE1Jo0UQbW9RYSqPgVRB-o70U96Ch_bFJrsw4FP9u1O8qgCMp4uy8QU4HAbdy_nvKYyj7nRmkE3wyJs8JPN7HHmfkebKqZfHiPlingj5A8-v-tFhUPdSLWKasYW42Fz65sDETGTqQHBzKsabtd0FQuZtM1xo_HEXirrkyqHgZg73r924EP87X_8dz6RCJo3o5LhFYi2lCmVK4qZ-b6EyCK3dXLCq-bYBeF4kYcNqL70flrF7WQ4g",
#       "uid": "KQv94Jc84bNNwcB39X0cvwZRBNu1",
#       "email": "phtn458@gmail.com",
#       "group_code": "APEX",
#       "refresh": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjkyODg2OGRjNDRlYTZhOThjODhiMzkzZDM2NDQ1MTM2NWViYjMwZDgiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiUmV2ZXJzZSBFbnRyb3B5IiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FDZzhvY0lZcjQyR1lPTUtVQXpJdDY5OExka1RFbm1ZUU1GbGtpcnFtdkxhTTRHUERkM29ZX0FqZkE9czk2LWMiLCJhZ2VudCI6dHJ1ZSwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL2Zhc3RpbnN1cmUtZjE4MDEiLCJhdWQiOiJmYXN0aW5zdXJlLWYxODAxIiwiYXV0aF90aW1lIjoxNzMyNDM0Mzc2LCJ1c2VyX2lkIjoiS1F2OTRKYzg0Yk5Od2NCMzlYMGN2d1pSQk51MSIsInN1YiI6IktRdjk0SmM4NGJOTndjQjM5WDBjdndaUkJOdTEiLCJpYXQiOjE3MzI0MzQzNzYsImV4cCI6MTczMjQzNzk3NiwiZW1haWwiOiJwaHRuNDU4QGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7Imdvb2dsZS5jb20iOlsiMTA4Mzc5MDAwNzQ0MTU5NDUwNjAzIl0sImVtYWlsIjpbInBodG40NThAZ21haWwuY29tIl19LCJzaWduX2luX3Byb3ZpZGVyIjoiZ29vZ2xlLmNvbSJ9fQ.b8K2promGQbx0yRZY0LXB5_yT6PTMM8HUeVSDLTPAIqw2LOsXSW0H9shVIS5vqKbUFhphOMRfpxfaVESlL_rU9ljcPDlRSJ-CJeE1Jo0UQbW9RYSqPgVRB-o70U96Ch_bFJrsw4FP9u1O8qgCMp4uy8QU4HAbdy_nvKYyj7nRmkE3wyJs8JPN7HHmfkebKqZfHiPlingj5A8-v-tFhUPdSLWKasYW42Fz65sDETGTqQHBzKsabtd0FQuZtM1xo_HEXirrkyqHgZg73r924EP87X_8dz6RCJo3o5LhFYi2lCmVK4qZ-b6EyCK3dXLCq-bYBeF4kYcNqL70flrF7WQ4g"
#     }'


# ACCOUNT_ID
# 39edf942-75e9-4bca-b71a-29c161be9b28
# GROUP_ID
# 4b224c94-dd5f-46c9-bf5a-178d3de50986
# USER_ID
#
