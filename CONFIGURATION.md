# Configuration
Olaris accepts the following configuration values. These values can be set via command line flags, environment
variables, or configuration files. Configuration sources are loaded in the following order (Each item takes precedence
over the item below it):

* Flags
* Environment Variables
* Configuration File
* Defaults

An `olaris.toml.sample` configuration file is included in the `doc/` folder; rename it to `olaris.toml` and place in
`$HOME/.config/olaris`. You can also override the configuration directory location with the `OLARIS_CONFIG_DIR`
environment variable or the `--config-dir` command-line flag.

You can run `olaris help` to see the complete documentation of all flags and subcommands. To see the documentation for a
specific sub-command, use that sub-command's help flag: `olaris server --help`.

## `config.dir`

* Flag: `--config-dir`
* Env: `OLARIS_CONFIG_DIR`
* Default: `~/.config/olaris`

The default configuration file directory.   

## `server.cacheDir`

* Flag: `--cache-dir`
* Env: `OLARIS_SERVER_CACHEDIR`
* Default: System dependent

The directory used to cache images and store temporary files from transcode sessions.

## `server.verbose`

* Flag: `--verbose`
* Env: `OLARIS_SERVER_VERBOSE`
* Default: `true`

Enables verbose logging.

## `server.port`

* Flag: `--port`, `-p`
* Env: `OLARIS_SERVER_PORT`
* Default: `8080`

HTTP port for the built-in web server.

## `server.dbLog`

* Flag: `--db-log`
* Env: `OLARIS_SERVER_DBLOG`
* Default: `false`

Enables verbose database logging.

## `server.zeroconf.enabled`

* Flag: `--zeroconf-enabled`
* Env: `OLARIS_SERVER_ZEROCONF_ENABLED`
* Default: `false`

Enables zeroconf for the Olaris server.

## `server.zeroconf.domain`

* Flag: `--zeroconf-domain`
* Env: `OLARIS_SERVER_ZEROCONF_DOMAIN`
* Default: `local.`

Sets the zeroconf domain, if enabled.

## `server.sqliteDir`

* Flag: `--sqlite-dir`
* Env: `OLARIS_SERVER_SQLITEDIR`
* Default: `<config_dir>/metadb`

Configures the default directory for SQLite databases.

## `server.directFileAccess`

* Flag: `--allow-direct-file-access`
* Env: `OLARIS_SERVER_DIRECTFILEACCESS`
* Default: `false`

Configures whether accessing files directly by path (without a valid JWT) is allowed.

## `metadata.scanHidden`

* Flag: `--scan-hidden`
* Env: `OLARIS_METADATA_SCANHIDDEN`
* Default: `false`

Tells Olaris whether it should scan hidden files and directories for metadata.

## `database.connection`

* Flag: `--db-conn`
* Env: `OLARIS_DATABASE_CONNECTION`
* Default: None. Uses SQLite.

The database connection string Olaris should use to store metadata for the libraries (default to the default SQLite file
path, overrides the `database.connection` configuration value). The connection string has to be in the following format:
`engine://<connection string data>`. The connection string data can be different for each database, please refer to
[GORM's documentation](https://gorm.io/docs/connecting_to_the_database.html) for more information about compatible
databases.

For example, `mysql://user:password@/dbname?charset=utf8&parseTime=True&loc=Local`

## `debug.streamingPages`

* Flag: `--enable-streaming-debug-pages`
* Env: `OLARIS_DEBUG_STREAMINGPAGES`
* Default: `false`

Whether to enable debug pages in the streaming server (default false, overrides the `debug.streamingPages` configuration
value).

## `debug.transcoderLog`

* Flag: `--write-transcoder-log`
* Env: `OLARIS_DEBUG_TRANSCODERLOG`
* Default: `true`

Whether to write transcoder output to logfile (default true, overrides the `debug.streamingPages` value from
configuration file).

## `rclone.configFile`

* Flag: `--rclone-config`
* Env: `OLARIS_RCLONE_CONFIGFILE`
* Default: `~/.config/rclone/rclone.conf`

The path to your RClone configuration file.
