package schema

// TODO(js) Is there some way to filter system tables (like mssql)? Or should we always just be using our own schema?

var oracle = dialect{
	escapeIdent: escapeWithBraces, // {tablename}
	queries: [4]string{
		// columnTypes query.
		`SELECT * FROM %s WHERE 1=0`,
		// tableNames query.
		pack(`
			SELECT table_name
			FROM
				all_tables
			WHERE
				owner IN (SELECT sys_context('userenv', 'current_schema') from dual)
		`),
		// viewNames query.
		pack(`
			SELECT view_name
			FROM
				all_views
			WHERE
				owner IN (SELECT sys_context('userenv', 'current_schema') from dual)
		`),
		// get object id
		`SELECT object_id FROM sys.objects WHERE name='%s'`,
	},
}
