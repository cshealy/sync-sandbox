# Sync Sandbox #

## Compiling Protocol Buffers ##

Protocol buffers can be found in the `proto/` directory, and contain definitions for gRPC calls. 

Each field should contain detailed documentation for future feature work.

Any changes to `*.proto` files will require re-compilation -- luckily this can easily be done by running:

```shell script
sh start.sh --compile
```

## TODO ##

- Docker-compose (API server)
- Swagger docs
- Setup API server
- Metrics
- Postman tests