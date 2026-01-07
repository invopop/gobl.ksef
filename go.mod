module github.com/invopop/gobl.ksef

go 1.24

toolchain go1.24.3

require (
	github.com/artemkunich/goxades v0.2.1
	github.com/beevik/etree v1.5.1
	github.com/go-resty/resty/v2 v2.16.5
	github.com/invopop/gobl v0.218.0
	github.com/invopop/xmldsig v0.11.0
	github.com/jarcoal/httpmock v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/russellhaering/goxmldsig v1.5.0
	github.com/spf13/cobra v1.9.1
	github.com/terminalstatic/go-xsd-validate v0.1.5
	software.sslmate.com/src/go-pkcs12 v0.5.0
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-jose/go-jose/v4 v4.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/invopop/yaml v0.3.1 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	cloud.google.com/go v0.121.3
	github.com/Masterminds/semver/v3 v3.3.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/invopop/jsonschema v0.13.0 // indirect
	github.com/invopop/validation v0.8.0 // indirect
	github.com/stretchr/testify v1.10.0
)

replace github.com/artemkunich/goxades => github.com/MieszkoGulinski/goxades v1.0.5-ksef

// replace github.com/invopop/gobl => ../gobl
