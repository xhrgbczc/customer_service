package filters

import (
	"kf_server/models"
	"kf_server/utils"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
)

// 判断用户是否带有token
func err401(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	ctx.ResponseWriter.WriteHeader(401)
	_ = ctx.Output.Body([]byte("{\"code\": \"401\", \"message\": \"未登录，或登录已失效\", \"data\": \"token expired\"}"))
}

// FilterToken token
var FilterToken = func(ctx *context.Context) {

	// 地址白名单
	whitelist := []string{
		"/api/auth/login",
		"/api/auth/logout",
		"/api/auth/token",
	}

	oldToken := ctx.Input.Header("Authorization")
	isExistInSlice := utils.InExistInSlice(ctx.Request.RequestURI, whitelist)
	isMatch, _ := regexp.MatchString(`^Bearer\s`, oldToken)
	if (isExistInSlice == false && oldToken == "") || !isMatch {
		err401(ctx)
		return
	}
	o := orm.NewOrm()
	auth := models.Auths{Token: oldToken}
	if isExistInSlice == false && oldToken != "" {
		token := strings.Split(oldToken, " ")[1]
		if err := o.Read(&auth, "Token"); err != nil {
			err401(ctx)
			return
		}

		if err := utils.ValidateToken(token); err != nil {
			err401(ctx)
			return
		}

	}
	// 判断是否需要刷新token
	token := strings.Split(oldToken, " ")[1]
	parmMap, err := utils.DecodeToken(token)
	if err != nil {
		err401(ctx)
		return
	}
	newNum := big.NewRat(1, 1)
	newNum.SetFloat64(parmMap["exp"].(float64))
	exp, _ := strconv.ParseInt(newNum.FloatString(0), 10, 64)

	// 该换token了
	if time.Now().Unix()+60*60*2 >= exp {
		if newToken, err := utils.RefreshToken(token); err == nil {
			auth.Token = newToken
			o.Update(&auth)
			ctx.Output.Header("Authorization", newToken)
			return
		}
	}
	ctx.Output.Header("Authorization", oldToken)
}
