package gapi

import (
	"fmt"

	db "github.com/joekingsleyMukundi/backend-intern-assesment/common/db/sqlc"
	"github.com/joekingsleyMukundi/backend-intern-assesment/common/token"
	"github.com/joekingsleyMukundi/backend-intern-assesment/common/util"
	"github.com/joekingsleyMukundi/backend-intern-assesment/payments/pb"
)

type Server struct {
	pb.UnimplementedPaymentsServer
	tokenMaker token.Maker
	store      db.Store
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMwker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create token: %d", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	return server, nil
}
