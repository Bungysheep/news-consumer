package configs

const (
	// DBCONNSTRING - Database connection string
	DBCONNSTRING = "postgres://news-local-pg:news-local-pg@localhost:5432/news-local-pg?sslmode=disable"

	// ELASTICSEARCHURL - Elasticsearch url
	ELASTICSEARCHURL = "http://localhost:9200/"

	// REDISURL - Redis url
	REDISURL = "localhost:6379"

	// REDISAUTH - Redis url
	REDISAUTH = ""

	// REDISNEWSPOSTCHANNEL - Redis news post channel
	REDISNEWSPOSTCHANNEL = "NEWS_POST_CHANNEL"
)
