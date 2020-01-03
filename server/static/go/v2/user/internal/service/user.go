package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"

	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/dao"
	modelconst "github.com/liov/hoper/go/v2/user/internal/model"
	"github.com/liov/hoper/go/v2/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/http/token"
	"github.com/liov/hoper/go/v2/utils/json"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/mail"
	"github.com/liov/hoper/go/v2/utils/strings2"
	"github.com/liov/hoper/go/v2/utils/time2"
	"github.com/liov/hoper/go/v2/utils/valid"
	"github.com/liov/hoper/go/v2/utils/verification/code"
	"github.com/liov/hoper/go/v2/utils/verification/luosimao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserService struct{}

func NewUserService(server model.UserServiceServer) *UserService {
	return &UserService{}
}

func (*UserService) VerifyCode(ctx context.Context, req *utils.Empty) (*model.VerifyRep, error) {
	var rep = &model.VerifyRep{}
	vcode := code.Generate()
	log.Info(vcode)
	rep.Details = vcode
	rep.Message = "字符串有问题吗啊"
	return rep, nil
}

func (*UserService) SignupVerify(ctx context.Context, req *model.SingUpVerifyReq) (*model.SingUpVerifyRep, error) {
	err := valid.Validate.Struct(req)
	if err != nil {
		return nil, errorcode.InvalidParams.ErrWithMessage(valid.Trans(err))
	}

	if req.Mail == "" && req.Phone == "" {
		return nil, errorcode.InvalidParams.ErrWithMessage("请填写邮箱或手机号")
	}

	if exist, _ := userDao.ExitByEmailORPhone(req.Mail, req.Phone); exist {
		if req.Mail != "" {
			return nil, errorcode.InvalidParams.ErrWithMessage("邮箱已被注册")
		} else {
			return nil, errorcode.InvalidParams.ErrWithMessage("手机号已被注册")
		}
	}
	vcode := code.Generate()
	log.Info(vcode)
	key := modelconst.VerificationCodeKey + req.Mail + req.Phone
	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	if _, err := RedisConn.Do("SET", key, vcode, "EX", modelconst.VerificationCodeDuration); err != nil {
		log.Error("UserService.Verify,RedisConn.Do: ", err)
		return nil, errorcode.ERROR.ErrWithMessage("新建出错")
	}
	var rep = &model.SingUpVerifyRep{}
	rep.Details = vcode
	rep.Message = "字符串有问题吗啊"
	return rep, err
}

