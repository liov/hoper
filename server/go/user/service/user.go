package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/hopeio/context/ginctx"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/pick"
	"github.com/hopeio/protobuf/request"
	"github.com/hopeio/protobuf/response"
	timepb "github.com/hopeio/protobuf/time"
	gormi "github.com/hopeio/utils/dao/database/gorm"
	"github.com/hopeio/utils/sdk/luosimao"
	stringsi "github.com/hopeio/utils/strings"
	model "github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/api/middle"
	"github.com/liov/hoper/server/go/user/confdao"
	"github.com/liov/hoper/server/go/user/data"
	"github.com/liov/hoper/server/go/user/data/redis"
	modelconst "github.com/liov/hoper/server/go/user/model"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net/http"
	"strconv"
	"time"

	"github.com/hopeio/protobuf/errcode"
	redisi "github.com/hopeio/utils/dao/redis"
	templatei "github.com/hopeio/utils/encoding/text/template"
	"github.com/hopeio/utils/log"
	httpi "github.com/hopeio/utils/net/http"

	"github.com/hopeio/utils/net/mail"
	"github.com/hopeio/utils/validation"
	"gorm.io/gorm"
)

type UserService struct {
	model.UnimplementedUserServiceServer
}

func (u *UserService) VerifyCode(ctx context.Context, req *emptypb.Empty) (*wrapperspb.StringValue, error) {
	device := httpctx.FromContextValue(ctx).DeviceInfo
	log.Debug(device)
	var rep = &wrappers.StringValue{}
	vcode := validation.GenerateCode()
	log.Info(vcode)
	rep.Value = vcode
	return rep, nil
}

// 验证码
func (u *UserService) SendVerifyCode(ctx context.Context, req *model.SendVerifyCodeReq) (*emptypb.Empty, error) {
	return nil, nil
}

func (*UserService) SignupVerify(ctx context.Context, req *model.SingUpVerifyReq) (*wrappers.StringValue, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()

	if req.Mail == "" && req.Phone == "" {
		return nil, errcode.InvalidArgument.Message("请填写邮箱或手机号")
	}

	userDao := data.GetDBDao(ctxi, confdao.Dao.GORMDB.DB)
	input := req.Mail
	if input == "" {
		input = req.Phone
	}
	checkUser, err := userDao.GetByEmailOrPhone(input)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errcode.DBError
	}
	if err == nil {
		if checkUser.Mail == req.Mail {
			return nil, errcode.InvalidArgument.Message("邮箱已被注册")
		}
		if checkUser.Phone == req.Phone {
			return nil, errcode.InvalidArgument.Message("手机号已被注册")
		}
	}
	vcode := validation.GenerateCode()
	log.Debug(vcode)
	key := modelconst.VerificationCodeKey + req.Mail + req.Phone
	if err := confdao.Dao.Redis.SetEX(ctx, key, vcode, modelconst.VerificationCodeDuration).Err(); err != nil {
		return nil, ctxi.ErrorLog(errcode.RedisErr.Message("新建出错"), err, "SetEX")
	}
	return new(wrappers.StringValue), nil
}

func (u *UserService) Signup(ctx context.Context, req *model.SignupReq) (*wrappers.StringValue, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()

	if req.Mail == "" && req.Phone == "" {
		return nil, errcode.InvalidArgument.Message("请填写邮箱或手机号")
	}
	if req.VCode != confdao.Conf.Customize.LuosimaoSuperPW {
		if err := LuosimaoVerify(req.VCode); err != nil {
			return nil, err
		}
	}

	userDao := data.GetDBDao(ctxi, confdao.Dao.GORMDB.DB)

	checkUser, err := userDao.GetByNameOrEmailOrPhone(req.Name, req.Mail, req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ctxi.ErrorLog(errcode.DBError.Message("查询出错"), err, "userDao.GetByNameOrEmailOrPhone")
	}
	if err == nil {
		if checkUser.Name == req.Name {
			return nil, errcode.InvalidArgument.Message("用户名已被注册")
		}
		if checkUser.Mail == req.Mail {
			return nil, errcode.InvalidArgument.Message("邮箱已被注册")
		}
		if checkUser.Phone == req.Phone {
			return nil, errcode.InvalidArgument.Message("手机号已被注册")
		}
	}

	var user = &model.User{
		Name:      req.Name,
		Account:   uuid.New().String(),
		Mail:      req.Mail,
		Phone:     req.Phone,
		Gender:    req.Gender,
		AvatarUrl: modelconst.DefaultAvatar,
		Role:      model.RoleNormal,
		Status:    model.UserStatusInActive,
	}

	user.Password = encryptPassword(req.Password)
	if err := userDao.Creat(user); err != nil {
		return nil, ctxi.ErrorLog(errcode.DBError.Message("新建出错"), err, "UserService.Creat")
	}

	activeUser := modelconst.ActiveTimeKey + strconv.FormatUint(user.Id, 10)

	curTime := ctxi.RequestAt.TimeStamp

	if err := confdao.Dao.Redis.SetEX(ctx, activeUser, curTime, modelconst.ActiveDuration).Err(); err != nil {
		return nil, ctxi.ErrorLog(errcode.RedisErr, err, "UserService.Signup,SetEX")
	}

	if req.Mail != "" {
		go sendMail(ctxi, model.ActionActive, curTime, user)
	}

	return &wrappers.StringValue{Value: "注册成功，注意查收邮件"}, nil
}

