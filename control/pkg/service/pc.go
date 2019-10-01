package service

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type PcObject struct {
	Name        string            `bson:"name,omitempty"`
	MacAddress  string            `bson:"mac_address,omitempty"`
	Remember    bool              `bson:"remember,omitempty"`
	AccessToken SecretTokenObject `bson:"access_token,omitempty"`
	CreateTime  time.Time         `bson:"create_time,omitempty"`
	UpdateTime  time.Time         `bson:"update_time,omitempty"`
}

func (p *PcObject) Read() (obj UserObject, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	readBson := bson.M{}
	readBson["pc_list"] = bson.M {
		"$elemMatch" : bson.M {
			"access_token.token" : p.AccessToken.Token,
		},
	}
	err = collection.Find(readBson).One(&obj)
	if err != nil {
		return
	}
	return
}
func (p *PcObject) ReadAll() (objList []UserObject, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)
	readBson := bson.M{}
	readBson["pc_list"] = bson.M {
		"$elemMatch" : bson.M {
			"mac_address" : p.MacAddress,
		},
	}
	err = collection.Find(readBson).Sort([]string{"-update_time"}...).Skip(0).Limit(9999).All(&objList)
	if err != nil {
		return
	}
	return
}
func (p *PcObject) AddedToUser(obj UserObject) (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	isFound := false
	for i, e := range obj.PcList {
		if e.MacAddress == p.MacAddress {
			obj.PcList[i].UpdateTime = time.Now().UTC()
			obj.PcList[i].Remember = false
			e.AccessToken.Refresh()
			p.AccessToken = e.AccessToken
			isFound = true
			break
		}
	}
	if isFound == false {
		_ = p.AccessToken.GenerateAccessToken()
		obj.PcList = append(obj.PcList, *p)
	}

	updateBson := bson.M{}
	updateBson["pc_list"] = obj.PcList

	obj.Time.Update()
	updateBson["time.update_time"] = obj.Time.UpdateTime

	err = collection.Update(bson.M{"_id": obj.Id}, bson.M{"$set": updateBson})
	if err != nil {
		return
	}
	return
}
func (p *PcObject) RememberDevice(obj UserObject, remember bool) (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	isFound := false
	for i, e := range obj.PcList {
		if e.MacAddress == p.MacAddress {
			obj.PcList[i].UpdateTime = time.Now().UTC()
			obj.PcList[i].Remember = remember
			isFound = true
			break
		}
	}
	if isFound == false {
		err = errors.New("not found")
		return
	}

	updateBson := bson.M{}
	updateBson["pc_list"] = obj.PcList

	obj.Time.Update()
	updateBson["time.update_time"] = obj.Time.UpdateTime

	err = collection.Update(bson.M{"_id": obj.Id}, bson.M{"$set": updateBson})
	if err != nil {
		return
	}
	return
}
func (p *PcObject) IsValidate() bool {
	return !p.AccessToken.IsExpiredLong()
}
func (p *PcObject) RefreshAccessToken() {
	p.AccessToken.Refresh()
}
func GetPcAccessToken(pcList []PcObject, macAddress string) (accessToken string, remember bool, err error) {
	for _, e := range pcList {
		if e.MacAddress == macAddress {
			//if e.IsValidate() && e.Remember {
			if e.IsValidate() {
				accessToken = e.AccessToken.Token
				remember = e.Remember
				return
			}
		}
	}
	err = errors.New("not found")
	return
}