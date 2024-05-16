package role

import "errors"

type Role string

const (
	Creator                  Role = "CREATOR" // The user who created the chat, no one can do anything with him
	Admin                    Role = "ADMIN"   // He can do anything with anyone except the Creator
	EditGroupImage           Role = "EDIT_GROUP_IMAGE"
	EditGroupName            Role = "EDIT_GROUP_NAME"
	EditGroupRoles           Role = "EDIT_GROUP_ROLES"
	EditGroupTheme           Role = "EDIT_GROUP_THEME"
	InviteUser               Role = "INVITE_USER"
	KickUser                 Role = "KICK_USER"
	SendMessage              Role = "SEND_MESSAGE"
	DeleteMessage            Role = "DELETE_MESSAGE"
	DeleteOtherMemberMessage Role = "DELETE_OTHER_MEMBER_MESSAGE"
	CreateIntegration        Role = "CREATE_INTEGRATION"
)

func New(s string) (Role, error) {
	switch Role(s) {
	case Creator:
		return Creator, nil
	case Admin:
		return Admin, nil
	case EditGroupImage:
		return EditGroupImage, nil
	case EditGroupName:
		return EditGroupName, nil
	case EditGroupRoles:
		return EditGroupRoles, nil
	case EditGroupTheme:
		return EditGroupTheme, nil
	case InviteUser:
		return InviteUser, nil
	case KickUser:
		return KickUser, nil
	case SendMessage:
		return SendMessage, nil
	case DeleteMessage:
		return DeleteMessage, nil
	case DeleteOtherMemberMessage:
		return DeleteOtherMemberMessage, nil
	case CreateIntegration:
		return CreateIntegration, nil
	}

	return "", errors.New("this role does not exist")
}

func All() *Roles {
	var roles Roles
	roles = append(roles, Creator, Admin, EditGroupImage, EditGroupName, EditGroupRoles, EditGroupTheme, InviteUser, KickUser, SendMessage, DeleteMessage, DeleteOtherMemberMessage, CreateIntegration)
	return &roles
}