// Salt 每个用户都有一个不同的盐
func salt(password string) string {
	return password[0:5]
}

// EncryptPassword 给密码加密
func encryptPassword(password string) string {
	hash := salt(password) + confdao.Conf.Customize.PassSalt + password[5:]
	return fmt.Sprintf("%x", md5.Sum(stringsi.ToBytes(hash)))
}

func sendMail(ctxi *httpctx.Context, action model.Action, curTime int64, user *model.User) {
	siteURL := "https://" + confdao.Conf.Customize.Domain
	title := action.String()
	secretStr := strconv.FormatInt(curTime, 10) + user.Mail + user.Password
	secretStr = fmt.Sprintf("%x", md5.Sum(stringsi.ToBytes(secretStr)))
	var ctiveOrRestPasswdValues = struct {
		UserName, SiteName, SiteURL, ActionURL, SecretStr string
	}{user.Name, "hoper", siteURL, "", secretStr}
	var templ string
	switch action {
	case model.ActionActive:
		ctiveOrRestPasswdValues.ActionURL = siteURL + "/#/user/active/" + strconv.FormatUint(user.Id, 10) + "/" + secretStr
		templ = modelconst.ActionActiveContent
	case model.ActionRestPassword:
		ctiveOrRestPasswdValues.ActionURL = siteURL + "/#/user/resetPassword/" + strconv.FormatUint(user.Id, 10) + "/" + secretStr
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
	addr := confdao.Conf.SendMail.Host + confdao.Conf.SendMail.Port
	m := &mail.Mail{
		FromName: ctiveOrRestPasswdValues.SiteName,
		From:     confdao.Conf.SendMail.From,
		Subject:  title,
		Content:  content,
		To:       []string{user.Mail},
	}
	log.Debug(content)
	err = m.SendMailTLS(addr, confdao.Dao.Mail.Auth)
	if err != nil {
		log.Error("sendMail:", err)
	}
}

// 验证密码是否正确
func checkPassword(password string, user *model.User) bool {
	if password == "" || user.Password == "" {
		return false
	}
	return encryptPassword(password) == user.Password
}

func (u *UserService) Active(ctx context.Context, req *model.ActiveReq) (*model.LoginRep, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()
	userDBDao := data.GetDBDao(ctxi, confdao.Dao.GORMDB.DB)

	user, err := userDBDao.GetByPrimaryKey(req.Id)
	if err != nil {
		return nil, errcode.DBError
	}
	redisKey := modelconst.ActiveTimeKey + strconv.FormatUint(req.Id, 10)
	emailTime, err := confdao.Dao.Redis.Get(ctx, redisKey).Int64()
	if err != nil {
		go sendMail(ctxi, model.ActionActive, ctxi.RequestAt.TimeStamp, user)
		return nil, ctxi.ErrorLog(errcode.InvalidArgument.Message("已过激活期限"), err, "Get")
	}
	if user.Status != model.UserStatusInActive {
		return nil, errcode.AlreadyExists.Message("已激活")
	}
	secretStr := strconv.Itoa((int)(emailTime)) + user.Mail + user.Password

	secretStr = fmt.Sprintf("%x", md5.Sum(stringsi.ToBytes(secretStr)))

	if req.Secret != secretStr {
		return nil, errcode.InvalidArgument.Message("无效的链接")
	}
	err = userDBDao.Active(user)
	if err != nil {
		return nil, errcode.DBError
	}
	confdao.Dao.Redis.Del(ctx, redisKey)
	return u.login(ctxi, user)
}

func (u *UserService) Edit(ctx context.Context, req *model.EditReq) (*emptypb.Empty, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()
	user, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	if user.Id != req.Id {
		return nil, errcode.PermissionDenied
	}
	device := ctxi.DeviceInfo

	if req.Details != nil {
		userDBDao := data.GetDBDao(ctxi, confdao.Dao.GORMDB.DB)

		originalIds, err := userDBDao.ResumesIds(user.Id)
		if err != nil {
			return nil, errcode.DBError.Message("更新失败")
		}
		var resumes []*model.Resume
		resumes = append(req.Details.EduExps, req.Details.WorkExps...)
		tx := userDBDao.Begin()
		defer tx.Rollback()
		userDBDao = data.GetDBDao(ctxi, tx)
		if len(resumes) > 0 {
			err = userDBDao.SaveResumes(req.Id, resumes, originalIds, model.ConvDeviceInfo(device))
			if err != nil {
				return nil, errcode.DBError.Message("更新失败")
			}
		}
		err = userDBDao.Update(req)
		if err != nil {
			return nil, errcode.DBError.Message("更新失败")
		}
		tx.Commit()
	}
	return new(emptypb.Empty), nil
}

func (u *UserService) Login(ctx context.Context, req *model.LoginReq) (*model.LoginRep, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()

	if req.VCode != confdao.Conf.Customize.LuosimaoSuperPW {
		if err := LuosimaoVerify(req.VCode); err != nil {
			return nil, err
		}
	}

	if req.Input == "" {
		return nil, errcode.InvalidArgument.Message("账号错误")
	}

	userDBDao := data.GetDBDao(ctxi, confdao.Dao.GORMDB.DB)
	user, err := userDBDao.UserInfoByAccount(req.Input)
	if err != nil {
		return nil, ctxi.ErrorLog(errcode.DBError.Message("账号不存在"), err, "Login")
	}

	if !checkPassword(req.Password, user) {
		return nil, errcode.InvalidArgument.Message("密码错误")
	}
	if user.Status == model.UserStatusInActive {
		//没看懂
		//encodedEmail := base64.StdEncoding.EncodeToString(stringsi.ToBytes(user.Mail))
		activeUser := modelconst.ActiveTimeKey + strconv.FormatUint(user.Id, 10)

		curTime := time.Now().Unix()
		if err := confdao.Dao.Redis.SetEX(ctx, activeUser, curTime, modelconst.ActiveDuration).Err(); err != nil {
			return nil, ctxi.ErrorLog(errcode.RedisErr, err, "SetEX")
		}
		go sendMail(ctxi, model.ActionActive, curTime, user)
		return nil, model.UserErrNoActive.Message("账号未激活,请进入邮箱点击激活")
	}

	return u.login(ctxi, user)
}

func (*UserService) login(ctxi *httpctx.Context, user *model.User) (*model.LoginRep, error) {
	authorization := Authorization{AuthInfo: &model.AuthInfo{
		Id:     user.Id,
		Name:   user.Name,
		Role:   user.Role,
		Status: user.Status,
	}}

	ctxi.AuthInfo = authorization.AuthInfo
	authorization.IssuedAt = &jwt.NumericDate{Time: ctxi.Time}
	authorization.ExpiresAt = &jwt.NumericDate{Time: ctxi.Time.Add(confdao.Conf.Customize.TokenMaxAge)}

	tokenString, err := authorization.GenerateToken(confdao.Conf.Customize.TokenSecretBytes)
	if err != nil {
		return nil, errcode.Internal
	}
	db := gormi.NewTraceDB(confdao.Dao.GORMDB.DB, ctxi.BaseContext(), ctxi.TraceID)

	db.Table(modelconst.UserExtTableName).Where(`id = ?`, user.Id).
		UpdateColumn("last_activated_at", ctxi.RequestAt.TimeString)
	userRedisDao := redis.GetUserDao(ctxi, confdao.Dao.Redis.Client)
	if err := userRedisDao.EfficientUserHashToRedis(); err != nil {
		return nil, errcode.RedisErr
	}
	resp := &model.LoginRep{
		Token: tokenString,
		User:  user,
	}

	cookie := (&http.Cookie{
		Name:  "token",
		Value: tokenString,
		Path:  "/",
		//Domain:   "hoper.xyz",
		Expires:  time.Now().Add(confdao.Conf.Customize.TokenMaxAge * time.Second),
		MaxAge:   int(confdao.Conf.Customize.TokenMaxAge),
		Secure:   false,
		HttpOnly: true,
	}).String()
	err = (*httpctx.HttpContext)(ctxi).SetCookie(cookie)
	if err != nil {
		return nil, errcode.Unavailable
	}
	return resp, nil
}

func (u *UserService) Logout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("Logout")()
	user, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	confdao.Dao.GORMDB.Table(modelconst.UserExtTableName).Where(`id = ?`, user.Id).UpdateColumn("last_activated_at", time.Now())

	if err := confdao.Dao.Redis.Del(ctx, redisi.CommandDEL, modelconst.LoginUserKey+strconv.FormatUint(user.Id, 10)).Err(); err != nil {
		return nil, ctxi.ErrorLog(errcode.RedisErr, err, "redisi.Del")
	}
	cookie := (&http.Cookie{
		Name:  httpi.HeaderCookieValueToken,
		Value: httpi.HeaderCookieValueDel,
		Path:  "/",
		//Domain:   "hoper.xyz",
		Expires:  time.Now().Add(-1),
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
	}).String()
	err = (*httpctx.HttpContext)(ctxi).SetCookie(cookie)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return new(emptypb.Empty), nil
}

