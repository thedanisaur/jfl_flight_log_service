# JFL Flight Log Service

Domain → Actual flight execution
* What actually happened

Insert Flight Log
```
USER_ID=5a1bc8b9-f502-11f0-a8a7-74563c2abceb
curl -i -k -X POST \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
http://127.0.0.1:8082/flight-logs/$USER_ID \
-d '{ "id": "1388db31-f4ef-11f0-a8a7-74563c2abceb", "user_id": "$USER_ID", "unit_id": "1388db31-f4ef-11f0-a8a7-74563c2abceb", "mds": "MDS", "flight_log_date": "2024-01-18T23:36:00Z", "serial_number": "Serial Number", "unit_charged": "Unit Charged", "harm_location": "Harm Location", "flight_authorization": "Flight Authorization", "issuing_unit": "Issuing Unit", "comments": [ { "user_id": "$USER_ID", "role_id": "1388db31-f4ef-11f0-a8a7-74563c2abceb", "comment": "this is my comment" } ], "missions": [ { "mission_symbol": "Mission Symbol", "mission_from": "Mission From", "mission_to": "Mission To" } ], "aircrew": [ { "user_id": "$USER_ID", "flying_origin": "Flying Origin", "flight_auth_code": "Flight Auth Code" } ] }'
```

Get Flight Log
```
USER_ID=3eb59016-f680-11f0-a8a7-74563c2abceb
FLIGHT_LOG_ID=d1da8ac1-eec9-4434-9e97-69460f9004d2
curl -i -k -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njg5NzgxMzgsImlhdCI6MTc2ODk2OTEzOCwiaXNzIjoiYXV0aC5qZmwuY29tIiwianRpIjoiMTc2NTJjOGUtMTlhYi00MjZkLWE1YmQtMTk5YzZlODBmNGE1IiwidXNlcl9pZCI6IjNlYjU5MDE2LWY2ODAtMTFmMC1hOGE3LTc0NTYzYzJhYmNlYiJ9.NN8w_J0nLj9T4cytYvIhhZ5RijlQWE-_NUFtnyieWHj3miX9Lb-DZySjuesK2Sma-DVRTBiAEmQomvmMjjF6uTnXh7V-7uMu80fEhmg7NHvzKvqfSUTMIJj0BxO5A5zKHIke_GvRlI6Qck5pBGr3sfayBzC4UoM8yK2oMk4crm2xNAPkStEj1vzkX7O1ZmBHJ-ENEj82LOCNmgEevUAcIgJMV26CsT_o54MPby8c3KtduF01UL1KsPdK3JDAyXNhEeh1sZFqbxwwZjV6VhjanCbVrcnGkeHW7L3x2jxVzC5ShpDOnmbR7d4ywxxO7rktzpEIbS7hOZNtYJ1buSHQ5g" \
http://127.0.0.1:8082/flight-logs/$USER_ID/$FLIGHT_LOG_ID
```