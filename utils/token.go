package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-template/models"
	"go-template/pkg"
	"time"
)

const TokenExpirationDuration = 24 * time.Hour // 设置 token 过期时间为24小时

// 检查数据库是否有token
func CheckToken(c *fiber.Ctx) error {
	// 从请求头中获取token 和 user ID
	token := c.Get("token")
	fmt.Println(token)
	user := models.UserBasic{} // 假设用户模型为 User

	db := pkg.DB

	if err := db.Where("Identity = ?", token).First(&user).Error; err != nil {
		// 这里处理数据库错误，例如返回 500 错误
		return c.JSON(pkg.JSONResponse{
			Code: pkg.CodeErrToken,
			JSON: struct {
				Message   string `json:"message"`
				MessageZh string `json:"message_zh"`
			}{"token错误! ", "token错误！"},
		})
	}

	if user.Identity != token || user.Identity == "" {
		// Token 不匹配，返回错误
		return c.JSON(pkg.JSONResponse{
			Code: pkg.CodeErrToken,
			JSON: struct {
				Message   string `json:"message"`
				MessageZh string `json:"message_zh"`
			}{"token错误! ", "token错误！"},
		})
	}

	if time.Since(user.UpdatedAt) > TokenExpirationDuration {
		// Token 过期，返回错误
		return c.JSON(pkg.JSONResponse{
			Code: pkg.CodeErrToken,
			JSON: struct {
				Message   string `json:"message"`
				MessageZh string `json:"message_zh"`
			}{"token过期! ", "token过期！"},
		})
	}

	// Token 验证成功
	return c.Next()
}
