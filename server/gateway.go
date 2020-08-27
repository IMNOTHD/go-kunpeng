package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	pb "go-kunpeng/api"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func ProvideHttp(endpoint string, grpcServer *grpc.Server) *http.Server {
	ctx := context.Background()
	gwmux := runtime.NewServeMux()
	dopts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterCacheUserHandlerFromEndpoint(ctx, gwmux, endpoint, dopts)
	if err != nil {
		log.Fatalf("Register Endpoint err: %v", err)
	}
	//新建mux，它是http的请求复用器
	mux := http.NewServeMux()
	//注册gwmux
	mux.Handle("/", gwmux)
	log.Println(endpoint + " HTTP.Listing...")
	return &http.Server{
		Addr:      endpoint,
		Handler:   grpcHandlerFunc(grpcServer, mux),
		TLSConfig: nil,
	}
}

// grpcHandlerFunc 根据不同的请求重定向到指定的Handler处理
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
