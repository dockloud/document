package config

import (
	"net/url"
	"os"

	"github.com/dockloud/config/pkg/config"
	"github.com/dockloud/config/pkg/secret"
	model "github.com/dockloud/document/internal/domain/config"
)

type viperConfig struct {
	Name      string `mapstructure:"name"`
	ID        string `mapstructure:"id"`
	Container string `mapstructure:"container"`
	Port      string `mapstructure:"port"`

	Database struct {
		User    string `mapstructure:"user"`
		Name    string `mapstructure:"name"`
		Address string `mapstructure:"address"`
	} `mapstructure:"database"`

	Redis struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"redis"`

	Dependencies []struct {
		Name string `mapstructure:"name"`
		ID   string `mapstructure:"id"`
	} `mapstructure:"dependencies"`
}

type secrets struct {
	DBPassword   string
	ServiceToken string
	GatewayToken string
	LinkToken    string
}

func LoadMergedConfig() (*model.MergedConfig, error) {
	confHandler := config.NewViperHandler[viperConfig](
		os.Getenv("CONFIG_PATH"),
		os.Getenv("CONFIG_NAME"),
		"yaml",
	)

	conf, err := confHandler.GetConfig()
	if err != nil {
		return nil, err
	}
	secr, err := loadSecrets()
	if err != nil {
		return nil, err
	}

	return conf.ToMergedConfig(*secr), nil
}

func loadSecrets() (*secrets, error) {
	secHandler := secret.NewDockerHandler()

	password, err := secHandler.ReadSecret("document_db_password")
	if err != nil {
		return nil, err
	}
	servJWT, err := secHandler.ReadSecret("document_jwt_secret")
	if err != nil {
		return nil, err
	}
	gatewayJWT, err := secHandler.ReadSecret("gateway_jwt_secret")
	if err != nil {
		return nil, err
	}
	linkJWT, err := secHandler.ReadSecret("link_jwt_secret")
	if err != nil {
		return nil, err
	}

	return &secrets{
		DBPassword:   password,
		ServiceToken: servJWT,
		GatewayToken: gatewayJWT,
		LinkToken:    linkJWT,
	}, nil
}

func (c viperConfig) ToMergedConfig(secr secrets) *model.MergedConfig {
	deps := make(map[string]model.DependencyConfig, len(c.Dependencies))
	for _, d := range c.Dependencies {
		deps[d.ID] = model.DependencyConfig{
			Name: d.Name,
			ID:   d.ID,
		}
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.Database.User, secr.DBPassword),
		Host:   c.Database.Address,
		Path:   c.Database.Name,
	}
	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()
	dbURL := u.String()

	return &model.MergedConfig{
		Name:      c.Name,
		ID:        c.ID,
		Container: c.Container,
		Port:      c.Port,
		JWT: model.JWTConfig{
			ServiceSecret: []byte(secr.ServiceToken),
			GatewaySecret: []byte(secr.GatewayToken),
			LinkSecret:    []byte(secr.LinkToken),
		},
		Database:     model.DatabaseConfig{URL: dbURL},
		Redis:        model.RedisConfig{URL: c.Redis.URL},
		Dependencies: deps,
	}
}
