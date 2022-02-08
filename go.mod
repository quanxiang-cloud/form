module github.com/quanxiang-cloud/form

go 1.16

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis/v8 v8.11.4
	github.com/quanxiang-cloud/cabin v0.0.4
	github.com/quanxiang-cloud/structor v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.8.1
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/gorm v1.22.5
)

replace github.com/quanxiang-cloud/structor => ../structor
