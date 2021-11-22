default: build async

perform:
	go run main.go -async=true

select:
	sqlite3 audit.db "SELECT SUM(error) AS "errors", COUNT(*) "total" FROM pdf;"

clear:
	sqlite3 audit.db "delete from pdf;"

build:
	go build -o audit main.go

sync: build
	./audit -async=false 2> errors.log

async: build
	./audit -async=true 2> errors.log

db:
	sqlite3 audit.db "create table pdf (id text primary key, code text, error integer, plain_text text);"

setup: db
	cp input.json.sample input.json
