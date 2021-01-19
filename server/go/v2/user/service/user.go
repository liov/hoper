package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	"github.com/liov/hoper/go/v2/protobuf/utils/response"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
	"github.com/liov/hoper/go/v2/user/middle"
	modelconst "github.com/liov/hoper/go/v2/user/model"
	templatei "github.com/liov/hoper/go/v2/utils/def/template"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/gateway"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	"github.com/liov/hoper/go/v2/utils/net/mail"
	"github.com/liov/hoper/go/v2/utils/strings"
	"github.com/liov/hoper/go/v2/utils/time"
	"github.com/liov/hoper/go/v2/utils/verification"
	"github.com/liov/hoper/go/v2/utils/verification/auth/jwt"
	"gorm.io/gorm"
)

type UserService struct {
	model.UnimplementedUserServiceServer
}

func (*UserService) Service() (string, string, []http.HandlerFunc) {
	return "用户相关", "/api/user", []http.HandlerFunc{middle.Log}
}

func (*UserService) FiberService() (string, string, []fiber.Handler) {
	return "用户相关", "/api/user", []fiber.Handler{middle.FiberLog}
}

func GetUserService() *UserService {
	if userSvc != nil {
		return userSvc
	}
	userSvc = new(UserService)
	return userSvc
}

func (u *UserService) VerifyCode(ctx context.Context, req *request.Empty) (*response.CommonRep, error) {
	device := CtxFromContext(ctx).UserDeviceInfo
	log.Debug(device)
	var rep = &response.CommonRep{}
	vcode := verification.GenerateCode()
	log.Info(vcode)
	rep.Details = vcode
	rep.Message = "字符串有问题吗啊"
	return rep, nil
}

func (*UserService) SignupVerify(ctx context.Context, req *model.SingUpVerifyReq) (*response.TinyRep, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	if req.Mail == "" && req.Phone == "" {
		return nil, errorcode.InvalidArgument.Message("请填写邮箱或手机号")
	}

	if exist, _ := userDao.ExitByEmailORPhone(nil, req.Mail, req.Phone); exist {
		if req.Mail != "" {
			return nil, errorcode.InvalidArgument.Message("邮箱已被注册")
		} else {
			return nil, errorcode.InvalidArgument.Message("手机号已被注册")
		}
	}
	vcode := verification.GenerateCode()
	log.Debug(vcode)
	key := modelconst.VerificationCodeKey + req.Mail + req.Phone
	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	if _, err := RedisConn.Do("SET", key, vcode, "EX", modelconst.VerificationCodeDuration); err != nil {
		log.Error("UserService.Verify,RedisConn.Do: ", err)
		return nil, errorcode.RedisErr.Message("新建出错")
	}
	return &response.TinyRep{Message: "验证码已发送"}, nil
}

func (*UserService) Signup(ctx context.Context, req *model.SignupReq) (*response.TinyRep, error) {
	if err := Validate(req); err != nil {
		return nil, err
	}
	if req.Mail == "" && req.Phone == "" {
		return nil, errorcode.InvalidArgument.Message("请填写邮箱或手机号")
	}

	if err := LuosimaoVerify(req.VCode); err != nil {
		return nil, err
	}
	if exist, _ := userDao.ExitByEmailORPhone(nil, req.Mail, req.Phone); exist {
		if req.Mail != "" {
			return nil, errorcode.InvalidArgument.Message("邮箱已被注册")
		} else {
			return nil, errorcode.InvalidArgument.Message("手机号已被注册")
		}
	}
	now := time.Now()
	formatNow := timei.Format(now)
	var user = &model.User{
		Name:            req.Name,
		Account:         uuid.New().String(),
		Mail:            req.Mail,
		Phone:           req.Phone,
		Gender:          req.Gender,
		AvatarURL:       modelconst.DefaultAvatar,
		LastActivatedAt: formatNow,
		Role:            model.Role_UserRoleNormal,
		CreatedAt:       formatNow,
		Status:          model.UserStatus_InActive,
	}

	user.Password = encryptPassword(req.Password)
	if err := userDao.Creat(nil, user); err != nil {
		log.Error(err)
		return nil, errorcode.RedisErr.Message("新建出错")
	}
	var rep = &response.TinyRep{Message: "新建成功,请前往邮箱激活"}

	activeUser := modelconst.ActiveTimeKey + strconv.FormatUint(user.Id, 10)
	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	curTime := time.Now().Unix()

	if _, err := RedisConn.Do("SET", activeUser, curTime, "EX", modelconst.ActiveDuration); err != nil {
		log.Error("UserService.Signup,RedisConn.Do: ", err)
	}

	if req.Mail != "" {
		go sendMail(model.Action_Active, curTime, user)
	}

	return rep, nil
}

