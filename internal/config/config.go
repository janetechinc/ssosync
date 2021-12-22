// Package config ...
package config

// Config ...
type Config struct {
	// Verbose toggles the verbosity
	Debug bool
	// LogLevel is the level with with to log for this config
	LogLevel string `mapstructure:"log_level"`
	// LogFormat is the format that is used for logging
	LogFormat string `mapstructure:"log_format"`
	// GoogleCredentials ...
	GoogleCredentials string `mapstructure:"google_credentials"`
	// GoogleAdmin ...
	GoogleAdmin string `mapstructure:"google_admin"`
	// UserMatch ...
	UserMatch string `mapstructure:"user_match"`
	// GroupFilter ...
	GroupMatch []string `mapstructure:"group_match"`
	// SCIMEndpoint ....
	SCIMEndpoint string `mapstructure:"scim_endpoint"`
	// SCIMAccessToken ...
	SCIMAccessToken string `mapstructure:"scim_access_token"`
	// IsLambda ...
	IsLambda bool
	// Ignore users ...
	IgnoreUsers []string `mapstructure:"ignore_users"`
	// Ignore groups ...
	IgnoreGroups []string `mapstructure:"ignore_groups"`
	// Include groups ...
	IncludeGroups []string `mapstructure:"include_groups"`
	// SyncMethod allow to defined the sync method used to get the user and groups from Google Workspace
	SyncMethod string `mapstructure:"sync_method"`
	// Type of datastore 
	DatastoreType string `mapstructure:"datastore_type"`
	// Prefix or bucket name for datastores
	DatastorePrefix string `mapstructure:"datastore_prefix"`
	// name of the datastore user object or file
	DatastoreUserObj string `mapstructure:"datastore_user_obj"`
	// name of the datastore group object or file
	DatastoreGroupObj string `mapstructure:"datastore_group_obj"`
}

const (
	// DefaultLogLevel is the default logging level.
	DefaultLogLevel = "info"
	// DefaultLogFormat is the default format of the logger
	DefaultLogFormat = "text"
	// DefaultDebug is the default debug status.
	DefaultDebug = false
	// DefaultGoogleCredentials is the default credentials path
	DefaultGoogleCredentials = "credentials.json"
	// DefaultSyncMethod is the default sync method to use.
	DefaultSyncMethod = "groups"
	// DefaultDatastoreType is the default datastore to use
	DefaultDatastoreType = "file"
	DefaultDatastorePrefix = "ssosync-"
	DefaultDatastoreUserObj = "Users.json"
	DefaultDatastoreGroupObj = "Groups.json"
)

// New returns a new Config
func New() *Config {
	return &Config{
		Debug:             DefaultDebug,
		LogLevel:          DefaultLogLevel,
		LogFormat:         DefaultLogFormat,
		SyncMethod:        DefaultSyncMethod,
		GoogleCredentials: DefaultGoogleCredentials,
		DatastoreType:     DefaultDatastoreType,
		DatastorePrefix:   DefaultDatastorePrefix,
		DatastoreUserObj:  DefaultDatastoreUserObj,
		DatastoreGroupObj: DefaultDatastoreGroupObj,
	}
}
