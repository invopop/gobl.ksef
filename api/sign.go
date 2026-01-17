package api

import (
	"net/url"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/regimes/pl"
)

// Sign attached QR code and other identification values to the envelope
func (c *Client) Sign(env *gobl.Envelope, nip string, uploadedInvoice *UploadedInvoice) error {
	env.Head.AddStamp(
		&head.Stamp{
			Provider: pl.StampProviderKSeFID,
			Value:    uploadedInvoice.KsefNumber,
		},
	)
	env.Head.AddStamp(
		&head.Stamp{
			Provider: pl.StampProviderKSeFHash,
			Value:    uploadedInvoice.InvoiceHash,
		},
	)
	// URL contains invoicing date in DD-MM-YYYY format
	env.Head.AddStamp(
		&head.Stamp{
			Provider: pl.StampProviderKSeFQR,
			Value:    c.qrUrl + "/" + nip + "/" + uploadedInvoice.InvoicingDate.Format("02-01-2006") + "/" + url.QueryEscape(uploadedInvoice.InvoiceHash),
		},
	)

	return nil
}
