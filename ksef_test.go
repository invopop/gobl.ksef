package ksef_test

import (
	"testing"

	"github.com/invopop/gobl.ksef/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	xsdvalidate "github.com/terminalstatic/go-xsd-validate"
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
	})

	t.Run("should return bytes of the KSeF document", func(t *testing.T) {
		doc, err := test.NewDocumentFrom("invoice-pl-pl.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		output, err := test.LoadOutputFile("invoice-pl-pl.xml")
		require.NoError(t, err)

		assert.Equal(t, output, data)
	})

	t.Run("should generate valid KSeF document", func(t *testing.T) {
		err := xsdvalidate.Init()
		require.NoError(t, err)
		defer xsdvalidate.Cleanup()

		xsdBuf, err := test.LoadSchemaFile("FA2.xsd")
		require.NoError(t, err)

		xsdhandler, err := xsdvalidate.NewXsdHandlerMem(xsdBuf, xsdvalidate.ParsErrVerbose)
		require.NoError(t, err)
		defer xsdhandler.Free()

		doc, err := test.NewDocumentFrom("invoice-pl-pl.json")
		require.NoError(t, err)

		data, err := doc.Bytes()
		require.NoError(t, err)

		validation := xsdhandler.ValidateMem(data, xsdvalidate.ParsErrDefault)
		assert.Nil(t, validation)
	})
}
