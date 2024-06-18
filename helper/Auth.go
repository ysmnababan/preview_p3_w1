package helper

import (
	"net/http"
	"preview/logger"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// do the auth here
		tokenString := ctx.Request().Header.Get("Auth")
		response := map[string]interface{}{}
		if tokenString == "" {
			logger.Logging(ctx).Warning("unable to get the token")
			response["message"] = "unautorized"
			return ctx.JSON(http.StatusUnauthorized, response)
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte("secret-key"), nil
		})

		if err != nil || !token.Valid {
			logger.Logging(ctx).Warning("token invalid")
			response["message"] = "unauthorized"
			return ctx.JSON(http.StatusUnauthorized, response)
		}

		// change token -> struct
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Set("user_id", claims["user_id"])
			ctx.Set("email", claims["email"])
		} else {
			logger.Logging(ctx).Warning("invalid claims")
			response["message"] = "invalid claims"
			return ctx.JSON(http.StatusUnauthorized, response)
		}

		return next(ctx)
	}
}
