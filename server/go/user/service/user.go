package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hopeio/zeta/context/http_context"
	"github.com/hopeio/zeta/pick"
	"github.com/hopeio/zeta/protobuf/request"
	contexti2 "github.com/hopeio/zeta/utils/context"
	dbi "github.com/hopeio/zeta/utils/dao/db/const"
	"github.com/hopeio/zeta/utils/sdk/luosimao"
	stringsi "github.com/hopeio/zeta/utils/strings"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"strconv"
	"time"

	model "github.com/actliboy/hoper/server/go/protobuf/user"
	"github.com/actliboy/hoper/server/go/user/confdao"
	"github.com/actliboy/hoper/server/go/user/dao"
	"github.com/actliboy/hoper/server/go/user/middle"
	modelconst "github.com/actliboy/hoper/server/go/user/model"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"

	"github.com/hopeio/zeta/protobuf/errorcode"
	"github.com/hopeio/zeta/protobuf/response"
	redisi "github.com/hopeio/zeta/utils/dao/redis"
	templatei "github.com/hopeio/zeta/utils/definition/template"
	"github.com/hopeio/zeta/utils/log"
	httpi "github.com/hopeio/zeta/utils/net/http"

	"github.com/hopeio/zeta/utils/net/mail"
	"github.com/hopeio/zeta/utils/verification"
	"gorm.io/gorm"
)

type UserService struct {
	model.UnimplementedUserServiceServer
}

func (u *UserService) VerifyCode(ctx context.Context, req *emptypb.Empty) (*wrappers.StringValue, error) {
	device := http_context.ContextFromContext(ctx).DeviceInfo
	log.Debug(device)
	var rep = &wrappers.StringValue{}
	vcode := verification.GenerateCode()
	log.Info(vcode)
	rep.Value = vcode
	return rep, nil
}

// 验证码
func (u *UserService) SendVerifyCode(ctx context.Context, req *model.SendVerifyCodeReq) (*emptypb.Empty, error) {
	return nil, nil
}

func (*UserService) SignupVerify(ctx context.Context, req *model.SingUpVerifyReq) (*wrappers.StringValue, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context
	if req.Mail == "" && req.Phone == "" {
		return nil, errorcode.InvalidArgument.Message("请填写邮箱或手机号")
	}
	userDao := dao.GetDao(ctxi)
	db := userDao.NewDB(confdao.Dao.GORMDB.DB)
	if req.Mail != "" {
		if exist, _ := userDao.ExitsCheck(db, "mail", req.Phone); exist {
			return nil, errorcode.InvalidArgument.Message("邮箱已被注册")
		}
	}
	if req.Phone != "" {
		if exist, _ := userDao.ExitsCheck(db, "phone", req.Phone); exist {
			return nil, errorcode.InvalidArgument.Message("手机号已被注册")
		}
	}
	vcode := verification.GenerateCode()
	log.Debug(vcode)
	key := modelconst.VerificationCodeKey + req.Mail + req.Phone
	if err := confdao.Dao.Redis.SetEX(ctx, key, vcode, modelconst.VerificationCodeDuration).Err(); err != nil {
		return nil, ctxi.ErrorLog(errorcode.RedisErr.Message("新建出错"), err, "SetEX")
	}
	return new(wrappers.StringValue), nil
}

