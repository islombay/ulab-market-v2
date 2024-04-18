package service

import (
	"app/pkg/logs"
	"app/storage"
)

type IServiceManager interface {

}

type Service struct {

}

func New(str storage.StoreInterface, log logs.LoggerInterface){
	
}