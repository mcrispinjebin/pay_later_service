package service

import "pay_later_service/models"

type SimpleService interface {
	CreateUser(userName, userEmail string, creditLimitOffered float32) (models.User, error)
	FetchUserDetails(userID int) (models.User, error)
	CreateMerchant(merchantName string, discountPercent float32) (models.Merchant, error)
	UpdateMerchantDiscount(merchantID int, discountPercent float32) error
	HandleUserOrder(userID, merchantID int, transactionAmount float32) (models.Order, error)
	HandleUserTransaction(userID int, paymentAmount float32) (models.User, error)
}
