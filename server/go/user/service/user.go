package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	contextx "github.com/hopeio/cherry"
	"github.com/hopeio/gox/math/rand"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/sdk/luosimao"
	stringsx "github.com/hopeio/gox/strings"
	"github.com/hopeio/pick"
	"github.com/hopeio/protobuf/request"
	"github.com/hopeio/protobuf/response"
	"github.com/hopeio/protobuf/time/timestamp"
	"github.com/hopeio/scaffold/errcode"
	jwtx "github.com/hopeio/scaffold/jwt"
	"github.com/liov/hoper/server/go/global"
	userpb "github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/api/middle"
	"github.com/liov/hoper/server/go/user/data"
	redisop "github.com/liov/hoper/server/go/user/data/redis"
	modelconst "github.com/liov/hoper/server/go/user/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/hopeio/gox/log"
	"github.com/hopeio/gox/net/mail"
	templatex "github.com/hopeio/gox/text/template"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
}

func (u *UserService) VerifyCode(ctx context.Context, req *userpb.VerifyCodeReq) (*emptypb.Empty, error) {
	if req.Mail != "" && req.Phone != "" {
		return nil, errcode.InvalidArgument.Msg("auth.err.onlyOneContact")
	}
	if req.Mail == "" && req.Phone == "" {
		return nil, errcode.InvalidArgument.Msg("auth.err.contactRequired")
	}
	_, err := u.SignupVerify(ctx, &userpb.SingUpVerifyReq{
		Mail: req.Mail,
		CountryCallingCode: req.CountryCallingCode,
		Phone: req.Phone,
	})
	if err != nil {
		return nil, err
	}
	vcode := rand.RandomNumber(6)
	log.Info(vcode)
	key := modelconst.VerificationCodeKey + req.Mail + req.CountryCallingCode +req.Phone
	if err = global.Dao.Redis.Set(ctx, key, vcode, modelconst.VerificationCodeDuration).Err(); err != nil {
		return nil, errcode.RedisErr.Wrap(err)
	}
	if req.Mail != "" {
		sendVcode(ctx, req.Action, vcode, req.Mail)
		return new(emptypb.Empty), nil
	}
	// 手机号：验证码已写入 Redis；下发走短信网关，接入前 Debug 下打印便于联调
	if global.Global.RootConfig.Debug {
		log.Infow("phone verify code (debug)", zap.String("phone", req.Phone), zap.String("code", vcode))
	}
	return new(emptypb.Empty), nil
}

