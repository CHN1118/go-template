package app

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go-template/models"
	"go-template/pkg"
	"go-template/utils"
	"log"
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

var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true // 注意：这样设置会允许所有的跨域请求，根据您的需求进行调整
	},
}

func SendMsg(c *fiber.Ctx) error {
	log.Println("试图升级到 WebSocket...")
	err := upgrader.Upgrade(c.Context(), func(netConn *websocket.Conn) {
		log.Println("WebSocket 连接建立!")
		defer netConn.Close() // 在函数返回时关闭 WebSocket 连接

		// 在此处订阅 Redis
		sub, err := utils.Subscribe(utils.PublishKey) // 订阅 Redis
		if err != nil {
			log.Printf("订阅 Redis 遇到错误: %v\n", err)
			return
		}
		defer sub.Close() // 在函数返回时关闭 Redis 订阅

		// 使用 goroutine 持续监听来自 Redis 的消息
		go func() {
			for {
				msg, err := sub.ReceiveMessage() // 从 Redis 订阅中读取
				if err != nil {
					log.Printf("从 Redis 订阅中读取时遇到错误: %v\n", err)
					break
				}

				// 输出接收到的消息
				log.Printf("从 Redis 接收到消息: %s\n", msg.Payload)

				// 将消息发送到 WebSocket 客户端
				err = netConn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
				if err != nil {
					log.Printf("写入 WebSocket 时遇到错误: %v\n", err)
					break
				}
			}
		}()

		// 在此处，您可以选择继续执行其他 WebSocket 相关的操作，或者只是阻塞直到连接关闭
		for {
			if _, _, err := netConn.NextReader(); err != nil {
				log.Println("WebSocket 连接关闭!")
				break
			}
		}
	})

	if err != nil {
		log.Printf("WebSocket 升级遇到错误: %v\n", err)
		return c.SendStatus(500)
	}
	return nil
}

//for {
//mt, msg, err := netConn.ReadMessage() // 从 WebSocket 中读取
//if err != nil {
//log.Println("read:", err)
//break
//}
//log.Printf("recv: %s", msg)
//err = utils.Publish(utils.PublishKey, string(msg)) // 发布到 Redis
//if err != nil {
//return
//}
//err = netConn.WriteMessage(mt, msg) // 写入 WebSocket
//if err != nil {
//log.Println("write:", err)
//break
//}
//}
