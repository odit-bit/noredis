WIP project

folder structure

(root)
    logic of noredis server to manage connection, parsing, and etc...
/app 
    entry point of application
/client (not implemented)
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

to use the app , it can build from app/server directory.