func (u *UserService) AuthInfo(ctx context.Context, req *emptypb.Empty) (*model.AuthInfoRep, error) {
	ctxi := httpctx.FromContextValue(ctx)
	user, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	return user.UserAuthInfo(), nil
}

func (u *UserService) Info(ctx context.Context, req *request.Id) (*model.UserRep, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	if req.Id == 0 {
		req.Id = auth.Id
	}
	userRedisDao := redis.GetUserDao(ctxi, confdao.Dao.Redis.Client)
	db := gormi.NewTraceDB(confdao.Dao.GORMDB.DB, ctxi.BaseContext(), ctxi.TraceID)
	var user1 model.User
	if err = db.First(&user1, req.Id).Error; err != nil {
		return nil, errcode.DBError.Message("账号不存在")
	}
	userExt, err := userRedisDao.GetUserExtRedis()
	if err != nil {
		return nil, err
	}
	return &model.UserRep{User: &user1, UerExt: userExt}, nil
}

func (u *UserService) ForgetPassword(ctx context.Context, req *model.LoginReq) (*wrappers.StringValue, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()
	if verifyErr := luosimao.Verify(confdao.Conf.Customize.LuosimaoVerifyURL, confdao.Conf.Customize.LuosimaoAPIKey, req.VCode); verifyErr != nil {
		return nil, errcode.InvalidArgument.Warp(verifyErr)
	}

	if req.Input == "" {
		return nil, errcode.InvalidArgument.Message("账号错误")
	}
	userDBDao := data.GetDBDao(ctxi, confdao.Dao.GORMDB.DB)

	user, err := userDBDao.GetByEmailOrPhone(req.Input, req.Input, "id", "name", "password")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if validation.PhoneOrMail(req.Input) != validation.Phone {
				return nil, errcode.InvalidArgument.Message("邮箱不存在")
			} else {
				return nil, errcode.InvalidArgument.Message("手机号不存在")
			}
		}
		log.Error(err)
		return nil, errcode.DBError
	}
	restPassword := modelconst.ResetTimeKey + strconv.FormatUint(user.Id, 10)

	curTime := time.Now().Unix()
	if err := confdao.Dao.Redis.SetEX(ctx, restPassword, curTime, modelconst.ResetDuration).Err(); err != nil {
		log.Error("redis set failed:", err)
		return nil, errcode.RedisErr
	}

	go sendMail(ctxi, model.ActionRestPassword, curTime, user)

	return &wrappers.StringValue{Value: "注意查收邮件"}, nil
}

