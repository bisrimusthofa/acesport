package user

type UserFormatter struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Token string `json:"token"`
}

func FormatUser(user User) UserFormatter {
	return UserFormatter{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}
}
