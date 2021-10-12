package dbmeta

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/smallnest/gen/schema"
)

type colInfo struct {
	Table   string           `json:"table"`
	ColInfo []D365ColumnInfo `json:"colInfo"`
}

type colInfos struct {
	ColInfos []colInfo `json:"colInfos"`
}

var ColInfos colInfos

// LoadD365Meta fetch db meta data for MS SQL database
func LoadD365Meta(db *sql.DB, sqlType, sqlDatabase, tableName string) (DbTableMeta, error) {
	m := &dbTableMeta{
		sqlType:     sqlType,
		sqlDatabase: sqlDatabase,
		tableName:   tableName,
	}

	cols, err := schema.Table(db, m.tableName)
	if err != nil {
		return nil, err
	}

	objectId, err := schema.ObjectId(db, m.tableName)
	if err != nil {
		return nil, err
	}

	m.columns = make([]*columnMeta, len(cols))

	colInfo, err := d365loadFromSysColumns(db, objectId, m.tableName)
	if err != nil {
		return nil, fmt.Errorf("unable to load ddl from ms sql: %v", err)
	}

	err = d365LoadPrimaryKey(db, objectId, colInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to load ddl from ms sql: %v", err)
	}

	infoSchema, err := LoadTableInfoFromD365InformationSchema(db, objectId, m.tableName)
	if err != nil {
		fmt.Printf("error calling LoadTableInfoFromMSSqlInformationSchema table: %s error: %v\n", tableName, err)
	}

	for i, v := range cols {

		nullable, ok := v.Nullable()
		if !ok {
			nullable = false
		}
		isAutoIncrement := false
		isPrimaryKey := i == 0
		var columnLen int64 = -1

		defaultVal := ""
		columnType := v.DatabaseTypeName()
		colDDL := v.DatabaseTypeName()

		colInfo, ok := colInfo[v.Name()]
		if ok {
			isPrimaryKey = colInfo.PrimaryKey
			nullable = colInfo.IsNullable
			isAutoIncrement = colInfo.IsIdentity
			dbType := strings.ToLower(v.DatabaseTypeName())

			if strings.Contains(dbType, "char") || strings.Contains(dbType, "text") {
				columnLen = colInfo.MaxLength
			}
		} else {
			MissingColumns[tableName] = append(MissingColumns[tableName], v.Name())
		}

		if infoSchema != nil {
			infoSchemaColInfo, ok := infoSchema[v.Name()]
			if ok {
				if infoSchemaColInfo.ColumnDefault != nil {
					defaultVal = fmt.Sprintf("%v", infoSchemaColInfo.ColumnDefault)
					defaultVal = cleanupDefault(defaultVal)
				}
			}
		}

		colMeta := &columnMeta{
			index:            i,
			name:             v.Name(),
			databaseTypeName: columnType,
			nullable:         nullable,
			isPrimaryKey:     isPrimaryKey,
			isAutoIncrement:  isAutoIncrement,
			colDDL:           colDDL,
			defaultVal:       defaultVal,
			columnType:       columnType,
			columnLen:        columnLen,
		}

		m.columns[i] = colMeta
	}

	m.ddl = BuildDefaultTableDDL(tableName, m.columns)
	m = updateDefaultPrimaryKey(m)
	return m, nil
}

func d365LoadPrimaryKey(db *sql.DB, tableName string, colInfo map[string]*D365ColumnInfo) error {

	primaryKeySQL := fmt.Sprintf(`
		SELECT a.name
          FROM sys.columns a 
		  LEFT OUTER JOIN (SELECT i.object_id, ic.column_id, i.is_primary_key
			FROM sys.indexes i
		  	LEFT JOIN sys.index_columns ic ON ic.object_id = i.object_id AND ic.index_id = i.index_id
			WHERE i.is_primary_key = 1) AS p on p.object_id = a.object_id AND p.column_id = a.column_id
          WHERE a.object_id=%s AND p.is_primary_key=1
		`, tableName)
	res, err := db.Query(primaryKeySQL)
	if err != nil {
		return fmt.Errorf("unable to load ddl from ms sql: %v", err)
	}
	defer res.Close()
	for res.Next() {

		var columnName string
		err = res.Scan(&columnName)
		if err != nil {
			return fmt.Errorf("unable to load identity info from ms sql Scan: %v", err)
		}

		//fmt.Printf("## PRIMARY KEY COLUMN_NAME: %s\n", columnName)
		colInfo, ok := colInfo[columnName]
		if ok {
			colInfo.PrimaryKey = true
			//fmt.Printf("name: %s primary_key: %t\n", colInfo.name, colInfo.primary_key)
		}
	}
	return nil
}

func d365loadFromSysColumns(db *sql.DB, objectId string, tableName string) (colInfo map[string]*D365ColumnInfo, err error) {
	colInfo = make(map[string]*D365ColumnInfo)

	identitySQL := fmt.Sprintf(`
SELECT name, is_identity, is_nullable, max_length 
FROM sys.columns 
WHERE  object_id = %s`, objectId)

	res, err := db.Query(identitySQL)
	if err != nil {
		return nil, fmt.Errorf("unable to load ddl from ms sql: %v", err)
	}

	defer res.Close()
	for res.Next() {
		var name string
		var isIdentity, isNullable bool
		var maxLength int64
		err = res.Scan(&name, &isIdentity, &isNullable, &maxLength)
		if err != nil {
			return nil, fmt.Errorf("unable to load identity info from ms sql Scan: %v", err)
		}

		colInfo[name] = &D365ColumnInfo{
			Name:       name,
			IsIdentity: isIdentity,
			IsNullable: isNullable,
			MaxLength:  maxLength,
		}
	}

	// Load missing columns info
	for _, c := range ColInfos.ColInfos {
		if c.Table == tableName {
			for _, ci := range c.ColInfo {
				colInfo[ci.Name] = &ci
			}
		}
	}

	return colInfo, err
}

type D365ColumnInfo struct {
	Name       string `json:"name"`
	IsIdentity bool   `json:"isIdentity"`
	IsNullable bool   `json:"isNullable"`
	PrimaryKey bool   `json:"primaryKey"`
	MaxLength  int64  `json:"maxLength"`
}

/*
https://www.mssqltips.com/sqlservertip/1512/finding-and-listing-all-columns-in-a-sql-server-database-with-default-values/
*/
