module github.com/zorneth/osg-providers

go 1.27.0

require (
	github.com/zorneth/osg-core v0.1.0-alpha.1
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/kr/text v0.2.0 // indirect

replace github.com/zorneth/osg-core => ../osg-core
