## GarageSale

This is my practice project, while covering the nice Ardanlabs' Ultimate Service course.

<br/>

### Setup

Run `./run-admin.sh migrate` to run the database migration, meaning populating the application's database objects.

Optionally, run `./run-admin.sh seed` to feed in some initial/testing data to play with.

Run `./run-admin keygen private.pem` to generate the `private.pem` file that will store the private key used for signing the JWT tokens returned as a result of a successful user authentication (see `/v1/users/token` operation for details).

<br/>

### Tests

The following tests are included:
- Business Logic tests, triggered using `go test -v ./internal/product`
- API Tests, triggered using `go test -v ./cmd/sales-api/tests`

<br/>

### Admin

The administrative features are accessible using `./run-admin.sh` script.<br/>
Besides the aforementioned (in Setup section above) database migration and seed capabilities, plus generation of the private key store, users can also be added using `useradd` command. Example:

```shell
$ ./run-admin.sh useradd joe@mail.com joe
Admin user will be created with email "joe@mail.com" and password "joe"
Continue? (1/0) 1
User created with id: c054c42e-bd6a-4236-a7d5-9395d76b4eff
$ 
```

<br/>

### Runtime Insights

#### Profiling

As part of _Getting Production Ready_, this service includes a Debug server, listening on a different port as this is not meant 
for public exposure.<br/>
For example, you can do CPU profiling by:
1. Triggering some stress:<br/>`hey -c 10 -n 15000 http://localhost:8000/v1/products`
2. Collect CPU profiling for some time:<br/>`go tool pprof "http://localhost:6060/debug/pprof/profile?seconds=8"`

#### Metrics

The internal [http://localhost:6060/debug/vars](http://localhost:6060/debug/vars) endpoint gives memstats insights, plus application specific metrics (see `internal/middleware/metrics.go`) such as number of errors, go routines, and requests.

Also, [expvarmon](https://github.com/divan/expvarmon) tool used like this: `expvarmon -ports=":6060" -endpoint="/debug/vars" -vars="requests,goroutines,errors,mem:memstats.Alloc"` gives a nice auto-refreshed TUI (text based ui).
