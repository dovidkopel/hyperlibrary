package common

type User struct {
	ClientId string   `json:"clientId"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
}

type UserWithFees struct {
	User
	FeesOwed  map[string]float64 `json:"feesOwed"`
	TotalOwed float64            `json:"totalOwed"`
}

func MakeEmptyUser() User {
	user := User{}
	user.Roles = make([]string, 0)
	return user
}
