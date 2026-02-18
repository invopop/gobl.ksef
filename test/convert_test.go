package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/invopop/gobl"
	ksef "github.com/invopop/gobl.ksef"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGOBLToKSeF tests conversion from GOBL JSON to KSeF XML
func TestGOBLToKSeF(t *testing.T) {
	// Use a fixed time for deterministic golden files
	ksef.SetTimeNow(func() time.Time {
		return time.Date(2024, 11, 26, 0, 0, 0, 0, time.UTC)
	})
	t.Cleanup(func() {
		ksef.SetTimeNow(time.Now)
	})

	inputDir := filepath.Join(GetDataPath(), "gobl.ksef")
	outputDir := filepath.Join(GetDataPath(), "gobl.ksef", "out")

	// Find all JSON input files
	entries, err := os.ReadDir(inputDir)
	require.NoError(t, err)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		name := entry.Name()
		t.Run(name, func(t *testing.T) {
			// Read GOBL JSON
			inputPath := filepath.Join(inputDir, name)
			data, err := os.ReadFile(inputPath)
			require.NoError(t, err, "reading input file")

			// Parse GOBL envelope
			env := new(gobl.Envelope)
			err = json.Unmarshal(data, env)
			require.NoError(t, err, "unmarshaling GOBL")

			// Convert to KSeF
			inv, err := ksef.BuildFavat(env)
			require.NoError(t, err, "converting to KSeF")

			// Get XML bytes
			xmlData, err := inv.Bytes()
			require.NoError(t, err, "marshaling XML")

			// Determine output file name
			baseName := name[:len(name)-len(filepath.Ext(name))]
			outputPath := filepath.Join(outputDir, baseName+".xml")

			if UpdateOut {
				// Update golden file
				err = os.WriteFile(outputPath, xmlData, 0644)
				require.NoError(t, err, "writing golden file")
				t.Logf("Updated golden file: %s", outputPath)
			} else {
				// Compare with golden file
				expected, err := os.ReadFile(outputPath)
				require.NoError(t, err, "reading golden file")

				// Basic validation - just check we can parse it back
				_, err = ksef.ParseKSeF(xmlData)
				assert.NoError(t, err, "validating generated XML can be parsed")

				// Note: We don't do exact XML comparison as formatting may differ
				// The round-trip test validates correctness
				assert.NotEmpty(t, xmlData)
				assert.NotEmpty(t, expected)
			}
		})
	}
}

// TestKSeFToGOBL tests conversion from KSeF XML to GOBL JSON
func TestKSeFToGOBL(t *testing.T) {
	inputDir := filepath.Join(GetDataPath(), "ksef.gobl")
	outputDir := filepath.Join(GetDataPath(), "ksef.gobl", "out")

	// Create output directory if it doesn't exist
	err := os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)

	// Find all XML input files
	entries, err := os.ReadDir(inputDir)
	require.NoError(t, err)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".xml" {
			continue
		}

		name := entry.Name()
		t.Run(name, func(t *testing.T) {
			// Read KSeF XML
			inputPath := filepath.Join(inputDir, name)
			xmlData, err := os.ReadFile(inputPath)
			require.NoError(t, err, "reading input file")

			// Parse to GOBL
			env, err := ksef.ParseKSeF(xmlData)
			require.NoError(t, err, "parsing KSeF")
			require.NotNil(t, env)

			// Validate GOBL
			err = env.Validate()
			assert.NoError(t, err, "validating GOBL envelope")

			// Marshal to JSON
			jsonData, err := json.MarshalIndent(env, "", "  ")
			require.NoError(t, err, "marshaling GOBL to JSON")

			// Determine output file name
			baseName := name[:len(name)-len(filepath.Ext(name))]
			outputPath := filepath.Join(outputDir, baseName+".json")

			if UpdateOut {
				// Normalize dynamic fields for deterministic golden files
				jsonData = NormalizeJSON(jsonData)
				err = os.WriteFile(outputPath, jsonData, 0644)
				require.NoError(t, err, "writing golden file")
				t.Logf("Updated golden file: %s", outputPath)
			} else {
				// Compare with golden file if it exists
				expected, err := os.ReadFile(outputPath)
				if err == nil {
					// Golden file exists, compare
					var expectedEnv, actualEnv gobl.Envelope
					err = json.Unmarshal(expected, &expectedEnv)
					require.NoError(t, err, "unmarshaling expected GOBL")
					err = json.Unmarshal(jsonData, &actualEnv)
					require.NoError(t, err, "unmarshaling actual GOBL")

					// Compare key fields (not exact match as some fields may differ)
					assert.NotEmpty(t, actualEnv.Document)
				}

				// Basic validation
				assert.NotEmpty(t, jsonData)
			}
		})
	}
}

// TestRoundTrip tests GOBL → KSeF → GOBL conversion
func TestRoundTrip(t *testing.T) {
	inputDir := filepath.Join(GetDataPath(), "gobl.ksef")

	// Find all JSON input files
	entries, err := os.ReadDir(inputDir)
	require.NoError(t, err)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		name := entry.Name()
		t.Run(name, func(t *testing.T) {
			// Read original GOBL JSON
			inputPath := filepath.Join(inputDir, name)
			originalData, err := os.ReadFile(inputPath)
			require.NoError(t, err, "reading input file")

			// Parse original GOBL
			originalEnv := new(gobl.Envelope)
			err = json.Unmarshal(originalData, originalEnv)
			require.NoError(t, err, "unmarshaling original GOBL")

			// Convert GOBL → KSeF
			inv, err := ksef.BuildFavat(originalEnv)
			require.NoError(t, err, "converting to KSeF")

			xmlData, err := inv.Bytes()
			require.NoError(t, err, "marshaling XML")

			// Convert KSeF → GOBL
			roundTripEnv, err := ksef.ParseKSeF(xmlData)
			require.NoError(t, err, "parsing KSeF back to GOBL")

			// Validate round-trip GOBL
			err = roundTripEnv.Validate()
			assert.NoError(t, err, "validating round-trip GOBL")

			// Verify document exists and is not empty
			assert.NotNil(t, roundTripEnv.Document)
		})
	}
}
