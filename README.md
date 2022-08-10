# freshrss-tenant-manager

## Usage

### Create Tenant

```
$ cat examples/tenant.json 
{
    "Name": "Alice"
}

$ curl -i -X POST -H "Content-Type: application/json" -d "@examples/tenant.json" http://localhost:8080/v1alpha1/tenants/
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Location: http://localhost:8080/v1alpha1/tenants/a174905e-95f5-48b1-a35f-1c61077acb92
Date: Wed, 10 Aug 2022 18:15:47 GMT
Content-Length: 335

{
    "ID": "a174905e-95f5-48b1-a35f-1c61077acb92",
    "Name": "Alice",
    "Subscriptions": null,
    "links": {
        "self": "http://localhost:8080/v1alpha1/tenants/a174905e-95f5-48b1-a35f-1c61077acb92",
        "subscriptions": "http://localhost:8080/v1alpha1/tenants/a174905e-95f5-48b1-a35f-1c61077acb92/subscriptions/"
    }
}
```

### Create Subscription

POST to the URL found in the tenant's `links` sub-document.

```
$ cat examples/subscription.json 
{
    "Title": "My Fav Feeds",
    "Username": "alice",
    "Service": "enabled"
}

$ curl -i -X POST -H "Content-Type: application/json" -d "@examples/subscription.json" http://localhost:8080/v1alpha1/tenants/a174905e-95f5-48b1-a35f-1c61077acb92/subscriptions/
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Location: http://localhost:8080/v1alpha1/tenants/a174905e-95f5-48b1-a35f-1c61077acb92/subscriptions/b7331e7c-99fc-41de-87df-3a7d5d9020c5
Date: Wed, 10 Aug 2022 18:20:29 GMT
Content-Length: 467

{
    "ID": "b7331e7c-99fc-41de-87df-3a7d5d9020c5",
    "TenantID": "a174905e-95f5-48b1-a35f-1c61077acb92",
    "Service": "enabled",
    "Title": "My Fav Feeds",
    "Username": "alice",
    "URL": "",
    "links": {
        "self": "http://localhost:8080/v1alpha1/tenants/a174905e-95f5-48b1-a35f-1c61077acb92/subscriptions/b7331e7c-99fc-41de-87df-3a7d5d9020c5",
        "tenant": "http://localhost:8080/v1alpha1/tenants/a174905e-95f5-48b1-a35f-1c61077acb92"
    }
}
```

### Get a Tenant with Subscriptions

```
$ curl -i http://localhost:8080/v1alpha1/tenants/a174905e-95f5-48b1-a35f-1c61077acb92
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 10 Aug 2022 18:22:48 GMT
Content-Length: 902

{
    "ID": "a174905e-95f5-48b1-a35f-1c61077acb92",
    "Name": "Alice",
    "Subscriptions": [
        {
            "ID": "b7331e7c-99fc-41de-87df-3a7d5d9020c5",
            "TenantID": "a174905e-95f5-48b1-a35f-1c61077acb92",
            "Service": "enabled",
            "Title": "My Fav Feeds",
            "Username": "alice",
            "URL": "",
            "links": {
                "self": "http://localhost:8080/v1alpha1/tenants/a174905e-95f5-48b1-a35f-1c61077acb92/subscriptions/b7331e7c-99fc-41de-87df-3a7d5d9020c5",
                "tenant": "http://localhost:8080/v1alpha1/tenants/a174905e-95f5-48b1-a35f-1c61077acb92"
            }
        }
    ],
    "links": {
        "self": "http://localhost:8080/v1alpha1/tenants/a174905e-95f5-48b1-a35f-1c61077acb92",
        "subscriptions": "http://localhost:8080/v1alpha1/tenants/a174905e-95f5-48b1-a35f-1c61077acb92/subscriptions/"
    }
}
```