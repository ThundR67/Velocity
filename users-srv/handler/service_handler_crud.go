package handler

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	proto "github.com/SonicRoshan/Velocity/users-srv/proto"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//Add is used to handle Add function
func (usersService UsersService) Add(
	ctx context.Context, request *proto.AddRequest, response *proto.AddResponse) error {

	log.Debug(
		"Adding User",
		zap.Any("Request", request),
	)

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

	log.Debug(
		"Adding User",
		zap.Any("Main Data", mainData),
		zap.Any("Extra Data", extraData),
	)

	userID, msg, err := usersService.users.Add(mainData, extraData)
	if err != nil {
		log.Error(
			"Adding User Returned Error",
			zap.Any("Main Data", mainData),
			zap.Any("Extra Data", extraData),
			zap.Error(err),
		)
		return errors.Wrap(err, "Error While Adding User")
	}

	log.Info(
		"Added User",
		zap.Any("Main Data", mainData),
		zap.Any("Extra Data", extraData),
	)
	response.Message = msg
	response.UserID = userID

	return nil
}

//Get is used to handle Get function
func (usersService UsersService) Get(
	ctx context.Context, request *proto.GetRequest, response *proto.UserMain) error {

	log.Debug("Getting User's  Main Data", zap.String("UserID", request.UserID))

	var data config.UserMain
	err := usersService.users.Get(request.UserID, config.DBConfigUserMainDataCollection, &data)

	if err != nil {
		log.Error(
			"Getting User's Main Data Returned Error",
			zap.String("UserID", request.UserID),
			zap.Error(err),
		)
		return errors.Wrap(err, "Error While Getting User Data")
	}

	err = usersService.copy(&response, &data)
	if err != nil {
		return err
	}

	log.Debug(
		"Got User's Main Data",
		zap.String("UserID", request.UserID),
		zap.Any("Data", response),
	)
	return nil
}

//GetExtra is used to handle GetExtra function
func (usersService UsersService) GetExtra(
	ctx context.Context, request *proto.GetRequest, response *proto.UserExtra) error {

	log.Debug("Getting User's Extra Data", zap.String("UserID", request.UserID))

	var data config.UserMain
	err := usersService.users.Get(request.UserID, config.DBConfigUserMainDataCollection, &data)

	if err != nil {
		log.Error(
			"Getting User's Extra Data Returned Error",
			zap.String("UserID", request.UserID),
			zap.Error(err),
		)
		return errors.Wrap(err, "Error While Getting User Data")
	}

	err = usersService.copy(&response, &data)
	if err != nil {
		return err
	}

	log.Debug(
		"Got User's Extra Data",
		zap.String("UserID", request.UserID),
		zap.Any("Data", response),
	)
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

	log.Debug(
		"Updating User Main Data",
		zap.String("UserID", request.UserID),
		zap.Any("Update", request.Update),
	)

	err = usersService.users.Update(request.UserID,
		update,
		config.DBConfigUserMainDataCollection)

	if err != nil {
		log.Error(
			"Updating User Main Data Returned Error",
			zap.String("UserID", request.UserID),
			zap.Any("Update", request.Update),
			zap.Error(err),
		)
		return errors.Wrap(err, "Error While Updating Data")
	}

	log.Info(
		"Updated User Main Data",
		zap.String("UserID", request.UserID),
		zap.Any("Update", request.Update),
	)

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

	log.Debug(
		"Updating User Extra Data",
		zap.String("UserID", request.UserID),
		zap.Any("Update", request.Update),
	)

	err = usersService.users.Update(request.UserID,
		update,
		config.DBConfigUserExtraDataCollection)

	if err != nil {
		log.Error(
			"Updating User Extra Data Returned Error",
			zap.String("UserID", request.UserID),
			zap.Any("Update", request.Update),
			zap.Error(err),
		)
		return errors.Wrap(err, "Error While Updating Data")
	}

	log.Info(
		"Updated User Extra Data",
		zap.String("UserID", request.UserID),
		zap.Any("Update", request.Update),
	)

	return nil
}

//Delete is used to handle Delete function
func (usersService UsersService) Delete(
	ctx context.Context,
	request *proto.DeleteRequest,
	response *proto.DeleteResponse) error {

	log.Debug("Deleting User", zap.String("Username", request.Username))

	msg, err := usersService.users.Delete(request.UserID, request.Username, request.Password)

	if err != nil {
		log.Error(
			"Deleting User Returned Error",
			zap.String("Username", request.Username),
			zap.Error(err),
		)
		return errors.Wrap(err, "Error While Deleting User")
	}

	log.Info("Deleted User", zap.String("Username", request.Username))
	response.Message = msg
	return err
}
