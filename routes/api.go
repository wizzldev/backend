package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/pkg/middlewares"
	"github.com/wizzldev/chat/pkg/utils/role"
)

func RegisterAPI(r fiber.Router) {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./pkg/configs/build.json")
	})

	{
		r.Post("/login", requests.Use[requests.Login](), middlewares.NewAuthLimiter(), handlers.Auth.Login).Name("auth.login")
		r.Get("/allow-ip/:token", handlers.Auth.AllowIP).Name("auth.allow-ip")
		r.Post("/register", requests.Use[requests.Register](), handlers.Auth.Register).Name("auth.register")
		r.Get("/logout", handlers.Auth.Logout).Name("auth.logout")
		r.Post("/request-new-password", requests.Use[requests.NewPassword](), handlers.Auth.RequestNewPassword).Name("auth.request-password")
		r.Get("/set-new-password/:token", handlers.Auth.IsPasswordResetExists).Name("auth.is-set-password-exists")
		r.Post("/set-new-password/:token", requests.Use[requests.SetNewPassword](), handlers.Auth.SetNewPassword).Name("auth.set-password")
		r.Get("/verify-email/:token", handlers.Auth.VerifyEmail).Name("auth.verify-email")
		r.Post("/request-new-email-verification", requests.Use[requests.Email](), handlers.Auth.RequestNewEmailVerification).Name("auth.request-email-verification")
	}

	auth := r.Group("/", middlewares.AnyAuth).Name("main.")
	{
		auth.Get("/me", handlers.Me.Hello).Name("me")
		auth.Put("/me", middlewares.NoBots, requests.Use[requests.UpdateMe](), middlewares.NewSimpleLimiter(3, 10*time.Minute, "Too many modifications, try again later"), handlers.Me.Update).Name("me")
		auth.Get("/me/ip-check", handlers.Me.SwitchIPCheck).Name("ip-switch")
		auth.Post("/me/profile-image", middlewares.NoBots, handlers.Me.UploadProfileImage).Name("profile-image")
		// wait a week before deletion
		// auth.Delete("/me", middlewares.NoBots, handlers.Me.Delete)
	}

	mobile := r.Group("/mobile", middlewares.NoBots).Name("mobile.")
	{
		mobile.Post("/register-push-notification", requests.Use[requests.PushToken](), handlers.Mobile.RegisterPushNotification).Name("register-push")
	}

	security := auth.Group("/security", middlewares.NoBots).Name("security.")
	{
		security.Get("/sessions", handlers.Security.Sessions).Name("sessions")
		security.Delete("/sessions", handlers.Security.DestroySessions).Name("sessions")
		security.Delete("/sessions/:id<int>", handlers.Security.DestroySession).Name("sessions.single")
		security.Get("/ips", handlers.Security.IPs).Name("ips")
		security.Delete("/ips/:id<int>", handlers.Security.DestroyIP).Name("ips.single")
		security.Use(HandleNotFoundError)
	}

	users := auth.Group("/users", middlewares.NoBots).Name("users.")
	{
		users.Post("/findByEmail", requests.Use[requests.Email](), handlers.Users.FindByEmail).Name("find:email")
		users.Use(HandleNotFoundError)
	}

	auth.Get("/themes", handlers.Theme.Paginate).Name("themes")
	auth.Get("/chat/contacts", handlers.Chat.Contacts).Name("contacts")
	auth.Get("/chat/user/:id<int>", middlewares.GroupAccess("id"), handlers.Group.GetInfo).Name("chat.user")
	auth.Get("/chat/private/:id<int>", handlers.Chat.PrivateMessage).Name("chat.user:private")
	auth.Post("/chat/search", requests.Use[requests.SearchContacts](), handlers.Chat.Search).Name("chat.search")

	auth.Post("/chat/group", requests.Use[requests.NewGroup](), handlers.Group.New).Name("chat.group")
	auth.Get("/chat/roles", handlers.Group.GetAllRoles).Name("chat.roles")

	chat := auth.Group("/chat/:id<int>", middlewares.GroupAccess("id")).Name("chat.")
	{
		chat.Get("/", handlers.Chat.Find).Name("find")
		chat.Put("/", middlewares.NewRoleMiddleware(role.EditGroupName), requests.Use[requests.EditGroupName](), handlers.Group.EditName).Name("name")
		chat.Get("/paginate", handlers.Chat.Messages).Name("paginate")
		chat.Put("/invite", middlewares.NewRoleMiddleware(role.Creator), requests.Use[requests.CustomInvite](), middlewares.NewSimpleLimiter(3, 15*time.Minute, "too many requests"), handlers.Group.CustomInvite).Name("invite")
		chat.Put("/emoji", middlewares.NewRoleMiddleware(role.Admin), requests.Use[requests.Emoji](), handlers.Group.Emoji).Name("emoji")
		chat.Get("/leave", handlers.Group.Leave).Name("leave")
		chat.Delete("/", middlewares.NewRoleMiddleware(role.Creator), handlers.Group.Delete).Name("delete")
		chat.Put("/roles", middlewares.NewRoleMiddleware(role.Admin), requests.Use[requests.ModifyRoles](), handlers.Group.ModifyRoles).Name("roles")
		chat.Get("/message/:messageID", handlers.Chat.FindMessage)
		chat.Post("/new-invite", middlewares.NewRoleMiddleware(role.InviteUser), requests.Use[requests.NewInvite](), middlewares.NewSimpleLimiter(3, 10*time.Minute, "Try again later before creating another"), handlers.Invite.Create).Name("invite")

		chat.Post("/file", middlewares.NewRoleMiddleware(role.AttachFile), handlers.Chat.UploadFile).Name("upload-file")
		chat.Post("/group-image", middlewares.NewRoleMiddleware(role.EditGroupImage), handlers.Group.UploadGroupImage).Name("image")

		chat.Get("/users", handlers.Group.Users).Name("users")
		chat.Get("/user_count", handlers.Group.UserCount).Name("user-count")

		chat.Put("/theme/:themeID", middlewares.NewRoleMiddleware(role.EditGroupTheme), handlers.Group.SetTheme).Name("theme")
		chat.Delete("/theme", middlewares.NewRoleMiddleware(role.EditGroupTheme), handlers.Group.RemoveTheme).Name("theme")

		chat.Post("/nickname/:userID", middlewares.NewRoleMiddleware(role.Admin), requests.Use[requests.Nickname](), handlers.GroupUser.EditNickName).Name("nickname")
		chat.Delete("/nickname/:userID", middlewares.NewRoleMiddleware(role.Admin), handlers.GroupUser.RemoveNickName).Name("nickname")

		chat.Use(HandleNotFoundError)
	}

	auth.Get("/invite/:code", handlers.Invite.Describe).Name("invite")
	auth.Get("/invite/:code/use", handlers.Invite.Use).Name("invite:use")

	// bot := r.Group("/bots", middlewares.Auth)
	{
		// bot.Get("/")
		// bot.Post("/")
		// bot.Post("/image")
		// bot.Put("/")
		// bot.Delete("/")
		// bot.Use(HandleNotFoundError)
	}

	auth.Use(HandleNotFoundError)
	r.Use(HandleNotFoundError)
}
