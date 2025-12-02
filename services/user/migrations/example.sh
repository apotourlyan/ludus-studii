# create up/down sql migration files
migrate create "create_users_table"
# migrate to latest version
migrate up "postgres://ludus:ludus@db:5432/ludus_studii_user_db?sslmode=disable"
# migrate to specific version
migrate up "postgres://ludus:ludus@db:5432/ludus_studii_user_db?sslmode=disable" 20251124050929_create_users_table
# rollback to specific version
migrate down "postgres://ludus:ludus@db:5432/ludus_studii_user_db?sslmode=disable" 20251124050929_create_users_table