func (u *UserService) ResetPassword(ctx context.Context, req *model.ResetPasswordReq) (*wrappers.StringValue, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("ResetPassword")()

	redisKey := modelconst.ResetTimeKey + strconv.FormatUint(req.Id, 10)
	emailTime, err := confdao.Dao.Redis.Get(ctx, redisKey).Int64()
	if err != nil {
		return nil, ctxi.ErrorLog(errcode.InvalidArgument.Message("无效的链接"), err, "Redis.Get")
	}
	userDBDao := data.GetDBDao(ctxi, confdao.Dao.GORMDB.DB)

	user, err := userDBDao.GetByPrimaryKey(req.Id)
	if err != nil {
		return nil, err
	}
	if user.Status != 1 {
		return nil, errcode.FailedPrecondition.Message("无效账号")
	}
	secretStr := strconv.Itoa(int(emailTime)) + user.Mail + user.Password

	secretStr = fmt.Sprintf("%x", md5.Sum(stringsi.ToBytes(secretStr)))

	if req.Secret != secretStr {
		return nil, errcode.InvalidArgument.Message("无效的链接")
	}
	db := gormi.NewTraceDB(confdao.Dao.GORMDB.DB, ctxi.BaseContext(), ctxi.TraceID)
	if err := db.Table(modelconst.UserTableName).
		Where(`id = ?`, user.Id).Update("password", req.Password).Error; err != nil {
		log.Error("UserService.ResetPassword,DB.Update", err)
		return nil, errcode.DBError
	}
	confdao.Dao.Redis.Del(ctx, redisKey)
	return &wrappers.StringValue{Value: "重置成功，请重新登录"}, nil
}

