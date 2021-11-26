default: build async

perform:
	go run main.go

select:
	sqlite3 audit.db "SELECT SUM(error) AS "errors", COUNT(*) "total" FROM pdf;"

clear:
	sqlite3 audit.db "delete from pdf;"

build:
	go build -o audit main.go

sync: build
	./audit -method=sync 2> errors.log

async: build
	./audit -method=async 2> errors.log

seller: build
	./audit -method=seller 2> errors.log

db:
	sqlite3 audit.db "create table pdf (id text primary key, code text, error integer, plain_text text);"

setup: db
	cp input.json.sample input.json

export:
	sqlite3 audit.db "SELECT id, code, COALESCE(NULLIF(TRIM(phone, ' '), ''), '(00) 000000000'), \"order\", seller, buyer FROM pdf WHERE error > 0;" > id.txt
