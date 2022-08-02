# Motivation

This is an example for GORM with postgresql inspired by:

* https://github.com/harranali/gorm-relationships-examples/tree/main/has-many
* https://stackoverflow.com/questions/73169934/using-relational-queries-with-gorm-go


 
# License

The code is public domain (educational purpose). Use it for whatever.

# How to get the database running

    docker compose build
    docker compose up -d

The postgres database is initialized automatically. No manual work needed, see also:
* https://herewecode.io/blog/create-a-postgresql-database-using-docker-compose/ so 

# How to get the example running

    docker compose build && docker compose run go-example

Then you can hack on the main.go and do this over and over.

    