// Salt 每个用户都有一个不同的盐
func salt(password string) string {
	return password[0:5]
}

// EncryptPassword 给密码加密
func encryptPassword(password string) string {
	hash := salt(password) + conf.Conf.Customize.PassSalt + password[5:]
	return fmt.Sprintf("%x", md5.Sum(stringsi.ToBytes(hash)))
}

func sendMail(action model.Action, curTime int64, user *model.User) {
	siteURL := "https://" + conf.Conf.Server.Domain
	title := action.String()
	secretStr := strconv.FormatInt(curTime, 10) + user.Mail + user.Password
	secretStr = fmt.Sprintf("%x", md5.Sum(stringsi.ToBytes(secretStr)))
	var ctiveOrRestPasswdValues = struct {
		UserName, SiteName, SiteURL, ActionURL, SecretStr string
	}{user.Name, "hoper", siteURL, "", secretStr}
	var templ string
	switch action {
	case model.Action_Active:
		ctiveOrRestPasswdValues.ActionURL = siteURL + "/user/active/" + strconv.FormatUint(user.Id, 10) + "/" + secretStr
		templ = modelconst.ActionActiveContent
	case model.Action_RestPassword:
		ctiveOrRestPasswdValues.ActionURL = siteURL + "/user/resetPassword/" + strconv.FormatUint(user.Id, 10) + "/" + secretStr
		templ = modelconst.ActionRestPasswordContent
	}
	log.Debug(ctiveOrRestPasswdValues.ActionURL)
	var buf = new(bytes.Buffer)
	err := templatei.Execute(buf, templ, &ctiveOrRestPasswdValues)
	if err != nil {
		log.Error("executing template:", err)
	}
	//content += "<p><img src=\"" + siteURL + "/images/logo.png\" style=\"height: 42px;\"/></p>"
	//fmt.Println(content)
	content := buf.String()
	addr := conf.Conf.Mail.Host + conf.Conf.Mail.Port
	m := &mail.Mail{
		FromName: ctiveOrRestPasswdValues.SiteName,
		From:     conf.Conf.Mail.From,
		Subject:  title,
		Content:  content,
		To:       []string{user.Mail},
	}
	log.Debug(content)
	err = m.SendMailTLS(addr, dao.Dao.MailAuth)
	if err != nil {
		log.Error("sendMail:", err)
	}
}

//验证密码是否正确
func checkPassword(password string, user *model.User) bool {
	if password == "" || user.Password == "" {
		return false
	}
	return encryptPassword(password) == user.Password
}

func (u *UserService) Active(ctx context.Context, req *model.ActiveReq) (*model.LoginRep, error) {
	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	redisKey := modelconst.ActiveTimeKey + strconv.FormatUint(req.Id, 10)
	emailTime, err := redis.Int64(RedisConn.Do("GET", redisKey))
	if err != nil {
		log.Error("UserService.Active,redis.Int64", err)
		return nil, errorcode.InvalidArgument.Message("无效的链接")
	}

	user, err := userDao.GetByPrimaryKey(nil, req.Id)
	if err != nil {
		return nil, errorcode.DBError
	}
	if user.Status != 0 {
		return nil, errorcode.AlreadyExists.Message("已激活")
	}
	secretStr := strconv.Itoa((int)(emailTime)) + user.Mail + user.Password

	secretStr = fmt.Sprintf("%x", md5.Sum(stringsi.ToBytes(secretStr)))

	if req.Secret != secretStr {
		return nil, errorcode.InvalidArgument.Message("无效的链接")
	}

	dao.Dao.GORMDB.Model(user).Updates(map[string]interface{}{"activated_at": time.Now(), "status": 1})
	RedisConn.Do("DEL", redisKey)
	return u.login(ctx, user)
}

func (u *UserService) Edit(ctx context.Context, req *model.EditReq) (*response.TinyRep, error) {
	defer timei.TimeCost(time.Now())
	c := CtxFromContext(ctx)
	user, err := c.GetAuthInfo()
	if err != nil || user.Id != req.Id {
		return nil, err
	}
	device := c.UserDeviceInfo

	originalIds, err := userDao.ResumesIds(nil, user.Id)
	if err != nil {
		return nil, errorcode.DBError.Message("更新失败")
	}
	var resumes []*model.Resume
	resumes = append(req.Details.EduExps, req.Details.WorkExps...)

	tx := dao.Dao.GORMDB.Begin()
	err = userDao.SaveResumes(tx, req.Id, resumes, originalIds, device)
	if err != nil {
		tx.Rollback()
		return nil, errorcode.DBError.Message("更新失败")
	}
	var nUser model.User
	tx.First(&nUser, user.Id)
	err = tx.Model(&user).Updates(nUser).Error
	if err != nil {
		tx.Rollback()
		return nil, errorcode.DBError.Message("更新失败")
	}
	tx.Commit()
	return &response.TinyRep{Message: "修改成功"}, nil
}

