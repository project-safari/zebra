package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/spf13/cobra"
	"gojini.dev/web"
	"gopkg.in/yaml.v3"
)

const ReadWriteOnly = 0o600

func NewInitCmd() *cobra.Command {
	initCmd := new(cobra.Command)

	initCmd.Use = "init"
	initCmd.Short = "create default zebra server configuration"
	initCmd.RunE = initServer
	initCmd.SilenceUsage = true

	initCmd.Flags().StringP("store", "s", cwd("zebra-store"),
		"zebra server store (default: $PWD/zebra-store)")
	initCmd.Flags().StringP("address", "a", "tcp://127.0.0.1:443",
		"zebra server address (default: tcp://127.0.0.1:443")
	initCmd.Flags().StringP("cert", "t", cwd("zebra-server.crt"),
		"zebra server certificate (default: $PWD/zebra-server.crt)")
	initCmd.Flags().StringP("key", "k", cwd("zebra-server.key"),
		"zebra server key (default: $PWD/zebra-server.key)")
	initCmd.Flags().StringP("user", "u", "", "admin user configuration file")
	_ = initCmd.MarkFlagRequired("user")
	initCmd.Flags().StringP("password", "p", "", "admin user password")
	_ = initCmd.MarkFlagRequired("password")
	initCmd.Flags().StringP("auth-key", "j", "", "zebra server auth key")
	_ = initCmd.MarkFlagRequired("auth-key")

	return initCmd
}

type ServerConfig struct {
	Store struct {
		Root string `json:"rootDir"`
	} `json:"store"`

	Server struct {
		Address string   `json:"address"`
		TLS     *web.TLS `json:"tls"`
	} `json:"server"`

	AuthKey string `json:"authKey"`

	Admin *auth.User `json:"admin"`
}

func initServer(cmd *cobra.Command, args []string) error {
	cfgFile := cmd.Flag("config").Value.String()

	admin, err := makeAdminConfig(cmd)
	if err != nil {
		return err
	}

	serverCfg := new(ServerConfig)
	serverCfg.Store.Root = cmd.Flag("store").Value.String()
	serverCfg.Server.Address = cmd.Flag("address").Value.String()
	serverCfg.Server.TLS = new(web.TLS)
	serverCfg.Server.TLS.CertFile = cmd.Flag("cert").Value.String()
	serverCfg.Server.TLS.KeyFile = cmd.Flag("key").Value.String()

	serverCfg.Admin = admin
	serverCfg.AuthKey = cmd.Flag("auth-key").Value.String()

	data, err := json.MarshalIndent(serverCfg, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgFile, data, ReadWriteOnly)
}

func makeAdminConfig(cmd *cobra.Command) (*auth.User, error) {
	userConfig := cmd.Flag("user").Value.String()
	cfg := &struct {
		User  string            `yaml:"user"`
		Email string            `yaml:"email"`
		Key   *auth.RsaIdentity `yaml:"key"`
	}{}

	cfgData, err := ioutil.ReadFile(userConfig)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(cfgData, cfg); err != nil {
		return nil, err
	}

	p, err := auth.NewPriv("", true, true, true, true)
	if err != nil {
		return nil, err
	}

	user := auth.NewUser(cfg.User, cfg.Email,
		cmd.Flag("password").Value.String(),
		cfg.Key.Public(), zebra.Labels{})
	user.Role = &auth.Role{
		Name:       "admin",
		Privileges: []*auth.Priv{p},
	}
	user.Status = nil

	return user, nil
}
