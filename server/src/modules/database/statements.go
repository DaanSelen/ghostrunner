package database

type Statements struct {
	SetupDatabase string

	AdminTokenCreate   string
	GetTokenID         string
	CreateToken        string
	DeleteToken        string
	RetrieveTokens     string
	RetrieveTokenNames string

	CreateTask   string
	DeleteTask   string
	ListAllTasks string
	CountTasks   string
}

var declStat = Statements{
	SetupDatabase: `
	CREATE TABLE IF NOT EXISTS tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		token TEXT UNIQUE NOT NULL
	);
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		command TEXT NOT NULL,
		nodeid TEXT NOT NULL,
		creation TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS completed (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		command TEXT NOT NULL,
		nodeid TEXT NOT NULL,
		completion TEXT NOT NULL,
		result TEXT NOT NULL
	);`,

	AdminTokenCreate: `
	INSERT INTO tokens (id, name, token)
		VALUES (0, ?, ?)
		ON CONFLICT (id) DO NOTHING;`,
	GetTokenID: `
	SELECT id FROM tokens WHERE name = ?`,
	CreateToken: `
	INSERT INTO tokens(name, token)
		VALUES(?, ?);`,
	DeleteToken: `
	DELETE FROM tokens WHERE
		name = ?;`,
	RetrieveTokens: `
	SELECT token FROM tokens`,
	RetrieveTokenNames: `
	SELECT name FROM tokens`,

	CreateTask: `
	INSERT INTO tasks (name, command, nodeid, creation)
		VALUES (?, ?, ?, ?);`,
	DeleteTask: `
	DELETE FROM tasks WHERE name = ? AND nodeid = ?;`,
	ListAllTasks: `
	Select name, command, nodeid, creation from tasks;`,
	CountTasks: `
	SELECT COUNT(*)
	FROM tasks
	WHERE name = ?;`,
}
