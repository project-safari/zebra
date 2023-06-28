package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/project-safari/zebra/model/lease"
	"github.com/spf13/cobra"
)

var ErrCreateLease = errors.New("error creating resource")

const (
	DefaultResourceCount = 3
	DefaultDuration      = 240
)

func NewLease() *cobra.Command {
	leaseCmd := &cobra.Command{
		Use:          "lease",
		Short:        "request a lease",
		RunE:         leaseRequest,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
	}

	leaseCmd.Flags().StringP("group", "g", "global", "resource group")
	leaseCmd.Flags().IntP("count", "k", DefaultResourceCount, "number of resources")
	leaseCmd.Flags().IntP("duration", "d", DefaultDuration, "length of lease(in minutes)")
	leaseCmd.Flags().StringP("name", "n", "", "name of lease")

	return leaseCmd
}

func leaseRequest(cmd *cobra.Command, args []string) error {
	cfg, req, err := makeLeaseReq(cmd, args)
	if err != nil {
		return err
	}

	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	dur, err := cmd.Flags().GetInt("duration")
	if err != nil {
		return err
	}

	leasereq := &struct {
		Email    string               `json:"email"`
		Duration time.Duration        `json:"duration"`
		Request  []*lease.ResourceReq `json:"request"`
	}{
		Email:    cfg.Email,
		Duration: time.Duration(dur) * time.Minute,
		Request:  []*lease.ResourceReq{req},
	}

	// Create a new lease
	resCode, err := client.Post("/api/v1/lease", leasereq, nil)
	if resCode != http.StatusOK {
		return ErrCreateLease
	}

	fmt.Println("Request - Type:", req.Type, "Group:", req.Group, "Count:", req.Count)
	fmt.Println("Lease request successfully created")

	return err
}

func makeLeaseReq(cmd *cobra.Command, args []string) (*Config, *lease.ResourceReq, error) {
	cfgFile := cmd.Flag("config").Value.String()

	cfg, err := Load(cfgFile)
	if err != nil {
		return nil, nil, err
	}

	resCount, err := cmd.Flags().GetInt("count")
	if err != nil {
		return nil, nil, err
	}

	req := &lease.ResourceReq{
		Type:  args[0],
		Group: cmd.Flag("group").Value.String(),
		Count: resCount,
		Name:  cmd.Flag("name").Value.String(),
	}

	return cfg, req, nil
}
