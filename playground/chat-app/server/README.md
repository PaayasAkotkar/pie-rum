


## in-case of grpc error use this
``` 
go mod edit -replace go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc=go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc@v0.45.0

go mod tidy

```

## just for testing we are not using the text embeddding


## Walkthrough
📁 admin -> pie-rum init flow
📁 server/services -> in ask-chess-coach the profile is being pushed 
📁 rag -> ingredients to make the rag complete
flow -> contains the main func where the model will be placed
google|nvidea -> models that are used


## any tips?
1. study the admin that would teach you one way of using the pie-rum
2. study the flow how I have created it and used
3. ignore store, marlin
4. do changes and make it work & also please do comment you guys; where and what's actually you been stuck on and why this or that been added

## NOTE
in case if something went wrong simply add the ```go work init ./pie-rum .``` and ``` go work sync ```