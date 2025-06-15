package model

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

func IsValidGender(g string) bool {
	return g == string(Male) || g == string(Female)
}
