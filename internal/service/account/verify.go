package account

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

// ROLE 角色定义
type ROLE = int

const (
	STUDENT ROLE = 1 // 学生
	SCHOOL  ROLE = 2 // 学校
	COMPANY ROLE = 3 // 企业
	AUTH    ROLE = 4 // 认证机构
)

func verifyRole(role int) error {
	switch role {
	case STUDENT, SCHOOL, COMPANY, AUTH:
		return nil
	default:
		return status.Error(codes.InvalidArgument, "无效的角色")
	}
}

func verifyNick(nick string) error {
	if nick == "" {
		return status.Error(codes.InvalidArgument, "用户名不能为空")
	}
	if len(nick) <= 4 {
		return status.Error(codes.InvalidArgument, "用户名长度需大于4")
	}
	return nil
}

func verifyPhone(phone string) error {
	if phone == "" {
		return status.Error(codes.InvalidArgument, "手机号不能为空")
	}
	return nil
}

func verifyEmail(email string) error {
	if email == "" {
		return status.Error(codes.InvalidArgument, "邮箱不能为空")
	}
	return nil
}

func verifyPassword(password string) error {
	if password == "" {
		return status.Error(codes.InvalidArgument, "密码不能为空")
	}
	if len(password) <= 6 {
		return status.Error(codes.InvalidArgument, "密码长度需大于6")
	}
	return nil
}

func verifyPlatform(platform string) error {
	if platform == "" {
		return status.Error(codes.InvalidArgument, "未指名来源平台")
	}

	switch strings.ToLower(platform) {
	case "iphone", "android", "web", "unknown":
		return nil
	default:
		return status.Errorf(codes.InvalidArgument, "无法识别的平台: %s", platform)
	}
}
