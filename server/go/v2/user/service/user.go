package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	contexti "github.com/liov/hoper/go/v2/tailmon/context"
	"github.com/liov/hoper/go/v2/tailmon/pick"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/empty"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/protobuf/utils/response"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
	"github.com/liov/hoper/go/v2/user/middle"
	modelconst "github.com/liov/hoper/go/v2/user/model"
	redisi "github.com/liov/hoper/go/v2/utils/dao/redis"
	templatei "github.com/liov/hoper/go/v2/utils/def/template"
	"github.com/liov/hoper/go/v2/utils/log"
	httpi "github.com/liov/hoper/go/v2/utils/net/http"

	"github.com/liov/hoper/go/v2/utils/net/mail"
	"github.com/liov/hoper/go/v2/utils/strings"
	"github.com/liov/hoper/go/v2/utils/verification"
	"gorm.io/gorm"
)

type UserService struct {
	model.UnimplementedUserServiceServer
}

func (u *UserService) VerifyCode(ctx context.Context, req *empty.Empty) (*wrappers.StringValue, error) {
	device := contexti.CtxFromContext(ctx).DeviceInfo
	log.Debug(device)
	var rep = &wrappers.StringValue{}
	vcode := verification.GenerateCode()
	log.Info(vcode)
	rep.Value = vcode
	return rep, nil
}

func (*UserService) SignupVerify(ctx context.Context, req *model.SingUpVerifyReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context
	if req.Mail == "" && req.Phone == "" {
		return nil, errorcode.InvalidArgument.Message("请填写邮箱或手机号")
	}
	userDao := dao.GetDao(ctxi)
	db := dao.Dao.GetDB(ctxi.Logger)
	if exist, _ := userDao.ExitByEmailORPhone(db, req.Mail, req.Phone); exist {
		if req.Mail != "" {
			return nil, errorcode.InvalidArgument.Message("邮箱已被注册")
		} else {
			return nil, errorcode.InvalidArgument.Message("手机号已被注册")
		}
	}
	vcode := verification.GenerateCode()
	log.Debug(vcode)
	key := modelconst.VerificationCodeKey + req.Mail + req.Phone
	if err := dao.Dao.Redis.SetEX(ctx, key, vcode, modelconst.VerificationCodeDuration).Err(); err != nil {
		return nil, ctxi.ErrorLog(errorcode.RedisErr.Message("新建出错"), err, "SetEX")
	}
	return new(empty.Empty), nil
}

func (u *UserService) Signup(ctx context.Context, req *model.SignupReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context
	if err := Validate(req); err != nil {
		return nil, err
	}
	if req.Mail == "" && req.Phone == "" {
		return nil, errorcode.InvalidArgument.Message("请填写邮箱或手机号")
	}

	if err := LuosimaoVerify(req.VCode); err != nil {
		return nil, err
	}
	userDao := dao.GetDao(ctxi)
	db := dao.Dao.GetDB(ctxi.Logger)
	if exist, _ := userDao.ExitByEmailORPhone(db, req.Mail, req.Phone); exist {
		if req.Mail != "" {
			return nil, errorcode.InvalidArgument.Message("邮箱已被注册")
		} else {
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
		log.Error(err)
		return nil, errorcode.RedisErr.Message("新建出错")
	}

	activeUser := modelconst.ActiveTimeKey + strconv.FormatUint(user.Id, 10)

	curTime := time.Now().Unix()

	if err := dao.Dao.Redis.SetEX(ctx, activeUser, curTime, modelconst.ActiveDuration).Err(); err != nil {
		log.Error("UserService.Signup,RedisConn.Do: ", err)
	}

	if req.Mail != "" {
		go sendMail(model.ActionActive, curTime, user)
	}

	return nil, nil
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
	case model.ActionActive:
		ctiveOrRestPasswdValues.ActionURL = siteURL + "/user/active/" + strconv.FormatUint(user.Id, 10) + "/" + secretStr
		templ = modelconst.ActionActiveContent
	case model.ActionRestPassword:
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
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("Active")
	defer span.End()
	ctx = ctxi.Context
	redisKey := modelconst.ActiveTimeKey + strconv.FormatUint(req.Id, 10)
	emailTime, err := dao.Dao.Redis.Get(ctx, redisKey).Int64()
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.InvalidArgument.Message("无效的链接"), err, "Get")
	}
	userDao := dao.GetDao(ctxi)
	db := dao.Dao.GetDB(ctxi.Logger)
	user, err := userDao.GetByPrimaryKey(db, req.Id)
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

	db.Model(user).Updates(map[string]interface{}{"activated_at": time.Now(), "status": 1})
	dao.Dao.Redis.Del(ctx, redisKey)
	return u.login(ctxi, user)
}

