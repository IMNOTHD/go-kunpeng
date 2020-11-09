package service

import (
	"database/sql"
	"errors"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
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
		zap.L().Warn("no such user: " + userId)
		return errors.New("no such user: " + userId)
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

	zap.L().Info("cache user " + userId + " succeed")

	return nil
}

func CacheUserByGrade(db *sql.DB, rc *redis.Client, grade string) (int32, error) {
	var err error

	ul, err := QueryUserByGrade(db, grade)
	if err != nil {
		return 0, err
	}
	if *ul == nil {
		zap.L().Warn("no such grade: " + grade)
		return 0, errors.New("no such grade: " + grade)
	}

	var successCount int32 = 0

	for _, u := range *ul {
		err = CleanUserAllInfo(rc, u.UserInfo.UserID)
		if err != nil {
			continue
		}

		err = SetUserInfoRedis(rc, &u.UserInfo)
		if err != nil {
			continue
		}
		err = AddRoleInfoRedis(rc, u.UserInfo.UserID, &u.RoleInfo)
		if err != nil {
			continue
		}
		err = AddJobInfoRedis(rc, u.UserInfo.UserID, &u.JobInfo)
		if err != nil {
			continue
		}
		err = SetAvatarUrlRedis(rc, u.UserInfo.UserID, &u.AvatarUrl)
		if err != nil {
			continue
		}
		successCount++
	}

	if int(successCount) != len(*ul) {
		zap.L().Warn("some cache failed")
		return successCount, errors.New("some cache failed")
	} else {
		zap.L().Info("cache grade " + grade + " succeed")

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
		zap.L().Error("no such class: " + class)
		return 0, errors.New("no such class: " + class)
	}

	var successCount int32 = 0

	for _, u := range *ul {
		err = CleanUserAllInfo(rc, u.UserInfo.UserID)
		if err != nil {
			continue
		}

		err = SetUserInfoRedis(rc, &u.UserInfo)
		if err != nil {
			continue
		}
		err = AddRoleInfoRedis(rc, u.UserInfo.UserID, &u.RoleInfo)
		if err != nil {
			continue
		}
		err = AddJobInfoRedis(rc, u.UserInfo.UserID, &u.JobInfo)
		if err != nil {
			continue
		}
		err = SetAvatarUrlRedis(rc, u.UserInfo.UserID, &u.AvatarUrl)
		if err != nil {
			continue
		}
		successCount++
	}

	if int(successCount) != len(*ul) {
		zap.L().Warn("some cache failed")
		return successCount, errors.New("some cache failed")
	} else {
		zap.L().Info("cache class " + class + " succeed")

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
		zap.L().Error("db query all user error")
		return 0, errors.New("db query all user error")
	}

	var successCount int32 = 0

	for _, u := range *ul {
		err = CleanUserAllInfo(rc, u.UserInfo.UserID)
		if err != nil {
			continue
		}

		err = SetUserInfoRedis(rc, &u.UserInfo)
		if err != nil {
			continue
		}
		err = AddRoleInfoRedis(rc, u.UserInfo.UserID, &u.RoleInfo)
		if err != nil {
			continue
		}
		err = AddJobInfoRedis(rc, u.UserInfo.UserID, &u.JobInfo)
		if err != nil {
			continue
		}
		err = SetAvatarUrlRedis(rc, u.UserInfo.UserID, &u.AvatarUrl)
		if err != nil {
			continue
		}
		successCount++
	}

	if int(successCount) != len(*ul) {
		zap.L().Warn("some cache failed")
		return successCount, errors.New("some cache failed")
	} else {
		zap.L().Info("cache all user succeed")

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

	err = AddActivityRecordRedis(rc, ar, false)
	if err != nil {
		return err
	}

	zap.L().Info("cache user " + userId + " activity record succeed")

	return nil
}

func CacheActivityRecordByGrade(db *sql.DB, rc *redis.Client, grade string) (int32, error) {
	var err error

	ul, err := QueryUserByGrade(db, grade)
	if err != nil {
		return 0, err
	}
	if *ul == nil {
		zap.L().Error("no such grade: " + grade)
		return 0, errors.New("no such grade: " + grade)
	}

	var successCount int32 = 0

	for _, u := range *ul {
		err = CacheSingleUserAllActivityRecord(db, rc, u.UserInfo.UserID)
		if err != nil {
			continue
		}

		successCount++
	}

	if int(successCount) != len(*ul) {
		zap.L().Warn("some cache failed")
		return successCount, errors.New("some cache failed")
	} else {
		zap.L().Info("cache grade " + grade + " activity record succeed")

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
		zap.L().Error("no such class: " + class)
		return 0, errors.New("no such class: " + class)
	}

	var successCount int32 = 0

	for _, u := range *ul {
		err = CacheSingleUserAllActivityRecord(db, rc, u.UserInfo.UserID)
		if err != nil {
			continue
		}

		successCount++
	}

	if int(successCount) != len(*ul) {
		zap.L().Error("some cache failed")
		return successCount, errors.New("some cache failed")
	} else {
		zap.L().Info("cache class " + class + " activity record succeed")

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
		zap.L().Error("db query all user error")
		return 0, errors.New("db query all user error")
	}

	var successCount int32 = 0

	for _, u := range *ul {
		err = CacheSingleUserAllActivityRecord(db, rc, u.UserInfo.UserID)
		if err != nil {
			continue
		}

		successCount++
	}

	if int(successCount) != len(*ul) {
		zap.L().Warn("some cache failed")
		return successCount, errors.New("some cache failed")
	} else {
		zap.L().Info("cache all user activity record succeed")

		return successCount, nil
	}
}

func CacheAllActivity(db *sql.DB, rc *redis.Client) (int32, error) {
	var successCount int32 = 0

	// TODO

	return successCount, nil
}
