## GarageSale

This is my practice project, while covering the nice Ardanlabs' Ultimate Service course.

<br/>

### Tests

The following tests are included:
- Business Logic tests, triggered using `go test -v ./internal/product`
- API Tests, triggered using `go test -v ./cmd/sales-api/tests`

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
