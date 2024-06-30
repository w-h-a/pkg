package broker

type HandlerWrapper func(Handler) Handler

type Handler func(Publication) error
