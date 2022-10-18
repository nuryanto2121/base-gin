package util

import (
	"fmt"
	"time"

	"app/pkg/setting"

	"github.com/dgrijalva/jwt-go"
)

// Claims :
type Claims struct {
	UserID   string `json:"id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Role     string `json:"role,omitempty"`
	OutletId string `json:"outlet_id,omitempty"`
	jwt.StandardClaims
}

// GenerateToken :
func GenerateToken(id string, user_name string, role string) (string, error) {

	screet := setting.AppSetting.JwtSecret
	expired_time := setting.AppSetting.ExpiredJwt
	issuer := setting.AppSetting.Issuer
	var jwtSecret = []byte(screet)
	// Set custom claims
	// Ids,_ :=strconv.I(id)
	claims := &Claims{
		UserID:   id,
		UserName: user_name,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(expired_time)).Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(jwtSecret)
}

// ParseToken :
func ParseToken(token string) (*Claims, error) {
	var screet = setting.AppSetting.JwtSecret
	var jwtSecret = []byte(screet)

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// GetEmailToken :
func GetEmailToken(email string) string {
	var screet = setting.AppSetting.JwtSecret
	var jwtSecret = []byte(screet)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

// ParseEmailToken :
func ParseEmailToken(token string) (string, error) {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if err, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return err, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}
	claims, _ := tkn.Claims.(jwt.MapClaims)
	return fmt.Sprintf("%s", claims["email"]), nil
}
