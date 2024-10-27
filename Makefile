run:
	go run cmd/fast/main.go

b:
	go build -o ./build/fast cmd/fast/main.go

tidy:
	go mod tidy

entf:
	go run -mod=mod entgo.io/ent/cmd/ent new $(f)

entg:
	go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema

entd:
	go run -mod=mod entgo.io/ent/cmd/ent describe ./ent/schema

entm:
	atlas migrate diff \
	--dir "file://ent/migrate/migrations" \
	--to "ent://ent/schema?globalindex=1" \
	--dev-url "postgres://xpriori:phtn458@localhost:5432/dpqb?search_path=public&sslmode=disable"

entma:
	atlas migrate apply \
	--dir "file://ent/migrate/migrations" \
	--url "postgres://xpriori:phtn458@localhost:5432/dpqb?search_path=public&sslmode=disable"

entms:
	atlas migrate status \
	--dir "file://ent/migrate/migrations" \
	--url "postgres://xpriori:phtn458@localhost:5432/dpqb?search_path=public&sslmode=disable"

pqd:
	psql postgres://xpriori:phtn458@localhost:5432/dpqb?sslmode=disable

atlasg:
	atlas migrate new $(n) \
	--dir "file://ent/migrate/migrations" \


atlasp:
	atlas migrate push $(n) \
	--dir "file://ent/migrate" \
	--dev-url "postgres://xpriori:phtn458@localhost:5432/dpqb?search_path=public&sslmode=disable"

atlasi:
	atlas schema inspect -u "postgres://xpriori:phtn458@localhost:5432/dpqb?search_path=public&sslmode=disable"

atlasl:
	atlas migrate lint \
	--dev-url "postgres://xpriori:phtn458@localhost:5432/dpqb?search_path=public&sslmode=disable" \
	--latest 1 \
	-w


clean:
	rm -rf ./build/fast
