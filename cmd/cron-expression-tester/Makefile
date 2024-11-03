
runeverysecond:
	@go run -mod=vendor main.go \
	-cron-expression='* * * * * ?' \
	-v=2

runsundayonly:
	@go run -mod=vendor main.go \
	-cron-expression='* * * * * 0' \
	-v=2

runevery1s:
	@go run -mod=vendor main.go \
	-cron-expression='@every 1s' \
	-v=2

runevery1m:
	@go run -mod=vendor main.go \
	-cron-expression='@every 1m' \
	-v=2

runevery1h:
	@go run -mod=vendor main.go \
	-cron-expression='@every 1h' \
	-v=2
