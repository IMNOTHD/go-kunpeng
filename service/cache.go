package service

import (
	"database/sql"
	"errors"

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
