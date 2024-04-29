package sqlstore

import (
	"github.com/SerFiLiuZ/MEDODS/internal/app/models"
	"github.com/SerFiLiuZ/MEDODS/internal/app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Receive(guid string, jwtKey []byte) (*models.AuthUser, error) {
	logger := *utils.NewLogger()
	//logger.EnableDebug()

	logger.Debugf("UserRepository: Receive: start GetAccessRefreshTokens")
	au, err := r.store.tokenRepository.GetAccessRefreshTokens(guid, jwtKey)
	if err != nil {
		logger.Errorf("UserRepository: Receive: GetAccessRefreshTokens: " + err.Error())
		return nil, err
	}

	logger.Debugf("UserRepository: Receive: start GetHashedToken")
	logger.Debugf("UserRepository: Receive: GetHashedToken: au.RtSigned: %v", au.RtSigned)
	logger.Debugf("UserRepository: Receive: GetHashedToken: au.RtSigned (size): %v", len([]byte(au.RtSigned)))

	rtHashed, err := r.store.tokenRepository.GetHashedToken(au.RtSigned)

	if err != nil {
		logger.Errorf("UserRepository: Receive: GetHashedToken: " + err.Error())
		return nil, err
	}

	logger.Debugf("UserRepository: Receive: start StartSession")

	var session mongo.Session

	if session, err = r.store.db.StartSession(); err != nil {
		logger.Errorf("UserRepository: Receive: StartSession: " + err.Error())
		return nil, err
	}
	if err = session.StartTransaction(); err != nil {
		logger.Errorf("UserRepository: Receive: StartTransaction: " + err.Error())
		return nil, err
	}

	logger.Debugf("UserRepository: Receive: start WithSession")

	if err = mongo.WithSession(*r.store.ctx, session, func(sc mongo.SessionContext) error {
		collection := r.store.db.Database("test_medods").Collection("users")
		findFilter := bson.M{"guid": guid}
		var result models.User

		err = collection.FindOne(*r.store.ctx, findFilter).Decode(&result)
		isNewUser := err == mongo.ErrNoDocuments
		if err != nil && !isNewUser {
			logger.Errorf("UserRepository: Receive: FindOne: " + err.Error())
			return err
		}

		if isNewUser {
			newRts := make([]string, 0, 1)
			newRts = append(newRts, rtHashed)
			newUser := models.User{GUID: guid, Rts: newRts}

			_, err = collection.InsertOne(*r.store.ctx, newUser)
			if err != nil {
				logger.Errorf("UserRepository: Receive: InsertOne: " + err.Error())
				return err
			}
		} else {
			newRts := make([]string, len(result.Rts), len(result.Rts)+1)
			copy(newRts, result.Rts)
			newRts = append(newRts, rtHashed)
			newUser := models.User{GUID: guid, Rts: newRts}
			updateFilter := bson.D{
				primitive.E{Key: "$set", Value: bson.D{
					primitive.E{Key: "rts", Value: newUser.Rts}}}}
			_, err = collection.UpdateOne(*r.store.ctx, findFilter, updateFilter)
			if err != nil {
				logger.Errorf("UserRepository: Receive: UpdateOne: " + err.Error())
				return err
			}
		}
		return nil

	}); err != nil {
		return nil, err
	}

	return au, nil
}