func (*UserService) SignupVerify(ctx context.Context, req *userpb.SingUpVerifyReq) (*emptypb.Empty, error) {

	if req.Mail != "" && req.Phone != "" {
		return nil, errcode.InvalidArgument.Msg("auth.err.onlyOneContact")
	}
	if req.Mail == "" && req.Phone == "" {
		return nil, errcode.InvalidArgument.Msg("auth.err.contactRequired")
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	userDao := data.GetDBDao(db)
	input := req.Mail
	if input == "" {
		input = req.Phone
	}
	checkUser, err := userDao.GetByEmailOrPhone(ctx, req.Mail, req.CountryCallingCode, req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errcode.DBError
	}
	if err == nil {
		if checkUser.Mail == req.Mail {
			return nil, errcode.InvalidArgument.Msg("auth.err.mailRegistered")
		}
		if checkUser.Phone == req.Phone {
			return nil, errcode.InvalidArgument.Msg("auth.err.phoneRegistered")
		}
	}
	return new(emptypb.Empty), nil
}

func (u *UserService) Signup(ctx context.Context, req *userpb.SignupReq) (*wrappers.StringValue, error) {

	if req.Mail != "" && req.Phone != "" {
		return nil, errcode.InvalidArgument.Msg("auth.err.onlyOneContact")
	}
	if req.Mail == "" && req.Phone == "" {
		return nil, errcode.InvalidArgument.Msg("auth.err.contactRequired")
	}
	if req.VCode != global.Conf.User.LuosimaoSuperPW {
		if err := LuosimaoVerify(req.VCode); err != nil {
			return nil, err
		}
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	userDao := data.GetDBDao(db)

	_, err := u.SignupVerify(ctx, &userpb.SingUpVerifyReq{
		Mail: req.Mail,
		CountryCallingCode: req.CountryCallingCode,
		Phone: req.Phone,
	})
	if err != nil {
		return nil, err
	}

	if req.Name == "" {
		req.Name = rand.RandomChars(10)
	}

	var user = &userpb.User{
		Name:    req.Name,
		Account: uuid.New().String(),
		Mail:    req.Mail,
		Phone:   req.Phone,
		Gender:  req.Gender,
		Avatar:  modelconst.DefaultAvatar,
		Role:    userpb.RoleNormal,
		Status:  userpb.UserStatusInActive,
	}

	if req.VCode != "" {
		vcode, err := global.Dao.Redis.Get(ctx, modelconst.VerificationCodeKey + req.Mail + req.CountryCallingCode +req.Phone).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return nil, errcode.RedisErr.Wrap(err)
		}
		if vcode != req.VCode && vcode != "" {
			return nil, errcode.InvalidArgument.Msg("auth.err.invalidCode")
		}
		if err = global.Dao.Redis.Del(ctx, modelconst.VerificationCodeKey + req.Mail + req.CountryCallingCode +req.Phone).Err(); err != nil {
			return nil, errcode.RedisErr.Wrap(err)
		}
		user.Status = userpb.UserStatusActivated
	}


	user.Password = encryptPassword(req.Password)
	if err := userDao.Create(ctx, user); err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	if req.VCode != "" {
		return &wrappers.StringValue{Value: "注册成功"}, nil
	}

	activeUser := modelconst.ActiveTimeKey + strconv.FormatUint(user.Id, 10)

	curTime := time.Now().UnixMilli()

	if err := global.Dao.Redis.Set(ctx, activeUser, curTime, modelconst.ActiveDuration).Err(); err != nil {
		return nil, errcode.RedisErr.Wrap(err)
	}

	if req.Mail != "" {
		go sendMail(ctx, userpb.ActionActive, curTime, user)
	}

	return &wrappers.StringValue{Value: "注册成功，注意查收邮件"}, nil
}

// Salt 每个用户都有一个不同的盐
func salt(password string) string {
	return password[0:5]
}

// EncryptPassword 给密码加密
func encryptPassword(password string) string {
	hash := salt(password) + global.Conf.User.PassSalt + password[5:]
	return fmt.Sprintf("%x", md5.Sum(stringsx.ToBytes(hash)))
}

func sendMail(ctx context.Context, action userpb.Action, curTime int64, user *userpb.User) {
	siteURL := global.Conf.SiteURL
	title := action.String()
	secretStr := strconv.FormatInt(curTime, 10) + user.Mail + user.Password
	secretStr = fmt.Sprintf("%x", md5.Sum(stringsx.ToBytes(secretStr)))
	var activeOrRestPasswdValues = struct {
		UserName, SiteName, SiteURL, ActionURL, SecretStr string
	}{user.Name, global.Conf.SiteName, global.Conf.SiteURL, "", secretStr}
	var templ string
	switch action {
	case userpb.ActionActive:
		activeOrRestPasswdValues.ActionURL = siteURL + "/api/user/active/" + strconv.FormatUint(user.Id, 10) + "/" + secretStr
		templ = modelconst.ActionActiveContent
	case userpb.ActionRestPassword:
		activeOrRestPasswdValues.ActionURL = siteURL + "/api/user/resetPassword/" + strconv.FormatUint(user.Id, 10) + "/" + secretStr
		templ = modelconst.ActionRestPasswordContent
	}
	log.Debug(activeOrRestPasswdValues.ActionURL)
	var buf = new(bytes.Buffer)
	err := templatex.Execute(buf, templ, &activeOrRestPasswdValues)
	if err != nil {
		log.Error("executing template:", err)
	}
	//content += "<p><img src=\"" + siteURL + "/images/logo.png\" style=\"height: 42px;\"/></p>"
	//fmt.Println(content)
	content := buf.String()

	m := &mail.Mail{
		Addr:     global.Dao.Mail.Conf.Host + global.Dao.Mail.Conf.Port,
		FromName: activeOrRestPasswdValues.SiteName,
		From:     global.Dao.Mail.Conf.UserName,
		Subject:  title,
		Content:  content,
		To:       []string{user.Mail},
		Auth:     global.Dao.Mail.Auth,
	}
	log.Debug(content)
	err = m.SendMailTLS()
	if err != nil {
		log.Error("sendMail:", err)
	}
}

func sendVcode(ctx context.Context, action userpb.Action, vcode string, mailAddr string) {
	content := global.LocalizerMap["zh-Hans"].MustLocalize(&i18n.LocalizeConfig{
		MessageID: "auth.mail.verifyCodeContent",
		TemplateData: map[string]interface{}{"Action": action.String(), "Vcode": vcode},
	})
	m := &mail.Mail{
		Addr:     global.Dao.Mail.Conf.Host + global.Dao.Mail.Conf.Port,
		FromName: global.Conf.SiteName,
		From:     global.Dao.Mail.Conf.UserName,
		Subject:  global.LocalizerMap["zh-Hans"].MustLocalize(&i18n.LocalizeConfig{
			MessageID: "auth.mail.verifyCodeSubject",
		}),
		Content:  content,
		To:       []string{mailAddr},
		Auth:     global.Dao.Mail.Auth,
	}
	log.Debug(content)
	err := m.SendMailTLS()
	if err != nil {
		log.Error("sendMail:", err)
	}
}

// 验证密码是否正确
func checkPassword(password string, user *userpb.User) bool {
	if password == "" || user.Password == "" {
		return false
	}
	return encryptPassword(password) == user.Password
}

func (u *UserService) Active(ctx context.Context, req *userpb.ActiveReq) (*userpb.LoginResp, error) {

	db := global.Dao.GORMDB.DB.WithContext(ctx)
	userDBDao := data.GetDBDao(db)

	user, err := userDBDao.GetByPrimaryKey(ctx, req.Id)
	if err != nil {
		return nil, errcode.DBError
	}

	if user.Status != userpb.UserStatusInActive {
		return nil, errcode.AlreadyExists.Msg("auth.err.activated")
	}
	redisKey := modelconst.ActiveTimeKey + strconv.FormatUint(req.Id, 10)
	emailTime, err := global.Dao.Redis.Get(ctx, redisKey).Int64()
	if err != nil {
		go sendMail(ctx, userpb.ActionActive, time.Now().UnixMilli(), user)
		return nil, errcode.InvalidArgument.Msg("auth.err.activationExpired")
	}
	secretStr := strconv.Itoa((int)(emailTime)) + user.Mail + user.Password

	secretStr = fmt.Sprintf("%x", md5.Sum(stringsx.ToBytes(secretStr)))

	if req.Secret != secretStr {
		return nil, errcode.InvalidArgument.Msg("auth.err.invalidLink")
	}
	err = userDBDao.Active(ctx, user)
	if err != nil {
		return nil, errcode.DBError
	}
	global.Dao.Redis.Del(ctx, redisKey)
	return u.login(ctx, user)
}

func (u *UserService) Edit(ctx context.Context, req *userpb.EditReq) (*emptypb.Empty, error) {

	user, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	if user.Id != req.Id {
		return nil, errcode.PermissionDenied
	}

	device := Device(ctx)
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	if req.Detail != nil {
		userDBDao := data.GetDBDao(db)

		originalIds, err := userDBDao.ResumesIds(ctx, user.Id)
		if err != nil {
			return nil, errcode.DBError.Msg("auth.err.updateFailed")
		}
		var resumes []*userpb.Resume
		resumes = append(req.Detail.EduExps, req.Detail.WorkExps...)
		tx := db.Begin()
		defer tx.Rollback()
		userDBDao = data.GetDBDao(tx)
		if len(resumes) > 0 {
			err = userDBDao.SaveResumes(ctx, req.Id, resumes, originalIds, userpb.ConvDeviceInfo(device))
			if err != nil {
				return nil, errcode.DBError.Msg("auth.err.updateFailed")
			}
		}
		err = userDBDao.Update(ctx, req)
		if err != nil {
			return nil, errcode.DBError.Msg("auth.err.updateFailed")
		}
		tx.Commit()
	}
	return new(emptypb.Empty), nil
}

func (u *UserService) Login(ctx context.Context, req *userpb.LoginReq) (*userpb.LoginResp, error) {

	if req.VCode != global.Conf.User.LuosimaoSuperPW {
		if err := LuosimaoVerify(req.VCode); err != nil {
			return nil, err
		}
	}

	if req.Mail == "" && req.Phone == "" {
		return nil, errcode.InvalidArgument.Msg("auth.err.invalidAccount")
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	userDBDao := data.GetDBDao(db)
	user, err := userDBDao.UserInfoByAccount(ctx, req.Mail, req.CountryCallingCode, req.Phone)
	if err != nil {
		return nil, errcode.DBError.Msg("auth.err.accountNotFound")
	}

	if !checkPassword(req.Password, user) {
		return nil, errcode.InvalidArgument.Msg("auth.err.passwordWrong")
	}
	if user.Status == userpb.UserStatusInActive {
		//没看懂
		//encodedEmail := base64.StdEncoding.EncodeToString(stringsx.ToBytes(user.Mail))
		activeUser := modelconst.ActiveTimeKey + strconv.FormatUint(user.Id, 10)

		curTime := time.Now().UnixMilli()
		if err := global.Dao.Redis.Set(ctx, activeUser, curTime, modelconst.ActiveDuration).Err(); err != nil {
			return nil, errcode.RedisErr.Wrap(err)
		}
		go sendMail(ctx, userpb.ActionActive, curTime, user)
		return nil, userpb.UserErrNoActive.Msg("auth.err.notActivated")
	}

	return u.login(ctx, user)
}

func (*UserService) login(ctx context.Context, user *userpb.User) (*userpb.LoginResp, error) {
	md := contextx.GetMetadata(ctx)
	authorization := jwtx.Claims[*userpb.AuthInfo]{Auth: &userpb.AuthInfo{
		Id:   user.Id,
		Name: user.Name,
		Role: user.Role,
	}}
	now := time.Now()
	md.Data = authorization.Auth
	authorization.IssuedAt = &jwt.NumericDate{Time: now}
	authorization.ExpiresAt = &jwt.NumericDate{Time: now.Add(global.Conf.User.TokenMaxAge)}

	tokenString, err := authorization.GenerateToken(global.Conf.User.TokenSecretBytes)
	if err != nil {
		return nil, errcode.Internal
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)

	db.Table(modelconst.TableNameUserExt).Where(`id = ?`, user.Id).
		UpdateColumn("last_activated_at", now)
	userRedisDao := redisop.GetUserDao(global.Dao.Redis.Client)
	if err := userRedisDao.EfficientUserHashToRedis(ctx, authorization.Auth); err != nil {
		return nil, errcode.RedisErr
	}
	resp := &userpb.LoginResp{
		Token: tokenString,
		User:  user,
	}

	cookie := (&http.Cookie{
		Name:  "token",
		Value: tokenString,
		Path:  "/",
		//Domain:   "hoper.xyz",
		Expires:  time.Now().Add(global.Conf.User.TokenMaxAge * time.Second),
		MaxAge:   int(global.Conf.User.TokenMaxAge),
		Secure:   false,
		HttpOnly: true,
	}).String()

	serverTransportStream := grpc.ServerTransportStreamFromContext(ctx)
	if serverTransportStream != nil {
		err = serverTransportStream.SetHeader(metadata.MD{httpx.HeaderSetCookie: []string{cookie}})
		if err != nil {
			return nil, errcode.Unavailable
		}
	} else {
		md.ResponseWriter.Header().Set(httpx.HeaderSetCookie, cookie)
	}
	return resp, nil
}

func (u *UserService) Logout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	md := contextx.GetMetadata(ctx)

	user, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}

	global.Dao.GORMDB.Table(modelconst.TableNameUserExt).Where(`id = ?`, user.Id).UpdateColumn("last_activated_at", time.Now())

	if err := global.Dao.Redis.Del(ctx, modelconst.LoginUserKey+strconv.FormatUint(user.Id, 10)).Err(); err != nil {
		return nil, errcode.RedisErr.Wrap(err)
	}
	cookie := (&http.Cookie{
		Name:  httpx.HeaderCookieValueToken,
		Value: httpx.HeaderCookieValueDel,
		Path:  "/",
		//Domain:   "hoper.xyz",
		Expires:  time.Now().Add(-1),
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
	}).String()
	md.ResponseWriter.Header().Set(httpx.HeaderSetCookie, cookie)
	serverTransportStream := grpc.ServerTransportStreamFromContext(ctx)
	if serverTransportStream != nil {
		err = serverTransportStream.SetHeader(metadata.MD{httpx.HeaderSetCookie: []string{cookie}})
		if err != nil {
			return nil, errcode.Unavailable
		}
	}
	return new(emptypb.Empty), nil
}

