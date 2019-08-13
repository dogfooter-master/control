package service

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type HospitalObject struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	CreatedBy bson.ObjectId `bson:"created_by,omitempty"`
	Name      string        `bson:"name,omitempty"`
	Address   string        `bson:"address,omitempty"`
	Logo      FileObject    `bson:"logo,omitempty"`
	Status    string        `bson:"status,omitempty"`
	Time      TimeLogObject `bson:"time,omitempty"`
}

func (d *HospitalObject) Create() (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.HospitalCollection)

	//index := mgo.Index{
	//	Key:    []string{"login.account"},
	//	Unique: true,
	//}
	//if err = collection.EnsureIndex(index); err != nil {
	//	err = fmt.Errorf("collection.EnsureIndex: %v", err)
	//	return
	//}

	if len(d.Id) == 0 {
		d.Id = bson.NewObjectId()
	}
	d.Time.Initialize()

	err = collection.Insert(d)
	if err != nil {
		err = fmt.Errorf("insert: %v", err)
		return
	}

	return
}

func (d *HospitalObject) Read() (obj HospitalObject, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.HospitalCollection)

	readBson := bson.M{}
	if len(d.Id) > 0 {
		readBson["_id"] = d.Id
	}

	err = collection.Find(readBson).One(&obj)
	if err != nil {
		return
	}
	return
}
func (d *HospitalObject) ReadAll(skip int, limit int) (objList []HospitalObject, err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.HospitalCollection)

	readBson := bson.M{}
	if len(d.Name) > 0 {
		readBson["name"] = d.Name
	}

	err = collection.Find(readBson).Skip(skip).Limit(limit).All(&objList)
	if err != nil {
		return
	}
	return
}
func (d *HospitalObject) Update() (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.HospitalCollection)

	object := HospitalObject{
		Id: d.Id,
	}
	err = collection.Find(bson.M{"_id": object.Id}).One(&object)
	if err != nil {
		return
	}

	updateBson := bson.M{}

	if len(d.Name) > 0 {
		updateBson["name"] = d.Name
	}
	if len(d.Address) > 0 {
		updateBson["address"] = d.Address
	}
	if len(d.CreatedBy) > 0 {
		updateBson["created_by"] = d.CreatedBy
	}

	d.Time.Update()
	updateBson["time.update_time"] = d.Time.UpdateTime

	err = collection.Update(bson.M{"_id": object.Id}, bson.M{"$set": updateBson})
	if err != nil {
		return
	}

	return
}
func (d *HospitalObject) Delete() (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.HospitalCollection)

	err = collection.Remove(bson.M{"_id": d.Id})
	if err != nil {
		return
	}

	return
}
