package server

import (
	"fmt"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	inn "github.com/tit/go-inn-validator"
)

type BodyEmail struct {
	Emails []string `json:"emails"`
}

type BodyInn struct {
	Numbers []string `json:"inn_numbers"`
}

func (s *Server) chekEmail(ctx *gin.Context) {
	body := BodyEmail{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	emails := make([]string, 0)

	for _, v := range body.Emails {
		e, ok := isEmailValid(v)
		if !ok {
			fmt.Println(e, "not valid")
			continue
		}
		fmt.Println(e)
		emails = append(emails, e)

	}

	ctx.JSON(http.StatusOK, gin.H{"valid_emails": emails})
}

func isEmailValid(email string) (string, bool) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", false
	}
	return addr.Address, true
}

func (s *Server) checkINN(ctx *gin.Context) {
	body := BodyInn{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	validInn := make([]string, 0)

	for _, v := range body.Numbers {
		if ok := isValid(v); !ok {
			fmt.Println(v, "is not valid")
			continue
		}
		validInn = append(validInn, v)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"valid_inn": validInn,
	})
}

func isValid(numbers string) bool {
	ok, err := inn.IsPrivatePersonInnValid(numbers)
	if err != nil {
		return ok
	}
	return ok
}
