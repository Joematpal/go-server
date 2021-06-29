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

// host             string
func WithHost(host string) Option {
	return optionApplyFunc(func(s *Server) error {
		s.host = host
		return nil
	})
}

// port             string
func WithPort(port string) Option {
	return optionApplyFunc(func(s *Server) error {
		s.port = port
		return nil
	})
}

// tls              bool
func WithTLS(tls bool) Option {
	return optionApplyFunc(func(s *Server) error {
		s.tls = tls
		return nil
	})
}

// pubCert          string
func WithPubCert(pubCert string) Option {
	return optionApplyFunc(func(s *Server) error {
		s.pubCert = pubCert
		return nil
	})
}

// privCert         string
func WithPrivCert(privCert string) Option {
	return optionApplyFunc(func(s *Server) error {
		s.privCert = privCert
		return nil
	})
}

func WithLogger(logr logger) Option {
	return optionApplyFunc(func(s *Server) error {
		s.logger = logr
		return nil
	})
}
