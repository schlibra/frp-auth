package constant

const (
	CreateUserTable = "" +
		"CREATE TABLE IF NOT EXISTS `user`  (" +
		"`id` int NOT NULL AUTO_INCREMENT," +
		"`username` varchar(255) NOT NULL," +
		"`nickname` varchar(255) NOT NULL," +
		"`password` varchar(255) NOT NULL," +
		"`admin` int NULL DEFAULT 0," +
		"`enable` int NULL DEFAULT 0," +
		"`token_version` varchar(255) NULL," +
		"PRIMARY KEY (`id`)" +
		");"
	CreateTokenTable = "" +
		"CREATE TABLE IF NOT EXISTS `token`  (" +
		"`id` int NOT NULL AUTO_INCREMENT," +
		"`name` varchar(255) NOT NULL," +
		"`token` varchar(255) NOT NULL," +
		"`user` varchar(255) NOT NULL," +
		"`enable` int NULL DEFAULT 1," +
		"PRIMARY KEY (`id`)" +
		");"
	CreatePortTable = "" +
		"CREATE TABLE IF NOT EXISTS `port`  (" +
		"`id` int NOT NULL AUTO_INCREMENT," +
		"`min` int NOT NULL," +
		"`max` int NOT NULL," +
		"`token` varchar(255) NOT NULL," +
		"`user` varchar(255) NOT NULL," +
		"PRIMARY KEY (`id`)" +
		");"
	DropUserTable             = "DROP TABLE IF EXISTS `user`;"
	DropTokenTable            = "DROP TABLE IF EXISTS `token`;"
	DropPortTable             = "DROP TABLE IF EXISTS `port`;"
	SelectUserTableByUsername = "" +
		"SELECT `id`, `username`, `password`, `nickname`, `admin`, `enable`, `token_version`" +
		"FROM `user`" +
		"WHERE `username` = ?;"
	SelectUserTableById = "" +
		"SELECT `id`, `username`, `password`, `nickname`, `admin`, `enable`, `token_version`" +
		"FROM `user`" +
		"WHERE `id` = ?;"
	SelectUserTable = "" +
		"SELECT `id`, `username`, `password`, `nickname`, `admin`, `enable`, `token_version`" +
		"FROM `user`;"
	UpdateUserTableNicknameById = "" +
		"UPDATE `user`" +
		"SET `nickname` = ? " +
		"WHERE `id` = ?;"
	UpdateUserTablePasswordById = "" +
		"UPDATE `user`" +
		"SET `password` = ? " +
		"WHERE `id` = ?;"
	UpdateUserTableById = "" +
		"UPDATE `user` " +
		"SET `password` = ?, `nickname` = ?, `admin` = ?, `enable` = ? " +
		"WHERE `id` = ?;"
	UpdateUserTableTokenVersionById = "" +
		"UPDATE `user`" +
		"SET `token_version` = ? " +
		"WHERE `id` = ?;"
	InsertUserTable = "" +
		"INSERT INTO `user`" +
		"(`username`, `password`, `nickname`, `admin`, `enable`)" +
		"VALUES (?, ?, ?, ?, ?)"
	DeleteUserTableById = "" +
		"DELETE FROM `user` " +
		"WHERE `id` = ?;"
	SelectTokenTableByUser = "" +
		"SELECT `id`, `name`, `token`, `user`, `enable`" +
		"FROM `token`" +
		"WHERE `user` = ?;"
	SelectTokenTableById = "" +
		"SELECT `id`, `name`, `token`, `user`, `enable`" +
		"FROM `token`" +
		"WHERE `id` = ?;"
	SelectTokenTableByName = "" +
		"SELECT `id`, `name`, `token`, `user`, `enable`" +
		"FROM `token`" +
		"WHERE `name` = ?;"
	InsertTokenTable = "" +
		"INSERT INTO `token`" +
		"(`name`, `token`, `user`, `enable`)" +
		"VALUES (?, ?, ?, 1);"
	DeleteTokenTableRowById = "" +
		"DELETE FROM `token`" +
		"WHERE `id` = ?;"
	UpdateTokenTableEnableById = "" +
		"UPDATE `token` " +
		"SET `enable` = ? " +
		"WHERE `id` = ?; "
	UpdateTokenTableTokenById = "" +
		"UPDATE `token` " +
		"SET `token` = ? " +
		"WHERE `id` = ?;"
	UpdateTokenTableById = "" +
		"UPDATE `token` " +
		"SET `enable` = ? , `user` = ? " +
		"WHERE `id` = ?;"
	SelectTokenTable = "" +
		"SELECT `id`, `name`, `token`, `user`, `enable` FROM `token`"
	SelectPortTable = "" +
		"SELECT `id`, `min`, `max`, `token`, `user` " +
		"FROM `port`"
	SelectPortTableByToken = "" +
		"SELECT `id`, `min`, `max`, `token`, `user` " +
		"FROM `port` " +
		"WHERE `token` = ?;"
	SelectPortTableByUser = "" +
		"SELECT `id`, `min`, `max`, `token`, `user` " +
		"FROM `port` " +
		"WHERE `user` = ?;"
	SelectPortTableById = "" +
		"SELECT `id`, `min`, `max`, `token`, `user` " +
		"FROM `port` " +
		"WHERE `id` = ?;"
	InsertPortTable = "" +
		"INSERT INTO `port` " +
		"(`min`, `max`, `token`, `user`) " +
		"VALUES (?, ?, ?, ?)"
	DeletePortTableById = "" +
		"DELETE FROM `port` " +
		"WHERE `id` = ?;"
	UpdatePortTableById = "" +
		"UPDATE `port` SET `min` = ? , `max` = ? , `token` = ? , `user` = ? WHERE `id` = ?;"
)
