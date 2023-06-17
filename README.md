# Motivation

This is an example for GORM with postgresql inspired by:

* https://github.com/harranali/gorm-relationships-examples/tree/main/has-many
* https://stackoverflow.com/questions/73169934/using-relational-queries-with-gorm-go

This is a simple cross-platform usable example on GORM and has-many relations.

# License

The code is public domain (educational purpose). Use it for whatever.

# How to get the database running

    docker compose build
    docker compose up -d
    docker compose run go-example

The postgres database is initialized automatically. No manual work needed, see also:
* https://herewecode.io/blog/create-a-postgresql-database-using-docker-compose/ so 

The schema is created by gorm, see this lines in main.go:

	db.Migrator().DropTable(&CreditCardUser{}, &CreditCard{})
	db.AutoMigrate(&CreditCardUser{})
	db.AutoMigrate(&CreditCard{})

So your examples will always have a correct schema.

## Use psql to access the database

    docker compose up -d
    docker exec -it gorm-examples-postgresql-1 sh -c 'psql -U postgres -d gorm-example'
    \dt
    \l
    gorm-example=# \dt
               List of relations
    Schema |       Name        | Type  |  Owner
    --------+-------------------+-------+----------
    public | credit_card_users | table | postgres
    public | credit_cards      | table | postgres
    (2 rows)
    
    gorm-example=# \d+
    List of relations
    Schema |           Name           |   Type   |  Owner   | Persistence |    Size    | Description
    --------+--------------------------+----------+----------+-------------+------------+-------------
    public | credit_card_users        | table    | postgres | permanent   | 16 kB      |
    public | credit_card_users_id_seq | sequence | postgres | permanent   | 8192 bytes |
    public | credit_cards             | table    | postgres | permanent   | 16 kB      |
    public | credit_cards_id_seq      | sequence | postgres | permanent   | 8192 bytes |
    (4 rows)
    
    gorm-example=# select * from credit_cards;
    id |          created_at           |          updated_at           | deleted_at |    number    |       bank       | user_id
    ----+-------------------------------+-------------------------------+------------+--------------+------------------+---------
    1 | 2022-08-02 09:01:02.187082+00 | 2022-08-02 09:01:02.187082+00 |            | 1234567898   | FinFisher        |       1
    2 | 2022-08-02 09:01:02.187082+00 | 2022-08-02 09:01:02.187082+00 |            | 345657881    | MaxedOut Limited |       1
    4 | 2022-08-02 09:01:02.188971+00 | 2022-08-02 09:01:02.188971+00 |            | 2342         | Bankxter         |       2
    7 | 2022-08-02 09:01:02.199361+00 | 2022-08-02 09:01:02.199361+00 |            | 666666666666 | happyBank        |       4
    8 | 2022-08-02 09:01:02.201586+00 | 2022-08-02 09:01:02.201586+00 |            | 666666666666 | happyhappyBank   |       4
    (5 rows)
    
    gorm-example=# select * from credit_card_users;
    id |          created_at           |          updated_at           | deleted_at |      name
    ----+-------------------------------+-------------------------------+------------+----------------
    1 | 2022-08-02 09:01:02.186323+00 | 2022-08-02 09:01:02.186323+00 |            | mrFlux
    2 | 2022-08-02 09:01:02.188735+00 | 2022-08-02 09:01:02.188735+00 |            | sirTuxedo
    3 | 2022-08-02 09:01:02.190299+00 | 2022-08-02 09:01:02.190299+00 |            | missFraudinger
    4 | 2022-08-02 09:01:02.191714+00 | 2022-08-02 09:01:02.201393+00 |            | happyUser
    (4 rows)

## Use pgadmin to access the database

You need to enable it in the docker-compose.yaml and then:

    docker compose up -d

pgadmin might be nice to build bigger DB requests. In general, I prefer psql for its simplicity as pgadmin's session
frequently break with the docker compose restarts from the goland IDE.

## Flush the contents of the postgresql state

If you want to recreate the database, basically delete the volume and then the database gets created automatically:

If gorm-example is still running:

    docker compose down

Let's see the volume, should be something like this:

    docker volumes ls
    local     gorm-examples_postgresql

Then kill it with fire:

    docker volume rm --force gorm-examples_postgresql
    gorm-examples_postgresql

Afterwards basically start working with it as usual, i.e.:

    docker compose up -d

Which should create a completely new database. Introspect the logs with:

    docker container logs  -f -n 200  gorm-examples-postgresql-1

Note: Normally this should not be needed.

# How to get the example running

    docker compose build && docker compose run go-example

Then you can hack on the main.go and do this over and over.

# sqlite3

The sqlite3 folder contains an example for an articles/tags example with lots of calls to the database which
might be interesting. It is basically a test but it is also a good example for how to use sqlite3 with gorm.