package db

import "context"

type Client interface {
	Create(context.Context, string, interface{}) error
	List(context.Context, string, interface{}, func(interface{}) error) error
	Close(context.Context) error
}
