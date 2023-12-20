package user

import "time"

const (
	AvatarPath = "images/avatar/"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password" gorm:"column:password"`
	Role     string `json:"role"`
	// EmailVerifiedAt time.Time `json:"email_verified_at"`
	Phone          string    `json:"phone"`
	RefferalCode   string    `json:"refferal_code"`
	RefferalCodeIn string    `json:"refferal_code_in"`
	Bank           string    `json:"bank"`
	Norek          string    `json:"norek"`
	RememberToken  string    `json:"remember_token"`
	FcmToken       string    `json:"fcm_token"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Balance        int       `json:"balance"`
	Avatar         string    `json:"avatar"`
}
