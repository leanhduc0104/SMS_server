package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int        `json:"id" gorm:"column:id;type:uuid;uuid_generate_v4();primary_key"`
	Username  string     `gorm:"column:username;uniqueIndex;not null" json:"username"`
	Password  string     `json:"-"`
	Role      string     `json:"role"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

// HashPassword hashes the user's password using bcrypt
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password is correct
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
