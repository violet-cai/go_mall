package logic

import (
	"code/service/order/model"
	"code/service/user/rpc/types/user"
	"context"
	"google.golang.org/grpc/status"

	"code/service/order/rpc/internal/svc"
	"code/service/order/rpc/types/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLogic) List(in *order.ListRequest) (*order.ListResponse, error) {
	_, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{Id: in.Uid})
	if err != nil {
		return nil, err
	}
	list, err := l.svcCtx.OrderModel.FindAllByUid(l.ctx, in.Uid)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "订单不存在")
		}
		return nil, status.Error(500, err.Error())
	}
	orderList := make([]*order.DetailResponse, 0)
	for _, item := range list {
		orderList = append(orderList, &order.DetailResponse{
			Id:     item.Id,
			Uid:    item.Uid,
			Pid:    item.Pid,
			Amount: item.Amount,
			Status: item.Status,
		})
	}
	return &order.ListResponse{
		Data: orderList,
	}, nil
}
