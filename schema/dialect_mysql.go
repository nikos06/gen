package schema

var mysql = dialect{
	escapeIdent: escapeWithBackticks, // `tablename`
	queries: [4]string{
		// columnTypes query.
		`SELECT * FROM %s LIMIT 0`,
		// tableNames query.
		pack(`
			SELECT table_name
			FROM
				information_schema.tables
			WHERE
				table_type = 'BASE TABLE' AND
				table_schema = database()
		`),
		// viewNames query.
		pack(`
			SELECT table_name
			FROM
				information_schema.tables
			WHERE
				table_type = 'VIEW' AND
				table_schema = database()
		`),
		// get object id
		`SELECT object_id FROM sys.objects WHERE name='%s'`,
	},
}
