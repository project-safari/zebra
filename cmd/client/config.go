package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const ReadOnly = 0o600

// ErrLeaseDuration is an error if the lease duration is grater than 4.
var ErrLeaseDuration = errors.New("lease duration cannot be more that 4 hours")

func NewConfigure() *cobra.Command {
	configCmd := &cobra.Command{
		Use:          "config",
		Short:        "customize zebra client",
		RunE:         showConfig,
		SilenceUsage: true,
	}

	// Add config command.
	addConfigCommands(configCmd)

	configCmd.AddCommand(&cobra.Command{
		Use:          "public-key",
		Short:        "show current user public key",
		RunE:         showLocalKey,
		SilenceUsage: true,
	})

	defaultCmd := &cobra.Command{
		Use:          "defaults",
		Short:        "configure default values",
		RunE:         configDefaults,
		SilenceUsage: true,
	}
	defaultCmd.Flags().IntP("duration", "t", zebra.DefaultMaxDuration, "duration in hours")
	configCmd.AddCommand(defaultCmd)

	return configCmd
}

// Configuration commands.
func addConfigCommands(configCmd *cobra.Command) {
	configCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "initialize default configuration, with server address",
		RunE:  initConfig,
		Args:  cobra.ExactArgs(1),
	})

	configCmd.AddCommand(&cobra.Command{
		Use:          "user",
		Short:        "zebra user name, default: $USER",
		RunE:         configUser,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	})

	configCmd.AddCommand(&cobra.Command{
		Use:          "ca-cert",
		Short:        "zebra CA cert file",
		RunE:         configCACert,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	})

	configCmd.AddCommand(&cobra.Command{
		Use:          "server",
		Short:        "zebra server address",
		RunE:         configServer,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	})

	configCmd.AddCommand(&cobra.Command{
		Use:          "email",
		Short:        "zebra user email address",
		RunE:         configEmail,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	})
}

// Default configuration for duration.
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

// Function to create new configuration.
func NewConfig() *Config {
	return &Config{
		User:          os.Getenv("USER"),
		Email:         "",
		ServerAddress: "",
		Key:           nil,
		CACert:        "",
		Defaults: ConfigDefaults{
			Duration: zebra.DefaultMaxDuration,
		},
	}
}

// Function to load config file, returns the config and a(n) (potential) error.
func Load(cfgFile string) (*Config, error) {
	data, e := ioutil.ReadFile(cfgFile)
	if e != nil {
		return nil, e
	}

	c := NewConfig()
	e = yaml.Unmarshal(data, c)

	return c, e
}

// Function to save the config to a file, returns the file with the config.
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

// Function to show the config.
func showConfig(cmd *cobra.Command, args []string) error {
	cfg := cmd.Flag("config").Value.String()

	return show(cfg)
}

// Function to show the config file.
func show(cfg string) error {
	b, e := ioutil.ReadFile(cfg)
	if e != nil {
		return e
	}

	fmt.Println("Zebra Config File: ", cfg)
	fmt.Println(string(b))

	return nil
}

// Initialize configuration.
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

// Configuration for users.
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

// Configuration with certificate authority for files that will help with keys.
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

// Configuration for email.
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

// Configuration for server.
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

// Configuration for default values (duration isnthe only default).
// It ensures that the duration is no morethan 4 hours.
func configDefaults(cmd *cobra.Command, args []string) error {
	cfgFile := cmd.Flag("config").Value.String()

	duration, e := cmd.Flags().GetInt("duration")
	if e != nil {
		return e
	}

	if duration > zebra.DefaultMaxDuration {
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

// Function to show local public key as string.
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
