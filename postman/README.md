# Postman Integration Tests #

## Importing ##

Postman integration tests can easily be imported through the UI by performing the following commands:

```
File > Import... > Upload Files > sync-sandbox-integration-tests.json
```

## Testing ##

Tests can be run manually when sending a request through postman and will eventually be automated whenever a new push is sent to github.

## Exporting New Tests ##

Tests are exported when a collection is exported. This can be be achieved by clicking:

```
sync-sandbox collection > ... > Export > <postman directory>
```

* Note: make sure this is committed to the sync-sandbox repo *