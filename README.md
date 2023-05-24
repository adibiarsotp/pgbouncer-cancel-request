Cancel requests outage pgbouncer
================================

forked from https://github.com/chobostar/pgbouncer-cancel-request, with modification put webserver and context timeout passed from client 

Usage:

```
$ go get github.com/adibiarsotp/pgbouncer-cancel-request
$ cd ~/go/src/github.com/adibiarsotp/pgbouncer-cancel-request
$ docker-compose up -d
$ make run_pgbouncer_tester
```

Check `used_clients`:

```
$ psql -h localhost -p 6432 -U pgbouncer -c "show lists" | grep 'used_clients'
```

After several minutes (example 5 minutes), we can stop the script, and recheck used clients. the number will not reset,
and we will unable to enter the db.

```
$ psql -h localhost -p 6432 -U postgres -d db
psql: ERROR:  no more connections allowed (max_client_conn)
```

We can change the pgbouncer version in docker-compose.yml to test fixed version pgbouncer >= 1.16. dont forget to rebuild the docker image.