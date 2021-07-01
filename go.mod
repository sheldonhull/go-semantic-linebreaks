module github.com/sheldonhull/go-semantic-sentences

go 1.16

require (
	github.com/matryer/is v1.4.0
	github.com/peterbourgon/ff/v3 v3.0.0
	github.com/rs/zerolog v1.23.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	github.com/sheldonhull/go-semantic-sentences/internal/logger v0.0.0
)

replace github.com/sheldonhull/go-semantic-sentences/internal/logger v0.0.0 => ./internal/logger
