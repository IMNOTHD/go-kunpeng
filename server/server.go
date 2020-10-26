package server

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
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
		zap.L().Fatal(fmt.Sprintf("net.Listen err: %v", err))
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCacheUserServer(grpcServer, &CacheUserService{})
	pb.RegisterCacheActivityRecordServer(grpcServer, &CacheActivityService{})
	zap.L().Info(Address + " net.Listening...")

	httpServer := ProvideHttp(Address, grpcServer)

	//用服务器 Serve() 方法以及端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	if err = httpServer.Serve(listener); err != nil {
		zap.L().Fatal(fmt.Sprintf("ListenAndServe: %v", err))
	}
}

// 由于要兼容grpc-gateway, 这里所有的Handler的error均返回nil, 请直接读取Code获取状态
// -------------------------------------

func (c CacheUserService) CacheSingleUser(ctx context.Context, request *pb.CacheSingleUserRequest) (*pb.CacheSingleUserResponse, error) {
	db, err := service.CreateMysqlWorker()
	if err != nil {
		return &pb.CacheSingleUserResponse{Code: pb.CacheUserResponseCode_FAIL, Msg: err.Error()}, nil
	}
	defer db.Close()

	rc, err := service.CreateRedisClient()
	if err != nil {
		return &pb.CacheSingleUserResponse{Code: pb.CacheUserResponseCode_FAIL, Msg: err.Error()}, nil
	}
	defer rc.Close()

	err = service.CacheSingleUserAllInfo(db, rc, request.GetUserId())
	if err != nil {
		return &pb.CacheSingleUserResponse{Code: pb.CacheUserResponseCode_FAIL, Msg: err.Error()}, nil

	}

	return &pb.CacheSingleUserResponse{Code: pb.CacheUserResponseCode_SUCCESS}, nil
}

func (c CacheUserService) CacheMultiSingleUser(s pb.CacheUser_CacheMultiSingleUserServer) error {
	const (
		// 允许channel缓冲的长度
		_commodityLoad = 2048
		// 消费者数量
		_consumerNumber = 32
	)

	var err error
	db, err := service.CreateMysqlWorker()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	defer db.Close()

	rc, err := service.CreateRedisClient()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	defer rc.Close()

	var wg sync.WaitGroup

	ch := make(chan string, _commodityLoad)

	consumer := func(c chan string, db *sql.DB, rc *redis.Client) {
		defer wg.Done()

		for {
			userId, ok := <-c

			if !ok {
				return
			}

			err = service.CacheSingleUserAllInfo(db, rc, userId)
			if err != nil {
				xErr := s.Send(&pb.CacheMultiSingleUserResponse{
					Code:   pb.CacheUserResponseCode_FAIL,
					Msg:    err.Error(),
					UserId: userId,
				})
				if xErr != nil {
					zap.L().Error(xErr.Error())
				}
			} else {
				xErr := s.Send(&pb.CacheMultiSingleUserResponse{
					Code:   pb.CacheUserResponseCode_SUCCESS,
					UserId: userId,
				})
				if xErr != nil {
					zap.L().Error(xErr.Error())
				}
			}
		}
	}

	// 启动消费者
	wg.Add(_consumerNumber)
	for i := 0; i < _consumerNumber; i++ {
		go consumer(ch, db, rc)
	}

	for {
		r, err := s.Recv()

		if err == io.EOF {
			break
		}

		if r == nil {
			continue
		}

		if err != nil {
			zap.L().Error(err.Error())
			_ = s.Send(&pb.CacheMultiSingleUserResponse{
				Code:   pb.CacheUserResponseCode_FAIL,
				Msg:    err.Error(),
				UserId: r.GetUserId(),
			})
			continue
		}

		ch <- r.GetUserId()
	}

	// 等待通道中的数据全部被读取, 关闭channel
	for {
		if len(ch) == 0 {
			break
		}
	}
	close(ch)

	// 等待消费者关闭
	wg.Wait()

	return nil
}

func (c CacheUserService) CacheUserByGrade(ctx context.Context, request *pb.CacheUserByGradeRequest) (*pb.MultiCacheUserResponse, error) {
	db, err := service.CreateMysqlWorker()
	if err != nil {
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}
	defer db.Close()

	rc, err := service.CreateRedisClient()
	if err != nil {
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}
	defer rc.Close()

	sc, err := service.CacheUserByGrade(db, rc, request.GetGrade())
	if err != nil {
		if sc == 0 {
			zap.L().Error(err.Error())
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
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}
	defer db.Close()

	rc, err := service.CreateRedisClient()
	if err != nil {
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}
	defer rc.Close()

	sc, err := service.CacheUserByClass(db, rc, request.GetClass())
	if err != nil {
		if sc == 0 {
			zap.L().Error(err.Error())
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
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}
	defer db.Close()

	rc, err := service.CreateRedisClient()
	if err != nil {
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}
	defer rc.Close()

	sc, err := service.CacheAllUser(db, rc)
	if err != nil {
		if sc == 0 {
			zap.L().Error(err.Error())
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
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserResponse{Code: pb.CacheUserResponseCode_FAIL}, nil
	}
	defer rc.Close()

	sc, err := service.CleanAllUserAllInfo(rc)
	if err != nil {
		if sc == 0 {
			zap.L().Error(err.Error())
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
		zap.L().Error(err.Error())
		return &pb.CacheSingleUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL, Msg: err.Error()}, nil
	}
	defer db.Close()

	rc, err := service.CreateRedisClient()
	if err != nil {
		zap.L().Error(err.Error())
		return &pb.CacheSingleUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL, Msg: err.Error()}, nil
	}
	defer rc.Close()

	err = service.CacheSingleUserAllActivityRecord(db, rc, request.GetUserId())
	if err != nil {
		zap.L().Error(err.Error())
		return &pb.CacheSingleUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL, Msg: err.Error()}, nil
	}

	return &pb.CacheSingleUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_SUCCESS}, nil
}

