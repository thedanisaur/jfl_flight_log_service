# JFL Flight Log Service

```
USER_ID=5a1bc8b9-f502-11f0-a8a7-74563c2abceb
curl -i -k -X POST \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
http://127.0.0.1:8082/flight-logs/$USER_ID \
-d '{ "id": "1388db31-f4ef-11f0-a8a7-74563c2abceb", "user_id": "$USER_ID", "unit_id": "1388db31-f4ef-11f0-a8a7-74563c2abceb", "mds": "MDS", "flight_log_date": "2024-01-18T23:36:00Z", "serial_number": "Serial Number", "unit_charged": "Unit Charged", "harm_location": "Harm Location", "flight_authorization": "Flight Authorization", "issuing_unit": "Issuing Unit", "comments": [ { "user_id": "$USER_ID", "role_id": "1388db31-f4ef-11f0-a8a7-74563c2abceb", "comment": "this is my comment" } ], "missions": [ { "mission_symbol": "Mission Symbol", "mission_from": "Mission From", "mission_to": "Mission To" } ], "aircrew": [ { "user_id": "$USER_ID", "flying_origin": "Flying Origin", "flight_auth_code": "Flight Auth Code" } ] }'
```