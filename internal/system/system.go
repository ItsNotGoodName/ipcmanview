package system

func GetConfig() (Config, error) {
	return app.CP.GetConfig()
}

type UpdateConfigParams struct {
	SiteName     string
	EnableSignUp bool
}

func UpdateConfig(arg UpdateConfigParams) error {
	return app.CP.UpdateConfig(func(cfg Config) (Config, error) {
		cfg.SiteName = arg.SiteName
		cfg.EnableSignUp = arg.EnableSignUp

		return cfg, nil
	})
}
