package dao

import "gorm.io/gorm"

type Issuer struct {
	gorm.Model
	Uid       int64
	Did       string
	Website   string
	Endpoint  string
	ShortDesc string
	LongDesc  string
	Type      int
}

func (d *Dao) AddIssuer(uid int64, did, website, endpoint, shortDesc, longDesc string, typ int) error {
	return d.mysql.Create(&Issuer{
		Uid:       uid,
		Did:       did,
		Website:   website,
		Endpoint:  endpoint,
		ShortDesc: shortDesc,
		LongDesc:  longDesc,
		Type:      typ,
	}).Error
}

func (d *Dao) GetIssued(uid int64) ([]Issuer, error) {
	issuers := make([]Issuer, 0)
	err := d.mysql.Where("uid = ?", uid).Find(&issuers).Error
	return issuers, err
}

func (d *Dao) DelIssued(id uint64) error {
	return d.mysql.Where("id = ?", id).Delete(&Issuer{}).Error
}

func (d *Dao) GetIssuers() ([]Issuer, error) {
	issuers := make([]Issuer, 0)
	err := d.mysql.Find(&issuers).Error
	return issuers, err
}
