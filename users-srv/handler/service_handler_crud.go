package handler

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	proto "github.com/SonicRoshan/Velocity/users-srv/proto"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

//Add is used to handle Add function
func (usersService UsersService) Add(
	ctx context.Context, request *proto.AddRequest, response *proto.AddResponse) error {
	log.Debugf("Adding User With Request %+v", request)

	var mainData config.UserMain
	var extraData config.UserExtra

	err := usersService.copy(&mainData, &request.MainData)
	if err != nil {
		return err
	}
	err = usersService.copy(&extraData, &request.ExtraData)
	if err != nil {
		return err
	}

	log.Debugf("Adding User With Main Data %+v \n  Extra Data %v",
		mainData, extraData)

	userID, msg, err := usersService.users.Add(mainData, extraData)
	if err != nil {
		log.Errorf("Adding User With Main Data %+v\n  Extra Data %v\n Returned Error %+v",
			mainData, extraData, err)
		return errors.Wrap(err, "Error While Adding User")
	}

	log.Infof("Added User With Main Data %+v  \n Extra Data %+v",
		mainData, extraData)
	response.Message = msg
	response.UserID = userID

	return nil
}

//Get is used to handle Get function
func (usersService UsersService) Get(
	ctx context.Context, request *proto.GetRequest, response *proto.UserMain) error {

	log.Debugf("Getting User Main Data With User ID %s", request.UserID)

	var data config.UserMain
	err := usersService.users.Get(request.UserID, config.DBConfigUserMainDataCollection, &data)

	if err != nil {
		log.Errorf("Getting User Main Data With User ID %s Returned Error %+v",
			request.UserID, err)
		return errors.Wrap(err, "Error While Getting User Data")
	}

	err = usersService.copy(&response, &data)
	if err != nil {
		return err
	}

	log.Infof("Got User Data With Data %+v", response)
	return nil
}

//GetExtra is used to handle GetExtra function
func (usersService UsersService) GetExtra(
	ctx context.Context, request *proto.GetRequest, response *proto.UserExtra) error {

	log.Debugf("Getting User Main Data With User ID %s", request.UserID)

	var data config.UserMain
	err := usersService.users.Get(request.UserID, config.DBConfigUserMainDataCollection, &data)

	if err != nil {
		log.Errorf("Getting User Main Data With User ID %s Returned Error %+v",
			request.UserID, err)
		return errors.Wrap(err, "Error While Getting User Data")
	}

	err = usersService.copy(&response, &data)
	if err != nil {
		return err
	}

	log.Infof("Got User Data With Data %+v", response)
	return nil
}

//Update is used to handle Update function
func (usersService UsersService) Update(
	ctx context.Context,
	request *proto.UpdateRequest,
	response *proto.UpdateResponse) error {

	log.Debug("Updating User Main Data")
	var update config.UserMain

	err := copier.Copy(&update, &request.Update)
	if err != nil {
		return err
	}

	log.Debugf("Updating User Main Data With ID %s With Update %+v",
		request.UserID, update)

	err = usersService.users.Update(request.UserID,
		update,
		config.DBConfigUserMainDataCollection)

	if err != nil {
		log.Debugf("Updating User Main Data With ID %s With Update %+v Returned Error %+v",
			request.UserID, update, err)
		return errors.Wrap(err, "Error While Updating Data")
	}
	log.Infof("Updated User Main Data With ID %s With Update %+v",
		request.UserID, update)

	return nil
}

//UpdateExtra is used to handle UpdateExtra function
func (usersService UsersService) UpdateExtra(
	ctx context.Context,
	request *proto.UpdateRequest,
	response *proto.UpdateResponse) error {

	log.Debug("Updating User Extra Data")
	var update config.UserExtra

	err := copier.Copy(&request.Update, &update)
	if err != nil {
		return err
	}

	log.Debugf("Updating User Extra Data With ID %s With Update %+v",
		request.UserID, update)

	err = usersService.users.Update(request.UserID,
		update,
		config.DBConfigUserExtraDataCollection)

	if err != nil {
		log.Debugf("Updating User Extra Data With ID %s With Update %+v Returned Error %+v",
			request.UserID, update, err)
		return errors.Wrap(err, "Error While Updating Data")
	}

	log.Infof("Updated User Extra Data With ID %s With Update %+v",
		request.UserID, update)

	return nil
}

//Delete is used to handle Delete function
func (usersService UsersService) Delete(
	ctx context.Context,
	request *proto.DeleteRequest,
	response *proto.DeleteResponse) error {

	log.Debugf("Deleting User %s", request.Username)
	msg, err := usersService.users.Delete(request.UserID, request.Username, request.Password)

	if err != nil {
		log.Errorf("Delete User %s Returned Error %+v", request.Username, err)
		return errors.Wrap(err, "Error While Deleting User")
	}

	log.Infof("Deleted User %s", request.Username)
	response.Message = msg
	return err
}
