module github.com/smallnest/gen

go 1.14

require (
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc
	github.com/bxcodec/faker/v3 v3.3.1
	github.com/davecgh/go-spew v1.1.1
	github.com/denisenkom/go-mssqldb v0.0.0-20191128021309-1d7a30a10f73
	github.com/droundy/goopt v0.0.0-20170604162106-0b8effe182da
	github.com/gobuffalo/packd v1.0.0
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/jimsmart/schema v0.0.4
	github.com/jinzhu/gorm v1.9.11
	github.com/jinzhu/inflection v1.0.0
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/lib/pq v1.3.0
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/nikos06/odbc v0.9.2
	github.com/ompluscator/dynamic-struct v1.2.0
	github.com/rogpeppe/go-internal v1.6.2 // indirect
	github.com/serenize/snaker v0.0.0-20171204205717-a683aaf2d516
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/swaggo/swag v1.7.3
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
)

replace github.com/smallnest/gen/packrd => ./packerd
