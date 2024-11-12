# Quotes Service

Quotes service implements following functionality:

- Accepts quotes requests.
- Updates quotes during negotiation.
- Informs Orders service when quote is accepted.
- Interacts with Catalog service for information about Items.

It is assumed that:

- The Customer is already registered with integer ID.
- The Orders service implements mechanism to create the order based on accepted quote.

Each quote consists of Items. Quotes service uses [Catalog service](system_design_catalog.md) to retrieve requested items.

On quote creation all prices and tax rates are frozen. And should not change on changes in Catalog service. Prices
can be changed only on quote update. This also means that any item removed in Catalog service should not affect
quote and can be changed only during negotiation.

All quotes descriptions for one customer should be unique.

This leads to following data structures:

`quotes`

- `id`, int, sequence, PK
- `customer_id`, int, not null
- `description`, text (unique for customer), not null
- `status`, text (pending, accepted, rejected), not null

Unique: `customer_id`, `description`

`items`

- `quote_id`, int, not null
- `item_id`, text, not null
- `segment`, enum (product, service), not null
- `price`, money, not null
- `tax`, money, not null

PK: `quote_id`, `item_id`, `segment`

## Technologies

- AWS RDS
- Containers can be deployed in any environment like AWS ECS or Kubernetes
- Supplementary AWS services like IAM for Service Account
- (optional) AWS SNS/SQS or Kafka

## API

### `POST /quotes`

Create new quote. Returns the ID of created quote. The price and/or tax rate for items are
not required. In this case they will be retrieved from Catalog service. The request fails if requested
items are not exists in Catalog.

```http request
### Expect 201 { "id": 10 }
POST /quotes

{
  "customer_id": 100,
  "description": "Quote for customer",
  "status": "pending",
  "items": [
    {
      "id": "basic",
      "segment": "product",
      "price": 100,
      "tax": 0.1
    },
    {
      "id": "basic",
      "segment": "service"
    }
  ]
}
```

### `GET /quotes`

Retrieve quotes. The features can be filtered by Item IDs (`?id=10&id=100`) and/or customer
IDs (`?customer_id=20`). `skip` and `limit` parameters can be used for pagination.

All quotes are sorted by their id in ascending order. Use `order` parameter to change the order.
The valid `order` values are `asc` (default) and `desc`.

```http request
### Expect 200
GET /quotes?id=100
```
```json
[
  {
    "id": 1000,
    "customer_id": 100,
    "description": "Quote for customer",
    "status": "pending",
    "items": [
      {
        "id": "basic",
        "segment": "product",
        "price": 100,
        "tax": 0.1
      },
      {
        "id": "basic",
        "segment": "service",
        "price": 100,
        "tax": 0.1
      }
    ]
  }
]
```

### `GET /quotes/{id}`

Retrieve specific quote.

```http request
### Expect 200
GET /quotes/1000

{
  "id": 1000,
  "customer_id": 100,
  "description": "Quote for customer",
  "status": "pending",
  "items": [
    {
      "id": "basic",
      "segment": "product",
      "price": 100,
      "tax": 0.1
    },
    {
      "id": "basic",
      "segment": "service",
      "price": 100,
      "tax": 0.1
    }
  ]
}
```

### `PUT /quotes/{id}`

Update existing quote. In case if quote is "accepted" the Quotes service will send message to Orders.

```http request
### Expect 200
PUT /quotes/1000

{
  "id": 1000,
  "customer_id": 100,
  "description": "Updated quote for customer",
  "status": "accepted",
  "items": [
    {
      "id": "basic",
      "segment": "product",
      "price": 200,
      "tax": 0.2
    },
    {
      "id": "basic",
      "segment": "service",
      "price": 100,
      "tax": 0.1
    }
  ]
}
```

### `DELETE /quotes/{id}`

Delete quote.

```http request
### Expect 200
DELETE /quotes/1000
```

## Bulk API

Use `/bulk/...` API for mass operations. All bulk operations imply transactional integrity. This means that
if an error occurs, no changes will be made. For example no quotes will be updated if one of them is not exists.

### `POST /bulk/items`

Create multiple quotes. Returns IDs of created quotes.

```http request
### Expect 201 { "ids": [ 10, 11 ] } 
POST /bulk/quotes

[
  {
    "customer_id": 100,
    "description": "First quote",
    "status": "pending",
    "items": [
      {
        "id": "basic",
        "segment": "product"
      }
    ]
  },
  {
    "customer_id": 100,
    "description": "Second quote",
    "status": "pending",
    "items": [
      {
        "id": "advanced",
        "segment": "product"
      }
    ]
  }
]
```

### `PUT /bulk/quotes`
Update multiple items. For all accepted quotes the Quote service will inform Orders service.

```http request
### Expect 200
PUT /bulk/quotes

[
  {
    "id": 10,
    "customer_id": 100,
    "description": "First updated quote",
    "status": "pending",
    "items": [
      {
        "id": "basic",
        "segment": "product",
        "price": 200,
        "tax": 0.2
      }
    ]
  },
  {
    "id": 11,
    "customer_id": 100,
    "description": "Second updated quote",
    "status": "accepted",
    "items": [
      {
        "id": "basic",
        "segment": "product",
        "price": 200,
        "tax": 0.2
      }
    ]
  }
]
```

### `DELETE /bulk/quotes`

Delete multiple quotes.

```http request
### Expect 200
DELETE /bulk/quotes

{
  "ids": [
    10,
    11
  ]
}
```
