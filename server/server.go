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

	//用服务器 Serve() 方法以及端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	if err = httpServer.Serve(listener); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 由于要兼容grpc-gateway, 这里所有的Handler的error均返回nil, 请直接读取Code获取状态
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

	defer db.Close()
	defer rc.Close()

	err = service.CacheSingleUserAllInfo(db, rc, request.GetUserId())
	if err != nil {
		log.Println(err.Error())
		return &pb.CacheSingleUserResponse{Code: pb.CacheUserResponseCode_FAIL, Msg: err.Error()}, nil

	}

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

		u, err := service.QueryUserByUserId(db, r.UserId)
		if err != nil {
			log.Println(err.Error())
			_ = s.Send(&pb.CacheMultiSingleUserResponse{
				Code:   pb.CacheUserResponseCode_FAIL,
				Msg:    err.Error(),
				UserId: r.UserId,
			})
		}
		err = service.SetUserInfoRedis(rc, &u.UserInfo)
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
	db, err := service.CreateMysqlWorker()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}
	rc, err := service.CreateRedisClient()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}

	defer db.Close()
	defer rc.Close()

	sc, err := service.CacheUserByGrade(db, rc, request.GetGrade())
	if err != nil {
		if sc == 0 {
			log.Println(err.Error())
			return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
		}

		return &pb.MultiCacheUserResponse{
			Code:         pb.CacheUserResponseCode_PARTIAL_SUCCESS,
			SuccessCount: sc,
		}, nil
	}

	return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_SUCCESS}, nil
}

func (c CacheUserService) CacheUserByClass(ctx context.Context, request *pb.CacheUserByClassRequest) (*pb.MultiCacheUserResponse, error) {
	db, err := service.CreateMysqlWorker()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}
	rc, err := service.CreateRedisClient()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}

	defer db.Close()
	defer rc.Close()

	sc, err := service.CacheUserByClass(db, rc, request.GetClass())
	if err != nil {
		if sc == 0 {
			log.Println(err.Error())
			return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
		}

		return &pb.MultiCacheUserResponse{
			Code:         pb.CacheUserResponseCode_PARTIAL_SUCCESS,
			SuccessCount: sc,
		}, nil
	}

	return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_SUCCESS}, nil
}

func (c CacheUserService) CacheAllUser(ctx context.Context, e *empty.Empty) (*pb.MultiCacheUserResponse, error) {
	db, err := service.CreateMysqlWorker()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}
	rc, err := service.CreateRedisClient()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}

	defer db.Close()
	defer rc.Close()

	sc, err := service.CacheAllUser(db, rc)
	if err != nil {
		if sc == 0 {
			log.Println(err.Error())
			return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
		}

		return &pb.MultiCacheUserResponse{
			Code:         pb.CacheUserResponseCode_PARTIAL_SUCCESS,
			SuccessCount: sc,
		}, nil
	}

	return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_SUCCESS}, nil
}

func (c CacheUserService) RemoveAllUserCache(ctx context.Context, e *empty.Empty) (*pb.MultiCacheUserResponse, error) {
	rc, err := service.CreateRedisClient()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}

	defer rc.Close()

	sc, err := service.CleanAllUserAllInfo(rc)
	if err != nil {
		if sc == 0 {
			log.Println(err.Error())
			return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
		}

		return &pb.MultiCacheUserResponse{
			Code:         pb.CacheUserResponseCode_PARTIAL_SUCCESS,
			SuccessCount: sc,
		}, nil
	}

	return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_SUCCESS}, nil
}

// -------------------------------------

