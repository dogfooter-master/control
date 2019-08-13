package service

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type DogfooterPrivate struct {
}

func (s *DogfooterPrivate) Service(ctx context.Context, req Payload) (res Payload, err error) {
	//TODO: 인증토큰 체크 로직
	accessToken := SecretTokenObject{
		Token: req.AccessToken,
	}
	var user UserObject
	if user, err = accessToken.Authenticate(); err != nil {
		return
	}
	switch req.Service {
	case "UpdateUserInformation":
		res, err = s.UpdateUserInformation(ctx, req, user)
	case "GetPoint":
		res, err = s.GetPoint(ctx, req, user)
	case "GetLoginPoint":
		res, err = s.GetLoginPoint(ctx, req, user)

	default:
		err = fmt.Errorf("unknown service '%v' in category: '%v'", req.Service, req.Category)
	}
	return
}
func (s *DogfooterPrivate) GetLoginPoint(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Point: GetLoginPoint(),
	}
	return
}
func (s *DogfooterPrivate) GetPoint(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Point: do.Point,
	}
	return
}
func (s *DogfooterPrivate) GetWebsocketHost(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Host: GetConfigWebsocketHosts(),
	}
	return
}
func (s *DogfooterPrivate) SetDefaultTag(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomTag(req.Tag); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) DelTag(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomTag(req.Tag); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) AddTag(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomTag(req.Tag); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) ListTag(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Tag.Apply(GetDefaultTag(), GetTagList())
	return
}

func (s *DogfooterPrivate) SetDefaultGender(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomGender(req.Gender); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) DelGender(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomGender(req.Gender); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) AddGender(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomGender(req.Gender); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) ListGender(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Gender.Apply(GetDefaultGender(), GetGenderList())
	return
}

func (s *DogfooterPrivate) SetDefaultCountry(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomCountry(req.Country); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) DelCountry(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomCountry(req.Country); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) AddCountry(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomCountry(req.Country); err != nil {
		return
	}
	return
}

