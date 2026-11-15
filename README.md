# Fence

Cedar based auth system

## run modes

1. local

- all policies and entities are managed through local client
- policies/entities are on local files

2. cached local

- policies/entities are managed via local files
- state manager has a function that can be run as a go routine to regularly update entities/policies from filesystem

2. proxied mode

- all policies and entities are manged on server, local client calls server for authz

3. cached mode

- local client or agent caches info from server and regularly updates client

## Ideas

- crdt with memberlist or serf

## TODO

- switch entity management to just be for roles and users
- IsAllowed takes in an additional set of entities that can be merged with the users and roles stored in the provider

## Resources

- https://docs.cedarpolicy.com/bestpractices/bp-separate-principals.html
- https://docs.cedarpolicy.com/bestpractices/bp-implementing-roles-templates.html
- https://cedarpolicy.com/en/playground
