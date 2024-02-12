package web

func (s *impl) PreConfiguration() {
	s.StringVar(&s.Params.ModeValue, "mode", "", "mode value is https / http")
	s.UintVar(&s.Params.Port, "port", 0, "port using listing service")
	//s.X509KeyPairVar(&s.Params.Certs, "tls", "certificado server")*/
}

func (s *impl) PostConfiguration() error {
	/*if common.StringIsRequiered(&s.mode) != nil {
		return errors.New("mode is requeried")
	}

	vmode, err := common.ParseMode(&s.mode)
	s.Params.Mode = vmode
	if err != nil {
		return errors.New("mode value must be http/https")
	}

	if *vmode == commons.TLS {
		if s.Params.Certs.PrivateKey == nil {
			return errors.New("no se puedo cargar el certificado 2")
		}
	}*/
	return nil
}
