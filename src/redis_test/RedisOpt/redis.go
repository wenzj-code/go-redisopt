package RedisOpt

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	redis "gopkg.in/redis.v5"
)

//RedisOpt .
type RedisOpt struct {
	redisClusterClient *redis.ClusterClient
	clusterAddr        []string

	redisSingleClient *redis.Client
	signleAddr        string

	redisPasswd string
	//redis mode: 1 single mode, 2 cluster mode
	redisMothed int
}

//InitSingle .
//Init single client
func (opt *RedisOpt) InitSingle(redisAddr, redisPasswd string) error {
	opt.redisMothed = 1
	opt.signleAddr = redisAddr
	opt.redisPasswd = redisPasswd
	opt.redisSingleClient = redis.NewClient(&redis.Options{
		Addr:     opt.signleAddr,
		Password: opt.redisPasswd,
		DB:       0,
	})

	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := opt.redisSingleClient.Ping().Result()
	return err
}

//InitCluster .
//Init cluster client
func (opt *RedisOpt) InitCluster(ClusterAddr []string, redisPasswd string) error {
	opt.redisMothed = 2
	opt.clusterAddr = ClusterAddr
	opt.redisPasswd = redisPasswd
	opt.redisClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    opt.clusterAddr,
		Password: opt.redisPasswd,
	})

	statusCmd := opt.redisClusterClient.Ping()
	_, err := statusCmd.Result()
	return err
}

func (opt *RedisOpt) reconnectRedis() {
	log.Error("reconnect RedisServer")
	var err error
	if opt.redisMothed == 2 {
		err = opt.InitCluster(opt.clusterAddr, opt.redisPasswd)
	} else {
		err = opt.InitSingle(opt.signleAddr, opt.redisPasswd)
	}
	if err != nil {
		log.Error("reconnect Redis failed")
		return
	}
	log.Info("reconnect Redis success")
}

//Set .
func (opt *RedisOpt) Set(key string, data interface{}, sec int) error {
	var cmd *redis.StatusCmd
	if opt.redisMothed == 2 {
		cmd = opt.redisClusterClient.Set(key, data, time.Duration(sec)*time.Second)
	} else {
		cmd = opt.redisSingleClient.Set(key, data, time.Duration(sec)*time.Second)
	}
	_, err := cmd.Result()
	if err != nil {
		opt.reconnectRedis()
	}
	return err
}

//Get .
func (opt *RedisOpt) Get(key string) ([]byte, error) {
	var cmd *redis.StringCmd
	if opt.redisMothed == 2 {
		cmd = opt.redisClusterClient.Get(key)
	} else {
		cmd = opt.redisSingleClient.Get(key)
	}
	data, err := cmd.Bytes()
	if err != nil {
		opt.reconnectRedis()
		return nil, err
	}
	return data, nil
}

//Delete .
func (opt *RedisOpt) Delete(key string) error {
	var cmd *redis.IntCmd
	if opt.redisMothed == 2 {
		cmd = opt.redisClusterClient.Del(key)
	} else {
		cmd = opt.redisSingleClient.Del(key)
	}
	if cmd.Err() != nil {
		opt.reconnectRedis()
		return cmd.Err()
	}
	return nil
}

//-------------------------------
//HExists .
func (opt *RedisOpt) HExists(key, field string) bool {
	if opt.redisMothed == 2 {
		return false
	}
	cmd := opt.redisSingleClient.HExists(key, field)

	status, err := cmd.Result()
	if err != nil {
		opt.reconnectRedis()
	}
	return status
}

//HSet .
func (opt *RedisOpt) HSet(key, field string, value interface{}, second int) error {
	if opt.redisMothed == 2 {
		return errors.New("it's single methed not cluster")
	}
	cmd := opt.redisSingleClient.HSet(key, field, value)
	_, err := cmd.Result()
	if err != nil {
		opt.reconnectRedis()
	}
	if second > 0 {
		opt.redisSingleClient.Expire(key, time.Duration(second)*time.Second)
	}
	return err
}

//HMSet .
func (opt *RedisOpt) HMSet(key string, fields map[string]string, second int) error {
	if opt.redisMothed == 2 {
		return errors.New("it's single methed not cluster")
	}
	cmd := opt.redisSingleClient.HMSet(key, fields)
	_, err := cmd.Result()
	if err != nil {
		opt.reconnectRedis()
	}
	if second > 0 {
		opt.redisSingleClient.Expire(key, time.Duration(second)*time.Second)
	}
	return err
}

//HMGet .
func (opt *RedisOpt) HGetAll(key string) (map[string]string, error) {
	if opt.redisMothed == 2 {
		return nil, errors.New("it's single methed not cluster")
	}
	cmd := opt.redisSingleClient.HGetAll(key)

	data, err := cmd.Result()
	if err != nil {
		opt.reconnectRedis()
		return nil, err
	}
	return data, nil
}

//HGet .
func (opt *RedisOpt) HGet(key, field string) (string, error) {
	if opt.redisMothed == 2 {
		return "", errors.New("it's single methed not cluster")
	}
	cmd := opt.redisSingleClient.HGet(key, field)

	data, err := cmd.Result()
	if err != nil {
		opt.reconnectRedis()
		return "", err
	}
	return data, nil
}

//HDelete .
func (opt *RedisOpt) HDelete(key string) error {
	if opt.redisMothed == 2 {
		return errors.New("it's single methed not cluster")
	}
	cmd := opt.redisSingleClient.HDel(key)
	if cmd.Err() != nil {
		opt.reconnectRedis()
		return cmd.Err()
	}
	return nil
}
