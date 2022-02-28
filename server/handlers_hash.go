package server

import (
	"fmt"
	"hash/crc64"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type Uri struct {
	Id string `uri:"id" binding:"required"`
}

// generate hash
func (s *Server) genHash(ctx *gin.Context) {

	input := Body{}
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var (
		//new uuid for client
		uid = uuid.New()
		//timer for 1 minute
		mintute = time.Minute * 1
		seconds = time.Second * 5
	)

	go func() {
		timer := time.NewTimer(mintute)
		ticker := time.NewTicker(seconds)
		//decimal
		var dec int

		for {
			select {
			case <-timer.C:
				log.Println("timer is done")

				err := s.store.Client.Set(fmt.Sprint(uid), dec, 0).Err()
				if err != nil {
					fmt.Println("Set uuid err:", err)
					return
				}

				// log.Printf("dec=%d\n", dec)
				return

			case <-ticker.C:
				crcTable := crc64.MakeTable(crc64.ECMA)
				checksum64 := crc64.Checksum([]byte(input.Text), crcTable)

				fmt.Printf("%d\n%d\n", checksum64, time.Now().UnixNano())

				var conjunc int = int(time.Now().UnixNano()) & int(checksum64)
				// fmt.Println("dec", conjunc)

				dec += conjunc

				fmt.Println("dec", dec)
			}
		}
	}()

	_, err := s.store.Client.Set(fmt.Sprint(uid), 0, 0).Result()
	if err != nil {
		log.Fatal("Set uuid err:", err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": uid,
	})

}

func (s *Server) getHash(ctx *gin.Context) {
	uri := Uri{}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	get, err := s.store.Client.Get(uri.Id).Int()
	if err != nil {
		if err.Error() == redis.Nil.Error() {
			ctx.JSON(http.StatusBadRequest, "ID is not valid")
			return
		}

		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if get == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "PENDING",
		})
		return
	}

	bin := strconv.FormatInt(int64(get), 2)

	result := ""

	for _, v := range bin {
		if string(v) == "1" {
			result += string(v)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"hash": result,
	})
}
