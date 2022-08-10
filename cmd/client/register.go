package main

import (
	"net/http"

	"github.com/project-safari/zebra/auth"
	"github.com/spf13/cobra"
)

func NewRegistration() *cobra.Command {
	registerCmd := &cobra.Command{ //nolint:exhaustivestruct,exhaustruct
		Use:          "registration",
		Short:        "register for zebra",
		RunE:         registerReq,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
	}

	registerCmd.Flags().String("password", "", "Create passowrd for zebra")

	return registerCmd
}

func Approve() *cobra.Command {
	approvecmd := &cobra.Command{
		Use:          "approve",
		Short:        "Approve user to use zebra",
		RunE:         approveReq,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
	}

	approvecmd.Flags().String("user", "", "User that as been approved for zebra")

	return approvecmd
}

func Reject() *cobra.Command {
	rejectcmd := &cobra.Command{
		Use:          "reject",
		Short:        "Reject user to use zebra",
		RunE:         rejectReq,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
	}

	rejectcmd.Flags().String("user", "", "User that as been rejected for zebra")

	return rejectcmd
}

func registerReq(cmd *cobra.Command, arg []string) error {
	cfg, client, err := createClient(cmd)
	if err != nil {
		return err
	}

	password, _ := cmd.Flags().GetString("password")

	user := auth.UserType().Constructor()
	reqBody := &struct {
		Name     string            `json:"name"`
		Password string            `json:"password"`
		Email    string            `json:"email"`
		Key      *auth.RsaIdentity `json:"key"`
	}{
		Name:     cfg.User,
		Password: password,
		Email:    cfg.Email,
		Key:      cfg.Key,
	}

	resCode, err := client.Post("api/v1/resources", reqBody, user)
	if resCode != http.StatusOK {
		return err
	}

	return nil
}

func approveReq(cmd *cobra.Command, arg []string) error {
	_, client, err := createClient(cmd)
	if err != nil {
		return err
	}

	user, _ := cmd.Flags().GetString("user")

	resCode, err := client.Post("api/v1/resources", user, nil)
	if resCode != http.StatusOK {
		return err
	}

	return nil
}

func rejectReq(cmd *cobra.Command, arg []string) error {
	_, client, err := createClient(cmd)
	if err != nil {
		return err
	}

	user, _ := cmd.Flags().GetString("user")

	resCode, err := client.Post("/register", user, nil)
	if resCode != http.StatusOK {
		return err
	}

	return nil
}

func createClient(cmd *cobra.Command) (*Config, *Client, error) {
	cfgFile := cmd.Flag("config").Value.String()

	cfg, err := Load(cfgFile)
	if err != nil {
		return nil, nil, err
	}

	client, err := NewClient(cfg)
	if err != nil {
		return nil, nil, err
	}

	return cfg, client, nil
}
