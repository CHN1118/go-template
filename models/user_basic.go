package models

import (
	"fmt"
	"go-template/pkg"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model    `json:"-"`
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Name          string         `json:"name" valid:"required"`
	PassWord      string         `json:"password" valid:"required"`
	Phone         string         `json:"phone" valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string         `json:"email" valid:"email"`
	Identity      string         `json:"identity"`
	ClientIp      string         `json:"client_ip"`
	ClientPort    string         `json:"client_port"`
	LoginTime     time.Time      `json:"login_time"`
	HeartbeatTime time.Time      `json:"heartbeat_time"`
	LoginOutTime  time.Time      `json:"login_out_time"`
	IsLogout      bool           `json:"is_logout"`
	DeviceInfo    string         `json:"device_info"`
	Salt          string         `json:"salt"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	pkg.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func CreateUser(user UserBasic) interface{} {
	// 检查是否有相同用户名
	var count int64
	pkg.DB.Model(&UserBasic{}).Where("name = ?", user.Name).Count(&count)
	if count > 0 {
		fmt.Println("用户名已存在！")
		return pkg.JSONResponse{
			Code: pkg.CodeErr,
			JSON: struct {
				Message   string `json:"message"`
				MessageZh string `json:"message_zh"`
			}{"用户名已存在！", "用户名已存在！"},
		}
	}
	// 创建用户
	pkg.DB.Create(&user)
	return pkg.SuccessResponse("创建成功！")
}

func DeleteUser(user UserBasic) interface{} {
	// 检查是否有相同用户名
	var count int64
	pkg.DB.Model(&UserBasic{}).Where("id = ?", user.ID).Count(&count)
	if count == 0 {
		fmt.Println("用户不存在！")
		return pkg.JSONResponse{
			Code: pkg.CodeErr,
			JSON: struct {
				Message   string `json:"message"`
				MessageZh string `json:"message_zh"`
			}{"用户不存在！", "用户不存在！"},
		}
	}
	// 删除用户
	pkg.DB.Delete(&user)
	return pkg.SuccessResponse("删除成功！")
}

func UpdateUser(user UserBasic) interface{} {
	// 检查是否有相同用户名
	var count int64
	pkg.DB.Model(&UserBasic{}).Where("id = ?", user.ID).Count(&count)
	if count == 0 {
		fmt.Println(user.ID)
		fmt.Println("用户不存在！")
		return pkg.JSONResponse{
			Code: pkg.CodeErr,
			JSON: struct {
				Message   string `json:"message"`
				MessageZh string `json:"message_zh"`
			}{"用户不存在！", "用户不存在！"},
		}
	}
	// 更新用户
	pkg.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord})
	return pkg.SuccessResponse("更新成功！")
}

// 根据用户名查找用户 (返回 UserBasic)
func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	pkg.DB.Where("name = ?", name).First(&user)
	return user
}

// 根据手机号查找用户 (返回 *gorm.DB)
func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return pkg.DB.Where("Phone = ?", phone).First(&user)
}

// 根据邮箱查找用户 (返回 *gorm.DB)
func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return pkg.DB.Where("email = ?", email).First(&user)
}

// 根据用户名和密码查找用户 (返回 UserBasic)
func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}
	pkg.DB.Where("name = ? and pass_word=?", name,
		password).First(&user)
	temp := pkg.GenerateLongHash()
	pkg.DB.Model(&user).Where("id = ?", user.ID).Update("identity",
		temp) //更新token
	return user
}

// 判断token是否有效
func CheckToken(token string) bool {
	user := UserBasic{}
	pkg.DB.Where("identity = ?", token).First(&user)
	if user.Name == "" {
		return false
	}
	return true
}
