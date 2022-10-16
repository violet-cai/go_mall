package logic

import (
	"code/service/order/rpc/types/order"
	"code/service/pay/model"
	"code/service/user/rpc/types/user"
	"context"
	"google.golang.org/grpc/status"

	"code/service/pay/rpc/internal/svc"
	"code/service/pay/rpc/types/pay"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackLogic {
	return &CallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallbackLogic) Callback(in *pay.CallbackRequest) (*pay.CallbackResponse, error) {
	_, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{Id: in.Uid})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.OrderRpc.Detail(l.ctx, &order.DetailRequest{Id: in.Oid})
	if err != nil {
		return nil, err
	}
	res, err := l.svcCtx.PayModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "支付不存在")
		}
		return nil, err
	}
	// 支付金额与订单金额不符
	if in.Amount != res.Amount {
		return nil, status.Error(100, "支付金额与订单金额不符")
	}
	res.Source = in.Source
	res.Status = in.Status
	_, err = l.svcCtx.OrderRpc.Paid(l.ctx, &order.PaidRequest{Id: in.Oid})
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	return &pay.CallbackResponse{}, nil
}
