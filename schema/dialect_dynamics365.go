package schema

// See https://stackoverflow.com/questions/8774928/how-to-exclude-system-table-when-querying-sys-tables

var dynamics365 = dialect{
	escapeIdent: escapeWithBrackets, // [tablename]
	queries: [4]string{
		// columnTypes query.
		`SELECT * FROM %s WHERE 1=0`,
		// tableNames query.
		`SELECT name from sys.objects where type='U'`,
		// viewNames query.
		`SELECT NULL`,
		// get object id
		`SELECT object_id FROM sys.objects WHERE name='%s'`,
	},
}
