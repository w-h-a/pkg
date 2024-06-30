package server

import "github.com/google/uuid"

var (
	defaultNamespace = "default"
	defaultName      = "server"
	defaultID        = uuid.New().String()
	defaultVersion   = "v0.1.0"
	defaultAddress   = ":0"
)

type Server interface {
	Options() ServerOptions
	NewController(c interface{}, opts ...ControllerOption) Controller
	RegisterController(c Controller) error
	NewSubscriber(t string, sub interface{}, opts ...SubscriberOption) Subscriber
	RegisterSubscriber(sub Subscriber) error
	Run() error
	String() string
}
