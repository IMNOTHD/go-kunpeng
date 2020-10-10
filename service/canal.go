package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/golang/protobuf/proto"

	"github.com/withlin/canal-go/client"
	protocol "github.com/withlin/canal-go/protocol"

	cc "go-kunpeng/config/canal"
	"go-kunpeng/model"
)

type EntryData struct {
	// 字段mysql类型, 拿到的是切掉()之后转小写的文本, 原始的type来自desc命令
	MysqlType string
	// 字段值,timestamp,datetime是一个时间格式的文本
	Value string
	// 是否为null
	IsNull bool
}

type M map[string]EntryData

func StartCanalClient() {
	connector := client.NewSimpleCanalConnector(cc.Address, cc.Port, cc.Username, cc.Password, cc.Destination, cc.SoTimeOut, cc.IdleTimeOut)
	err := connector.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// https://github.com/alibaba/canal/wiki/AdminGuide
	//mysql 数据解析关注的表，Perl正则表达式.
	//
	//多个正则之间以逗号(,)分隔，转义符需要双斜杠(\\)
	//
	//常见例子：
	//
	//  1.  所有表：.*   or  .*\\..*
	//	2.  canal schema下所有表： canal\\..*
	//	3.  canal下的以canal打头的表：canal\\.canal.*
	//	4.  canal schema下的一张表：canal\\.test1
	//  5.  多个规则组合使用：canal\\..*,mysql.test1,mysql.test2 (逗号分隔)

	err = connector.Subscribe(".*\\..*")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for {
		message, err := connector.Get(100, nil, nil)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		batchId := message.Id
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(cc.PollingInterval * time.Millisecond)
			// log.Println("===暂时没有数据更新===")
			continue
		}

		// 打数据变更log
		go printEntry(message.Entries)
		// 数据处理协程
		go handleData(message.Entries)
	}
}

