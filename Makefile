.PHONY: perform

perform:
	go run main.go -async=true

select:
	sqlite3 audit.db "SELECT SUM(error) AS "errors", COUNT(*) "total" FROM pdf;"

clear:
	sqlite3 audit.db "delete from pdf;"

build:
	go build -o audit main.go
