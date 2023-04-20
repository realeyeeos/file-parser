module fileparser

go 1.16

require (
	baliance.com/gooxml v0.0.0-00010101000000-000000000000
	github.com/BurntSushi/toml v1.0.0 // indirect
	github.com/beevik/etree v1.1.0
	github.com/gogf/gf v1.16.9
	github.com/henrylee2cn/pholcus v1.3.4
	github.com/ledongthuc/pdf v0.0.0-20220302134840-0c2507a12d80
	github.com/saintfish/chardet v0.0.0-20230101081208-5e3ef4b5456d
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/text v0.8.0
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/bytedance/godlp => ./godlp

replace baliance.com/gooxml => ./gooxml
