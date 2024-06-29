package server

import "context"

type ControllerWrapper func(ControllerFunc) ControllerFunc

type ControllerFunc func(ctx context.Context, req Request, rsp interface{}) error

type SubscriberWrapper func(SubscriberFunc) SubscriberFunc

type SubscriberFunc func(ctx context.Context, pub Publication) error
