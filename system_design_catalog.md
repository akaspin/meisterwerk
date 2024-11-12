# Catalog Service

The Catalog service operates products and services. This service accepts requests from "Orders" and "Quotes" services.
The API endpoints allows data management from Client or CLI.

> It's recommended to implement changes broadcast to avoid redundant calls from "Orders" and "Quotes" services.

Both product and service are encapsulated in Item. The Item object has required properties:

- ID
- Name
- Description
- Segment
- Price
- Tax Rate

The Segment value can be "product" or "service".

> Note that depends on service the Item behaviour can vary. This can be implemented with `oneOf` mechanism in 
> OpenAPI 3.1 and polymorphic constraints in relational DB. This functionality not implemented to reduce complexity.

> In case of significant behaviour discrepancy it's recommended to split Catalog service to separate services. 

The products and features should be named uniquely. This allows semi-natural indexing. In basic cases the ID can be
calculated from name in kebab-case. Like "Time planner" to "time-planner".

This leads to the following data structure:

`items`

- `id`, text
- `segment`, enum (product, service)
- `name`, text, unique
- `description`, text
- `price`, money, not null
- `tax`, double, not null

Unique: `id`, `segment`

## Technologies

- AWS RDS
- Containers can be deployed in any environment like AWS ECS or Kubernetes
- Supplementary AWS services like IAM for Service Account

## API

### `POST /items`

Create new item. Returns the ID of created product calculated from name 
as kebab-case. Fails in case if item already exists.

```http request
### Expect 201 { "id": "basic" }
POST /items

{
  "segment": "product",
  "name": "Basic",
  "description": "Basic package"
  "price": 100,
  "tax": 0.5
}
```

### `GET /items`

Retrieve items. The features can be filtered by Item IDs (`?id=basic&id=advanced`) and/or segments 
(`?segment=service`). `skip` and `limit` parameters can be used for pagination.

By default all items are sorted by segment and their IDs in ascending lexicographic order. Use `sort` parameter
to change behaviour. For example `?sort=name` will sort the output by name. Use `order` parameter to change the order.
The valid `order` values are `asc` (default) and `desc`.

```http request
### Expect 200
GET /items?id=basic&segment=product
```
```json
[
  {
      "id": "basic",
      "segment": "product",
      "name": "Basic",
      "description": "Basic package"
      "price": 10,
      "tax": 0.3
  }
]
```

### `GET /segments/{segment}/items/{id}`

Retrieve specific item.

```http request
### Expect 200
GET /segments/product/item/basic

{
  "id": "basic",
  "segment": "product",
  "name": "Basic",
  "description": "Basic package",
  "price": 200,
  "tax": 0.1
}
```

### `PUT /segments/{segment}/items/{id}`

Update existing item.

```http request
### Expect 200
PUT /segments/product/item/basic

{
  "description": "Basic package",
  "price": 200,
  "tax": 0.1
}
```

### `DELETE /segments/{segment}/items/{id}`

Delete existing item.

```http request
### Expect 200
DELETE /segments/product/item/basic
```

## Bulk API

Use `/bulk/...` API for mass operations. All bulk operations imply transactional integrity. This means that 
if an error occurs, no changes will be made. For example no items will be created if one of them already exists.

### `POST /bulk/items`

Create multiple items. Returns pairs of segment and ID of created items.

```http request
### Expect 201 [ { "segment": "product", "id": "basic" }, { "segment": "service", "id": "basic" }]
POST /bulk/items

[
  {
    "name": "Basic",
    "segment": "product",
    "description": "Basic package"
    "price": 10,
    "tax": 0.3
  },
  {
    "name": "Basic",
    "segment": "service",
    "description": "Basic package"
    "price": 10,
    "tax": 0.3
  },
]
```

### `PUT /bulk/items`

Update multiple items.

```http request
### Expect 200
PUT /bulk/items

[
  {
    "id": "basic",
    "name": "Basic",
    "segment": "product",
    "description": "Basic package"
    "price": 10,
    "tax": 0.3
  },
  {
    "id": "basic",
    "name": "Basic",
    "segment": "service",
    "description": "Basic package"
    "price": 10,
    "tax": 0.3
  },
]
```

### `DELETE /bulk/items`

Delete multiple items.

```http request
### Expect 200
DELETE /bulk/items

[ 
  { 
    "segment": "product",
    "id": "basic" 
  },
  { 
    "segment": "service",
    "id": "basic"
  }
]
```
