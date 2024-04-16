package role

type Role string

const (
	EditGroupImage           Role = "edit_group_image"
	EditGroupName            Role = "edit_group_name"
	EditGroupRoles           Role = "edit_group_roles"
	EditGroupTheme           Role = "edit_group_theme"
	DeleteGroup              Role = "delete_group"
	InviteUser               Role = "invite_user"
	KickUser                 Role = "kick_user"
	BanUser                  Role = "ban_user"
	SendMessage              Role = "send_message"
	DeleteMessage            Role = "delete_message"
	DeleteOtherMemberMessage Role = "delete_other_member_message"
	CreateIntegration        Role = "create_integration"
)
