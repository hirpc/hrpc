package option

type certs struct {
	PubKey []byte
	PriKey []byte
}

func WithServerCerts(pubkey []byte, prikey []byte) Option {
	return func(o *Options) {
		o.ServerCerts = &certs{
			PubKey: pubkey,
			PriKey: prikey,
		}
	}
}
