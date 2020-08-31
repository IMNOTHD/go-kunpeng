package server

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	pb "go-kunpeng/api"
	"go-kunpeng/service"
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

// 由于一些奇怪的设计, 这里所有的Handler的error均返回nil, 请读取Code
// -------------------------------------

func (c CacheUserService) CacheSingleUser(ctx context.Context, request *pb.CacheSingleUserRequest) (*pb.CacheSingleUserResponse, error) {
	db, err := service.CreateMysqlWorker()
	if err != nil {
		log.Println(err.Error())
		return &pb.CacheSingleUserResponse{Code: pb.CacheUserResponseCode_FAIL, Msg: err.Error()}, nil
	}
	rc, err := service.CreateRedisClient()
	if err != nil {
		log.Println(err.Error())
		return &pb.CacheSingleUserResponse{Code: pb.CacheUserResponseCode_FAIL, Msg: err.Error()}, nil
	}

	u, err := service.QueryUserInfoByUserId(db, request.UserId)
	if err != nil {
		log.Println(err.Error())
		return &pb.CacheSingleUserResponse{Code: pb.CacheUserResponseCode_FAIL, Msg: err.Error()}, nil
	}
	err = service.SetUserInfoRedis(rc, u)
	if err != nil {
		log.Println(err.Error())
		return &pb.CacheSingleUserResponse{Code: pb.CacheUserResponseCode_FAIL, Msg: err.Error()}, nil
	}
	// todo 只做了userinfo

	return &pb.CacheSingleUserResponse{Code: pb.CacheUserResponseCode_SUCCESS}, nil
}

func (c CacheUserService) CacheMultiSingleUser(s pb.CacheUser_CacheMultiSingleUserServer) error {
	// todo 该服务需完全的重构, 做成多线程模式, 这里仅是做个示例
	var err error
	db, err := service.CreateMysqlWorker()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	rc, err := service.CreateRedisClient()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for {
		r, err := s.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Println(err.Error())
			_ = s.Send(&pb.CacheMultiSingleUserResponse{
				Code:   pb.CacheUserResponseCode_FAIL,
				Msg:    err.Error(),
				UserId: r.UserId,
			})
		}
		if err != nil {
			log.Println(err.Error())
			_ = s.Send(&pb.CacheMultiSingleUserResponse{
				Code:   pb.CacheUserResponseCode_FAIL,
				Msg:    err.Error(),
				UserId: r.UserId,
			})
		}

		u, err := service.QueryUserInfoByUserId(db, r.UserId)
		if err != nil {
			log.Println(err.Error())
			_ = s.Send(&pb.CacheMultiSingleUserResponse{
				Code:   pb.CacheUserResponseCode_FAIL,
				Msg:    err.Error(),
				UserId: r.UserId,
			})
		}
		err = service.SetUserInfoRedis(rc, u)
		if err != nil {
			log.Println(err.Error())
			_ = s.Send(&pb.CacheMultiSingleUserResponse{
				Code:   pb.CacheUserResponseCode_FAIL,
				Msg:    err.Error(),
				UserId: r.UserId,
			})
		}

		_ = s.Send(&pb.CacheMultiSingleUserResponse{
			Code: pb.CacheUserResponseCode_SUCCESS,
		})

	}

	return nil
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
