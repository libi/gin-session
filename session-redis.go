package gsession

import (
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"fmt"
	"sync"
)

var Redis_Prefix = "lbcache:"
var Redis_life = 20000
type redisStore struct{
	redis redis.Conn
	lock sync.RWMutex
}

func (r *redisStore) Get(k string,sid string) interface{}{
	r.lock.RLock()
	defer r.lock.RUnlock()
	jsondata,err := r.getSessionData(sid)

	if(err != nil){
		return nil
	}
	return jsondata[k]
}
func (r *redisStore) Set(k string,v interface{},sid string) error{
	r.lock.Lock()
	defer r.lock.Unlock()
	jsondata,err := r.getSessionData(sid)
	if(err != nil){
		fmt.Println(err)
		jsondata = make(map[string]interface{})
	}

	jsondata[k] = v
	jsonString,err := json.Marshal(jsondata)
	if(err != nil){
		fmt.Println(err)
		return err
	}
	redisstring := string(jsonString)

	r.redis.Do("SETEX", Redis_Prefix+sid, Redis_life, redisstring)
	return nil
}
func (r *redisStore)getSessionData(sid string) (jsondata map[string]interface{},err error){
	if(r.redis == nil){
		r.initRedisConn()
	}
	data,err := redis.String(r.redis.Do("GET",Redis_Prefix+sid))
	if(err != nil){
		fmt.Println(err)
		return nil,err
	}

	if err := json.Unmarshal([]byte(data), &jsondata); err != nil {
		fmt.Println(err)
		return nil,err
	}
	return jsondata,nil
}
func (r *redisStore)initRedisConn() {
	fmt.Println("链接redis服务器")
	c,err := redis.DialURL("redis://192.168.6.108")

	if(err != nil){
		fmt.Println("redis初始化失败")
		fmt.Println(err)
	}
	r.redis = c
}
func init(){
	driver := new(redisStore)
	Register("redis",driver)
}

