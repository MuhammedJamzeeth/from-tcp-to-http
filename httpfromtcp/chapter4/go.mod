module chapter4

go 1.26.4

require (
	chapter5 v0.0.0
	github.com/stretchr/testify v1.12.1
)

replace chapter5 => ../chapter5

require go.yaml.in/yaml/v3 v3.0.5 // indirect
