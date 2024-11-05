package account

import (
	"context"
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
	resp = &types.LoginResp{
		Token: token,
	}
	return

}
