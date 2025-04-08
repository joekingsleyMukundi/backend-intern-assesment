package gapi

import (
	"github.com/joekingsleyMukundi/backend-intern-assesment/auth/pb"
	db "github.com/joekingsleyMukundi/backend-intern-assesment/common/db/sqlc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertUser(user db.User) *pb.User {
	return &pb.User{
		Username:  user.Username,
		FullNamr:  user.FullName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
