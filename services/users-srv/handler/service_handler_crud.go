package handler

import (
	"context"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	"github.com/SonicRoshan/Velocity/global/config"
	proto "github.com/SonicRoshan/Velocity/services/users-srv/proto"
)

//Add is used to handle Add function
func (usersService UsersService) Add(
	ctx context.Context,
	request *proto.AddRequest,
	response *proto.AddResponse) error {

	log.Debug(
		"Adding User",
		zap.Any("Request", request),
	)

	var mainData config.UserMain
	var extraData config.UserExtra

	copier.Copy(&mainData, &request.MainData)
	copier.Copy(&extraData, &request.ExtraData)

	log.Debug(
		"Adding User",
		zap.Any("Main Data", mainData),
		zap.Any("Extra Data", extraData),
	)

	userID, msg := usersService.users.Add(mainData, extraData)
	if msg != "" {
		response.Message = msg
		return nil
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

	log.Debug("Getting User's Main Data", zap.String("UserID", request.UserID))

	var data config.UserMain
	err := usersService.users.Get(request.UserID, config.DBConfigUserMainDataCollection, &data)

	if err != nil {
		return usersService.errHandler.Check(
			err,
			log,
			"Getting User's Main Data Returned Error",
			zap.String("UserID", request.UserID),
		).(error)
	}

	err = copier.Copy(&response, &data)
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

	var data config.UserExtra
	err := usersService.users.Get(request.UserID, config.DBConfigUserExtraDataCollection, &data)

	if err != nil {
		return usersService.errHandler.Check(
			err,
			log,
			"Getting User's Extra Data Returned Error",
			zap.String("UserID", request.UserID),
		).(error)
	}

	log.Info(
		"Got User's Extra Data From Low Level",
		zap.String("UserID", request.UserID),
		zap.Any("Data", data),
	)

	err = copier.Copy(&response, &data)
	if err != nil {
		return err
	}

	log.Info(
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

	usersService.users.Update(request.UserID,
		update,
		config.DBConfigUserMainDataCollection)

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
	request *proto.UpdateExtraRequest,
	response *proto.UpdateResponse) error {

	log.Debug("Updating User Extra Data")
	var update config.UserExtra

	err := copier.Copy(&update, &request.Update)
	if err != nil {
		return err
	}

	log.Debug(
		"Updating User Extra Data",
		zap.String("UserID", request.UserID),
		zap.Any("Update", request.Update),
	)

	usersService.users.Update(request.UserID,
		update,
		config.DBConfigUserExtraDataCollection)

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

	msg := usersService.users.Delete(request.UserID, request.Username, request.Password)

	log.Info("Deleted User", zap.String("Username", request.Username))
	response.Message = msg
	return nil
}

//Activate is used handle Activate function
func (usersService UsersService) Activate(
	ctx context.Context,
	request *proto.ActivateRequest,
	response *proto.ActivateResponse) error {

	log.Debug("Activating Account", zap.String("Email", request.Email))

	msg := usersService.users.Activate(request.Email)

	log.Info("Activated Account", zap.String("Email", request.Email))
	response.Message = msg
	return nil
}
