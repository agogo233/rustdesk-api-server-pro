package app

import (
	"rustdesk-api-server-pro/app/controller/admin"
	"rustdesk-api-server-pro/app/controller/api"
	"rustdesk-api-server-pro/app/middleware"
	"rustdesk-api-server-pro/config"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func SetRoute(app *iris.Application, cfg *config.ServerConfig) {
	apiParty := app.Party("/api")
	apiMvc := mvc.New(apiParty)
	apiMvc.Handle(new(api.SystemController))
	apiMvc.Handle(new(api.LoginController))
	apiMvc.Handle(new(api.AuditController))

	apiWithAuthParty := app.Party("/api")
	apiWithAuthParty.Use(middleware.ApiAuth(app))
	{
		apiWithAuthMvc := mvc.New(apiWithAuthParty)
		apiWithAuthMvc.Handle(new(api.UserController))
		apiWithAuthMvc.Handle(new(api.PeerController))
		apiWithAuthMvc.Handle(new(api.AddressBookController))
		apiWithAuthMvc.Handle(new(api.AddressBookPeerController))
		apiWithAuthMvc.Handle(new(api.AddressBookTagController))
	}

	adminPath := cfg.HttpConfig.AdminPath
	if adminPath == "" || adminPath == "/" {
		adminPath = "/admin"
	}

	adminParty := app.Party(adminPath)
	adminParty.Use(middleware.RateLimit(cfg.SecurityConfig.IpRateLimitPerMinute))
	adminMvc := mvc.New(adminParty)
	adminMvc.Handle(new(admin.AuthController))

	adminWithAuthParty := app.Party(adminPath)
	adminWithAuthParty.Use(middleware.AdminAuth(app))
	{
		adminWithAuthMvc := mvc.New(adminWithAuthParty)
		adminWithAuthMvc.Handle(new(admin.IndexController))
		adminWithAuthMvc.Handle(new(admin.DashboardController))
		adminWithAuthMvc.Handle(new(admin.UsersController))
		adminWithAuthMvc.Handle(new(admin.SessionsController))
		adminWithAuthMvc.Handle(new(admin.AuditController))
		adminWithAuthMvc.Handle(new(admin.MailTemplateController))
		adminWithAuthMvc.Handle(new(admin.MaiLogsController))
		adminWithAuthMvc.Handle(new(admin.DevicesController))
	}
}
