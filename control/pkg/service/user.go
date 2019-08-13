package service

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserObject struct {
	Id               bson.ObjectId                `bson:"_id,omitempty"`
	Name             string                       `bson:"name,omitempty"`
	Nickname         string                       `bson:"nickname,omitempty"`
	Login            LoginObject                  `bson:"login,omitempty"`
	Relation         RelationObject               `bson:"relation,omitempty"`
	SecretToken      SecretTokenObject            `bson:"secret_token,omitempty"`
	Client           ClientObject                 `bson:"client,omitempty"`
	Avatar           FileObject                   `bson:"avatar,omitempty"`
	Status           string                       `bson:"status,omitempty"`
	Time             TimeLogObject                `bson:"time,omitempty"`
	CustomConfig     CustomConfigObject           `bson:"custom_config,omitempty"`
}

func (d *UserObject) Create() (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	index := mgo.Index{
		Key:    []string{"login.account"},
		Unique: true,
	}
	if err = collection.EnsureIndex(index); err != nil {
		err = fmt.Errorf("collection.EnsureIndex: %v", err)
		return
	}

	d.Time.Initialize()

	err = collection.Insert(d)
	if err != nil {
		err = fmt.Errorf("insert: %v", err)
		return
	}

	return
}

func (d *UserObject) Read() (obj UserObject, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	readBson := bson.M{}
	if len(d.Id) > 0 {
		readBson["_id"] = d.Id
	}
	if len(d.Login.Account) > 0 {
		readBson["login.account"] = d.Login.Account
	}

	err = collection.Find(readBson).One(&obj)
	if err != nil {
		return
	}
	return
}
func (d *UserObject) ReadAll(skip int, limit int) (objList []UserObject, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	readBson := bson.M{}
	if len(d.Login.Account) > 0 {
		readBson["login.account"] = d.Login.Account
	}

	err = collection.Find(readBson).Skip(skip).Limit(limit).All(&objList)
	if err != nil {
		return
	}
	return
}
func (d *UserObject) Update() (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	//object := UserObject{
	//	Id: d.Id,
	//}
	//err = collection.Find(bson.M{"_id": object.Id}).One(&object)
	//if err != nil {
	//	return
	//}

	updateBson := bson.M{}

	if len(d.Login.Password) > 0 {
		updateBson["login.password"] = d.Login.Password
	}
	if len(d.Name) > 0 {
		updateBson["name"] = d.Name
	}
	if len(d.Nickname) > 0 {
		updateBson["nickname"] = d.Nickname
	}
	if len(d.Avatar.Path) > 0 {
		updateBson["avatar"] = d.Avatar
	}
	if len(d.SecretToken.Token) > 0 {
		updateBson["secret_token"] = d.SecretToken
	}
	if len(d.Status) > 0 {
		updateBson["status"] = d.Status
	}
	if d.Time.LoginTime.IsZero() == false {
		updateBson["time.login_time"] = d.Time.LoginTime
	}
	if len(d.Relation.HospitalId) > 0 {
		updateBson["relation.hospital_id"] = d.Relation.HospitalId
	}
	if len(d.Relation.PcUserId) > 0 {
		updateBson["relation.pc_user_id"] = d.Relation.PcUserId
	}

	d.Time.Update()
	updateBson["time.update_time"] = d.Time.UpdateTime

	err = collection.Update(bson.M{"_id": d.Id}, bson.M{"$set": updateBson})
	if err != nil {
		return
	}

	return
}
func (d *UserObject) UpdatePcUser() (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	updateBson := bson.M{}
	updateBson["relation.pc_user_id"] = d.Relation.PcUserId

	d.Time.Update()
	updateBson["time.update_time"] = d.Time.UpdateTime

	err = collection.Update(bson.M{"_id": d.Id}, bson.M{"$set": updateBson})
	if err != nil {
		return
	}

	return
}
func (d *UserObject) Delete() (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	err = collection.Remove(bson.M{"_id": d.Id})
	if err != nil {
		return
	}

	return
}
func (d *UserObject) DeleteDuplicatedAccount() (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	err = collection.Remove(bson.M{"login.account": d.Login.Account})
	if err != nil {
		return
	}

	return
}

