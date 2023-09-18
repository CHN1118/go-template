package app

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"go-template/models"
	"go-template/pkg"
	"strconv"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /api/user/getUserList [get]
func GetUserList(c *fiber.Ctx) error {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	return c.JSON(pkg.SuccessResponse(data))
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /api/user/createUser [get]
func CreateUser(c *fiber.Ctx) error {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")
	if password != repassword {
		return c.JSON(pkg.JSONResponse{
			Code: pkg.CodeErr,
			JSON: struct {
				Message   string `json:"message"`
				MessageZh string `json:"message_zh"`
			}{"两次密码不一致！", "两次密码不一致！"},
		})
	}
	user.Salt = pkg.RandomMaxString(10)
	user.PassWord = pkg.MakePassword(password, user.Salt)
	data := models.CreateUser(user)
	return c.JSON(data)
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /api/user/deleteUser [get]
func DeleteUser(c *fiber.Ctx) error {
	user := models.UserBasic{}           // 传入空结构体
	id, _ := strconv.Atoi(c.Query("id")) // string 转 int
	user.ID = uint(id)                   // 赋值
	data := models.DeleteUser(user)
	return c.JSON(data)
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /api/user/updateUser [post]
func UpdateUser(c *fiber.Ctx) error {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.FormValue("id")) // string 转 int
	user.ID = uint(id)                       // 赋值
	user.Name = c.FormValue("name")          // 赋值
	user.PassWord = c.FormValue("password")  // 赋值
	user.Phone = c.FormValue("phone")        // 赋值
	user.Email = c.FormValue("email")        // 赋值
	fmt.Println("update :", user)            // 赋值

	_, err := govalidator.ValidateStruct(user) // 验证参数
	if err != nil {
		fmt.Println(err)
		return c.JSON(pkg.JSONResponse{
			Code: pkg.CodeErr,
			JSON: struct {
				Message   string `json:"message"`
				MessageZh string `json:"message_zh"`
			}{"修改参数不匹配！", "修改参数不匹配！"},
		})
	}

	data := models.UpdateUser(user)
	return c.JSON(data)
}

// FindUserByNameAndPwd
// @Summary 根据用户名和密码查找用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /api/user/findUserByNameAndPwd [get]
func FindUserByNameAndPwd(c *fiber.Ctx) error {
	data := models.UserBasic{}
	name := c.Query("name")
	password := c.Query("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		return c.JSON(pkg.JSONResponse{
			Code: pkg.CodeErr,
			JSON: struct {
				Message   string `json:"message"`
				MessageZh string `json:"message_zh"`
			}{"用户不存在！", "用户不存在！"},
		})

	}
	flag := pkg.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		return c.JSON(pkg.JSONResponse{
			Code: pkg.CodeErr,
			JSON: struct {
				Message   string `json:"message"`
				MessageZh string `json:"message_zh"`
			}{"密码错误！", "密码错误！"},
		})

	}
	pwd := pkg.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	fmt.Println(data)
	fmt.Println(models.CheckToken(data.Identity))
	return c.JSON(pkg.SuccessResponse(data))
}
