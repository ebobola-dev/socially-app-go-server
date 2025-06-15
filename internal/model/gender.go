package model

type Gender string

const (
	GMale   Gender = "male"
	GFemale Gender = "female"
)

func IsValidGender(g string) bool {
	return g == string(GMale) || g == string(GFemale)
}

func GenderFromString(sg *string) *Gender {
	if sg != nil {
		switch *sg {
		case "male":
			g := GMale
			return &g
		case "female":
			g := GFemale
			return &g
		}
	}
	return nil
}
