echo -e '\n####### Remove Database #######\n'
rm database-dev.db
echo -e '\n####### Create Table #######\n'
go run sql/sqlite-migrate/main.go database-dev.db sql/migration
echo -e '\n####### Insert test data #######\n'
go run sql/sqlite-migrate/main.go database-dev.db sql/test-data
