brew services start postgresql@14
go run main.go 
migrate create -ext sql -dir database/migrations -seq create_board_member_table migrate -path database/migrations -database "postgres://postgres:superrahasia123@localhost:5432/management_app?sslmode=disable" up
migrate -path database/migrations -database "postgres://postgres:superrahasia123@localhost:5432/management_app?sslmode=disable" up