func (u *UserService) AuthInfo(ctx context.Context, req *emptypb.Empty) (*userpb.Auth, error) {

	user, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	return user.Proto(), nil
}

func (u *UserService) Info(ctx context.Context, req *request.Id) (*userpb.UserResp, error) {

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	if req.Id == 0 {
		req.Id = auth.Id
	}

	userRedisDao := redisop.GetUserDao(global.Dao.Redis.Client)
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	var user1 userpb.User
	if err = db.First(&user1, req.Id).Error; err != nil {
		return nil, errcode.DBError.Msg("auth.err.accountNotFound")
	}
	userExt, err := userRedisDao.GetUserExtRedis(ctx, auth.Id)
	if err != nil {
		return nil, err
	}
	return &userpb.UserResp{User: &user1, UerExt: userExt}, nil
}

func (u *UserService) ForgetPassword(ctx context.Context, req *userpb.LoginReq) (*wrappers.StringValue, error) {

	if verifyErr := luosimao.Verify(global.Conf.User.LuosimaoVerifyURL, global.Conf.User.LuosimaoAPIKey, req.VCode); verifyErr != nil {
		return nil, errcode.InvalidArgument.Wrap(verifyErr)
	}

	if req.Mail == "" && req.Phone == "" {
		return nil, errcode.InvalidArgument.Msg("auth.err.invalidAccount")
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	userDBDao := data.GetDBDao(db)

	user, err := userDBDao.GetByEmailOrPhone(ctx, req.Mail, req.CountryCallingCode, req.Phone, "id", "name", "password")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if req.Mail != "" {
				return nil, errcode.InvalidArgument.Msg("auth.err.mailNotFound")
			} else {
				return nil, errcode.InvalidArgument.Msg("auth.err.phoneNotFound")
			}
		}
		log.Error(err)
		return nil, errcode.DBError
	}
	restPassword := modelconst.ResetTimeKey + strconv.FormatUint(user.Id, 10)

	curTime := time.Now().Unix()
	if err := global.Dao.Redis.Set(ctx, restPassword, curTime, modelconst.ResetDuration).Err(); err != nil {
		log.Error("redis set failed:", err)
		return nil, errcode.RedisErr
	}

	go sendMail(ctx, userpb.ActionRestPassword, curTime, user)

	return &wrappers.StringValue{Value: "注意查收邮件"}, nil
}

