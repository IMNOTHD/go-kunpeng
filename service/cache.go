package service

import (
	"database/sql"
	"errors"
	"log"

	"github.com/go-redis/redis/v8"
)

func CacheSingleUserAllInfo(db *sql.DB, rc *redis.Client, userId string) error {
	var err error

	err = CleanUserAllInfo(rc, userId)
	if err != nil {
		return err
	}

	u, err := QueryUserByUserId(db, userId)
	if err != nil {
		return err
	}
	if &u == nil {
		return errors.New("no such user")
	}

	err = SetUserInfoRedis(rc, &u.UserInfo)
	if err != nil {
		return err
	}
	err = AddRoleInfoRedis(rc, userId, &u.RoleInfo)
	if err != nil {
		return err
	}
	err = AddJobInfoRedis(rc, userId, &u.JobInfo)
	if err != nil {
		return err
	}
	err = SetAvatarUrlRedis(rc, userId, &u.AvatarUrl)
	if err != nil {
		return err
	}

	return nil
}

func CacheUserByGrade(db *sql.DB, rc *redis.Client, grade string) (int32, error) {
	var err error

	ul, err := QueryUserByGrade(db, grade)
	if err != nil {
		return 0, err
	}
	if *ul == nil {
		return 0, errors.New("no such grade")
	}

	var successCount int32 = 0

	for _, u := range *ul {
		err = CleanUserAllInfo(rc, u.UserInfo.UserID)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		err = SetUserInfoRedis(rc, &u.UserInfo)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = AddRoleInfoRedis(rc, u.UserInfo.UserID, &u.RoleInfo)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = AddJobInfoRedis(rc, u.UserInfo.UserID, &u.JobInfo)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = SetAvatarUrlRedis(rc, u.UserInfo.UserID, &u.AvatarUrl)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		successCount++
	}

	if int(successCount) != len(*ul) {
		return successCount, errors.New("some cache failed")
	} else {
		return successCount, nil
	}
}

func CacheUserByClass(db *sql.DB, rc *redis.Client, class string) (int32, error) {
	var err error

	ul, err := QueryUserByClass(db, class)
	if err != nil {
		return 0, err
	}
	if *ul == nil {
		return 0, errors.New("no such class")
	}

	var successCount int32 = 0

	for _, u := range *ul {
		err = CleanUserAllInfo(rc, u.UserInfo.UserID)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		err = SetUserInfoRedis(rc, &u.UserInfo)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = AddRoleInfoRedis(rc, u.UserInfo.UserID, &u.RoleInfo)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = AddJobInfoRedis(rc, u.UserInfo.UserID, &u.JobInfo)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = SetAvatarUrlRedis(rc, u.UserInfo.UserID, &u.AvatarUrl)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		successCount++
	}

	if int(successCount) != len(*ul) {
		return successCount, errors.New("some cache failed")
	} else {
		return successCount, nil
	}
}

func CacheAllUser(db *sql.DB, rc *redis.Client) (int32, error) {
	var err error

	ul, err := QueryUserAll(db)
	if err != nil {
		return 0, err
	}
	if *ul == nil {
		return 0, errors.New("db query all user error")
	}

	var successCount int32 = 0

	for _, u := range *ul {
		err = CleanUserAllInfo(rc, u.UserInfo.UserID)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		err = SetUserInfoRedis(rc, &u.UserInfo)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = AddRoleInfoRedis(rc, u.UserInfo.UserID, &u.RoleInfo)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = AddJobInfoRedis(rc, u.UserInfo.UserID, &u.JobInfo)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = SetAvatarUrlRedis(rc, u.UserInfo.UserID, &u.AvatarUrl)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		successCount++
	}

	if int(successCount) != len(*ul) {
		return successCount, errors.New("some cache failed")
	} else {
		return successCount, nil
	}
}

func CacheSingleUserAllActivityRecord(db *sql.DB, rc *redis.Client, userId string) error {
	var err error

	err = CleanUserAllActivityRecord(rc, userId)
	if err != nil {
		return err
	}

	ar, err := QueryActivityRecordByUserId(db, userId)
	if err != nil {
		return err
	}
	if &ar == nil {
		return nil
	}

	err = AddActivityRecordRedis(rc, ar)
	if err != nil {
		return err
	}

	return nil
}

func CacheActivityRecordByGrade(db *sql.DB, rc *redis.Client, grade string) (int32, error) {
	var err error

	ul, err := QueryUserByGrade(db, grade)
	if err != nil {
		return 0, err
	}
	if *ul == nil {
		return 0, errors.New("no such grade")
	}

	var successCount int32 = 0

	for _, u := range *ul {
		err = CacheSingleUserAllActivityRecord(db, rc, u.UserInfo.UserID)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		successCount++
	}

	if int(successCount) != len(*ul) {
		return successCount, errors.New("some cache failed")
	} else {
		return successCount, nil
	}
}

func CacheActivityRecordByClass(db *sql.DB, rc *redis.Client, class string) (int32, error) {
	var err error

	ul, err := QueryUserByClass(db, class)
	if err != nil {
		return 0, err
	}
	if *ul == nil {
		return 0, errors.New("no such class")
	}

	var successCount int32 = 0

	for _, u := range *ul {
		err = CacheSingleUserAllActivityRecord(db, rc, u.UserInfo.UserID)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		successCount++
	}

	if int(successCount) != len(*ul) {
		return successCount, errors.New("some cache failed")
	} else {
		return successCount, nil
	}
}

func CacheAllUserActivityRecord(db *sql.DB, rc *redis.Client) (int32, error) {
	var err error

	ul, err := QueryUserAll(db)
	if err != nil {
		return 0, err
	}
	if *ul == nil {
		return 0, errors.New("db query all user error")
	}

	var successCount int32 = 0

	for _, u := range *ul {
		err = CacheSingleUserAllActivityRecord(db, rc, u.UserInfo.UserID)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		successCount++
	}

	if int(successCount) != len(*ul) {
		return successCount, errors.New("some cache failed")
	} else {
		return successCount, nil
	}
}
