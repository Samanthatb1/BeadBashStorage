# Bead Bash Storage
Database to store a history of users and oders for the jewelry business Bead Beash Studio

# Docker Usage
## 1) Clone for personal development with Docker
1. Have [Docker](https://www.docker.com/) installed 
2. clone this repository with `git clone git@github.com:Samanthatb1/beadBashStorage.git`
3. cd into the root directory with `cd beadBashOrders` and create an env file titled `app.env`
4. Add to the env file: 
```
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:secret@localhost:5432/BB-DB?sslmode=disable
SERVER_ADDRESS=0.0.0.0:8080
MIGRATION_URL=file://db/migration
```
5. Run with **either** `docker compose up` **or** run this series of commands:
```
$ make startPostgresContainer
$ make createDB
$ make migrate_up
$ make server
```
--
## 2) Run without cloning code with Docker Hub
1. Grab the docker image from [Docker Hub](https://hub.docker.com/repository/docker/samanthatb1/bead-bash-orders/general) by running :
```
$ docker pull samanthatb1/bead-bash-orders:latest
```
2. Run database and api service:

```
$ docker run --name postgres --network bb-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=BB-DB -d postgres

$ docker run --name api --network bb-network -e DB_SOURCE=postgresql://root:secret@postgres:5432/BB-DB?sslmode=disable -e SERVER_ADDRESS=0.0.0.0:8080 -e DB_DRIVER=postgres -e ENVIRONMENT=production -p 8080:8080 -e GIN_MODE=release  -d samanthatb1/bead-bash-orders
```

# API Endpoints
## Data Layout
User:

      {
          "id": number
          "username": "your username",
          "full_name": "your name",
          "total_orders": number,
          "created_at": date
      }

Order:

      {
          "order_id": number,
          "account_id": number,
          "username": "your username",
          "full_name": "your name",
          "purchased_item": "purchased item",
          "purchase_amount": "item cost",
          "shipping_location": "shipping location",
          "currency": "currency code",
          "date_ordered": "date ordered",
      }
## Endpoints

Get an existing User

    GET /users/:identifier
    -> where identifier can be a user id or username
    -> returns the user


Get all Users

    GET /users/all?page_id={number}&page_size={number}
    -> returns an array of users based on the page and amount requested
Add a new User

    POST /users
    -> returns the new created User

    Body Params:
      {
          "full_name": "your name",
          "username": "your username",
      }

Delete User

    DELETE /users/:username
    -> returns deletion status

Get A Users Orders

    GET /orders/:username
    -> returns an array of orders from that user

Get all Orders

    GET /orders/all?page_id={number}&page_size={number}
    -> returns an array of orders based on the page and amount requested

Create new Order


    POST /orders
    -> returns the user and the edited order

    Body Params:
      {
          "username": "your username",
          "full_name": "your name",
          "purchased_item": "purchased item",
          "purchase_amount": "item cost",
          "shipping_location": "shipping location",
          "currency": "currency code",
          "date_ordered": "date ordered",
      }

Delete Order


    DELETE /orders/:order_id
    -> returns deletion status and the item that was deleted

Edit Order

    PATCH /orders
    -> returns the updated order

    Body Params:
      {
          "order_id": "order id",
          "purchase_amount": "updated amount ", OPTIONAL
          "purchased_item": "updated item", OPTIONAL
          "shipping_location": "updated shipping location", OPTIONAL
      }

## DB Schema
  ![Database Image](./images/DB_Tables.png?raw=true)