func (s *DogfooterPrivate) SetDefaultDisease(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomDisease(req.Disease); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) DelDisease(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomDisease(req.Disease); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) AddDisease(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomDisease(req.Disease); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) SetDefaultLocation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomLocation(req.Location); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) DelLocation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomLocation(req.Location); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) AddLocation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomLocation(req.Location); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) SetDefaultSkin(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomSkin(req.Skin); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) DelSkin(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomSkin(req.Skin); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) AddSkin(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomSkin(req.Skin); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) SetDefaultEthnicity(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.SetDefaultCustomEthnicity(req.Ethnicity); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) DelEthnicity(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.ExcludeCustomEthnicity(req.Ethnicity); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) AddEthnicity(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	if res.Default, res.List, err = do.AddCustomEthnicity(req.Ethnicity); err != nil {
		return
	}
	return
}
func (s *DogfooterPrivate) ListLocation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Location.Apply(GetDefaultLocation(), GetLocationList())
	return
}
func (s *DogfooterPrivate) ListDisease(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Disease.Apply(GetDefaultDisease(), GetDiseaseList())
	return
}
func (s *DogfooterPrivate) ListSkin(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Skin.Apply(GetDefaultSkin(), GetSkinList())
	return
}
func (s *DogfooterPrivate) ListCountry(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Country.Apply(GetDefaultCountry(), GetCountryList())
	return
}
func (s *DogfooterPrivate) ListEthnicity(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {

	res = Payload{
		Account: do.Login.Account,
	}
	res.Default, res.List, _ = do.CustomConfig.Ethnicity.Apply(GetDefaultEthnicity(), GetEthnicityList())

	return
}
func (s *DogfooterPrivate) GetUserInformation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	uri := do.Avatar.GetFileUri(do.SecretToken.Token)

	req.HospitalId = do.Relation.HospitalId.Hex()
	if res, err = s.GetHospital(ctx, req, do); err != nil {
		return
	}
	hospitalName := res.HospitalName
	hospitalCreatedBy := res.HospitalCreatedBy
	res = Payload{
		AccessToken:       req.AccessToken,
		Account:           do.Login.Account,
		SecretToken:       do.Login.Password,
		UserId:            do.Id.Hex(),
		HospitalId:        do.Relation.HospitalId.Hex(),
		HospitalName:      hospitalName,
		HospitalCreatedBy: hospitalCreatedBy,
		Name:              do.Name,
		Nickname:          do.Nickname,
	}
	if len(uri) > 0 {
		res.Uri = &UriObject{
			Avatar: uri,
		}
	}
	return
}
func (s *DogfooterPrivate) SignUpComplete(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	if len(req.HospitalName) == 0 {
		err = errors.New("'hospital_name' is mandatory")
		return
	}
	if len(req.Name) == 0 {
		err = errors.New("'name' is mandatory")
		return
	}
	if len(req.Nickname) == 0 {
		err = errors.New("'nickname' is mandatory")
		return
	}
	if res, err = s.CreateHospital(ctx, req, do); err != nil {
		return
	}
	req.HospitalId = res.HospitalId

	return s.UpdateUserInformation(ctx, req, do)
}
func (s *DogfooterPrivate) GetHospital(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	if len(req.HospitalId) == 0 {
		err = errors.New("'hospital_id' is mandatory")
		return
	}
	if bson.IsObjectIdHex(req.HospitalId) == false {
		err = errors.New("'hospital_id' is invalid")
		return
	}
	ho := HospitalObject{
		Id: bson.ObjectIdHex(req.HospitalId),
	}
	var ro HospitalObject
	if ro, err = ho.Read(); err != nil {
		return
	}
	res = Payload{
		HospitalId:        ro.Id.Hex(),
		HospitalName:      ro.Name,
		HospitalCreatedBy: ro.CreatedBy.Hex(),
	}
	return
}
func (s *DogfooterPrivate) CreateHospital(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	if len(req.HospitalName) == 0 {
		err = errors.New("'hospital_name' is mandatory")
		return
	}
	ho := HospitalObject{
		CreatedBy: do.Id,
		Name:      req.HospitalName,
	}
	if err = ho.Create(); err != nil {
		return
	}
	res = Payload{
		HospitalId: ho.Id.Hex(),
	}
	return
}
func (s *DogfooterPrivate) UpdateUserInformation(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	do.Nickname = req.Nickname

	if do.Status == "information" {
		do.Status = "active"
	}
	if err = do.Update(); err != nil {
		return
	}
	res = Payload{
		Account: do.Login.Account,
	}
	return
}
func (s *DogfooterPrivate) UpdateAccessToken(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	do.SecretToken.Refresh()
	if err = do.Update(); err != nil {
		return
	}
	res = Payload{
		Account: do.Login.Account,
	}
	return
}

// TODO: 아바타 파일 저장하는 순서
// 1. PrepareAvatar 를 콜해서 저장할 경로를 얻어온다.
// 2. 해당 경로에 POST 한다.
func (s *DogfooterPrivate) PrepareAvatar(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	if len(req.Avatar) == 0 {
		err = errors.New("'avatar' is mandatory")
		return
	}
	do.Avatar.PrepareAvatarPath(do.Id.Hex(), req.Avatar)
	if err = do.Update(); err != nil {
		return
	}
	uri := do.Avatar.GetFileUri(do.SecretToken.Token)
	res = Payload{
		Account: do.Login.Account,
		Avatar:  uri,
	}
	return
}
func (s *DogfooterPrivate) GetAvatarUri(ctx context.Context, req Payload, do UserObject) (res Payload, err error) {
	uri := do.Avatar.GetFileUri(do.SecretToken.Token)
	res = Payload{
		Account: do.Login.Account,
		Uri: &UriObject{
			Avatar: uri,
		},
	}
	return
}
