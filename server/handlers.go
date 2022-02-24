package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type addCounterReq struct {
	I int64 `uri:"i" binding:"required,min=1"`
}

type Count struct {
	Counter int64 `json:"counter"`
}

type Substr struct {
	Text string `json:"text"`
}
type RespBody struct {
	Text string `json:"text"`
	Len  int    `json:"len"`
}

type FindString struct {
	Str string `uri:"str" binding:"required"`
}

// incrby i
func (s *Server) addCounter(ctx *gin.Context) {
	var req addCounterReq

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incr, err := s.store.Client.IncrBy("counter", req.I).Result()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"counter": incr})
}

// decrby i
func (s *Server) subCounter(ctx *gin.Context) {
	var req addCounterReq

	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decr, err := s.store.Client.DecrBy("counter", req.I).Result()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"counter": decr})
}

// get value of counter
func (s *Server) getCounter(ctx *gin.Context) {

	get, err := s.store.Client.Get("counter").Int64()
	if err == redis.Nil {
		ctx.JSON(http.StatusOK, "counter doesn't exist")
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"counter": get})
}

//substr
func (s *Server) getSubstr(ctx *gin.Context) {

	substr := Substr{}

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = json.Unmarshal(body, &substr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	str := substr.Text

	longStr := ""
	curStr := ""

	for i := 0; i < len(str); i++ {
		char := string(str[i])

		if strings.Contains(curStr, char) {
			if curStr > longStr {
				longStr = curStr
			}
			curStr = strings.Split(curStr, char)[1]
		}
		curStr = curStr + char
	}

	ans := ""
	if len(curStr) > len(longStr) {
		ans = curStr
	} else {
		ans = longStr
	}

	resp := RespBody{Text: ans, Len: len(ans)}

	ctx.JSON(http.StatusOK, resp)
}

// find string
func (s *Server) findStr(ctx *gin.Context) {

	var uri FindString
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	routes := s.router.Routes()

	paths := make([]string, 0)

	for _, route := range routes {
		paths = append(paths, route.Path)
	}

	result := make([]string, 0)

	for _, v := range paths {
		str := strings.Split(v, "/")

		for _, s := range str {
			if uri.Str == s {
				result = append(result, v)
				break
			}
		}
	}

	ctx.JSON(http.StatusOK, result)
}
