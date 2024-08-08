### Ok: Ping-pong
GET {{url_base}}/ping
Origin: http://any-origin.io
Content-type: any-content-type
HTTP 200
content-type: application/json; charset=utf-8
content-length: 18
{"message":"pong"}

### Fail: 404: Route not found
GET {{url_base}}/unknown
HTTP 404
[Asserts]
body == "404 page not found"

### Ok: Add customer
POST {{url_base}}/customer
Content-type: application/json
{"id":"customer-id-1","name":"customer-name-1","created_at":1000001,"address":{"id":"shipping-address-id-1","city":"Riga","street":"Brivibas 123","longitude":24.11,"latitude":56.01}}
HTTP 201
content-type: application/json; charset=utf-8
content-length: 257
{"created":true,"customer":{"id":"customer-id-1","name":"customer-name-1","address":{"id":"shipping-address-id-1","city":"Riga","street":"Brivibas 123","longitude":24.11,"latitude":56.01},"created_at":1000001},"message":"Customer info created successfully"}

POST {{url_base}}/customer
Content-type: application/json
{"id":"customer-id-2","name":"customer-name-2","created_at":1000002}
HTTP 201
content-type: application/json; charset=utf-8
content-length: 158
{"created":true,"customer":{"id":"customer-id-2","name":"customer-name-2","address":null,"created_at":1000002},"message":"Customer info created successfully"}

### Ok: Get customer by id
GET {{url_base}}/customer/customer-id-1
Content-type: application/json
HTTP 200
content-type: application/json; charset=utf-8
content-length: 195
{"customer":{"id":"customer-id-1","name":"customer-name-1","address":{"id":"shipping-address-id-1","city":"Riga","street":"Brivibas 123","longitude":24.11,"latitude":56.01},"created_at":1000001}}

### Ok: Get customer list
GET {{url_base}}/customers
Content-type: application/json
HTTP 200
content-type: application/json; charset=utf-8
content-length: 282
{"customers":[{"id":"customer-id-1","name":"customer-name-1","address":{"id":"shipping-address-id-1","city":"Riga","street":"Brivibas 123","longitude":24.11,"latitude":56.01},"created_at":1000001},{"id":"customer-id-2","name":"customer-name-2","address":null,"created_at":1000002}]}

### Fail: Delete customer by id
DELETE {{url_base}}/customer/absent-customer-id
Content-type: application/json
HTTP 404
content-type: application/json; charset=utf-8
content-length: 105
{"deleted":false,"id":"absent-customer-id","message":"Customer info not found by id: absent-customer-id"}

### Ok: Delete customer by id
DELETE {{url_base}}/customer/customer-id-1
Content-type: application/json
HTTP 200
content-type: application/json; charset=utf-8
content-length:84
{"deleted":true,"id":"customer-id-1","message":"Customer info successfully deleted"}

### Fail: Get customer by id
GET {{url_base}}/customer/customer-id-1
Content-type: application/json
HTTP 404
content-type: application/json; charset=utf-8
content-length: 74
{"id":"customer-id-1","message":"Customer not found by id: customer-id-1"}
