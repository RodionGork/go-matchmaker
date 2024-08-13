# go-matchmaker

Prototype for "matchmaker" service to combine people in groups according to their skills etc. Tech details
are to be added later.

### Build and Run

    make build

    GROUP_SIZE=3 TCP_PORT=8085 build/server

    curl -d '{"name":"Kitty","skill":14,"latency":2}' http://localhost:8085/users
    curl -d '{"name":"Patty","skill":13,"latency":3}' http://localhost:8085/users
    curl -d '{"name":"Clown","skill":12,"latency":4}' http://localhost:8085/users

Initial implementation just combines group when there are enough users in the queue.

### Extra info

_Grfg Cebwrpg sbhaq ba UU sbe Yrfgn Tnzrf_
