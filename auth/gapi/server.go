package gapi

import (
	"fmt"

	"github.com/joekingsleyMukundi/backend-intern-assesment/auth/pb"
	db "github.com/joekingsleyMukundi/backend-intern-assesment/common/db/sqlc"
	"github.com/joekingsleyMukundi/backend-intern-assesment/common/token"
	"github.com/joekingsleyMukundi/backend-intern-assesment/common/util"
)

type Server struct {
	pb.UnimplementedAuthServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMwker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token:%d", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	return server, nil
}
