package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"
	"github.com/Prabandham/expense_tracker/config"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUuid string
	UserId     string
}

// Generate a token for a given userId
func CreateToken(userid string) (*TokenDetails, error) {
	var err error
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix() // Expire in 60 minutes
	td.AccessUuid = uuid.NewV4().String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix() // Expire in 7 days
	td.RefreshUuid = uuid.NewV4().String()

	// Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(config.GetEnv("ACCESS_SECRET", "")))
	if err != nil {
		return nil, err
	}

	// Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(config.GetEnv("REFRESH_SECRET", "")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

// CreateAuth will persist the data to redis, this will help us invalidate the token as soon as the user logs out.
func CreateAuth(userid string, redisClient *redis.Client) (*TokenDetails, error) {
	td, err := CreateToken(userid)
	if err != nil {
		panic(err)
	}
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := redisClient.Set(td.AccessUuid, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return nil, errAccess
	}
	errRefresh := redisClient.Set(td.RefreshUuid, userid, rt.Sub(now)).Err()
	if errRefresh != nil {
		return nil, errRefresh
	}
	return td, nil
}

func ExtractToken(r *http.Request) string {
  bearToken := r.Header.Get("Authorization")
  //normally Authorization the_token_xxx
  strArr := strings.Split(bearToken, " ")
  if len(strArr) == 2 {
     return strArr[1]
  }
  return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
  tokenString := ExtractToken(r)
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
     //Make sure that the token method conform to "SigningMethodHMAC"
     if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
     }
     return []byte(config.GetEnv("ACCESS_SECRET", "")), nil
  })
  if err != nil {
     return nil, err
  }
  return token, nil
}

func TokenValid(r *http.Request) error {
  token, err := VerifyToken(r)
  if err != nil {
     return err
  }
  if !token.Valid {
     return err
  }
  return nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
  token, err := VerifyToken(r)
  if err != nil {
     return nil, err
  }
  claims, ok := token.Claims.(jwt.MapClaims)
  if ok && token.Valid {
     accessUuid, ok := claims["access_uuid"].(string)
     if !ok {
        return nil, err
     }
     userId, err := claims["user_id"].(string)
     if !err {
        return nil, errors.New("could not fetch user Id from claims")
     }
     return &AccessDetails{
        AccessUuid: accessUuid,
        UserId:   userId,
     }, nil
  }
  return nil, err
}

func FetchAuth(authD *AccessDetails, redisClient *redis.Client) (string, error) {
  userId, err := redisClient.Get(authD.AccessUuid).Result()
  if err != nil {
     return "", err
  }
  return userId, nil
}

func DeleteAuth(givenUuid string, redisClient *redis.Client) (int64,error) {
  deleted, err := redisClient.Del(givenUuid).Result()
  if err != nil {
     return 0, err
  }
  return deleted, nil
}