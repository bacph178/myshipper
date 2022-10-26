package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"myshipper/infrastructure"
	"myshipper/models"
	"myshipper/utils/token"
	"net/http"
	"os"
	"strings"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := token.TokenValid(ctx)
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func EnforceAuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("currentUser")
		if exists && user.(models.User).ID != 0 {
			return
		} else {
			err, _ := c.Get("authErr")
			_ = c.AbortWithError(http.StatusUnauthorized, err.(error))
			return
		}
	}
}

func UserLoaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		if bearer != "" {
			jwtPath := strings.Split(bearer, " ")
			if len(jwtPath) == 2 {
				jwtEncode := jwtPath[1]
				token, err := jwt.Parse(jwtEncode, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
					}
					secret := []byte(os.Getenv("JWT_SECRET"))
					return secret, nil
				})
				if err != nil {
					println(err.Error())
					return
				}
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					userId := uint(claims["user_id"].(float64))
					fmt.Printf("[+] Authenticated request, authenticated user id is %d\n", userId)
					var user models.User
					if userId != 0 {
						database := infrastructure.GetDb()
						database.Preloads("Roles").First(&user, userId)
					}
					c.Set("currentUser", user)
					c.Set("currentUserId", user.ID)
				}
			}
		}
	}
}
