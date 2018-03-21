package Utils

import (
	"time"
	"strconv"
	"log"
	"github.com/gin-gonic/gin"
)
const DISCORD_EPOCH int64 = 1420070400000

func GetTimeFromSnowflake(id string) (time.Time, error){
	iid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return time.Now(), err
	}

	return time.Unix(((iid>>22)+DISCORD_EPOCH)/1000, 0).UTC(), nil
}


func GetTimeFromSnowflakeInt(id int64) (time.Time, error){
	return time.Unix(((id>>22)+DISCORD_EPOCH)/1000, 0).UTC(), nil
}

func GetGuildIdString(id string)(int64, bool){
	if id == ""{
		return 0, false
	} else{
		intId, err :=strconv.ParseInt(id, 10, 64)
		if err!= nil{
			log.Println("Failed parse of GuildId, ", err)
			return 0, false
		} else{
			return intId, true
		}
	}
}

func GetGuildId(params *gin.Params) (int64, bool){
	id := params.ByName("id")
	return GetGuildIdString(id)
}

func ValidIdInt(id int64) bool{
	t, err := GetTimeFromSnowflakeInt(id)
	if err != nil{
		return false
	}
	if t.After(time.Now()){
		return false
	}
	if t.Unix() < DISCORD_EPOCH/1000{
		return false
	}
	return true
}

func ValidId(id string) bool{
	t, err := GetTimeFromSnowflake(id)
	if err != nil{
		return false
	}
	if t.After(time.Now()){
		return false
	}
	if t.Unix() < DISCORD_EPOCH/1000{
		return false
	}
	return true
}
