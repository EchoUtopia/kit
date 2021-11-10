# opentelemetry golang trace sdk


## example


firstly, init the tracer for once:

```
if err := trace.InitTracerWithJaegerExporter(serviceName); err != nil {
    handle(err)
}
```

in this init function, we use the jaeger exporter, 

to configure jaeger endpoint,
set enviroment `OTEL_EXPORTER_JAEGER_ENDPOINT` or pass `jaeger.WithEndpoint(url)` to this function


### grpc example


for now only unary interceptor is provided

for server, just pass `grpctrace.UnaryServerInterceptor()` to server's unaryServerInterceptor

for client, just pass `grpctrace.UnaryClientInterceptor()` to grpc.Dial function


### http example


for now, only gin middleware is provided

for server, use `engine.Use(httptrace.NewGinHandler())`

for client:
```
req := http.NewRequestWithContext(ctx, method, url, body)
resp, err := httptrace.DefaultClient.Do(req)
handle(resp, err)
```





