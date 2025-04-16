package service

import (

)

type Server struct {
	pb.AccountServer
	Queris *postgresql.Queries
}

func New(query *postgresql.Queries) *Server {
	return &Server{
		Queris: query,
	}
}
