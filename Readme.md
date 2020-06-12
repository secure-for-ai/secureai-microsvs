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
