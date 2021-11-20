.PHONY: perform

perform:
	go run main.go

select:
	sqlite3 audit.db "SELECT SUM(error) AS "errors", COUNT(*) "total" FROM pdf;"

clear:
	sqlite3 audit.db "delete from pdf;"
