# Sync Sandbox #

## Quickstart ##

**Work in progress** - currently I'm in the process of working on the server, but have added the following command to spin up our docker container:

```shell script
sh start.sh --build
```

*Note:* currently you would be required to build the image locally, in the future I plan to add an ECR repo with terraform that will automatically push images in a Jenkins pipeline fashion.

## gRPC Gateway ##

The [gRPC gateway](https://github.com/grpc-ecosystem/grpc-gateway) translates a RESTful HTTP API into gRPC. Sync-sanbox accepts both gRPC and REST calls, which means you can easily perform the following to receive a response:

```shell script
# start the API
sh start.sh

# curl our test endpoint
curl http://localhost:8080/test?name=echo

# response
{"name":"echo"}
```

This is just a placeholder for now, but will hold interesting info in the future.

## Compiling Protocol Buffers ##

Protocol buffers can be found in the `proto/` directory, and contain definitions for gRPC calls. 

Each field should contain detailed documentation for future feature work.

Any changes to `*.proto` files will require re-compilation -- luckily this can easily be done by running:

```shell script
sh start.sh --compile
```

## Mocks & Testing ##

Inevitably mocks will be required when testing. To add a new mock, ensure your interface is structured in a way that will allow unit tests to run successfully.
 
This requires your interface to contain any functions that would contact an external resource, such as an API.

Once completed, update `start.sh` with your mock source, destination, and package, then execute the following command:

```shell script
sh start.sh --mocks
```

Your new mock files should be found in the destination, and can now be used in unit tests.
## TODO ##

- Docker-compose (API server)
- Swagger docs
- Setup API server
- Metrics
- Postman tests