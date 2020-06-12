package logging

import (
	"context"
	"github.com/sirupsen/logrus"
	"gophr.v2/user"
)

func Apply(svc user.Service) user.Service {
	return &loggingDecorator{svc}
}

type loggingDecorator struct {
	svc user.Service
}

func (l *loggingDecorator) GetByID(ctx context.Context, id interface{}) (*user.User, error) {
	logrus.Infof("METHOD: GetByID ID: %v\n", id)
	return l.svc.GetByID(ctx, id)
}

func (l *loggingDecorator) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	logrus.Infof("METHOD: GetByEmail EMAIL: %v\n", email)
	return l.svc.GetByEmail(ctx, email)
}

func (l *loggingDecorator) GetByUserID(ctx context.Context, id string) (*user.User, error) {
	logrus.Infof("METHOD: GetByUserID USERID: %v\n", id)
	return l.svc.GetByUserID(ctx, id)
}

func (l *loggingDecorator) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	logrus.Infof("METHOD: GetByUsername USERNAME: %v\n", username)
	return l.svc.GetByUsername(ctx, username)
}

func (l *loggingDecorator) Save(ctx context.Context, usr *user.User) error {
	logrus.Infof("METHOD: Save "+
		"USERID: %v USERNAME: %v EMAIL: %v\n",
		usr.UserID, usr.Username, usr.Email)
	return l.svc.Save(ctx, usr)
}

func (l *loggingDecorator) GetAll(ctx context.Context, cursor string, num int) (user []*user.User, nextCursor string, err error) {
	logrus.Infof("METHOD: GetAll CURSOR: %v NUM: %v\n", cursor, num)
	return l.svc.GetAll(ctx, cursor, num)
}

func (l *loggingDecorator) Delete(ctx context.Context, id interface{}) error {
	logrus.Infof("METHOD: Delete ID: %v\n", id)
	return l.svc.Delete(ctx, id)
}

func (l *loggingDecorator) Update(ctx context.Context, usr *user.User) error {
	logrus.Infof("METHOD: Update "+
		"USERID: %v USERNAME: %v EMAIL: %v\n",
		usr.UserID, usr.Username, usr.Email)
	return l.svc.Update(ctx, usr)
}

func (l *loggingDecorator) Register(ctx context.Context, usr *user.User) error {
	logrus.Infof("METHOD: Register "+
		"USERID: %v USERNAME: %v EMAIL: %v\n",
		usr.UserID, usr.Username, usr.Email)
	return l.svc.Register(ctx, usr)
}

func (l *loggingDecorator) Login(ctx context.Context, usr *user.User) error {
	logrus.Infof("METHOD: Login "+
		"USERID: %v USERNAME: %v EMAIL: %v\n",
		usr.UserID, usr.Username, usr.Email)
	return l.svc.Login(ctx, usr)

}
