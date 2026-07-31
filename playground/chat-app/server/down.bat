@echo off
:: required for milvus client
go get go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc@v0.46.0
go mod tidy