func (u *UserService) Edit(ctx context.Context, req *model.EditReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("Edit")
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
		db := dao.Dao.GetDB(ctxi.Logger)
		originalIds, err := userDao.ResumesIds(db, user.Id)
		if err != nil {
			return nil, errorcode.DBError.Message("更新失败")
		}
		var resumes []*model.Resume
		resumes = append(req.Details.EduExps, req.Details.WorkExps...)
		if len(resumes) > 0 {
			tx := dao.Dao.GORMDB.Begin()
			err = userDao.SaveResumes(tx, req.Id, resumes, originalIds, model.ConvDeviceInfo(device))
			if err != nil {
				tx.Rollback()
				return nil, errorcode.DBError.Message("更新失败")
			}
			err = tx.Model(req.Details).UpdateColumns(req.Details).Error
			if err != nil {
				tx.Rollback()
				return nil, errorcode.DBError.Message("更新失败")
			}
			tx.Commit()
		}
	}
	return new(empty.Empty), nil
}

func (u *UserService) Login(ctx context.Context, req *model.LoginReq) (*model.LoginRep, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("Login")
	defer span.End()
	ctx = ctxi.Context
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
	db := dao.Dao.GetDB(ctxi.Logger)
	var user model.User
	if err := db.Table(modelconst.UserTableName).
		Where(sql, req.Input).Find(&user).Error; err != nil {
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
		if err := dao.Dao.Redis.SetEX(ctx, activeUser, curTime, modelconst.ActiveDuration).Err(); err != nil {
			return nil, ctxi.ErrorLog(errorcode.RedisErr, err, "SetEX")
		}
		go sendMail(model.ActionActive, curTime, &user)
		return nil, model.UserErrNoActive.Message("账号未激活,请进入邮箱点击激活")
	}

	return u.login(ctxi, &user)
}

func (*UserService) login(ctxi *contexti.Ctx, user *model.User) (*model.LoginRep, error) {
	auth := &model.AuthInfo{
		Id:     user.Id,
		Name:   user.Name,
		Role:   user.Role,
		Status: user.Status,
	}

	ctxi.AuthInfo = auth
	ctxi.LoginAt = ctxi.TimeStamp
	ctxi.ExpiredAt = ctxi.TimeStamp + int64(conf.Conf.Customize.TokenMaxAge)

	tokenString, err := ctxi.GenerateToken(stringsi.ToBytes(conf.Conf.Customize.TokenSecret))
	if err != nil {
		return nil, errorcode.Internal
	}
	db := dao.Dao.GetDB(ctxi.Logger)

	db.Table(modelconst.UserExtTableName).Where(`id = ?`, user.Id).
		UpdateColumn("last_activated_at", ctxi.RequestAt.TimeString)
	userDao := dao.GetDao(ctxi)
	if err := userDao.EfficientUserHashToRedis(); err != nil {
		return nil, errorcode.RedisErr
	}
	resp := &model.LoginRep{
		Token: tokenString,
		User: &model.UserBaseInfo{
			Id:        user.Id,
			Name:      user.Name,
			Gender:    user.Gender,
			AvatarUrl: user.AvatarUrl,
		}}

	cookie := (&http.Cookie{
		Name:  "token",
		Value: tokenString,
		Path:  "/",
		//Domain:   "hoper.xyz",
		Expires:  time.Now().Add(conf.Conf.Customize.TokenMaxAge * time.Second),
		MaxAge:   int(conf.Conf.Customize.TokenMaxAge),
		Secure:   false,
		HttpOnly: true,
	}).String()
	err = ctxi.SetCookie(cookie)
	if err != nil {
		return nil, errorcode.Unavailable
	}
	return resp, nil
}

