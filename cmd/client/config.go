package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/lease"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const ReadOnly = 0o600

var ErrLeaseDuration = errors.New("lease duration cannot be more that 4 hours")

func NewConfigure() *cobra.Command {
	configCmd := &cobra.Command{ //nolint:exhaustivestruct,exhaustruct
		Use:          "config",
		Short:        "customize zebra client",
		RunE:         showConfig,
		SilenceUsage: true,
	}

	// add config command
	addConfigCommands(configCmd)

	configCmd.AddCommand(&cobra.Command{ //nolint:exhaustivestruct,exhaustruct
		Use:          "public-key",
		Short:        "show current user public key",
		RunE:         showLocalKey,
		SilenceUsage: true,
	})

	defaultCmd := &cobra.Command{ //nolint:exhaustivestruct,exhaustruct
		Use:          "defaults",
		Short:        "configure default values",
		RunE:         configDefaults,
		SilenceUsage: true,
	}
	defaultCmd.Flags().IntP("duration", "t", lease.DefaultMaxDuration, "duration in hours")
	configCmd.AddCommand(defaultCmd)

	return configCmd
}

func addConfigCommands(configCmd *cobra.Command) {
	configCmd.AddCommand(&cobra.Command{ //nolint:exhaustivestruct,exhaustruct
		Use:   "init",
		Short: "initialize default configuration, with server address",
		RunE:  initConfig,
		Args:  cobra.ExactArgs(1),
	})

	configCmd.AddCommand(&cobra.Command{ //nolint:exhaustivestruct,exhaustruct
		Use:          "user",
		Short:        "zebra user name, default: $USER",
		RunE:         configUser,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	})

	configCmd.AddCommand(&cobra.Command{ //nolint:exhaustivestruct,exhaustruct
		Use:          "ca-cert",
		Short:        "zebra CA cert file",
		RunE:         configCACert,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	})

	configCmd.AddCommand(&cobra.Command{ //nolint:exhaustivestruct,exhaustruct
		Use:          "server",
		Short:        "zebra server address",
		RunE:         configServer,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	})

	configCmd.AddCommand(&cobra.Command{ //nolint:exhaustivestruct,exhaustruct
		Use:          "email",
		Short:        "zebra user email address",
		RunE:         configEmail,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	})
}

type ConfigDefaults struct {
	Duration int `yaml:"duration"`
}

type Config struct {
	User          string            `yaml:"user"`
	Email         string            `yaml:"email"`
	ServerAddress string            `yaml:"zebraServer"`
	Key           *auth.RsaIdentity `yaml:"key"`
	CACert        string            `yaml:"caCert"`
	Defaults      ConfigDefaults    `yaml:"defaults,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		User:          os.Getenv("USER"),
		Email:         "",
		ServerAddress: "",
		Key:           nil,
		CACert:        "",
		Defaults: ConfigDefaults{
			Duration: lease.DefaultMaxDuration,
		},
	}
}

func Load(cfgFile string) (*Config, error) {
	data, e := ioutil.ReadFile(cfgFile)
	if e != nil {
		return nil, e
	}

	c := NewConfig()
	e = yaml.Unmarshal(data, c)

	return c, e
}

func (cfg *Config) Save(cfgFile string) error {
	if cfg.Key == nil {
		k, e := auth.Generate()
		if e != nil {
			return e
		}

		cfg.Key = k
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgFile, data, ReadOnly)
}

func showConfig(cmd *cobra.Command, args []string) error {
	cfg := cmd.Flag("config").Value.String()

	return show(cfg)
}

func show(cfg string) error {
	b, e := ioutil.ReadFile(cfg)
	if e != nil {
		return e
	}

	fmt.Println("Zebra Config File: ", cfg)
	fmt.Println(string(b))

	return nil
}

func initConfig(cmd *cobra.Command, args []string) error {
	cfg := cmd.Flag("config").Value.String()
	server := args[0]
	c := NewConfig()
	c.ServerAddress = server

	if e := c.Save(cfg); e != nil {
		return e
	}

	return show(cfg)
}

func configUser(cmd *cobra.Command, args []string) error {
	cfgFile := cmd.Flag("config").Value.String()
	user := args[0]

	cfg, e := Load(cfgFile)
	if e != nil {
		return e
	}

	cfg.User = user
	if e := cfg.Save(cfgFile); e != nil {
		return e
	}

	return show(cfgFile)
}

func configCACert(cmd *cobra.Command, args []string) error {
	cfgFile := cmd.Flag("config").Value.String()
	caCert := args[0]

	cfg, e := Load(cfgFile)
	if e != nil {
		return e
	}

	cfg.CACert = caCert
	if e := cfg.Save(cfgFile); e != nil {
		return e
	}

	return show(cfgFile)
}

func configEmail(cmd *cobra.Command, args []string) error {
	cfgFile := cmd.Flag("config").Value.String()
	email := args[0]

	cfg, e := Load(cfgFile)
	if e != nil {
		return e
	}

	cfg.Email = email
	if e := cfg.Save(cfgFile); e != nil {
		return e
	}

	return show(cfgFile)
}

func configServer(cmd *cobra.Command, args []string) error {
	cfgFile := cmd.Flag("config").Value.String()
	server := args[0]

	cfg, e := Load(cfgFile)
	if e != nil {
		return e
	}

	cfg.ServerAddress = server
	if e := cfg.Save(cfgFile); e != nil {
		return e
	}

	return show(cfgFile)
}

func configDefaults(cmd *cobra.Command, args []string) error {
	cfgFile := cmd.Flag("config").Value.String()

	duration, e := cmd.Flags().GetInt("duration")
	if e != nil {
		return e
	}

	if duration > lease.DefaultMaxDuration {
		return ErrLeaseDuration
	}

	cfg, e := Load(cfgFile)
	if e != nil {
		return e
	}

	cfg.Defaults.Duration = duration
	if e := cfg.Save(cfgFile); e != nil {
		return e
	}

	return show(cfgFile)
}

func showLocalKey(cmd *cobra.Command, args []string) error {
	cfgFile := cmd.Flag("config").Value.String()
	cfg, e := Load(cfgFile)

	if e != nil {
		return e
	}

	pk := cfg.Key.Public()
	b, e := pk.MarshalText()

	if e != nil {
		return e
	}

	fmt.Println(string(b))

	return nil
}
