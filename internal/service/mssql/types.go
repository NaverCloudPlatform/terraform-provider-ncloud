package cloudmssql

type CloudMssqlInstance struct {
	CloudMssqlInstanceNo         *string   `json:"cloudMssqlInstanceNo,omitempty"`
	CloudMssqlServiceName        *string   `json:"cloudMssqlServiceName,omitempty"`
	CloudMssqlInstanceStatusName *string   `json:"cloudMssqlInstanceStatusName,omitempty"`
	CloudMssqlImageProductCode   *string   `json:"cloudMssqlImageProductCode,omitempty"`
	IsHa                         *bool     `json:"isHa,omitempty"`
	CloudMssqlPort               *int32    `json:"cloudMssqlPort,omitempty"`
	BackupFileRetentionPeriod    *int32    `json:"backupFileRetentionPeriod,omitempty"`
	BackupTime                   *string   `json:"backupTime,omitempty"`
	ConfigGroupNo                *string   `json:"configGroupNo,omitempty"`
	EngineVersion                *string   `json:"engineVersion,omitempty"`
	CreateDate                   *string   `json:"createDate,omitempty"`
	DbCollation                  *string   `json:"dbCollation,omitempty"`
	AccessControlGroupNoList     []*string `json:"accessControlGroupNoList,omitempty"`
}