func (u *UserService) ResetPassword(ctx context.Context, req *userpb.ResetPasswordReq) (*wrappers.StringValue, error) {

	redisKey := modelconst.ResetTimeKey + strconv.FormatUint(req.Id, 10)
	emailTime, err := global.Dao.Redis.Get(ctx, redisKey).Int64()
	if err != nil {
		log.Errorw("Get faild", zap.Error(err))
		return nil, errcode.InvalidArgument.Msg("auth.err.invalidLink")
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	userDBDao := data.GetDBDao(db)
	user, err := userDBDao.GetByPrimaryKey(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if user.Status != 1 {
		return nil, errcode.FailedPrecondition.Msg("auth.err.invalidAccountStatus")
	}
	secretStr := strconv.Itoa(int(emailTime)) + user.Mail + user.Password

	secretStr = fmt.Sprintf("%x", md5.Sum(stringsx.ToBytes(secretStr)))

	if req.Secret != secretStr {
		return nil, errcode.InvalidArgument.Msg("auth.err.invalidLink")
	}

	if err := db.Table(modelconst.TableNameUser).
		Where(`id = ?`, user.Id).Update("password", req.Password).Error; err != nil {
		log.Error("UserService.ResetPassword,DB.Update", err)
		return nil, errcode.DBError
	}
	global.Dao.Redis.Del(ctx, redisKey)
	return &wrappers.StringValue{Value: "重置成功，请重新登录"}, nil
}

func (*UserService) ActionLogList(ctx context.Context, req *userpb.ActionLogListReq) (*userpb.ActionLogListResp, error) {
	resp := &userpb.ActionLogListResp{}
	var logs []*userpb.ActionLog
	err := global.Dao.GORMDB.Table(modelconst.TableNameActionLog).
		Offset(0).Limit(10).Find(&logs).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}

	resp.List = logs
	return resp, nil
}

func (*UserService) BaseList(ctx context.Context, req *userpb.BaseListReq) (*userpb.BaseListResp, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errcode.InvalidArgument
	}

	if md.Get(httpx.HeaderGrpcInternal) == nil || md.Get(httpx.HeaderGrpcInternal)[0] == "" {
		return nil, errcode.PermissionDenied
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	userDBDao := data.GetDBDao(db)

	count, users, err := userDBDao.GetBaseListDB(ctx, req.Ids, int(req.PageNo), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	return &userpb.BaseListResp{
		Total: count,
		List:  users,
	}, nil
}

func (*UserService) GetTest(ctx context.Context, req *request.Id) (*userpb.User, error) {
	return &userpb.User{Id: req.Id, Name: "测试"}, nil
}

func (*UserService) Service() (string, string, []gin.HandlerFunc) {
	return "用户相关", "/api/user", []gin.HandlerFunc{middle.GinLog}
}

func (*UserService) PickAdd(ctx *gin.Context, req *userpb.SignupReq) (*wrappers.StringValue, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() { pick.Get("/add").Title("用户注册").End() })
	client := global.Dao.Redis
	cmd, _ := client.Do(ctx, "HGETALL", modelconst.LoginUserKey+"1").Result()
	log.Debug(cmd)

	return &wrappers.StringValue{Value: req.Name}, nil
}

func (*UserService) PickAddv(ctx *gin.Context, req *response.ErrResp) (*response.ErrResp, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() { pick.Post("/add").Title("用户注册").End() })
	return req, nil
}

