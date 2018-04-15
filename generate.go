package cas

//go:generate mockgen -destination ./mock_io/rwc.go io ReadWriteCloser
//go:generate mockgen -destination ./mock_cas/kv.go -source store.go
