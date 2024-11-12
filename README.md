# Meisterwerk

## Assumptions

I have no information about the business logic and overall existing architecture. 
Therefore, all proposed solutions are based on simplifications. 

- All services are deployed as containers in AWS ECS. 
- All container images are available in AWS ECR.
- Public API mapped in AWS API Gateway with Authorizer Lambda.
- It is assumed that the Authentication and Authorisation mechanisms have already been implemented 
  and can be integrated with new services.
- Any domain model data synchronization broadcast mechanisms are not taken into account. 
  This does not mean that streams (SNS, SQS, Kafka) are not used at all.

To avoid misunderstandings, notes have been added to the text and code. They may contain assumptions 
or descriptions of alternative solutions.

## Document structure

This document is divided to two parts: [System Design](#system-design) 
and [Service Implementation](#service-implementation). Each section contains own assumptions, notes and description. 

To reduce complexity several parts are extracted to separate documents.

## System Design

Two new services have been added to the [diagram](https://drive.google.com/file/d/1uX_FXl-aJi59jzw-SmgwyXXlbgy9KQud/view?usp=sharing).

## Assumptions and Notes

- For proper UI function "Orders", "Quotes" and "Catalog" services should be accessible from Client. As well as
  from "Integration" service.
- The "Catalog" service used by "Orders" and "Quotes" services for calculations and validations.
- Relational databases used for storage.
- All obvious information like extra fields or content-type in request examples is not documented to reduce complexity.

Refer to [Catalog Service](system_design_catalog.md) and [Quotes Service](system_design_quotes.md) for details. 

## Service Implementation

The Quotes service is selected for implementation. For behaviour description refer to [service documentation](system_design_quotes.md)

### Tools and Libraries

> GNU Make can be used but is not required. In command-line snippets `make` targets added as alternatives.  

- [OpenAPI Generator](https://openapi-generator.tech): OpenAPI 3.1 code generation.
- [Docker](https://www.docker.com) (and Docker Compose): Integration tests.
- [PostgreSQL](https://www.postgresql.org): Relational DB.
- [ArigaIO Atlas](https://github.com/ariga/atlas): Database migrations.
- [pflag](https://github.com/spf13/pflag): Command-line flags.
- [GORM](https://gorm.io): Go ORM to reduce complexity.
- [testify](github.com/stretchr/testify): Tests assertions.

### Required Tools

- Go
- [Docker](https://www.docker.com) (and Docker Compose)

### Implementation Notes

- Mock Catalog Client.
- For testing reasons when quote is accepted Quotes service sends requests to mock Orders service.
- All the optimisations removed to reduce complexity.
- Despite the ArigaIO Atlas offers automatic migrations generation. All related code is removed to reduce complexity.

### Basic usage

Deploy test stack. The Quotes server will be accessible on `8080` port. For API refer 
to [OpenAPI manifest](api/openapi/quotes.yaml) or [Design Document](system_design_quotes.md).

```shell
# make up
docker compose up -d --remove-orphans --wait

curl http://localhost:8080/quotes
```

In case of GNU Make use `UP_ARGS` environment variable to add additional parameters.

```shell
UP_ARGS=--build make up
```

Tests are run in the usual way.

```shell
# make test
go test ./...
```

All residue can be removed by following command.

```shell
# make clean
docker compose down -v --remove-orphans
```


