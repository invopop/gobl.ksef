package ksef_test

import (
	"testing"

	"github.com/invopop/gobl.ksef/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDocument(t *testing.T) {
	t.Run("should return a Document with KSeF data", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice-pl-pl.json")
		require.NoError(t, err)

		assert.Equal(t, "Faktura", doc.XMLName.Local)
		assert.Equal(t, "http://crd.gov.pl/wzor/2023/06/29/12648/", doc.XMLNamespace)
		assert.Equal(t, "http://www.w3.org/2001/XMLSchema", doc.XSDNamespace)
		assert.Equal(t, "http://www.w3.org/2001/XMLSchema-instance", doc.XSINamespace)
		assert.NotNil(t, doc.Header)
		assert.NotNil(t, doc.Buyer)
		assert.NotNil(t, doc.Seller)
		assert.NotNil(t, doc.Inv)
		assert.NotNil(t, doc.Footer)
	})

	t.Run("should return bytes of the KSeF document", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice-pl-pl.json")
		require.NoError(t, err)

		bytes, bytes_err := doc.Bytes()
		require.NoError(t, bytes_err)

		output, output_err := test.LoadOutputFile("invoice-pl-pl.xml")
		require.NoError(t, output_err)

		assert.Equal(t, output, bytes)
	})
}
