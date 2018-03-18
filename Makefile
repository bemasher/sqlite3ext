all: go.dll

goext.a goext.h: goext.go
	CGO_ENABLED=1 GOARCH=386 go build -buildmode=c-archive goext.go

go.dll: ext.c goext.h goext.a
	gcc -m32 -g -fPIC -shared -pthread ext.c goext.a -o go.dll -lWinMM -lntdll -lWS2_32 -lsqlite3

test: go.dll
	sqlite3 < test.sql

clean:
	rm goext.a goext.h go.dll