func (*UserService) ActionLogList(ctx context.Context, req *model.ActionLogListReq) (*model.ActionLogListRep, error) {
	rep := &model.ActionLogListRep{}
	var logs []*model.ActionLog
	err := confdao.Dao.GORMDB.Table(modelconst.UserActionLogTableName).
		Offset(0).Limit(10).Find(&logs).Error
	if err != nil {
		return nil, errcode.DBError.Warp(err)
	}
	rep.List = logs
	return rep, nil
}

func (*UserService) BaseList(ctx context.Context, req *model.BaseListReq) (*model.BaseListRep, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("BaseList")()
	if ctxi.Internal == "" {
		return nil, errcode.PermissionDenied
	}
	ctx = ctxi.BaseContext()
	userDBDao := data.GetDBDao(ctxi, confdao.Dao.GORMDB.DB)

	count, users, err := userDBDao.GetBaseListDB(req.Ids, int(req.PageNo), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	return &model.BaseListRep{
		Total: count,
		List:  users,
	}, nil
}

func (*UserService) GetTest(ctx context.Context, req *request.Id) (*model.User, error) {
	return &model.User{Id: req.Id, Name: "测试"}, nil
}

func (*UserService) Service() (string, string, []gin.HandlerFunc) {
	return "用户相关", "/api/user", []gin.HandlerFunc{middle.GinLog}
}

func (*UserService) Add(ctx *ginctx.Context, req *model.SignupReq) (*wrappers.StringValue, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() {
		pick.Get("/add").
			Title("用户注册").
			Version(2).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试").End()
	})
	client := confdao.Dao.Redis
	cmd, _ := client.Do(ctx.BaseContext(), "HGETALL", modelconst.LoginUserKey+"1").Result()
	log.Debug(cmd)

	return &wrappers.StringValue{Value: req.Name}, nil
}

func (*UserService) Addv(ctx *ginctx.Context, req *response.TinyRep) (*response.TinyRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() {
		pick.Post("/add").
			Title("用户注册").
			Version(1).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试").End()
	})
	return req, nil
}

func (u *UserService) EasySignup(ctx context.Context, req *model.SignupReq) (*model.LoginRep, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("EasySignup")()

	if req.Mail == "" && req.Phone == "" {
		return nil, errcode.InvalidArgument.Message("请填写邮箱或手机号")
	}

	userDBDao := data.GetDBDao(ctxi, confdao.Dao.GORMDB.DB)
	checkUser, err := userDBDao.GetByNameOrEmailOrPhone(req.Name, req.Mail, req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errcode.DBError
	}
	if err == nil {
		if checkUser.Name == req.Name {
			return nil, errcode.InvalidArgument.Message("用户名已被注册")
		}
		if checkUser.Mail == req.Mail {
			return nil, errcode.InvalidArgument.Message("邮箱已被注册")
		}
		if checkUser.Phone == req.Phone {
			return nil, errcode.InvalidArgument.Message("手机号已被注册")
		}
	}

	var user = &model.User{
		Name:        req.Name,
		Account:     uuid.New().String(),
		Mail:        req.Mail,
		Phone:       req.Phone,
		Gender:      req.Gender,
		AvatarUrl:   modelconst.DefaultAvatar,
		Role:        model.RoleNormal,
		ActivatedAt: timepb.NewTime(ctxi.RequestAt.Time),
		Status:      model.UserStatusActivated,
	}

	user.Password = encryptPassword(req.Password)
	if err := userDBDao.Creat(user); err != nil {
		return nil, ctxi.ErrorLog(errcode.DBError.Message("新建出错"), err, "UserService.Creat")
	}
	return u.login(ctxi, user)
}
