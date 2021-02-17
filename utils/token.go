package utils

import (
	"fmt"
	"kf_server/models"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
)

// KEY ...
// DEFAULTEXPIRESECONDS ...
const (
	KEY                  string = "JWT-ARY-STARK"
	DEFAULTEXPIRESECONDS int    = 259200
)

// MyCustomClaims JWT -- json web token
// HEADER PAYLOAD SIGNATURE
// This struct is the PAYLOAD
type MyCustomClaims struct {
	models.JwtKeyDto
	jwt.StandardClaims
}

// RefreshToken Refresh token
func RefreshToken(tokenString string) (string, error) {
	// first get previous token
	token, err := jwt.ParseWithClaims(
		tokenString, &MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(KEY), nil
		})
	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok || !token.Valid {
		return "", err
	}
	mySigningKey := []byte(KEY)
	expireAt := time.Now().Add(time.Second * time.Duration(DEFAULTEXPIRESECONDS)).Unix()
	newClaims := MyCustomClaims{
		claims.JwtKeyDto,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			Issuer:    claims.JwtKeyDto.UserName,
			IssuedAt:  time.Now().Unix(),
		},
	}
	// generate new token with new claims
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenStr, err := newToken.SignedString(mySigningKey)
	logs.Info(tokenStr)
	if err != nil {
		logs.Error("generate new fresh json web token failed !! error :", err)
		return "", err
	}
	return "Bearer " + tokenStr, err
}

// ValidateToken 验证token
func ValidateToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(KEY), nil
		})

	if err != nil && !token.Valid {
		logs.Error("validate tokenString failed !!!", err)
		return err
	}
	return nil
}

// DecodeToken 解析token
func DecodeToken(token string) (map[string]interface{}, error) {
	parseAuth, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(KEY), nil
	})
	//将token中的内容存入parmMap
	claim := parseAuth.Claims.(jwt.MapClaims)
	var parmMap map[string]interface{}
	parmMap = make(map[string]interface{})
	for key, val := range claim {
		parmMap[key] = val
	}
	return parmMap, err
}

// GenerateToken 生成token
func GenerateToken(jwtKeyDto models.JwtKeyDto) (tokenString string) {
	// Create the Claims
	mySigningKey := []byte(KEY)
	expireAt := time.Now().Add(time.Second * time.Duration(DEFAULTEXPIRESECONDS)).Unix()
	// pass parameter to this func or not
	claims := MyCustomClaims{
		jwtKeyDto,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			Issuer:    jwtKeyDto.UserName,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("generate json web token failed !! error :", err)
	}
	return "Bearer " + tokenStr

}

// return this result to client then all later request should have header "Authorization: Bearer <token> "
func getHeaderTokenValue(tokenString string) string {
	//Authorization: Bearer <token>
	return fmt.Sprintf("Bearer %s", tokenString)
}
