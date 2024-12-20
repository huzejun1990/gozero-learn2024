package account

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"user-api/internal/biz"
	"user-api/internal/model"

	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	//1.校验用户名密码
	userModel := model.NewUserModel(l.svcCtx.Mysql)
	user, err := userModel.FindByUsernameAndPwd(l.ctx, req.Username, req.Password)
	if err != nil {
		l.Logger.Error("查询用户失败:", err)
		return nil, biz.DBError
	}
	if user == nil {
		return nil, biz.UserNameAndPwdError
	}
	secret := l.svcCtx.Config.Auth.AccessSecret
	expire := l.svcCtx.Config.Auth.Expire
	//2.生成token
	token, err := biz.GetJwtToken(secret, time.Now().Unix(), expire, user.Id)
	if err != nil {
		return nil, biz.TokenError
	}
	//3.把token存入redis中，操作过期，踢掉线
	//把token存入redis
	err = l.svcCtx.RedisClient.SetexCtx(context.Background(), "token:"+token, strconv.FormatInt(user.Id, 10), int(expire))
	if err != nil {
		return nil, biz.RedisError
	}
	/*	err = l.svcCtx.Redis.SetexCtx(l.ctx, "token:"+token, strconv.FormatInt(user.Id, 10), int(expire))
		if err != nil {
			l.Logger.Error("存入redis失败:", err)
			return nil, biz.RedisError
		}*/

	resp = &types.LoginResp{
		Token: token,
	}
	//4.记录日志
	go func() {
		logData := map[string]any{
			"username": user.Username,
			"ip":       l.ctx.Value("ip"),
			"userId":   user.Id,
			"time":     time.Now().Format("2006-01-02 15:04:05"),
			"type":     "login",
			"msg":      "登录成功",
		}
		bytes, _ := json.Marshal(logData)
		err2 := l.svcCtx.KafkaPushCli.PushWithKey(
			context.Background(),
			fmt.Sprintf("log_%s", user.Id),
			string(bytes))
		if err2 != nil {
			l.Logger.Error("写入日志失败:", err2)
		}
	}()
	return

}
