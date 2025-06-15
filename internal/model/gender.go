package model

type Gender string

const (
	GMale   Gender = "male"
	GFemale Gender = "female"
)

func IsValidGender(g string) bool {
	return g == string(GMale) || g == string(GFemale)
}
