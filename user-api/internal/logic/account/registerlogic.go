package account

import (
	"context"
	"errors"
	"time"
	"user-api/internal/model"

	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	//1、要用用户名查询 是否已经存在，存在返回错误
	userModel := model.NewUserModel(l.svcCtx.Mysql)
	user, err := userModel.FindByUsername(l.ctx, req.Username)
	if err != nil {
		l.Logger.Error("查询用户失败：", err)
		return nil, err
	}
	if user != nil {
		return nil, errors.New("此用户已注册")
	}
	//2.如果不存在，插入用户数据 注册
	_, err = userModel.Insert(l.ctx, &model.User{
		Username:      req.Username,
		Password:      req.Password,
		RegisterTime:  time.Now(),
		LastLoginTime: time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return
}
