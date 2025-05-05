package gapi

import (
	"context"
	"database/sql"
	"time"

	"github.com/hibiken/asynq"
	db "github.com/joekingsleyMukundi/backend-intern-assesment/common/db/sqlc"
	"github.com/joekingsleyMukundi/backend-intern-assesment/payments/mpesa"
	"github.com/joekingsleyMukundi/backend-intern-assesment/payments/pb"
	"github.com/joekingsleyMukundi/backend-intern-assesment/payments/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) InitializePayment(ctx context.Context, req *pb.InitializePaymentRequest) (*pb.InitializePaymentResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetOwner())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "something went wrong")
	}
	arg := db.CreatePaymentParams{
		Owner:  user.Username,
		Amount: req.GetAmount(),
		Status: "Initiated",
	}
	payment, err := server.store.CreatePayment(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment record: %v", err)
	}
	tPayload := &worker.PayloadInitiatePayment{
		Username: arg.Owner,
		Phone:    req.Phone,
		Amount:   arg.Amount,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritiacal),
	}
	err = server.taskDistributor.DistributetaskInitiatePayment(ctx, tPayload, opts...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create  : %v", err)
	}
	accessToken, err := mpesa.GetAccessToken()
	if err != nil {
		updateArg := db.UpdatePaymentParams{
			ID:     payment.ID,
			Status: "Failed",
		}
		payment, err = server.store.UpdatePayment(ctx, updateArg)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update failed payment: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to initiate payment: %v", err)
	}
	stkResponse, err := mpesa.InitializeSTKPush(accessToken.AccessToken, req.GetPhone(), req.GetAmount())
	if err != nil {
		updateArg := db.UpdatePaymentParams{
			ID:     payment.ID,
			Status: "Failed",
		}
		payment, err = server.store.UpdatePayment(ctx, updateArg)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update failed payment: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to initiate payment: %v", err)
	}
	if stkResponse.ResponseDescription == "Success. Request accepted for processing" {
		updateArg := db.UpdatePaymentParams{
			ID:     payment.ID,
			Status: "Pending Confirmtion",
		}
		payment, err = server.store.UpdatePayment(ctx, updateArg)
	}
	rsp := &pb.InitializePaymentResponse{
		Payment: &pb.Payment{
			Id:        payment.ID,
			Owner:     arg.Owner,
			Status:    payment.Status,
			Amount:    payment.Amount,
			CreatedAt: timestamppb.New(payment.CreatedAt),
		},
	}
	return rsp, nil
}
func (server *Server) GetPaymentstatus(ctx context.Context, req *pb.GetPaymentStatusRequest) (*pb.GetPaymentStatusResponse, error) {
	payment, err := server.store.GetPayment(ctx, req.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "payment not found")
		}
		return nil, status.Errorf(codes.Internal, "something went wrong")
	}
	rsp := &pb.GetPaymentStatusResponse{
		Payment: &pb.Payment{
			Id:        payment.ID,
			Owner:     payment.Owner,
			Status:    payment.Status,
			Amount:    payment.Amount,
			CreatedAt: timestamppb.New(payment.CreatedAt),
		},
	}
	return rsp, nil
}
func (server *Server) UpdatePayment(ctx context.Context, req *pb.UpdatePaymentRequest) (*pb.UpdatePaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePayment not implemented")
}
