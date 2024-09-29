# protoc-gateway-impl

Automatically creates implementation for [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) generated interfaces.

## Result:

IN: 
```bash
message PingRequest {
    string msg = 1;
}

message PingResponse {
    string msg = 1;
}

service BasicService {
  rpc Ping(PingRequest) returns (PingResponse) {
    option (google.api.http) = {
      post: "/ping"
      body: "*"
    };
  }
}
```
OUT:
service.auto.go
```go
package basic_service

import (
	"context"
	"net/http"
	gw "github.com/EvgenyOvsov/protoc-gateway-impl/example/api/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type BasicService interface {
	Apply(ctx context.Context, httpMux *http.ServeMux) error
	Ping(ctx context.Context, req *gw.PingRequest) (*gw.PingResponse, error)
	PingV2(ctx context.Context, req *gw.PingRequest) (*gw.PingResponse, error)
}

func (s *basicService) Apply(ctx context.Context, httpMux *http.ServeMux) error {
	grpcMux := runtime.NewServeMux()
	if err := gw.RegisterBasicServiceHandlerServer(ctx, grpcMux, s); err != nil {
		return err
	}
	httpMux.Handle("/ping", grpcMux)
	httpMux.Handle("/v1/ping", grpcMux)
	return nil
}
```
service.go
```go
package basic_service

import (
	gw "github.com/EvgenyOvsov/protoc-gateway-impl/example/api/proto"
)

type basicService struct {
	gw.UnimplementedBasicServiceServer
}

func NewBasicService() BasicService {
	return &basicService{}
}
```
ping.go
```go
package basic_service

import (
	"context"
	gw "github.com/EvgenyOvsov/protoc-gateway-impl/example/api/proto"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *basicService) Ping(ctx context.Context, req *gw.PingRequest) (*gw.PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method not implemented")
}
```

For details use `make example` ([buf](https://buf.build/) must be installed) or read code, I don't care...
