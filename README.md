# omise-SCB

#### start service
````
go run cmd/main.go
````
####  Create new payment request

````
curl --location --request POST 'localhost:8080/omise-scb' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount":10000
}'
````
#### Create new payment response

````
{
    "object": "charge",
    "id": "chrg_test_5pw21ew426s6wzuyz3m",
    "livemode": false,
    "location": "/charges/chrg_test_5pw21ew426s6wzuyz3m",
    "created": "0001-01-01T00:00:00Z",
    "status": "pending",
    "amount": 10000,
    "currency": "THB",
    "description": null,
    "capture": true,
    "authorized": false,
    "reversed": false,
    "paid": false,
    "transaction": "",
    "card": null,
    "refunded": 0,
    "refunds": {
        "object": "list",
        "id": "",
        "livemode": false,
        "location": "/charges/chrg_test_5pw21ew426s6wzuyz3m/refunds",
        "created": "0001-01-01T00:00:00Z",
        "from": "1970-01-01T00:00:00Z",
        "to": "2021-11-19T06:18:29Z",
        "offset": 0,
        "limit": 20,
        "total": 0,
        "order": "chronological",
        "data": []
    },
    "failure_code": null,
    "failure_message": null,
    "customer": "",
    "ip": null,
    "dispute": null,
    "return_uri": "http://localhost:8080/omise-scb/1637302708655035000/complete",
    "authorize_uri": "https://pay.omise.co/offsites/ofsp_test_5pw21ew5xnoq11ch3ts/pay",
    "source_of_fund": "",
    "offsite": "",
    "source": {
        "object": "source",
        "id": "src_test_5pw21euspas96rvezzy",
        "livemode": false,
        "location": "/sources/src_test_5pw21euspas96rvezzy",
        "type": "internet_banking_scb",
        "flow": "redirect",
        "amount": 10000,
        "currency": "THB"
    },
    "metadata": {}
}
````

#### Get transaction status request
````
curl --location --request GET 'localhost:8080/omise-scb?orderID=1637295744829628000&status=successful'
````

#### Get transaction status response
```
{
    "status": "success",
    "Data": [
        {
            "ID": 1,
            "order_id": 1637295744829628000,
            "charge_id": "chrg_test_5pw21ew426s6wzuyz3m",
            "amount": 10000,
            "currency": "THB",
            "payment_status": "successful",
            "soure_id": "src_test_5pw21euspas96rvezzy",
            "soure_type": "internet_banking_scb",
            "paid_at": "2021-11-19T13:18:43.171114+07:00",
            "create_at": "0001-01-01T00:00:00Z"
        }
    ]
}
```
