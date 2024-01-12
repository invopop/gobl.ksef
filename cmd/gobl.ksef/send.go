package main

import (
	"fmt"
	"io"

	ksef_api "github.com/invopop/gobl.ksef/api"
	"github.com/spf13/cobra"
)

type sendOpts struct {
	*rootOpts
}

func send(o *rootOpts) *sendOpts {
	return &sendOpts{rootOpts: o}
}

func (c *sendOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [infile] [nip] [token] [keyPath]",
		Short: "Send a GOBL JSON to the KSeF API",
		RunE:  c.runE,
	}

	return cmd
}

func (c *sendOpts) runE(cmd *cobra.Command, args []string) error {
	// ctx := commandContext(cmd)
	token := inputToken(args)
	keyPath := inputKeyPath(args)
	nip := inputNip(args)

	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	inData, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}
	env := ksef_api.KSeFEnv{Url: ksef_api.KSeFTestingBaseURL, KeyPath: keyPath}

	_, err = ksef_api.SendInvoices(env, nip, token, []string{string(inData)})
	if err != nil {
		return fmt.Errorf("sending invoices: %w", err)
	}
	return nil
}
