package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/invopop/gobl"
	ksef "github.com/invopop/gobl.ksef"
	ksef_api "github.com/invopop/gobl.ksef/api"
	"github.com/invopop/gobl/bill"
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
		ksef_api.WithDebugClient(),
	)

	env, err := SendInvoice(client, data)
	if err != nil {
		return fmt.Errorf("sending invoices: %w", err)
	}

	data, err = json.MarshalIndent(env, "", "  ")
	if err != nil {
		return err
	}

	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return fmt.Errorf("invalid type %T", env.Document)
	}

	err = saveFile(filename(inv), data)
	if err != nil {
		return err
	}

	return nil
}

// SendInvoice sends invoices to KSeF
func SendInvoice(c *ksef_api.Client, data []byte) (*gobl.Envelope, error) {
	ctx := context.Background()

	err := ksef_api.FetchSessionToken(ctx, c)
	if err != nil {
		return nil, err
	}

	env := new(gobl.Envelope)
	if err := json.Unmarshal(data, env); err != nil {
		return nil, fmt.Errorf("parsing input as GOBL Envelope: %w", err)
	}

	doc, err := ksef.NewDocument(env)
	if err != nil {
		return nil, fmt.Errorf("building FA_VAT document: %w", err)
	}

	data, err = doc.Bytes()
	if err != nil {
		return nil, fmt.Errorf("generating FA_VAT xml: %w", err)
	}

	sendInvoiceResponse, err := ksef_api.SendInvoice(ctx, c, data)
	if err != nil {
		return nil, err
	}

	_, err = waitUntilInvoiceIsProcessed(ctx, c, sendInvoiceResponse.ElementReferenceNumber)
	if err != nil {
		return nil, err
	}

	res, err := waitUntilSessionIsTerminated(ctx, c)
	if err != nil {
		return nil, err
	}
	upoBytes, err := base64.StdEncoding.DecodeString(res.Upo)
	if err != nil {
		return nil, err
	}
	// saveFile(res.ReferenceNumber+".xml", upoBytes)

	err = ksef_api.Sign(env, upoBytes, c)
	if err != nil {
		return nil, err
	}

	return env, nil
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

func saveFile(name string, data []byte) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error when closing:", err)
		}
	}()
	_, err = file.Write(data)
	return err
}

func filename(inv *bill.Invoice) string {
	if inv.Series != "" {
		return sanitizeFilename(inv.Series + "_" + inv.Code + ".xml")
	}
	return sanitizeFilename(inv.Code + ".xml")
}

func sanitizeFilename(filename string) string {
	re := regexp.MustCompile(`[^\w\.-]`)
	sanitized := re.ReplaceAllString(filename, "_")

	return sanitized
}
