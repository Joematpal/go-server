package grpc

type Option interface {
	applyOption(*Server) error
}

type optionApplyFunc func(*Server) error

func (f optionApplyFunc) applyOption(opts *Server) error {
	return f(opts)
}

func WithRegisterService(rs RegisterService) Option {
	return optionApplyFunc(func(s *Server) error {
		s.registerServices = append(s.registerServices, rs)
		return nil
	})
}
