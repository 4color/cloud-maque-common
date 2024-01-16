package model

type SdpRole struct {
	Roleid    string `xorm:"not null pk unique VARCHAR(50)" json:"roleid" `
	Rolename  string `xorm:"VARCHAR(50)" json:"rolename" `
	Des       string `xorm:"VARCHAR(100)" json:"des" `
	Indexpage string `xorm:"VARCHAR(100)" json:"indexpage" `
	Isinner   int    `xorm:"default '0' SMALLINT" json:"isinner" `
	RoleCode  string `xorm:"VARCHAR(100)" json:"roleCode" ` //角色代码
}
