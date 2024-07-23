package sql


//java-db中描述indice表的类
type Indices struct{
	ArtifactId int `gorm:"column:artifact_id;type:INTEGER"`
	Version string `gorm:"column:version;type:TEXT"`
	Sha1 []byte `gorm:"column:sha1;type:BLOB"`
	ArchiveType string `gorm:"column:archive_type;type:TEXT"`
}

func (indices Indices) TableName() string{
	return "indices"
}
