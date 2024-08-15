# go-matchmaker

Prototype for "matchmaker" service to combine people in groups according to their skills etc.
The main "feature" of the implementation is that real matching code is plugged in as a
file written in Lua, this gives several advantages:

- matcher code could be developed by persons (data-scientists) not familiar with Go
- matcher code could be super-easy tested separately from the service ("driver" test file is included)
- matcher code could be replaced without rebuilding the project (we even can add ability to
change it in run-time, for example to better fit time of the day, load etc)

### Build and Run

    make build

    GROUP_SIZE=3 TCP_PORT=8082 MATCHER_FILE=euclid_matcher.lua build/server

    tests/simple10.sh

Configuration in env variables:

- `TCP_PORT` - port to listen to
- `GROUP_SIZE` - users per group
- `MATCHER_FILE` - matcher code file in Lua (by default built-in naive implementation is used)
- `USER_BUFFER` - how many users could be accepted before HTTP interface will start blocking
- `MATCHER_PERIOD` - how often matcher code is called (seconds)
- `DEBUG_MATCHER` - extra info on accepting users and grouping them (in stdout)

### Extra info

See [DETAILS_RU.md](./DETAILS_RU.md) for implementation details.

_Grfg Cebwrpg sbhaq ba UU sbe Yrfgn Tnzrf_