func (c CacheActivityService) CacheSingleUserActivityRecord(ctx context.Context, request *pb.CacheSingleUserActivityRecordRequest) (*pb.CacheSingleUserActivityRecordResponse, error) {
	db, err := service.CreateMysqlWorker()
	if err != nil {
		log.Println(err.Error())
		return &pb.CacheSingleUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL, Msg: err.Error()}, nil
	}
	rc, err := service.CreateRedisClient()
	if err != nil {
		log.Println(err.Error())
		return &pb.CacheSingleUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL, Msg: err.Error()}, nil
	}

	defer db.Close()
	defer rc.Close()

	err = service.CacheSingleUserAllActivityRecord(db, rc, request.GetUserId())
	if err != nil {
		log.Println(err.Error())
		return &pb.CacheSingleUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL, Msg: err.Error()}, nil
	}

	return &pb.CacheSingleUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_SUCCESS}, nil
}

func (c CacheActivityService) CacheMultiSingleUserActivityRecord(server pb.CacheActivityRecord_CacheMultiSingleUserActivityRecordServer) error {
	// todo
	panic("implement me")
}

func (c CacheActivityService) CacheUserActivityRecordByGrade(ctx context.Context, request *pb.CacheUserActivityRecordByGradeRequest) (*pb.MultiCacheUserActivityRecordResponse, error) {
	db, err := service.CreateMysqlWorker()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}
	rc, err := service.CreateRedisClient()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}

	defer db.Close()
	defer rc.Close()

	sc, err := service.CacheActivityRecordByGrade(db, rc, request.GetGrade())
	if err != nil {
		if sc == 0 {
			log.Println(err.Error())
			return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
		}

		return &pb.MultiCacheUserActivityRecordResponse{
			Code:         pb.CacheActivityRecordResponseCode_PARTIAL_SUCCESS,
			SuccessCount: sc,
		}, nil
	}

	return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_PARTIAL_SUCCESS}, nil
}

func (c CacheActivityService) CacheUserActivityRecordByClass(ctx context.Context, request *pb.CacheUserActivityRecordByClassRequest) (*pb.MultiCacheUserActivityRecordResponse, error) {
	db, err := service.CreateMysqlWorker()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}
	rc, err := service.CreateRedisClient()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}

	defer db.Close()
	defer rc.Close()

	sc, err := service.CacheActivityRecordByClass(db, rc, request.GetClass())
	if err != nil {
		if sc == 0 {
			log.Println(err.Error())
			return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
		}

		return &pb.MultiCacheUserActivityRecordResponse{
			Code:         pb.CacheActivityRecordResponseCode_PARTIAL_SUCCESS,
			SuccessCount: sc,
		}, nil
	}

	return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_PARTIAL_SUCCESS}, nil
}

func (c CacheActivityService) CacheAllUserActivityRecord(ctx context.Context, e *empty.Empty) (*pb.MultiCacheUserActivityRecordResponse, error) {
	db, err := service.CreateMysqlWorker()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}
	rc, err := service.CreateRedisClient()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}

	defer db.Close()
	defer rc.Close()

	sc, err := service.CacheAllUserActivityRecord(db, rc)
	if err != nil {
		if sc == 0 {
			log.Println(err.Error())
			return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
		}

		return &pb.MultiCacheUserActivityRecordResponse{
			Code:         pb.CacheActivityRecordResponseCode_PARTIAL_SUCCESS,
			SuccessCount: sc,
		}, nil
	}

	return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_PARTIAL_SUCCESS}, nil
}

func (c CacheActivityService) RemoveAllUserActivityRecordCache(ctx context.Context, e *empty.Empty) (*pb.MultiCacheUserActivityRecordResponse, error) {
	rc, err := service.CreateRedisClient()
	if err != nil {
		log.Println(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}

	defer rc.Close()

	sc, err := service.CleanAllUserAllActivityRecord(rc)
	if err != nil {
		if sc == 0 {
			log.Println(err.Error())
			return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
		}

		return &pb.MultiCacheUserActivityRecordResponse{
			Code:         pb.CacheActivityRecordResponseCode_PARTIAL_SUCCESS,
			SuccessCount: sc,
		}, nil
	}

	return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_PARTIAL_SUCCESS}, nil

}
