package schema

var sqlite = dialect{
	escapeIdent: escapeWithDoubleQuotes, // "tablename"
	queries: [4]string{
		// columnTypes query.
		`SELECT * FROM %s LIMIT 0`,
		// tableNames query.
		pack(`
			SELECT name
			FROM
				sqlite_master
			WHERE
				type = 'table'
		`),
		// viewNames query.
		pack(`
			SELECT name
			FROM
				sqlite_master
			WHERE
				type = 'view'
		`),
		// get object id
		`SELECT NULL`,
	},
}
