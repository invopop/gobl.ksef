package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

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
	nip := inputNip(args)
	token := inputToken(args)
	keyPath := inputKeyPath(args)

	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer func() {
		err = input.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	data, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	client := ksef_api.NewClient(
		ksef_api.WithID(nip),
		ksef_api.WithToken(token),
		ksef_api.WithKeyPath(keyPath),
	)

	_, err = SendInvoice(client, data)
	if err != nil {
		return fmt.Errorf("sending invoices: %w", err)
	}
	return nil
}

// SendInvoice sends invoices to KSeF
func SendInvoice(c *ksef_api.Client, data []byte) (string, error) {
	ctx := context.Background()

	err := ksef_api.FetchSessionToken(ctx, c)
	if err != nil {
		return "", err
	}

	sendInvoiceResponse, err := ksef_api.SendInvoice(ctx, c, data)
	if err != nil {
		return "", err
	}

	_, err = waitUntilInvoiceIsProcessed(ctx, c, sendInvoiceResponse.ElementReferenceNumber)
	if err != nil {
		return "", err
	}

	res, err := waitUntilSessionIsTerminated(ctx, c)
	if err != nil {
		return "", err
	}
	upoBytes, err := base64.StdEncoding.DecodeString(res.Upo)
	if err != nil {
		return "", err
	}
	file, err := os.Create(res.ReferenceNumber + ".xml")
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error when closing:", err)
		}
	}()
	_, err = file.Write(upoBytes)
	if err != nil {
		return "", err
	}

	return string(upoBytes), nil
}

func waitUntilInvoiceIsProcessed(ctx context.Context, c *ksef_api.Client, referenceNumber string) (*ksef_api.InvoiceStatusResponse, error) {
	for {
		status, err := ksef_api.FetchInvoiceStatus(ctx, c, referenceNumber)
		if err != nil {
			return nil, err
		}
		if status.ProcessingCode == 200 || status.ProcessingCode >= 400 {
			return status, nil
		}
		sleepContext(ctx, 5*time.Second)
	}
}

func waitUntilSessionIsTerminated(ctx context.Context, c *ksef_api.Client) (*ksef_api.SessionStatusByReferenceResponse, error) {
	_, err := ksef_api.TerminateSession(ctx, c)
	if err != nil {
		return nil, err
	}
	for {
		status, err := ksef_api.GetSessionStatusByReference(ctx, c)

		if err != nil {
			return nil, err
		}
		if status.ProcessingCode == 200 || status.ProcessingCode >= 400 {
			return status, nil
		}
		sleepContext(ctx, 5*time.Second)
	}
}

func sleepContext(ctx context.Context, delay time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(delay):
	}
}
