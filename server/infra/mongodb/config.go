package mongodb

type ClientOptions struct {
	Host        string
	Port        string
	User        string
	Password    string
	RetryWrites bool
}

type DatabaseConfig struct {
	DatabaseName   string
	CollectionName string
}

type MongoDBConfig struct {
	Client   ClientOptions
	DBConfig DatabaseConfig
}
