package account

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"user-api/internal/biz"
	"user-api/internal/model"

	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResp, err error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		return nil, biz.TokenError
	}
	userModel := model.NewUserModel(l.svcCtx.Mysql)
	user, err := userModel.FindOne(l.ctx, userId)
	if err != nil && (errors.Is(err, model.ErrNotFound) || errors.Is(err, sql.ErrNoRows)) {
		return nil, biz.TokenError
	}
	if err != nil {
		l.Logger.Error("查询用户失败:", err)
		return nil, biz.DBError
	}

	return &types.UserInfoResp{
		Username: user.Username,
		Id:       user.Id,
	}, nil
}
