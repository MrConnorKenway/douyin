package redis

import (
	"douyin/common/constant"
	"douyin/message/infra/dal/model"
	redisModel "douyin/message/infra/redis/model"
	"encoding/json"
	"strconv"
)

func AddMessage(userId int64, toUserId int64, messageRedis redisModel.MessageRedis) error {
	redisConn := redisPool.Get()
	defer redisConn.Close()
	var key string
	if userId < toUserId {
		key = constant.MessageRedisPrefix + strconv.FormatInt(userId, 10) + ":" + strconv.FormatInt(toUserId, 10)
	} else {
		key = constant.MessageRedisPrefix + strconv.FormatInt(toUserId, 10) + ":" + strconv.FormatInt(userId, 10)
	}

	ub, err := json.Marshal(messageRedis)
	if err != nil {
		return err
	}
	_, err = redisConn.Do("zadd", key, messageRedis.CreateTime, ub)
	if err != nil {
		redisConn.Do("del", key)
		return err
	}

	expireTime := expireTimeUtil.GetRandTime()
	_, err = redisConn.Do("expire", key, expireTime)
	if err != nil {
		return err
	}

	return nil
}

func AddMessageList(userId int64, toUserId int64, messageListp []*model.Message) error {
	redisConn := redisPool.Get()
	defer redisConn.Close()

	var key string
	if userId < toUserId {
		key = constant.MessageRedisPrefix + strconv.FormatInt(userId, 10) + ":" + strconv.FormatInt(toUserId, 10)
	} else {
		key = constant.MessageRedisPrefix + strconv.FormatInt(toUserId, 10) + ":" + strconv.FormatInt(userId, 10)
	}

	for i := range messageListp {
		messageRedis := redisModel.MessageRedis{
			MessageId:  messageListp[i].MessageUUId,
			FromUserId: messageListp[i].FromUserId,
			ToUserId:   messageListp[i].ToUserId,
			Content:    messageListp[i].Contents,
			CreateTime: messageListp[i].CreateTime,
		}
		ub, err := json.Marshal(messageRedis)
		if err != nil {
			return err
		}
		_, err = redisConn.Do("zadd", key, messageRedis.CreateTime, ub)
		if err != nil {
			redisConn.Do("del", key)
			return err
		}
	}

	expireTime := expireTimeUtil.GetRandTime()
	_, err := redisConn.Do("expire", key, expireTime)
	if err != nil {
		return err
	}

	return nil
}
