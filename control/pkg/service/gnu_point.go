package service

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type GNUPointObject struct {
	Point      int       `bson:"point,omitempty"`
	UpdateTime time.Time `bson:"update_time,omitempty"`
}

func (p *GNUPointObject) UpdatePoint(user UserObject, point int) (err error) {
	collection := mgoSession.DB(mgoConfig.Database).C(mgoConfig.UserCollection)

	updateBson := bson.M{}
	updateBson["gnu_point.point"] = point
	updateBson["gnu_point.update_time"] = time.Now().UTC().Local()

	err = collection.Update(bson.M{"_id": user.Id}, bson.M{"$set": updateBson})
	if err != nil {
		return
	}
	return
}