func (u *UserService) Signup(ctx context.Context, req *model.SignupReq) (*wrappers.StringValue, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context

	if req.Mail == "" && req.Phone == "" {
		return nil, errorcode.InvalidArgument.Message("请填写邮箱或手机号")
	}
	if req.VCode != confdao.Conf.Customize.LuosimaoSuperPW {
		if err := LuosimaoVerify(req.VCode); err != nil {
			return nil, err
		}
	}

	userDao := dao.GetDao(ctxi)
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	if exist, _ := userDao.ExitsCheck(db, "name", req.Name); exist {
		return nil, errorcode.InvalidArgument.Message("用户名已被注册")
	}
	if req.Mail != "" {
		if exist, _ := userDao.ExitsCheck(db, "mail", req.Mail); exist {
			return nil, errorcode.InvalidArgument.Message("邮箱已被注册")
		}
	}
	if req.Phone != "" {
		if exist, _ := userDao.ExitsCheck(db, "phone", req.Phone); exist {
			return nil, errorcode.InvalidArgument.Message("手机号已被注册")
		}
	}
	formatNow := ctxi.TimeString
	var user = &model.User{
		Name:      req.Name,
		Account:   uuid.New().String(),
		Mail:      req.Mail,
		Phone:     req.Phone,
		Gender:    req.Gender,
		AvatarUrl: modelconst.DefaultAvatar,
		Role:      model.RoleNormal,
		CreatedAt: formatNow,
		Status:    model.UserStatusInActive,
	}

	user.Password = encryptPassword(req.Password)
	if err := userDao.Creat(db, user); err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError.Message("新建出错"), err, "UserService.Creat")
	}

	activeUser := modelconst.ActiveTimeKey + strconv.FormatUint(user.Id, 10)

	curTime := ctxi.RequestAt.TimeStamp

	if err := confdao.Dao.Redis.SetEX(ctx, activeUser, curTime, modelconst.ActiveDuration).Err(); err != nil {
		return nil, ctxi.ErrorLog(errorcode.RedisErr, err, "UserService.Signup,SetEX")
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

func sendMail(ctxi *http_context.Context, action model.Action, curTime int64, user *model.User) {
	siteURL := "https://" + confdao.Conf.Server.Domain
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
	addr := confdao.Dao.Mail.Conf.Host + confdao.Dao.Mail.Conf.Port
	m := &mail.Mail{
		FromName: ctiveOrRestPasswdValues.SiteName,
		From:     confdao.Dao.Mail.Conf.From,
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
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("Active")
	defer span.End()
	ctx = ctxi.Context
	redisKey := modelconst.ActiveTimeKey + strconv.FormatUint(req.Id, 10)
	emailTime, err := confdao.Dao.Redis.Get(ctx, redisKey).Int64()
	if err != nil {
		userDao := dao.GetDao(ctxi)
		db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
		user, err := userDao.GetByPrimaryKey(db, req.Id)
		if err != nil {
			return nil, errorcode.DBError
		}
		go sendMail(ctxi, model.ActionActive, ctxi.RequestAt.TimeStamp, user)
		return nil, ctxi.ErrorLog(errorcode.InvalidArgument.Message("已过激活期限"), err, "Get")
	}
	userDao := dao.GetDao(ctxi)
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	user, err := userDao.GetByPrimaryKey(db, req.Id)
	if err != nil {
		return nil, errorcode.DBError
	}
	if user.Status != model.UserStatusInActive {
		return nil, errorcode.AlreadyExists.Message("已激活")
	}
	secretStr := strconv.Itoa((int)(emailTime)) + user.Mail + user.Password

	secretStr = fmt.Sprintf("%x", md5.Sum(stringsi.ToBytes(secretStr)))

	if req.Secret != secretStr {
		return nil, errorcode.InvalidArgument.Message("无效的链接")
	}

	db.Model(user).Updates(map[string]interface{}{"activated_at": time.Now(), "status": model.UserStatusActivated})
	confdao.Dao.Redis.Del(ctx, redisKey)
	return u.login(ctxi, user)
}

func (u *UserService) Edit(ctx context.Context, req *model.EditReq) (*emptypb.Empty, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("Edit")
	defer span.End()
	ctx = ctxi.Context
	user, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	if user.Id != req.Id {
		return nil, errorcode.PermissionDenied
	}
	device := ctxi.DeviceInfo

	if req.Details != nil {
		userDao := dao.GetDao(ctxi)
		db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
		originalIds, err := userDao.ResumesIds(db, user.Id)
		if err != nil {
			return nil, errorcode.DBError.Message("更新失败")
		}
		var resumes []*model.Resume
		resumes = append(req.Details.EduExps, req.Details.WorkExps...)
		tx := confdao.Dao.GORMDB.Begin()
		if len(resumes) > 0 {
			err = userDao.SaveResumes(tx, req.Id, resumes, originalIds, model.ConvDeviceInfo(device))
			if err != nil {
				tx.Rollback()
				return nil, errorcode.DBError.Message("更新失败")
			}
		}
		err = tx.Table(modelconst.UserTableName).Where(`id = ?`, req.Id).UpdateColumns(req.Details).Error
		if err != nil {
			tx.Rollback()
			return nil, errorcode.DBError.Message("更新失败")
		}
		tx.Commit()
	}
	return new(emptypb.Empty), nil
}

func (u *UserService) Login(ctx context.Context, req *model.LoginReq) (*model.LoginRep, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("Login")
	defer span.End()
	ctx = ctxi.Context

	if req.VCode != confdao.Conf.Customize.LuosimaoSuperPW {
		if err := LuosimaoVerify(req.VCode); err != nil {
			return nil, err
		}
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
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	var user model.User
	if err := db.Table(modelconst.UserTableName).
		Where(sql+` AND status != ?`+dbi.WithNotDeleted, req.Input, model.UserStatusDeleted).First(&user).Error; err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError.Message("账号不存在"), err, "Find")
	}

	if !checkPassword(req.Password, &user) {
		return nil, errorcode.InvalidArgument.Message("密码错误")
	}
	if user.Status == model.UserStatusInActive {
		//没看懂
		//encodedEmail := base64.StdEncoding.EncodeToString(stringsi.ToBytes(user.Mail))
		activeUser := modelconst.ActiveTimeKey + strconv.FormatUint(user.Id, 10)

		curTime := time.Now().Unix()
		if err := confdao.Dao.Redis.SetEX(ctx, activeUser, curTime, modelconst.ActiveDuration).Err(); err != nil {
			return nil, ctxi.ErrorLog(errorcode.RedisErr, err, "SetEX")
		}
		go sendMail(ctxi, model.ActionActive, curTime, &user)
		return nil, model.UserErrNoActive.Message("账号未激活,请进入邮箱点击激活")
	}

	return u.login(ctxi, &user)
}

func (*UserService) login(ctxi *http_context.Context, user *model.User) (*model.LoginRep, error) {
	auth := &model.AuthInfo{
		Id:     user.Id,
		Name:   user.Name,
		Role:   user.Role,
		Status: user.Status,
	}

	ctxi.AuthInfo = auth
	ctxi.IssuedAt = &jwt.NumericDate{Time: ctxi.Time}
	ctxi.ExpiresAt = &jwt.NumericDate{Time: ctxi.Time.Add(confdao.Conf.Customize.TokenMaxAge)}

	tokenString, err := ctxi.GenerateToken(confdao.Conf.Customize.TokenSecret)
	if err != nil {
		return nil, errorcode.Internal
	}
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)

	db.Table(modelconst.UserExtTableName).Where(`id = ?`, user.Id).
		UpdateColumn("last_activated_at", ctxi.RequestAt.TimeString)
	userDao := dao.GetDao(ctxi)
	if err := userDao.EfficientUserHashToRedis(); err != nil {
		return nil, errorcode.RedisErr
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
	err = (*contexti2.HttpContext)(ctxi.RequestContext).SetCookie(cookie)
	if err != nil {
		return nil, errorcode.Unavailable
	}
	return resp, nil
}

func (u *UserService) Logout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("Logout")
	defer span.End()
	ctx = ctxi.Context
	user, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	confdao.Dao.GORMDB.Table(modelconst.UserExtTableName).Where(`id = ?`, user.Id).UpdateColumn("last_activated_at", time.Now())

	if err := confdao.Dao.Redis.Del(ctx, redisi.DEL, modelconst.LoginUserKey+strconv.FormatUint(user.Id, 10)).Err(); err != nil {
		return nil, ctxi.ErrorLog(errorcode.RedisErr, err, "redisi.Del")
	}
	cookie := (&http.Cookie{
		Name:  httpi.HeaderCookieToken,
		Value: httpi.HeaderCookieDel,
		Path:  "/",
		//Domain:   "hoper.xyz",
		Expires:  time.Now().Add(-1),
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
	}).String()
	(*contexti2.HttpContext)(ctxi.RequestContext).SetCookie(cookie)
	return new(emptypb.Empty), nil
}