func (c CacheActivityService) CacheMultiSingleUserActivityRecord(s pb.CacheActivityRecord_CacheMultiSingleUserActivityRecordServer) error {
	const (
		// 允许channel缓冲的长度
		_commodityLoad = 1000
		// 消费者数量
		_consumerNumber = 16
	)

	var err error
	db, err := service.CreateMysqlWorker()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	defer db.Close()

	rc, err := service.CreateRedisClient()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	defer rc.Close()

	var wg sync.WaitGroup

	ch := make(chan string, _commodityLoad)

	consumer := func(c chan string, db *sql.DB, rc *redis.Client) {
		defer wg.Done()

		for {
			userId, ok := <-c

			if !ok {
				return
			}

			err = service.CacheSingleUserAllActivityRecord(db, rc, userId)
			if err != nil {
				zap.L().Error(err.Error())
				xErr := s.Send(&pb.CacheMultiSingleUserActivityRecordResponse{
					Code:   pb.CacheActivityRecordResponseCode_FAIL,
					Msg:    err.Error(),
					UserId: userId,
				})
				if xErr != nil {
					zap.L().Error(xErr.Error())
				}
			} else {
				xErr := s.Send(&pb.CacheMultiSingleUserActivityRecordResponse{
					Code:   pb.CacheActivityRecordResponseCode_PARTIAL_SUCCESS,
					UserId: userId,
				})
				if xErr != nil {
					zap.L().Error(xErr.Error())
				}
			}
		}
	}

	// 启动消费者
	wg.Add(_consumerNumber)
	for i := 0; i < _consumerNumber; i++ {
		go consumer(ch, db, rc)
	}

	for {
		r, err := s.Recv()

		if err == io.EOF {
			break
		}

		if r == nil {
			continue
		}

		if err != nil {
			zap.L().Error(err.Error())
			_ = s.Send(&pb.CacheMultiSingleUserActivityRecordResponse{
				Code:   pb.CacheActivityRecordResponseCode_FAIL,
				Msg:    err.Error(),
				UserId: r.GetUserId(),
			})
			continue
		}

		ch <- r.GetUserId()
	}

	// 等待通道中的数据全部被读取, 关闭channel
	for {
		if len(ch) == 0 {
			break
		}
	}
	close(ch)

	// 等待消费者关闭
	wg.Wait()

	return nil
}

func (c CacheActivityService) CacheUserActivityRecordByGrade(ctx context.Context, request *pb.CacheUserActivityRecordByGradeRequest) (*pb.MultiCacheUserActivityRecordResponse, error) {
	db, err := service.CreateMysqlWorker()
	if err != nil {
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}
	defer db.Close()

	rc, err := service.CreateRedisClient()
	if err != nil {
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}
	defer rc.Close()

	sc, err := service.CacheActivityRecordByGrade(db, rc, request.GetGrade())
	if err != nil {
		if sc == 0 {
			zap.L().Error(err.Error())
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
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}
	defer db.Close()

	rc, err := service.CreateRedisClient()
	if err != nil {
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}
	defer rc.Close()

	sc, err := service.CacheActivityRecordByClass(db, rc, request.GetClass())
	if err != nil {
		if sc == 0 {
			zap.L().Error(err.Error())
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
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}
	defer db.Close()

	rc, err := service.CreateRedisClient()
	if err != nil {
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}
	defer rc.Close()

	sc, err := service.CacheAllUserActivityRecord(db, rc)
	if err != nil {
		if sc == 0 {
			zap.L().Error(err.Error())
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
		zap.L().Error(err.Error())
		return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
	}

	defer rc.Close()

	sc, err := service.CleanAllUserAllActivityRecord(rc)
	if err != nil {
		if sc == 0 {
			zap.L().Error(err.Error())
			return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_FAIL}, nil
		}

		return &pb.MultiCacheUserActivityRecordResponse{
			Code:         pb.CacheActivityRecordResponseCode_PARTIAL_SUCCESS,
			SuccessCount: sc,
		}, nil
	}

	return &pb.MultiCacheUserActivityRecordResponse{Code: pb.CacheActivityRecordResponseCode_PARTIAL_SUCCESS}, nil
}
