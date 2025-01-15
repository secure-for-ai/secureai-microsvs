# Micro Service Framework
The framework is built for container based applications focusing on **scalability**
and **simplicity**. It features both
relational and non relational database as well as distributed cache servers.
In addition, the framework integrates all necessary security features to support
backend development. Using this framework, one can develop micro services
and quickly deploy them on cloud.


### Features
1. Cache: redis,
2. DB: mongodb, postgres/yugabyte
   * sql builder: using api calls to build sql automatically  
3. Web Crypto: crfs, jwt, and aes encryption
4. session: cookie based session storage. use redis as the cache server,
   mongodb or postgres as the persistent storage.
5. snowflake: implementation of twitter distributed unique id generator.

### Performance
Single instance > 500 TPS on a 2 core 4 GB AWS instance which
runs nginx+backend+mongo all in one.

### How to run
- Start docker compose
    ```bash
    docker-compose up
    ```

- Initial Mongo
    ```bash
    # login mongo
    mongo --port 27017 --host=localhost --authenticationDatabase=admin \
        -p password --username test
    ```

    ```bash
    # setup replication with single master
    rs.initiate({_id: "rs0", members: [{_id: 0, host: "localhost:27017"}] })
    ```


### PG data upgrade

```
# $POSTGRES_USER is the old data role
docker run --rm \
        -e PGUSER=$POSTGRES_USER \
        -e POSTGRES_INITDB_ARGS="-U $POSTGRES_USER" \
        -v [local old dir]:/var/lib/postgresql/$OLD/data \
        -v [local new dir]:/var/lib/postgresql/$NEW/data \
        "tianon/postgres-upgrade:$OLD-to-$NEW"
```

Suppose the `pg 13` data is in `_data/postgres_13` and `pg 17` data will be
in `_data/postgres_17`. The old data role is `test`. The command would be

```
OLD="13"
NEW="17"

OLD_DATA="./_data/postgres_13"
NEW_DATA="./_data/postgres_17"

POSTGRES_USER="test"
POSTGRES_PASSWORD="password"

docker run --rm \
        -e PGUSER=$POSTGRES_USER \
        -e PGPASSWORD=$POSTGRES_PASSWORD \
        -e POSTGRES_INITDB_ARGS="-U $POSTGRES_USER" \
        -v $OLD_DATA:/var/lib/postgresql/$OLD/data \
        -v $NEW_DATA:/var/lib/postgresql/$NEW/data \
        "tianon/postgres-upgrade:$OLD-to-$NEW"

# specify the auth-method for `host` connections for `all` databases,
# `all` users, and `all` addresses. If unspecified then `scram-sha-256` password
# authentication⁠ is used 
printf '\n' | sudo tee -a $NEW_DATA/pg_hba.conf > /dev/null
printf 'host all all all scram-sha-256\n' | sudo tee -a $NEW_DATA/pg_hba.conf > /dev/null
```

Change the password

```
docker exec -it secureai-dev-postgres bash
psql -U test
ALTER USER test WITH PASSWORD 'password';
\q
```