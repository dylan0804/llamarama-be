run:
	go run cmd/web/main.go

migrateup:
	migrate -path cmd/internal/db/migrations -database postgresql://postgres.srcvvvmpnoyvhbizjywu:golangchatroom@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres up

migratedown:
	migrate -path cmd/internal/db/migrations -database postgresql://postgres.srcvvvmpnoyvhbizjywu:golangchatroom@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres down