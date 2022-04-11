package user

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	UserName     string `json:"username" bson:"username"`
	PasswordHash string `json:"-" bson:"password"`
	Email        string `json:"email" bson:"email"`
}

//bson нет, потому что CreateUserDTO не имеет никакого отношения к BD
//Будем брать пользовательские данные, упаковывать в сруктуру CreateUserDTO
//Эту структуру спускать до сервиса
//В сервие превращать эту структуру в модель

type CreateUserDTO struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
