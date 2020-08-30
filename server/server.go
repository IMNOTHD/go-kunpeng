package server

import (
	"context"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	pb "go-kunpeng/api"
)

type CacheUserService struct{}

type CacheActivityService struct{}

const (
	// Address grpc监听地址
	Address string = ":7111"
	// Network grpc网络通信协议
	Network string = "tcp"
)

func Start() {
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen err: %v\n", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCacheUserServer(grpcServer, &CacheUserService{})
	pb.RegisterCacheActivityRecordServer(grpcServer, &CacheActivityService{})
	log.Println(Address + " net.Listening...")

	httpServer := ProvideHttp(Address, grpcServer)

	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	if err = httpServer.Serve(listener); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// -------------------------------------

func (c CacheUserService) CacheSingleUser(ctx context.Context, request *pb.CacheSingleUserRequest) (*pb.CacheSingleUserResponse, error) {
	panic("implement me")
}

func (c CacheUserService) CacheMultiSingleUser(server pb.CacheUser_CacheMultiSingleUserServer) error {
	panic("implement me")
}

func (c CacheUserService) CacheUserByGrade(ctx context.Context, request *pb.CacheUserByGradeRequest) (*pb.MultiCacheUserResponse, error) {
	panic("implement me")
}

func (c CacheUserService) CacheUserByClass(ctx context.Context, request *pb.CacheUserByClassRequest) (*pb.MultiCacheUserResponse, error) {
	panic("implement me")
}

func (c CacheUserService) CacheAllUser(ctx context.Context, e *empty.Empty) (*pb.MultiCacheUserResponse, error) {
	panic("implement me")
}

func (c CacheUserService) RemoveAllUserCache(ctx context.Context, e *empty.Empty) (*pb.MultiCacheUserResponse, error) {
	panic("implement me")
}

// -------------------------------------

func (c CacheActivityService) CacheSingleUserActivityRecord(ctx context.Context, request *pb.CacheSingleUserActivityRecordRequest) (*pb.CacheSingleUserActivityRecordResponse, error) {
	panic("implement me")
}

func (c CacheActivityService) CacheMultiSingleUserActivityRecord(server pb.CacheActivityRecord_CacheMultiSingleUserActivityRecordServer) error {
	panic("implement me")
}

func (c CacheActivityService) CacheUserActivityRecordByGrade(ctx context.Context, request *pb.CacheUserActivityRecordByGradeRequest) (*pb.MultiCacheUserActivityRecordResponse, error) {
	panic("implement me")
}

func (c CacheActivityService) CacheUserActivityRecordByClass(ctx context.Context, request *pb.CacheUserActivityRecordByClassRequest) (*pb.MultiCacheUserActivityRecordResponse, error) {
	panic("implement me")
}

func (c CacheActivityService) CacheAllUserActivityRecord(ctx context.Context, e *empty.Empty) (*pb.MultiCacheUserActivityRecordResponse, error) {
	panic("implement me")
}

func (c CacheActivityService) RemoveAllUserActivityRecordCache(ctx context.Context, e *empty.Empty) (*pb.MultiCacheUserActivityRecordResponse, error) {
	panic("implement me")
}
