package accounts

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
)

// GetAllCmds calls the GetAllCmds method and returns all the users commands.
func (s *userService) GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Var(APIKey, "uuid")
	if validateErr != nil {
		s.log.Errorf("could not validate GET ALL CMDS request: %v", validateErr)
		return nil, errors.NewBadRequestError("request format incorrect.")
	}
	cmds, err := s.cache.GetAllCmds(ctx, APIKey)
	if err != nil {
		s.log.Error("could not get user from cache")
	} else {
		s.log.Info("got user from cache")
		return cmds, nil
	}
	cmds, apierr := s.db.GetAllCmds(reqCtx, APIKey)
	return cmds, apierr
}

// AddCmd calls the AddCmd method and returns the number of updated commands.
func (s *userService) AddCmd(ctx context.Context, requestData request.AddCmd, APIKey string) (int, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Struct(requestData)
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("could not validate ADD CMD request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	numUpdated, err := s.db.AddCmd(reqCtx, requestData, APIKey)
	s.cache.DeleteCmds(ctx, APIKey)
	return numUpdated, err
}

// DeleteCmd calls the DelCmd method and returns the number of updated commands.
func (s *userService) DeleteCmd(ctx context.Context, requestData request.DeleteCmd, APIKey string) (int, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Struct(requestData)
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("could not validate DELETE CMD request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	numUpdated, err := s.db.DeleteCmd(reqCtx, requestData, APIKey)
	s.cache.DeleteCmds(ctx, APIKey)
	return numUpdated, err
}