func (u *UserService) Login(ctx context.Context, req *model.LoginReq) (*model.LoginRep, error) {

	if verifyErr := verification.LuosimaoVerify(conf.Conf.Customize.LuosimaoVerifyURL, conf.Conf.Customize.LuosimaoAPIKey, req.VCode); verifyErr != nil {
		return nil, errorcode.InvalidArgument.Message(verifyErr.Error())
	}

	if req.Input == "" {
		return nil, errorcode.InvalidArgument.Message("账号错误")
	}
	var sql string

	switch verification.PhoneOrMail(req.Input) {
	case verification.Mail:
		sql = "mail = ?"
	case verification.Phone:
		sql = "phone = ?"
	default:
		sql = "account = ?"
	}

	var user model.User
	if err := dao.Dao.GORMDB.Where(sql, req.Input).Find(&user).Error; err != nil {
		return nil, errorcode.DBError.Message("账号不存在")
	}

	if !checkPassword(req.Password, &user) {
		return nil, errorcode.InvalidArgument.Message("密码错误")
	}
	if user.Status == model.UserStatus_InActive {
		//没看懂
		//encodedEmail := base64.StdEncoding.EncodeToString(stringsi.ToBytes(user.Mail))
		activeUser := modelconst.ActiveTimeKey + strconv.FormatUint(user.Id, 10)
		RedisConn := dao.Dao.Redis.Get()
		defer RedisConn.Close()

		curTime := time.Now().Unix()
		if _, err := RedisConn.Do("SET", activeUser, curTime, "EX", modelconst.ActiveDuration); err != nil {
			log.Error("UserService.Signup,RedisConn.Do: ", err)
			return nil, errorcode.RedisErr
		}
		go sendMail(model.Action_Active, curTime, &user)
		return nil, model.UserErr_NoActive.Message("账号未激活,请进入邮箱点击激活")
	}

	return u.login(ctx, &user)
}

func (*UserService) login(ctx context.Context, user *model.User) (*model.LoginRep, error) {
	now := time.Now()
	nowStamp := now.Unix()
	userInfo := &model.UserAuthInfo{
		Id:           user.Id,
		LastActiveAt: nowStamp,
		Status:       user.Status,
		Role:         user.Role,
		LoginAt:      nowStamp,
	}
	claims := &jwt.Claims{
		UserId:         userInfo.Id,
		StandardClaims: jwt.NewStandardClaims(conf.Conf.Customize.TokenMaxAge, "hoper"),
	}

	tokenString, err := jwt.GenerateToken(claims, stringsi.ToBytes(conf.Conf.Customize.TokenSecret))
	if err != nil {
		return nil, errorcode.Internal
	}

	dao.Dao.GORMDB.Model(&user).UpdateColumn("last_activated_at", now)
	conn := dao.NewUserRedis()
	defer conn.Close()
	if err := conn.EfficientUserHashToRedis(userInfo); err != nil {
		return nil, errorcode.RedisErr
	}
	resp := &model.LoginRep{}
	resp.Details = &model.LoginRep_LoginDetails{Token: tokenString, User: &model.UserBaseInfo{
		Id:     user.Id,
		Score:  user.Score,
		Gender: user.Gender,
	}}
	resp.Message = "登录成功"

	cookie := (&http.Cookie{
		Name:  "token",
		Value: tokenString,
		Path:  "/",
		//Domain:   "hoper.xyz",
		Expires:  time.Now().Add(time.Duration(conf.Conf.Customize.TokenMaxAge) * time.Second),
		MaxAge:   int(time.Duration(conf.Conf.Customize.TokenMaxAge) * time.Second),
		Secure:   false,
		HttpOnly: true,
	}).String()
	gateway.GrpcSetCookie(ctx, cookie)
	resp.Cookie = cookie
	return resp, nil
}

func (u *UserService) Logout(ctx context.Context, req *request.Empty) (*model.LogoutRep, error) {
	c := CtxFromContext(ctx)
	user, err := c.GetAuthInfo()
	if err != nil {
		return nil, err
	}
	dao.Dao.GORMDB.Model(&model.UserAuthInfo{Id: user.Id}).UpdateColumn("last_activated_at", time.Now())

	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	if _, err := RedisConn.Do("DEL", modelconst.LoginUserKey+strconv.FormatUint(user.Id, 10)); err != nil {
		log.Error(err)
		return nil, errorcode.RedisErr
	}
	cookie := (&http.Cookie{
		Name:  "token",
		Value: "del",
		Path:  "/",
		//Domain:   "hoper.xyz",
		Expires:  time.Now().Add(-1),
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
	}).String()
	gateway.GrpcSetCookie(ctx, cookie)
	return &model.LogoutRep{Message: "已注销", Cookie: cookie}, nil
}

