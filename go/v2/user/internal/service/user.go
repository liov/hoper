package service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/liov/hoper/go/v2/protobuf/response"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/dao"
	modelconst "github.com/liov/hoper/go/v2/user/internal/model"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/mail"
	"github.com/liov/hoper/go/v2/utils/protobuf"
	"github.com/liov/hoper/go/v2/utils/strings2"
	"github.com/liov/hoper/go/v2/utils/time2"
	"github.com/liov/hoper/go/v2/utils/valid"
	"github.com/liov/hoper/go/v2/utils/verificationCode"
)

type UserService struct{}

func (*UserService) Verify(ctx context.Context, req *model.VerifyReq) (*response.AnyReply, error) {
	var rep = &response.AnyReply{Code: 10000}
	err := valid.Validate.Struct(req)
	if err != nil {
		rep.Message = valid.Trans(err)
		return rep, nil
	}

	if req.Mail == "" && req.Phone == "" {
		rep.Message = "请填写邮箱或手机号"
		return rep, err
	}

	if exist, err := userDao.ExitByEmailORPhone(req.Mail, req.Phone); exist {
		if req.Mail != "" {
			rep.Message = "邮箱已被注册"
			return rep, err
		} else {
			rep.Message = "手机号已被注册"
			return rep, err
		}
	}
	vcode:=verificationCode.Generate()
	log.Info(vcode)
	key := modelconst.VerificationCodeKey + req.Mail + req.Phone
	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	if _, err := RedisConn.Do("SET", key, vcode, "EX", modelconst.VerificationCodeDuration); err != nil {
		log.Error("UserService.Verify,RedisConn.Do: ", err)
		rep.Message = "新建出错"
		return rep, nil
	}
	rep.Details = []byte("\""+vcode+"\"")
	rep.Message = "字符串有问题吗啊"
	return rep, err
}

func (*UserService) Signup(ctx context.Context, req *model.SignupReq) (*response.AnyReply, error) {
	var rep = &response.AnyReply{Code: 10000}
	err := valid.Validate.Struct(req)
	if err != nil {
		rep.Message = valid.Trans(err)
		return rep, nil
	}

	if req.Mail == "" && req.Phone == "" {
		rep.Message = "请填写邮箱或手机号"
		return rep, err
	}

	if exist, err := userDao.ExitByEmailORPhone(req.Mail, req.Phone); exist {
		if req.Mail != "" {
			rep.Message = "邮箱已被注册"
			return rep, err
		} else {
			rep.Message = "手机号已被注册"
			return rep, err
		}
	}
	var user = &model.User{}
	user.Mail = req.Mail
	user.Gender = modelconst.UserSexNil
	user.CreatedAt = time2.Format(time.Now())
	user.LastActiveAt = user.CreatedAt
	user.Role = modelconst.UserRoleNormal
	user.Status = modelconst.UserStatusInActive
	user.AvatarURL = modelconst.DefaultAvatar
	user.Password = encryptPassword(req.Password)
	if err = userDao.Creat(user); err != nil {
		rep.Message = "新建出错"
		return rep, err
	}
	data, _ := json.Marshal(user)
	rep.Details = data
	rep.Message = "新建成功"

	activeUser := modelconst.ActiveTimeKey + strconv.FormatUint(user.Id, 10)
	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	curTime := time.Now().Unix()

	if _, err := RedisConn.Do("SET", activeUser, curTime, "EX", modelconst.ActiveDuration); err != nil {
		log.Error("UserService.Signup,RedisConn.Do: ", err)
		rep.Message = "新建出错"
		return rep, nil
	}
	go func() {
		if req.Mail != "" {
			sendMail("/active", "账号激活", curTime, user)
		}
	}()
	return rep, nil
}

// Salt 每个用户都有一个不同的盐
func salt(password string) string {
	return password[0:5]
}

// EncryptPassword 给密码加密
func encryptPassword(password string) string {
	hash := salt(password) + config.Conf.Server.PassSalt + password[5:]
	return fmt.Sprintf("%x", md5.Sum(strings2.ToBytes(hash)))
}

