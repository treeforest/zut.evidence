package dao

import "gorm.io/gorm"

// Challenge 挑战
type Challenge struct {
	gorm.Model
	SenderUid       int64
	SenderDid       string
	RecipientUid    int64
	RecipientDid    string
	RecipientPubKey string
	Ciphertext      string // 密文
	Plaintext       string // 明文
	Status          int    // 0: 待验证；1：验证失败；2：验证成功
}

func (d *Dao) AddChallenge(senderUid, recipientUid int64, senderDid, recipientDid, recipientPubKey,
	Ciphertext, Plaintext string) error {
	return d.mysql.Create(&Challenge{
		SenderUid:       senderUid,
		SenderDid:       senderDid,
		RecipientUid:    recipientUid,
		RecipientDid:    recipientDid,
		RecipientPubKey: recipientPubKey,
		Ciphertext:      Ciphertext,
		Plaintext:       Plaintext,
		Status:          0,
	}).Error
}

func (d *Dao) GetChallengeById(id uint) (*Challenge, error) {
	c := &Challenge{}
	err := d.mysql.Where("id = ?", id).Find(c).Error
	if err != nil {
		return nil, err
	}
	if c.SenderUid == 0 {
		return nil, nil
	}
	return c, nil
}

func (d *Dao) GetChallengeSent(uid int64) ([]Challenge, error) {
	challenges := make([]Challenge, 0)
	err := d.mysql.Where("sender_uid = ?", uid).Order("created_at desc").Find(&challenges).Error
	return challenges, err
}

func (d *Dao) SetChallengeStatus(id uint, status int) error {
	return d.mysql.Model(&Challenge{}).Where("id = ?", id).Update("status", status).Error
}

func (d *Dao) GetChallengeDoing(uid int64) ([]Challenge, error) {
	challenges := make([]Challenge, 0)
	err := d.mysql.Where("recipient_uid = ? AND status = ?", uid, 0).Order("created_at desc").Find(&challenges).Error
	return challenges, err
}

func (d *Dao) GetChallengeDone(uid int64) ([]Challenge, error) {
	challenges := make([]Challenge, 0)
	err := d.mysql.Where("recipient_uid = ? AND status <> ?", uid, 0).Order("created_at desc").Find(&challenges).Error
	return challenges, err
}
