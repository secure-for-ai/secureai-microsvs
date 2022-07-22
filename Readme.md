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