func (*UserService) Signup(ctx context.Context, req *model.SignupReq) (*model.SignupRep, error) {

	err := valid.Validate.Struct(req)
	if err != nil {
		return nil, errorcode.InvalidParams.ErrWithMessage(valid.Trans(err))
	}

	if req.Mail == "" && req.Phone == "" {
		return nil, errorcode.InvalidParams.ErrWithMessage("请填写邮箱或手机号")
	}

	if exist, _ := userDao.ExitByEmailORPhone(req.Mail, req.Phone); exist {
		if req.Mail != "" {
			return nil, errorcode.InvalidParams.ErrWithMessage("邮箱已被注册")
		} else {
			return nil, errorcode.InvalidParams.ErrWithMessage("手机号已被注册")
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
		log.Error(err)
		return nil, errorcode.ERROR.ErrWithMessage("新建出错")
	}
	var rep = &model.SignupRep{Message: "新建成功", Details: user}

	activeUser := modelconst.ActiveTimeKey + strconv.FormatUint(user.Id, 10)
	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	curTime := time.Now().Unix()

	if _, err := RedisConn.Do("SET", activeUser, curTime, "EX", modelconst.ActiveDuration); err != nil {
		log.Error("UserService.Signup,RedisConn.Do: ", err)
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

func sendMail(action string, title string, curTime int64, user *model.User) {
	siteName := "hoper"
	siteURL := "https://" + config.Conf.Server.Domain
	secretStr := strconv.FormatInt(curTime, 10) + user.Mail + user.Password
	secretStr = fmt.Sprintf("%x", md5.Sum(strings2.ToBytes(secretStr)))
	actionURL := siteURL + "/user" + action + "/"

	actionURL = actionURL + strconv.FormatUint(user.Id, 10) + "/" + secretStr
	log.Debug(actionURL)

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

	addr := config.Conf.Mail.Host + config.Conf.Mail.Port
	m := mail.Mail{
		FromName: siteName,
		From:     config.Conf.Mail.From,
		Subject:  title,
		Content:  content,
		To:       []string{user.Mail},
	}
	err := mail.SendMailTLS(addr, dao.Dao.MailAuth, &m)
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

func (*UserService) Active(ctx context.Context, req *model.ActiveReq) (*model.ActiveRep, error) {
	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	emailTime, err := redis.Int64(RedisConn.Do("GET", modelconst.ActiveTimeKey+strconv.FormatUint(req.Id, 10)))
	if err != nil {
		log.Error("UserService.Active,redis.Int64", err)
		return nil, errorcode.ErrorWithMessage(errorcode.InvalidParams, "无效的链接")
	}

	user, err := userDao.GetByPrimaryKey(req.Id)
	if err != nil {
		return nil, errorcode.Error(errorcode.ERROR)
	}
	if user.Status != 0 {
		return nil, errorcode.ERROR.ErrWithMessage("已激活")
	}
	secretStr := strconv.Itoa((int)(emailTime)) + user.Mail + user.Password

	secretStr = fmt.Sprintf("%x", md5.Sum(strings2.ToBytes(secretStr)))

	if req.Secret != secretStr {
		return nil, errorcode.ErrorWithMessage(errorcode.InvalidParams, "无效的链接")
	}
	user.Status = 1
	dao.Dao.GORMDB.Model(user).Update("activated_at", time.Now(), "status", 1)
	var rep = &model.ActiveRep{}
	rep.Message = "激活成功"
	return rep, nil
}

func (*UserService) Edit(context.Context, *model.EditReq) (*model.EditRep, error) {
	panic("implement me")
}

func (*UserService) Login(ctx context.Context, req *model.LoginReq) (*model.LoginRep, error) {
	resp := &model.LoginRep{}
	verifyErr := luosimao.LuosimaoVerify(config.Conf.Server.LuosimaoVerifyURL, config.Conf.Server.LuosimaoAPIKey, req.Luosimao)

	if verifyErr != nil {
		return nil, errorcode.InvalidParams.ErrWithMessage(verifyErr.Error())
	}

	if req.Input == "" {
		return nil, errorcode.InvalidParams.ErrWithMessage("账号错误")
	}
	var sql string
	emailMatch, _ := regexp.MatchString(`^([a-zA-Z0-9]+[_.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_.]?)*[a-zA-Z0-9]+.[a-zA-Z]{2,3}$`, req.Input)
	if emailMatch {
		sql = "mail = ?"
	} else {
		phoneMatch, _ := regexp.MatchString(`^1[0-9]{10}$`, req.Input)
		if phoneMatch {
			sql = "phone = ?"
		}
	}

	var user model.User
	if err := dao.Dao.GORMDB.Where(sql, req.Input).Find(&user).Error; err != nil {
		return nil, errorcode.ERROR.ErrWithMessage("账号不存在")
	}

	if !checkPassword(req.Password, &user) {
		return nil, errorcode.InvalidParams.ErrWithMessage("密码错误")
	}
	if user.Status == modelconst.UserStatusInActive {
		//没看懂
		//encodedEmail := base64.StdEncoding.EncodeToString(strings2.ToBytes(user.Mail))
		resp.Message = "账号未激活,请进去邮箱点击激活"
		go sendMail("/active", "账号激活", time.Now().Unix(), &user)
		return resp, nil
	}

	tokenString, err := token.GenerateToken(user.Id, config.Conf.Server.TokenMaxAge, config.Conf.Server.TokenSecret)
	if err != nil {
		return nil, errorcode.ERROR
	}

	dao.Dao.GORMDB.Model(&user).UpdateColumn("last_activated_at", time.Now())

	if err := UserHashToRedis(&model.UserMainInfo{
		Id:     user.Id,
		Score:  user.Score,
		Status: user.Status,
		Role:   user.Role,
	}); err != nil {
		return nil, errorcode.ERROR
	}

	resp.Details = &model.LoginRep_LoginDetails{Token: tokenString, User: &user}
	resp.Message = "登录成功"

	return resp, nil
}

func (u *UserService) Logout(ctx context.Context, req *model.LogoutReq) (*model.LogoutRep, error) {
	user, err := u.Auth(ctx)
	if err != nil {
		return nil, err
	}
	dao.Dao.GORMDB.Model(&model.UserMainInfo{Id: user.Id}).UpdateColumn("last_activated_at", time.Now())

	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	if _, err := RedisConn.Do("DEL", modelconst.LoginUserKey+strconv.FormatUint(user.Id, 10)); err != nil {
	}
	return &model.LogoutRep{Msg: "已注销"}, nil
}

func (u *UserService) AuthInfo(ctx context.Context, req *utils.Empty) (*model.UserMainInfo, error) {
	return u.Auth(ctx)
}

func (*UserService) Auth(ctx context.Context) (*model.UserMainInfo, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokens := md.Get("authorization")
	authErr := status.Error(codes.Code(errorcode.Auth), errorcode.Auth.Error())
	if len(tokens) == 0 || tokens[0] == "" {
		return nil, authErr
	}
	claims, err := token.ParseToken(tokens[0], config.Conf.Server.TokenSecret)
	if err != nil {
		return nil, authErr
	}
	user, err := UserHashFromRedis(claims.UserID)
	if err != nil {
		return nil, authErr
	}

	return user, nil
}

func (u *UserService) GetUser(ctx context.Context, req *model.GetReq) (*model.GetRep, error) {
	user, err := u.Auth(ctx)
	if err != nil {
		return &model.GetRep{Details: &model.User{Id: req.Id}}, nil
	}
	return &model.GetRep{Details: &model.User{Id: user.Id, Role: user.Role}}, nil
}

// UserToRedis 将用户信息存到redis
func UserToRedis(user *model.UserMainInfo) error {
	UserString, err := json.Json.MarshalToString(user)
	if err != nil {
		return err
	}

	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)

	conn := dao.Dao.Redis.Get()
	defer conn.Close()
	conn.Send("SELECT", modelconst.UserIndex)
	if _, redisErr := conn.Do("SET", loginUserKey, UserString, "EX", config.Conf.Server.TokenMaxAge); redisErr != nil {
		return redisErr
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func UserFromRedis(userID uint64) (*model.UserMainInfo, error) {
	loginUser := modelconst.LoginUserKey + strconv.FormatUint(userID, 10)

	conn := dao.Dao.Redis.Get()
	defer conn.Close()
	conn.Send("SELECT", modelconst.UserIndex)
	userString, err := redis.String(conn.Do("GET", loginUser))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var user model.UserMainInfo
	err = json.Json.UnmarshalFromString(userString, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UserLastActiveTime(userID uint64) error {
	conn := dao.Dao.Redis.Get()
	defer conn.Close()

	err := conn.Send("SELECT", modelconst.CronIndex)
	_, err = conn.Do("ZADD", modelconst.LoginUserKey+"ActiveTime",
		time.Now().Unix(), strconv.FormatUint(userID, 10))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func EditRedisUser(user *model.UserMainInfo) error {
	UserString, err := json.Json.MarshalToString(user)
	if err != nil {
		return err
	}
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)

	conn := dao.Dao.Redis.Get()
	defer conn.Close()
	conn.Send("SELECT", modelconst.UserIndex)
	if _, redisErr := conn.Do("SET", loginUserKey, UserString); redisErr != nil {
		return redisErr
	}
	return nil
}

// UserToRedis 将用户信息存到redis
func UserHashToRedis(user *model.UserMainInfo) error {
	var redisArgs []interface{}
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	redisArgs = append(redisArgs, loginUserKey)

	uValue := reflect.ValueOf(user).Elem()
	uType := uValue.Type()
	for i := 0; i < uValue.NumField(); i++ {
		redisArgs = append(redisArgs, uType.Field(i).Name, uValue.Field(i).Interface())
	}

	conn := dao.Dao.Redis.Get()
	defer conn.Close()
	conn.Send("SELECT", modelconst.UserIndex)
	conn.Send("HMSET", redisArgs...)
	if _, redisErr := conn.Do("EXPIRE", loginUserKey, config.Conf.Server.TokenMaxAge); redisErr != nil {
		return redisErr
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func UserHashFromRedis(userID uint64) (*model.UserMainInfo, error) {
	loginUser := modelconst.LoginUserKey + strconv.FormatUint(userID, 10)

	conn := dao.Dao.Redis.Get()
	defer conn.Close()
	conn.Send("SELECT", modelconst.UserIndex)
	userArgs, err := redis.Strings(conn.Do("HGETALL", loginUser))
	log.Debug(userArgs)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var user model.UserMainInfo
	uValue := reflect.ValueOf(&user).Elem()
	for i := 0; i < uValue.NumField(); i += 2 {
		fieldValue := uValue.FieldByName(userArgs[i])
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, _ := strconv.ParseInt(userArgs[i+1], 10, 64)
			fieldValue.SetInt(v)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, _ := strconv.ParseUint(userArgs[i+1], 10, 64)
			fieldValue.SetUint(v)
		case reflect.String:
			fieldValue.SetString(userArgs[i+1])
		case reflect.Float32, reflect.Float64:
			v, _ := strconv.ParseFloat(userArgs[i+1], 64)
			fieldValue.SetFloat(v)
		case reflect.Bool:
			v, _ := strconv.ParseBool(userArgs[i+1])
			fieldValue.SetBool(v)
		}

	}
	return &user, nil
}