func (u *UserService) AuthInfo(ctx context.Context, req *emptypb.Empty) (*model.UserAuthInfo, error) {
	ctxi := http_context.ContextFromContext(ctx)
	user, err := auth(ctxi, true)
	if err != nil {
		return nil, nil
	}
	return user.UserAuthInfo(), nil
}

func (u *UserService) Info(ctx context.Context, req *request.Id) (*model.UserRep, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	if req.Id == 0 {
		req.Id = auth.Id
	}
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	userDao := dao.GetDao(ctxi)
	var user1 model.User
	if err = db.Find(&user1, req.Id).Error; err != nil {
		return nil, errorcode.DBError.Message("账号不存在")
	}
	userExt, err := userDao.GetUserExtRedis()
	if err != nil {
		return nil, err
	}
	return &model.UserRep{User: &user1, UerExt: userExt}, nil
}

func (u *UserService) ForgetPassword(ctx context.Context, req *model.LoginReq) (*wrappers.StringValue, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context
	if verifyErr := luosimao.Verify(confdao.Conf.Customize.LuosimaoVerifyURL, confdao.Conf.Customize.LuosimaoAPIKey, req.VCode); verifyErr != nil {
		return nil, errorcode.InvalidArgument.Warp(verifyErr)
	}

	if req.Input == "" {
		return nil, errorcode.InvalidArgument.Message("账号错误")
	}
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	userDao := dao.GetDao(ctxi)
	user, err := userDao.GetByEmailORPhone(db, req.Input, req.Input, "id", "name", "password")
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

	curTime := time.Now().Unix()
	if err := confdao.Dao.Redis.SetEX(ctx, restPassword, curTime, modelconst.ResetDuration).Err(); err != nil {
		log.Error("redis set failed:", err)
		return nil, errorcode.RedisErr
	}

	go sendMail(ctxi, model.ActionRestPassword, curTime, user)

	return &wrappers.StringValue{Value: "注意查收邮件"}, nil
}