func (u *UserService) EasySignup(ctx context.Context, req *userpb.SignupReq) (*userpb.LoginResp, error) {

	if req.Mail != "" && req.Phone != "" {
		return nil, errcode.InvalidArgument.Msg("auth.err.onlyOneContact")
	}
	if req.Mail == "" && req.Phone == "" {
		return nil, errcode.InvalidArgument.Msg("auth.err.contactRequired")
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	userDBDao := data.GetDBDao(db)
	_, err := u.SignupVerify(ctx, &userpb.SingUpVerifyReq{
		Mail: req.Mail,
		CountryCallingCode: req.CountryCallingCode,
		Phone: req.Phone,
	})
	if err != nil {
		return nil, err
	}

	var user = &userpb.User{
		Name:        req.Name,
		Account:     uuid.New().String(),
		Mail:        req.Mail,
		Phone:       req.Phone,
		Gender:      req.Gender,
		Avatar:      modelconst.DefaultAvatar,
		Role:        userpb.RoleNormal,
		ActivatedAt: timestamp.New(time.Now()),
		Status:      userpb.UserStatusActivated,
	}

	user.Password = encryptPassword(req.Password)
	if err := userDBDao.Create(ctx, user); err != nil {
		log.Errorw("Create faild", zap.Error(err))
		return nil, errcode.DBError.Msg("auth.err.createFailed")
	}
	return u.login(ctx, user)
}
