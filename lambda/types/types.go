package types

import "golang.org/x/crypto/bcrypt"



type User struct{
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func NewUser(user User)(User,error){
	// hash password using bcrypt
	hashedPass,err:=bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)

	if err!=nil{
		return User{},err
	}

	return User{
		Username: user.Username,
		Email: user.Email,
		Password: string(hashedPass),
	},nil
}

// func ValidatePassword(hashedPassword,plain)