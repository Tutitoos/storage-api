package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"net/http"
	"storage-api/src/domain"
)

func IPWhitelistMiddleware(ctx fiber.Ctx) error {
	bypassWhitelist := domain.CONFIG.BypassWhitelist
	if bypassWhitelist == "" {
		return ctx.Status(http.StatusOK).Next()
	}

	headerBypassWhitelist := ctx.Get("Forwarded-For")
	if headerBypassWhitelist != bypassWhitelist {
		return ctx.Status(http.StatusOK).Next()
	}

	ctx.Request().Header.Set(fiber.HeaderXForwardedFor, "BYPASS-IP")

	return ctx.Status(http.StatusOK).Next()
}

// TODO: Revisar este middleware
//func IPWhitelistMiddleware(ctx fiber.Ctx) error {
//	allowedIPs := make(map[string]struct{}, len(domain.CONFIG.WhitelistIps))
//	for _, ip := range domain.CONFIG.WhitelistIps {
//		allowedIPs[ip] = struct{}{}
//	}
//
//	clientIPs := ctx.IPs()
//
//	for _, clientIP := range clientIPs {
//		println(clientIP)
//		if _, exists := allowedIPs[clientIP]; exists {
//			return ctx.Next()
//		}
//	}
//
//	result := domain.ResultData[string]()
//	result.AddMessage("Access Denied")
//	result.AddError(http.StatusForbidden, "IP not allowed")
//
//	return ctx.Status(http.StatusForbidden).JSON(result)
//}
