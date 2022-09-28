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

// ErrLeaseDuration occurs if the lease is requested for more than 4 hours.
var ErrLeaseDuration = errors.New("lease duration cannot be more that 4 hours")

// Function that initiates a new configuration command.
//
// Function uses cobra commands and returs a pointer to a cobra.Command.
//
// Contains commands for configuration of: public key and defaults.
func NewConfigure() *cobra.Command {
	configCmd := &cobra.Command{
		Use:          "config",
		Short:        "customize zebra client",
		RunE:         showConfig,
		SilenceUsage: true,
	}

	// add config command
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

// Function for init commands.
//
// Function uses cobra commands and returs a pointer to a cobra.Command.
//
// Contains commands for initialization of: user, CACertificate, server.
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

type ConfigDefaults struct {
	Duration int `yaml:"duration"`
}

// Config struct for configuration of user details.
// It contains the user, the email, the addr of the server, the key,
// the key, the CA certificate, and any config defaults.
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
			Duration: zebra.DefaultMaxDuration,
		},
	}
}

// Function to load the file and initiate a new config.
//
// It takes in cfgFile, a string with the path to the file,
// returns a pointer to the Config and an error or nil in the absence thereof.
func Load(cfgFile string) (*Config, error) {
	data, e := ioutil.ReadFile(cfgFile)
	if e != nil {
		return nil, e
	}

	c := NewConfig()
	e = yaml.Unmarshal(data, c)

	return c, e
}

// Function to save config info into a file.
//
// It executes a pointer to Config a struct with user details.
//
// It takes in cfgFile, a string with the path to the file,
// where the info will be saved, and
// returns an error or nil in the absence thereof.
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

// This shows the config by implementing the show function,
// it adds a flag to a given cobra command and returns an error or nil in the absence thereof.
func showConfig(cmd *cobra.Command, args []string) error {
	cfg := cmd.Flag("config").Value.String()

	return show(cfg)
}

// Function that is implememted in showConfig and initConfig.
//
// It takes in a string containing the path to a config file to print a byte array
// and returns an error or nil in the absence thereof.
func show(cfg string) error {
	b, e := ioutil.ReadFile(cfg)
	if e != nil {
		return e
	}

	fmt.Println("Zebra Config File: ", cfg)
	fmt.Println(string(b))

	return nil
}

// This function initializes the configuration, given a pointer to a cobra command
// and returns an error or nil in the absence thereof.
//
// It utilizes the show function.
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

// This function provides the configuration for the user,
// given a pointer to a cobra command
// and returns an error or nil in the absence thereof.
//
// It loads the config file, saves the config info to the file,
// and utilizes the show function.
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

// This function provides the configuration for the CA Cert,
// given a pointer to a cobra command
// and returns an error or nil in the absence thereof.
//
// It loads the config file, saves the config info to the file,
// and utilizes the show function.
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

// This function provides the configuration for the user's email,
// given a pointer to a cobra command
// and returns an error or nil in the absence thereof.
//
// It loads the config file, saves the config info to the file,
// and utilizes the show function.
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

// This function provides the configuration for the server,
// given a pointer to a cobra command
// and returns an error or nil in the absence thereof.
//
// It loads the config file, saves the config info to the file,
// and utilizes the show function.
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

// This function provides the configuration for any default values,
// given a pointer to a cobra command
// and returns an error or nil in the absence thereof.
//
// It loads the config file, saves the config info to the file,
// and utilizes the show function.
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

// Function for the local key,
// given a pointer to a cobra command
// and returns an error or nil in the absence thereof.
//
// It gets the public key, marshals it, and prints the result.
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
