package types

import "golang.org/x/crypto/bcrypt"



type User struct{
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	PhotoURL    string  `json:"photo_url"`
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
		PhotoURL: user.PhotoURL,
	},nil
}

// func ValidatePassword(hashedPassword,plain)