func handleData(es []protocol.Entry) {
	rc, err := CreateRedisClient()
	if err != nil {
		log.Println(err)
		return
	}
	defer rc.Close()

	db, err := CreateMysqlWorker()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	for _, entry := range es {
		if entry.GetEntryType() == protocol.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == protocol.EntryType_TRANSACTIONEND {
			continue
		}

		rowChange := new(protocol.RowChange)
		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		header := entry.GetHeader()
		eventType := rowChange.GetEventType()

		switch header.GetTableName() {
		case "common_user_info":
			for _, rowData := range rowChange.GetRowDatas() {

				if eventType == protocol.EventType_DELETE {
					mBefore := unmarshalData(rowData.GetBeforeColumns())
					_ = CleanUserInfoRedis(rc, mBefore["user_id"].Value)
				} else {
					mAfter := unmarshalData(rowData.GetAfterColumns())

					var d model.UserInfo

					err := Assign(&d, mAfter)
					if err != nil {
						log.Println(err)
					}

					err = SetUserInfoRedis(rc, &d)
					if err != nil {
						log.Println(err)
					}
				}
			}
		case "common_user_role_relation":
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == protocol.EventType_INSERT {
					mAfter := unmarshalData(rowData.GetAfterColumns())

					x := make(model.RoleInfo, 0)

					roleCode, err := QueryRoleCodeByRoleId(db, mAfter["user_role_id"].Value)

					x = append(x, roleCode)

					err = AddRoleInfoRedis(rc, mAfter["user_id"].Value, &x)
					if err != nil {
						log.Println(err)
					}
				} else {
					mBefore := unmarshalData(rowData.GetBeforeColumns())
					err := CleanRoleInfoRedis(rc, mBefore["user_id"].Value)
					if err != nil {
						log.Println(err)
						continue
					}

					r, err := QueryRoleInfoByUserId(db, mBefore["user_id"].Value)
					if err != nil {
						log.Println(err)
						continue
					}

					err = AddRoleInfoRedis(rc, mBefore["user_id"].Value, r)
					if err != nil {
						log.Println(err)
					}
				}
			}
		case "common_role":
			// 不用管, 修改了这个表需要重新跑全量缓存, 改起来太麻烦了
			log.Println("table common_role changed, you HAVE TO rerun full user info cache")
		case "organization_member":
			// 如果只是插入了一行不用管, 否则需要重新缓存全部的jobInfo
			for _, rowData := range rowChange.GetRowDatas() {
				mAfter := unmarshalData(rowData.GetAfterColumns())
				j, err := QueryJobInfoByUserId(db, mAfter["member_id"].Value)
				if err != nil {
					log.Println(err)
					continue
				}

				err = CleanJobInfoRedis(rc, mAfter["member_id"].Value)
				if err != nil {
					log.Println(err)
					continue
				}

				err = AddJobInfoRedis(rc, mAfter["member_id"].Value, j)
				if err != nil {
					log.Println(err)
				}
			}
		case "common_user":
			// 此表缓存的只有avatar_url字段, 其他不用管
			for _, rowData := range rowChange.GetRowDatas() {
				mAfter := unmarshalData(rowData.GetAfterColumns())

				var avatarUrl model.AvatarUrl
				if mAfter["avatar_url"].IsNull {
					avatarUrl = model.AvatarUrl{Url: ""}
				} else {
					avatarUrl = model.AvatarUrl{Url: mAfter["avatar_url"].Value}
				}

				err := SetAvatarUrlRedis(rc, mAfter["user_id"].Value, &avatarUrl)

				if err != nil {
					log.Println(err)
				}
			}
		case "activity_record":
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == protocol.EventType_INSERT {
					mAfter := unmarshalData(rowData.GetAfterColumns())

					var x model.ActivityRecord

					err := Assign(&x, mAfter)
					if err != nil {
						log.Println(err)
						continue
					}

					su, err := QueryUserInfoByUserId(db, x.ScannerUserID)
					if err != nil {
						log.Println(err)
						continue
					}

					x.ActivityTime = fmt.Sprintf("%.1f", float64(x.Time/10))
					x.ScannerName = su.RealName

					err = AddActivityRecordRedis(rc, &[]model.ActivityRecord{x})
					if err != nil {
						log.Println(err)
					}
				} else {
					mAfter := unmarshalData(rowData.GetAfterColumns())
					err := CacheSingleUserAllActivityRecord(db, rc, mAfter["user_id"].Value)
					if err != nil {
						log.Println(err)
					}
				}
			}
		case "activity":
			for _, rowData := range rowChange.GetRowDatas() {
				mAfter := unmarshalData(rowData.GetAfterColumns())
				if eventType == protocol.EventType_DELETE {
					err := CleanActivityRedis(rc, mAfter["activity_id"].Value)
					if err != nil {
						log.Println(err)
					}
				} else if eventType == protocol.EventType_INSERT || eventType == protocol.EventType_UPDATE {
					mAfter := unmarshalData(rowData.GetAfterColumns())

					var x model.Activity
					err := Assign(&x, mAfter)
					if err != nil {
						log.Println(err)
						continue
					}

					err = SetActivityRedis(rc, mAfter["activity_id"].Value, &x)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

func printEntry(entrys []protocol.Entry) {

	for _, entry := range entrys {
		if entry.GetEntryType() == protocol.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == protocol.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := new(protocol.RowChange)

		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		checkError(err)
		eventType := rowChange.GetEventType()
		header := entry.GetHeader()
		log.Println(fmt.Sprintf("================> binlog[%s : %d],name[%s,%s], eventType: %s", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType()))
		for _, rowData := range rowChange.GetRowDatas() {
			if eventType == protocol.EventType_DELETE {
				printColumn(rowData.GetBeforeColumns())
			} else if eventType == protocol.EventType_INSERT {
				printColumn(rowData.GetAfterColumns())
			} else {
				log.Println("-------> before")
				printColumn(rowData.GetBeforeColumns())
				log.Println("-------> after")
				printColumn(rowData.GetAfterColumns())
			}
		}
	}
}

func printColumn(columns []*protocol.Column) {
	for _, col := range columns {
		log.Println(fmt.Sprintf("%s : %s  update= %t", col.GetName(), col.GetValue(), col.GetUpdated()))
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %s", err.Error())
	}
}

func unmarshalData(columns []*protocol.Column) M {
	m := make(M)

	for _, col := range columns {
		m[col.GetName()] = EntryData{
			MysqlType: strings.Split(col.GetMysqlType(), "(")[0],
			Value:     col.GetValue(),
			IsNull:    col.GetIsNull(),
		}
	}

	return m
}

func Assign(ptr interface{}, m M) error {
	// 获取入参的类型
	reType := reflect.TypeOf(ptr)
	// 入参类型校验
	if reType.Kind() != reflect.Ptr || reType.Elem().Kind() != reflect.Struct {
		return errors.New("ptr not a struct ptr")
	}
	// 取指针指向的结构体变量
	v := reflect.ValueOf(ptr).Elem()

	for i := 0; i < v.NumField(); i++ {
		// 获取结构体字段信息
		structField := v.Type().Field(i)

		// 取db tag
		tag := structField.Tag
		// 获取tag值, 无tag则跳过
		dbField := tag.Get("db")
		if dbField == "" {
			continue
		}

		if e, ok := m[dbField]; ok {
			switch e.MysqlType {
			//.....如果用到更多的数据类型, 请自行添加
			case "datetime":
				if e.IsNull {
					v.Field(i).SetInt(0)
					continue
				}

				// 转化所需模板
				timeLayout := "2006-01-02 15:04:05"
				// 获取时区
				loc, _ := time.LoadLocation("Local")
				tmp, _ := time.ParseInLocation(timeLayout, e.Value, loc)
				// 转化为时间戳 类型是int64
				timestamp := tmp.Unix()
				v.Field(i).SetInt(timestamp)
			case "bigint", "int":
				if e.IsNull {
					v.Field(i).SetInt(0)
					continue
				}

				tmp, _ := strconv.ParseInt(e.Value, 10, 64)
				v.Field(i).SetInt(tmp)
			case "varchar":
				if v.Field(i).Kind() == reflect.String {
					if e.IsNull {
						v.Field(i).SetString("")
						continue
					}
					v.Field(i).SetString(e.Value)
				} else if v.Field(i).Kind() == reflect.Map {
					if e.IsNull {
						continue
					}

					tmp := make(map[string]interface{})
					err := json.Unmarshal([]byte(e.Value), &tmp)
					if err != nil {
						log.Println("Unmarshal failed:", err)
					}

					// 黑科技, 别问, 我也不知道怎么写出来的
					*(*map[string]interface{})(unsafe.Pointer(v.Field(i).Addr().Pointer())) = tmp

				} else {
					log.Println("who fucking did this stupid variable, do write this stupid reflect")
				}
			}
		}
	}

	return nil
}
