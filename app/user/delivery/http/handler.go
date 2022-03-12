package http

import "github.com/golang-jwt/jwt"

type jwtCustomClaims struct {
	ID int32 `json:"id"`
	jwt.StandardClaims
}

//// generate token
//jwtCfg := config.Get().Jwt
//claims := &jwtCustomClaims{
//user.ID,
//jwt.StandardClaims{
//ExpiresAt: time.Now().Add(jwtCfg.ExpireTime * time.Second).Unix(),
//},
//}
//
//unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//token, err = unsignedToken.SignedString([]byte(jwtCfg.SecretKey))
//if err != nil {
//return "", err
//}
//return token, err
