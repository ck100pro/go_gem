package main

import (
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claim struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var (
	mySigningKey string = "jwt_test"
)

func main() {
	r := gin.Default()
	r.GET("/get_jwt_token", getJwtToken)
	r.POST("/valid_jwt_token", validJwtToken)
	r.Run()
}

func getJwtToken(c *gin.Context) {
	var cl Claim
	now := time.Now()
	key := []byte(mySigningKey)

	if err := c.ShouldBindJSON(&cl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	cl.Id = "jwt_id"                                               // Id
	cl.ExpiresAt = now.Add(time.Duration(10) * time.Minute).Unix() // 權證到期時間 => 10 分鐘後
	// cl.NotBefore = now.Unix()                                      // 權證於設定的時間點前無效 => 發方當下時間點前都無效
	cl.Issuer = "jwt_user"       // 權證發行人
	cl.IssuedAt = now.Unix()     // 權證發行時間
	cl.Subject = "valid account" // 簽證主題

	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
	token, err := tokenClaim.SignedString(key)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func validJwtToken(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	token := strings.Split(auth, "Bearer ")[1]

	tokenClaim, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if claim, ok := tokenClaim.Claims.(*Claim); ok && tokenClaim.Valid {
		c.JSON(http.StatusOK,
			gin.H{"Password": claim.Password,
				"Accound":   claim.Account,
				"Issuer":    claim.Issuer,
				"ExpiresAt": claim.ExpiresAt,
				"IssuedAt":  claim.IssuedAt,
				"Subject":   claim.Subject,
				"ID":        claim.Id})
	} else {
		c.JSON(http.StatusOK, gin.H{"err": err.Error()})
	}
}
