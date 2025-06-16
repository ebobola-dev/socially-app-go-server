package model

type AvatarType string

const (
	ExternalAvatar AvatarType = "external"
	Avatar1        AvatarType = "avatar1"
	Avatar2        AvatarType = "avatar2"
	Avatar3        AvatarType = "avatar3"
	Avatar4        AvatarType = "avatar4"
	Avatar5        AvatarType = "avatar5"
	Avatar6        AvatarType = "avatar6"
	Avatar7        AvatarType = "avatar7"
	Avatar8        AvatarType = "avatar8"
	Avatar9        AvatarType = "avatar9"
	Avatar10       AvatarType = "avatar10"
)

func IsValidAvatarType(at string) bool {
	return at == string(ExternalAvatar) ||
		at == string(Avatar1) ||
		at == string(Avatar2) ||
		at == string(Avatar3) ||
		at == string(Avatar4) ||
		at == string(Avatar5) ||
		at == string(Avatar6) ||
		at == string(Avatar7) ||
		at == string(Avatar8) ||
		at == string(Avatar9) ||
		at == string(Avatar10)
}

func AvatarTypeFromString(sat *string) *AvatarType {
	if sat != nil {
		switch *sat {
		case "external":
			at := ExternalAvatar
			return &at
		case "avatar1":
			at := Avatar1
			return &at
		case "avatar2":
			at := Avatar2
			return &at
		case "avatar3":
			at := Avatar3
			return &at
		case "avatar4":
			at := Avatar4
			return &at
		case "avatar5":
			at := Avatar5
			return &at
		case "avatar6":
			at := Avatar6
			return &at
		case "avatar7":
			at := Avatar7
			return &at
		case "avatar8":
			at := Avatar8
			return &at
		case "avatar9":
			at := Avatar9
			return &at
		case "avatar10":
			at := Avatar10
			return &at
		}
	}
	return nil
}