/*
	letter:=`
From:{{.From}}
To: {{.To}}
Subject: {{.Title}}
Content-Type: text/html; charset=UTF-8

`
*/
func sendMail(action string, title string, curTime int64, user *model.User) {
	siteName := "hoper"
	siteURL := "https://" + config.Conf.Server.Domain
	secretStr := strconv.FormatInt(curTime, 10) + user.Mail + user.Password
	secretStr = fmt.Sprintf("%x", md5.Sum(strings2.ToBytes(secretStr)))
	actionURL := siteURL + "/user" + action + "/"

	actionURL = actionURL + strconv.FormatUint(user.Id, 10) + "/" + secretStr
	log.Info(actionURL)

	content := "<p><b>亲爱的" + user.Name + ":</b></p>" +
		"<p>我们收到您在 " + siteName + " 的注册信息, 请点击下面的链接, 或粘贴到浏览器地址栏来激活帐号.</p>" +
		"<a href=\"" + actionURL + "\">" + actionURL + "</a>" +
		"<p>如果您没有在 " + siteName + " 填写过注册信息, 说明有人滥用了您的邮箱, 请删除此邮件, 我们对给您造成的打扰感到抱歉.</p>" +
		"<p>" + siteName + " 谨上.</p>"

	if action == "/reset" {
		content = "<p><b>亲爱的" + user.Name + ":</b></p>" +
			"<p>你的密码重设要求已经得到验证。请点击以下链接, 或粘贴到浏览器地址栏来设置新的密码: </p>" +
			"<a href=\"" + actionURL + "\">" + actionURL + "</a>" +
			"<p>感谢你对" + siteName + "的支持，希望你在" + siteName + "的体验有益且愉快。</p>" +
			"<p>(这是一封自动产生的email，请勿回复。)</p>"
	}
	//content += "<p><img src=\"" + siteURL + "/images/logo.png\" style=\"height: 42px;\"/></p>"
	//fmt.Println(content)
	headers := make(map[string]string)
	headers["From"] = siteName + "<" + config.Conf.Mail.User + ">"
	headers["To"] = user.Mail
	headers["Subject"] = title
	headers["Content-Type"] = "text/html; charset=UTF-8"
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + content

	addr := config.Conf.Mail.Host + config.Conf.Mail.Port
	err := mail.SendMailTLS(addr, dao.Dao.MailAuth, config.Conf.Mail.User, []string{user.Mail}, []byte(message))
	if err != nil {
		log.Error("sendMail", err)
	}
}

//验证密码是否正确
func checkPassword(password string, user *model.User) bool {
	if password == "" || user.Password == "" {
		return false
	}
	return encryptPassword(password) == user.Password
}

func (*UserService) Active(ctx context.Context, req *model.ActiveReq) (*response.AnyReply, error) {
	var rep = &response.AnyReply{Code: 10000, Message: "无效的链接"}
	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	emailTime, err := redis.Int64(RedisConn.Do("GET", modelconst.ActiveTimeKey+strconv.FormatUint(req.Id, 10)))
	if err != nil {
		log.Error("UserService.Active,redis.Int64", err)
		return rep, err
	}

	user, err := userDao.GetByPrimaryKey(req.Id)
	if err != nil {
		return rep, err
	}

	secretStr := strconv.Itoa((int)(emailTime)) + user.Mail + user.Password

	secretStr = fmt.Sprintf("%x", md5.Sum(strings2.ToBytes(secretStr)))

	if req.Secret != secretStr {
		return rep, nil
	}
	rep.Message = "激活成功"
	return rep, nil
}

func (*UserService) Edit(context.Context, *model.EditReq) (*model.EditRep, error) {
	panic("implement me")
}

func (*UserService) Login(context.Context, *model.LoginReq) (*model.LoginRep, error) {
	panic("implement me")
}

func (*UserService) Logout(context.Context, *model.LogoutReq) (*model.LogoutRep, error) {
	panic("implement me")
}

func (*UserService) GetUser(ctx context.Context,req *model.GetReq) (*response.Reply, error) {
	user,_:=protobuf.GenGogoAny(&model.User{Id:req.Id})
	return &response.Reply{Details:user}, nil
}
