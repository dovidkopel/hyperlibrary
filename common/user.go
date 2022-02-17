package common

type User struct {
	ClientId string `json:"clientId"`
	Name     string `json:"name"`
}

type UserWithFees struct {
	User
	FeesOwed  map[string]float64 `json:"feesOwed"`
	TotalOwed float64            `json:"totalOwed"`
}
