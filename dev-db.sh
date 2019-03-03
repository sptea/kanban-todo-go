echo -e '\n\n####### Create Table #######\n'
go run sql/sqlite-migrate/main.go database-dev.db sql/migration
echo -e '\n\n####### Insert test data #######\n'
go run sql/sqlite-migrate/main.go database-dev.db sql/test-data
