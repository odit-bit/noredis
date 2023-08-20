WIP project

folder structure

(root)
    logic of noredis server to manage connection, parsing, and etc...

/app 
    entry point of application

/client 
    client API

/db
    implementation of redis data store, with go map as datastructure

/resp 
    implementation of RESP protocol.

motivation of this project is to develop understanding of the redis app design internals.
No depedency beside redis client package just for testing purpose.It implemented RESP protocol, 
so it compatible with redis client library for some limited feature or use client that available from this project

for now feature:
SET - (no options)
GET - (no options)
INCR

to try server locally , it can run from app/server directory.
`go run ./app/server`
or with flag
`go  run ./app/server --port {port} --password {password}`


and

try run example cli app, 
`go run ./example/cli` 
or with flag
`go  run ./example/cli --addr {addr} --password {password}`

