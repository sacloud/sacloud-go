module github.com/sacloud/sacloud-go

go 1.17

require (
	github.com/go-playground/validator/v10 v10.10.1
	github.com/sacloud/ftps v1.1.0
	github.com/sacloud/iaas-api-go v0.0.0-20220314063652-5eaa6e6cade6
	github.com/sacloud/sacloud-go/pkg v0.0.0-20220314055142-1db1c3d10889
	github.com/stretchr/testify v1.7.0
)

replace github.com/sacloud/sacloud-go/pkg => ./pkg

require (
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sacloud/api-client-go v0.0.0-20220311054319-f37467272e84 // indirect
	github.com/sacloud/go-http v0.0.4 // indirect
	go.uber.org/ratelimit v0.2.0 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sys v0.0.0-20210806184541-e5e7981a1069 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
