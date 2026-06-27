package utility

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const employeeJwtSecret = "bit303_shop_employee_secret"

type EmployeeClaims struct {
	EmployeeId uint   `json:"employee_id"`
	Username   string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateEmployeeToken(employeeId uint, username string, remember bool) (token string, expireAt time.Time, err error) {
	expireAt = time.Now().Add(24 * time.Hour)
	if remember {
		expireAt = time.Now().Add(7 * 24 * time.Hour)
	}
	claims := EmployeeClaims{
		EmployeeId: employeeId,
		Username:   username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(employeeJwtSecret))
	return
}

func ParseEmployeeToken(token string) (*EmployeeClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &EmployeeClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(employeeJwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsedToken.Claims.(*EmployeeClaims)
	if !ok || !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

const adminJwtSecret = "bit303_shop_admin_secret"

type AdminClaims struct {
	AdminId  uint   `json:"admin_id"`
	Username string `json:"username"`
	IsSuper  int    `json:"is_super"`
	jwt.RegisteredClaims
}

func GenerateAdminToken(adminId uint, username string, isSuper int, remember bool) (token string, expireAt time.Time, err error) {
	expireAt = time.Now().Add(24 * time.Hour)
	if remember {
		expireAt = time.Now().Add(7 * 24 * time.Hour)
	}
	claims := AdminClaims{
		AdminId:  adminId,
		Username: username,
		IsSuper:  isSuper,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(adminJwtSecret))
	return
}

func ParseAdminToken(token string) (*AdminClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(adminJwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsedToken.Claims.(*AdminClaims)
	if !ok || !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
