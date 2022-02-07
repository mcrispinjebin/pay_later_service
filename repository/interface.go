package repository

import "pay_later_service/models"

type RepoInterface interface {
	AddMerchant(merchant *models.Merchant) error
	UpdateMerchant(merchant models.Merchant) error
	GetMerchant(merchantID int) (models.Merchant, error)

	AddUser(userDetails *models.User) error
	GetUser(userID int) (models.User, error)

	HandleUserPayment(user *models.User, paymentAmount float32) error
	HandleUserOrder(user models.User, merchantID int, orderAmount, discountedAmount, newAvailableCreditLimit float32) (models.Order, error)

	GetMerchantDiscountsReport(merchantID int) (float32, error)
	GetUsersAtLowCreditLimit(thresholdLimit float32) ([]models.User, error)
	GetAllUsers() ([]models.User, error)
}
