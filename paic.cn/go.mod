module paic.cn

go 1.15

require (
	github.com/hashicorp/consul/api v1.8.1
	paic.cn/source v1.0.0
)

replace (
	github.com/Sirupsen/logrus v1.7.0 => github.com/sirupsen/logrus v1.7.0
	github.com/sirupsen/logrus v1.7.0 => github.com/Sirupsen/logrus v1.7.0
	paic.cn/source => ../paic.cn/source
)