func (u *UserService) Logout(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("Logout")
	defer span.End()
	ctx = ctxi.Context
	user, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	dao.Dao.GORMDB.Model(&model.UserAuthInfo{Id: user.Id}).UpdateColumn("last_activated_at", time.Now())

	if err := dao.Dao.Redis.Del(ctx, redisi.DEL, modelconst.LoginUserKey+strconv.FormatUint(user.Id, 10)).Err(); err != nil {
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
	ctxi.SetCookie(cookie)
	return nil, nil
}

func (u *UserService) AuthInfo(ctx context.Context, req *empty.Empty) (*model.UserAuthInfo, error) {
	ctxi := contexti.CtxFromContext(ctx)
	user, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	return user.UserAuthInfo(), nil
}

func (u *UserService) Info(ctx context.Context, req *request.Object) (*model.UserRep, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	if req.Id == 0 {
		req.Id = auth.Id
	}
	db := dao.Dao.GetDB(ctxi.Logger)
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

func (u *UserService) ForgetPassword(ctx context.Context, req *model.LoginReq) (*response.TinyRep, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context
	if verifyErr := verification.LuosimaoVerify(conf.Conf.Customize.LuosimaoVerifyURL, conf.Conf.Customize.LuosimaoAPIKey, req.VCode); verifyErr != nil {
		return nil, errorcode.InvalidArgument.Warp(verifyErr)
	}

	if req.Input == "" {
		return nil, errorcode.InvalidArgument.Message("账号错误")
	}
	db := dao.Dao.GetDB(ctxi.Logger)
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
	if err := dao.Dao.Redis.SetEX(ctx, restPassword, curTime, modelconst.ResetDuration).Err(); err != nil {
		log.Error("redis set failed:", err)
		return nil, errorcode.RedisErr
	}

	go sendMail(model.ActionRestPassword, curTime, user)

	return &response.TinyRep{}, nil
}

func (u *UserService) ResetPassword(ctx context.Context, req *model.ResetPasswordReq) (*response.TinyRep, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("Logout")
	defer span.End()
	ctx = ctxi.Context
	redisKey := modelconst.ResetTimeKey + strconv.FormatUint(req.Id, 10)
	emailTime, err := dao.Dao.Redis.Get(ctx, redisKey).Int64()
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.InvalidArgument.Message("无效的链接"), err, "Redis.Get")
	}
	db := dao.Dao.GetDB(ctxi.Logger)
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
	dao.Dao.Redis.Del(ctx, redisKey)
	return nil, nil
}

func (*UserService) ActionLogList(ctx context.Context, req *model.ActionLogListReq) (*model.ActionLogListRep, error) {
	rep := &model.ActionLogListRep{}
	var logs []*model.UserActionLog
	err := dao.Dao.GORMDB.Table(modelconst.UserActionLogTableName).
		Offset(0).Limit(10).Find(&logs).Error
	if err != nil {
		return nil, errorcode.DBError.Warp(err)
	}
	rep.List = logs
	return rep, nil
}

func (*UserService) BaseList(ctx context.Context, req *model.BaseListReq) (*model.BaseListRep, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	if ctxi.Internal == "" {
		return nil, errorcode.PermissionDenied
	}
	ctx = ctxi.Context
	db := dao.Dao.GetDB(ctxi.Logger)
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

func (*UserService) GetTest(ctx context.Context, req *request.Object) (*model.User, error) {
	return &model.User{Id: req.Id, Name: "测试"}, nil
}

func (*UserService) Service() (string, string, []http.HandlerFunc) {
	return "用户相关", "/api/user", []http.HandlerFunc{middle.Log}
}

func (*UserService) Add(ctx *contexti.Ctx, req *model.SignupReq) (*response.TinyRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() interface{} {
		return pick.Path("/add").
			Method(http.MethodGet).
			Title("用户注册").
			Version(2).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试")
	})
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Conf.Redis.Addr,
		Password: conf.Conf.Redis.Password,
	})
	cmd, _ := client.Do(ctx, "HGETALL", modelconst.LoginUserKey+"1").Result()
	log.Debug(cmd)

	return &response.TinyRep{Message: req.Name}, nil
}

func (*UserService) FiberService() (string, string, []fiber.Handler) {
	return "用户相关", "/api/user", []fiber.Handler{middle.FiberLog}
}

func (*UserService) Addv(ctx *contexti.Ctx, req *response.TinyRep) (*response.TinyRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.FiberApi(func() interface{} {
		return pick.Path("/add").
			Method(http.MethodGet).
			Title("用户注册").
			Version(1).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试")
	})
	return req, nil
}
