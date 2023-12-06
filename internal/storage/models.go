package storage

import "encoding/json"

type User struct {
	Login string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

func (u User) String() string {

	data, err := json.Marshal(&u)
	if err != nil {
		return ""
	}
	return string(data)
}