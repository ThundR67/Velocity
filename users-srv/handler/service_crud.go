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

	var mainData, extraData config.UserType
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
	ctx context.Context, request *proto.GetRequest, response *proto.User) error {

	log.Debugf("Getting User Data With User ID %s In Coll %s",
		request.UserID, request.Collection)

	data, err := usersService.users.Get(request.UserID, request.Collection)

	if err != nil {
		log.Errorf("Getting User Data With User ID %s In Coll %s Returned Error %+v",
			request.UserID, request.Collection, err)
		return errors.Wrap(err, "Error While Getting User Data")
	}

	err = usersService.copy(&response, &data)
	return err
}

//Update is used to handle Update function
func (usersService UsersService) Update(
	ctx context.Context,
	request *proto.UpdateRequest,
	response *proto.UpdateResponse) error {

	log.Debug("Updating User Data")
	var update config.UserType

	err := copier.Copy(&request.Update, &update)
	if err != nil {
		return err
	}

	log.Debugf("Updating User Data With ID %s With Update %+v In Coll %s",
		request.UserID, update, request.Collection)

	err = usersService.users.Update(request.UserID, update, request.Collection)

	if err != nil {
		log.Debugf("Updating User Data With ID %s With Update %+v In Coll %s Returned Error %+v",
			request.UserID, update, request.Collection, err)
		return errors.Wrap(err, "Error While Updating Data")
	}
	log.Infof("Updated User Data With ID %s With Update %+v In Coll %s",
		request.UserID, update, request.Collection)

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