func (u *UserService) ResetPassword(ctx context.Context, req *model.ResetPasswordReq) (*wrappers.StringValue, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("Logout")
	defer span.End()
	ctx = ctxi.Context
	redisKey := modelconst.ResetTimeKey + strconv.FormatUint(req.Id, 10)
	emailTime, err := confdao.Dao.Redis.Get(ctx, redisKey).Int64()
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.InvalidArgument.Message("无效的链接"), err, "Redis.Get")
	}
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	userDao := dao.GetDao(ctxi)
	user, err := userDao.GetByPrimaryKey(db, req.Id)
	if err != nil {
		return nil, err
	}
	if user.Status != 1 {
		return nil, errorcode.FailedPrecondition.Message("无效账号")
	}
	secretStr := strconv.Itoa(int(emailTime)) + user.Mail + user.Password

	secretStr = fmt.Sprintf("%x", md5.Sum(stringsi.ToBytes(secretStr)))

	if req.Secret != secretStr {
		return nil, errorcode.InvalidArgument.Message("无效的链接")
	}

	if err := db.Table(modelconst.UserTableName).
		Where(`id = ?`, user.Id).Update("password", req.Password).Error; err != nil {
		log.Error("UserService.ResetPassword,DB.Update", err)
		return nil, errorcode.DBError
	}
	confdao.Dao.Redis.Del(ctx, redisKey)
	return &wrappers.StringValue{Value: "重置成功，请重新登录"}, nil
}

func (*UserService) ActionLogList(ctx context.Context, req *model.ActionLogListReq) (*model.ActionLogListRep, error) {
	rep := &model.ActionLogListRep{}
	var logs []*model.UserActionLog
	err := confdao.Dao.GORMDB.Table(modelconst.UserActionLogTableName).
		Offset(0).Limit(10).Find(&logs).Error
	if err != nil {
		return nil, errorcode.DBError.Warp(err)
	}
	rep.List = logs
	return rep, nil
}

func (*UserService) BaseList(ctx context.Context, req *model.BaseListReq) (*model.BaseListRep, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	if ctxi.Internal == "" {
		return nil, errorcode.PermissionDenied
	}
	ctx = ctxi.Context
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	userDao := dao.GetDao(ctxi)
	count, users, err := userDao.GetBaseListDB(db, req.Ids, int(req.PageNo), int(req.PageSize))
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

func (*UserService) Service() (string, string, []http.HandlerFunc) {
	return "用户相关", "/api/user", []http.HandlerFunc{middle.Log}
}

func (*UserService) Add(ctx *http_context.Context, req *model.SignupReq) (*wrappers.StringValue, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() {
		pick.Get("/add").
			Title("用户注册").
			Version(2).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试").End()
	})
	client := redis.NewClient(&redis.Options{
		Addr:     confdao.Dao.Redis.Conf.Addr,
		Password: confdao.Dao.Redis.Conf.Password,
	})
	cmd, _ := client.Do(ctx, "HGETALL", modelconst.LoginUserKey+"1").Result()
	log.Debug(cmd)

	return &wrappers.StringValue{Value: req.Name}, nil
}

func (*UserService) Addv(ctx *http_context.Context, req *response.TinyRep) (*response.TinyRep, error) {
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
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context

	if req.Mail == "" && req.Phone == "" {
		return nil, errorcode.InvalidArgument.Message("请填写邮箱或手机号")
	}

	userDao := dao.GetDao(ctxi)
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	if exist, _ := userDao.ExitsCheck(db, "name", req.Name); exist {
		return nil, errorcode.InvalidArgument.Message("用户名已被注册")
	}
	if req.Mail != "" {
		if exist, _ := userDao.ExitsCheck(db, "mail", req.Phone); exist {
			return nil, errorcode.InvalidArgument.Message("邮箱已被注册")
		}
	}
	if req.Phone != "" {
		if exist, _ := userDao.ExitsCheck(db, "phone", req.Phone); exist {
			return nil, errorcode.InvalidArgument.Message("手机号已被注册")
		}
	}
	formatNow := ctxi.TimeString
	var user = &model.User{
		Name:        req.Name,
		Account:     uuid.New().String(),
		Mail:        req.Mail,
		Phone:       req.Phone,
		Gender:      req.Gender,
		AvatarUrl:   modelconst.DefaultAvatar,
		Role:        model.RoleNormal,
		CreatedAt:   formatNow,
		ActivatedAt: ctxi.RequestAt.TimeString,
		Status:      model.UserStatusActivated,
	}

	user.Password = encryptPassword(req.Password)
	if err := userDao.Creat(db, user); err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError.Message("新建出错"), err, "UserService.Creat")
	}
	return u.login(ctxi, user)
}