func (d *UserObject) Validate() (err error) {
	if len(d.Name) < 1 {
		err = errors.New("'name' is mandatory")
		return
	}
	if len(d.Nickname) == 0 {
		d.Nickname = d.Name
	}

	return
}

func (d *UserObject) AddCustomEthnicity(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.ethnicity.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.ethnicity.excluded": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Ethnicity.Apply(GetDefaultEthnicity(), GetEthnicityList())
}
func (d *UserObject) ExcludeCustomEthnicity(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.ethnicity.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.ethnicity.added": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Ethnicity.Apply(GetDefaultEthnicity(), GetEthnicityList())
}
func (d *UserObject) SetDefaultCustomEthnicity(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.ethnicity.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.ethnicity.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$set": bson.M{
				"custom_config.ethnicity.default": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Ethnicity.Apply(GetDefaultEthnicity(), GetEthnicityList())
}
func (d *UserObject) AddCustomCountry(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.country.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.country.excluded": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Country.Apply(GetDefaultCountry(), GetCountryList())
}
func (d *UserObject) ExcludeCustomCountry(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.country.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.country.added": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Country.Apply(GetDefaultCountry(), GetCountryList())
}
func (d *UserObject) SetDefaultCustomCountry(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.country.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.country.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$set": bson.M{
				"custom_config.country.default": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Country.Apply(GetDefaultCountry(), GetCountryList())
}
func (d *UserObject) AddCustomSkin(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.skin.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.skin.excluded": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Skin.Apply(GetDefaultSkin(), GetSkinList())
}
func (d *UserObject) ExcludeCustomSkin(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.skin.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.skin.added": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Skin.Apply(GetDefaultSkin(), GetSkinList())
}
func (d *UserObject) SetDefaultCustomSkin(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.skin.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.skin.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$set": bson.M{
				"custom_config.skin.default": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Skin.Apply(GetDefaultSkin(), GetSkinList())
}
func (d *UserObject) AddCustomDisease(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.disease.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.disease.excluded": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Disease.Apply(GetDefaultDisease(), GetDiseaseList())
}
func (d *UserObject) ExcludeCustomDisease(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.disease.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.disease.added": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Disease.Apply(GetDefaultDisease(), GetDiseaseList())
}
func (d *UserObject) SetDefaultCustomDisease(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.disease.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.disease.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$set": bson.M{
				"custom_config.disease.default": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Disease.Apply(GetDefaultDisease(), GetDiseaseList())
}
func (d *UserObject) AddCustomLocation(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.location.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.location.excluded": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Location.Apply(GetDefaultLocation(), GetLocationList())
}
func (d *UserObject) ExcludeCustomLocation(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.location.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.location.added": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Location.Apply(GetDefaultLocation(), GetLocationList())
}
func (d *UserObject) SetDefaultCustomLocation(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.location.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.location.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$set": bson.M{
				"custom_config.location.default": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Location.Apply(GetDefaultLocation(), GetLocationList())
}
func (d *UserObject) AddCustomGender(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.gender.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.gender.excluded": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Gender.Apply(GetDefaultGender(), GetGenderList())
}
func (d *UserObject) ExcludeCustomGender(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.gender.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.gender.added": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Gender.Apply(GetDefaultGender(), GetGenderList())
}
func (d *UserObject) SetDefaultCustomGender(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.gender.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.gender.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$set": bson.M{
				"custom_config.gender.default": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Gender.Apply(GetDefaultGender(), GetGenderList())
}
func (d *UserObject) AddCustomTag(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.tag.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.tag.excluded": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Tag.Apply(GetDefaultTag(), GetTagList())
}
func (d *UserObject) ExcludeCustomTag(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.tag.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.tag.added": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Tag.Apply(GetDefaultTag(), GetTagList())
}
func (d *UserObject) SetDefaultCustomTag(value string) (defaultValue string, dstList []string, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$addToSet": bson.M{
				"custom_config.tag.added": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$pull": bson.M{
				"custom_config.tag.excluded": value,
			},
		}); err != nil {
		return
	}
	if err = collection.Update(
		bson.M{"_id": d.Id},
		bson.M{
			"$set": bson.M{
				"custom_config.tag.default": value,
			},
		}); err != nil {
		return
	}
	var ro UserObject
	if ro, err = d.Read(); err != nil {
		return
	}

	return ro.CustomConfig.Tag.Apply(GetDefaultTag(), GetTagList())
}