func (u *UserService) AuthInfo(ctx context.Context, req *request.Empty) (*model.UserAuthInfo, error) {
	return CtxFromContext(ctx).GetAuthInfo()
}

func (u *UserService) GetUser(ctx context.Context, req *model.GetReq) (*model.GetRep, error) {
	/*	_, err := u.GetAuthInfo(ctx)
		if err != nil {
			return &model.GetRep{Details: &model.User{Id: req.Id}}, nil
		}*/
	var user1 model.User
	if err := dao.Dao.GORMDB.Find(&user1, req.Id).Error; err != nil {
		return nil, errorcode.DBError.Message("账号不存在")
	}
	return &model.GetRep{Details: &user1}, nil
}

func (u *UserService) ForgetPassword(ctx context.Context, req *model.LoginReq) (*response.TinyRep, error) {
	if verifyErr := verification.LuosimaoVerify(conf.Conf.Customize.LuosimaoVerifyURL, conf.Conf.Customize.LuosimaoAPIKey, req.VCode); verifyErr != nil {
		return nil, errorcode.InvalidArgument.Warp(verifyErr)
	}

	if req.Input == "" {
		return nil, errorcode.InvalidArgument.Message("账号错误")
	}
	user, err := userDao.GetByEmailORPhone(nil, req.Input, req.Input, "id", "name", "password")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if verification.PhoneOrMail(req.Input) != verification.Phone {
				return nil, errorcode.InvalidArgument.Message("邮箱不存在")
			} else {
				return nil, errorcode.InvalidArgument.Message("手机号不存在")
			}
		}
		log.Error(err)
		return nil, errorcode.DBError
	}
	restPassword := modelconst.ResetTimeKey + strconv.FormatUint(user.Id, 10)
	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	curTime := time.Now().Unix()
	if _, err := RedisConn.Do("SET", restPassword, curTime, "EX", modelconst.ResetDuration); err != nil {
		log.Error("redis set failed:", err)
		return nil, errorcode.RedisErr
	}

	go sendMail(model.Action_RestPassword, curTime, user)

	return &response.TinyRep{}, nil
}

func (u *UserService) ResetPassword(ctx context.Context, req *model.ResetPasswordReq) (*response.TinyRep, error) {
	RedisConn := dao.Dao.Redis.Get()
	defer RedisConn.Close()

	redisKey := modelconst.ResetTimeKey + strconv.FormatUint(req.Id, 10)
	emailTime, err := redis.Int64(RedisConn.Do("GET", redisKey))
	if err != nil {
		log.Error(model.UserserviceServicedesc.ServiceName, "ResetPassword,redis.Int64", err)
		return nil, errorcode.InvalidArgument.Message("无效的链接")
	}

	user, err := userDao.GetByPrimaryKey(nil, req.Id)
	if err != nil {
		return nil, errorcode.DBError
	}
	if user.Status != 1 {
		return nil, errorcode.FailedPrecondition.Message("无效账号")
	}
	secretStr := strconv.Itoa((int)(emailTime)) + user.Mail + user.Password

	secretStr = fmt.Sprintf("%x", md5.Sum(stringsi.ToBytes(secretStr)))

	if req.Secret != secretStr {
		return nil, errorcode.InvalidArgument.Message("无效的链接")
	}

	if err := dao.Dao.GORMDB.Model(user).Update("password", req.Password).Error; err != nil {
		log.Error("UserService.ResetPassword,DB.Update", err)
		return nil, errorcode.DBError
	}
	RedisConn.Do("DEL", redisKey)
	return nil, nil
}

func (*UserService) ActionLogList(ctx context.Context, req *model.ActionLogListReq) (*model.ActionLogListRep, error) {
	rep := &model.ActionLogListRep{}
	var logs []*model.UserActionLog
	err := dao.Dao.GORMDB.Offset(0).Limit(10).Find(&logs).Error
	if err != nil {
		return nil, errorcode.DBError.Log(err)
	}
	rep.Details = logs
	return rep, nil
}

func (*UserService) GetTest(ctx context.Context, req *model.GetReq) (*model.GetRep, error) {
	return &model.GetRep{Code: uint32(req.Id), Message: "测试"}, nil
}

func (*UserService) Add(ctx context.Context, req *model.SignupReq) (*response.TinyRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() interface{} {
		return pick.Method(http.MethodGet).
			Title("用户注册").
			Version(2).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试")
	})

	return &response.TinyRep{Message: req.Name}, nil
}

func (*UserService) AddV2(ctx context.Context, req *response.TinyRep) (*response.TinyRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.FiberApi(func() interface{} {
		return pick.Method(http.MethodGet).
			Title("用户注册").
			Version(2).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试")
	})

	return req, nil
}
