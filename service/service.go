package service

import (
	"fmt"
	"pay_later_service/models"
	"pay_later_service/repository"
	"pay_later_service/utils"
)

type simpleService struct {
	repo repository.RepoInterface
}

func NewServ(store repository.RepoInterface) SimpleService {
	return &simpleService{repo: store}
}

func (s simpleService) CreateUser(userName, userEmail string, creditLimitOffered float32) (models.User, error) {
	newUser := models.User{
		UserName:             userName,
		UserEmail:            userEmail,
		CreditLimitOffered:   creditLimitOffered,
		AvailableCreditLimit: creditLimitOffered,
		CreatedAt:            utils.CurrentMillis(),
	}
	if err := s.repo.AddUser(&newUser); err != nil {
		fmt.Println("error: ", err.Error())
		return newUser, err
	}

	return newUser, nil
}

func (s simpleService) FetchUserDetails(userID int) (models.User, error) {
	userDetails, err := s.repo.GetUser(userID)

	if err != nil {
		fmt.Println(err.Error())
		return userDetails, err
	}

	return userDetails, nil
}

func (s simpleService) CreateMerchant(merchantName string, discountPercent float32) (models.Merchant, error) {
	merchant := models.Merchant{
		MerchantName:    merchantName,
		DiscountPercent: discountPercent,
		CreatedAt:       utils.CurrentMillis(),
		UpdatedAt:       utils.CurrentMillis(),
	}
	if err := s.repo.AddMerchant(&merchant); err != nil {
		fmt.Println("error: ", err.Error())
		return merchant, err
	}

	return merchant, nil
}

func (s simpleService) UpdateMerchantDiscount(merchantID int, discountPercent float32) error {
	merchant := models.Merchant{
		MerchantID:      merchantID,
		DiscountPercent: discountPercent,
		UpdatedAt:       utils.CurrentMillis(),
	}

	if err := s.repo.UpdateMerchant(merchant); err != nil {
		return err
	}
	fmt.Println("Merchant updated")
	return nil
}

func (s simpleService) HandleUserOrder(userID, merchantID int, transactionAmount float32) (models.Order, error) {
	var order models.Order

	user, err := s.repo.GetUser(userID)

	fmt.Println("err", err, user)

	if err != nil {
		fmt.Println("Error in fetching user")
		return order, err
	}

	merchant, err := s.repo.GetMerchant(merchantID)

	if err != nil {
		fmt.Println("Error in fetching merchant")
		return order, err
	}

	// TODO: save it in order and mark status as failed and reason as insufficient credit
	if transactionAmount > user.AvailableCreditLimit {
		fmt.Println("insufficient credit limit")
		return order, fmt.Errorf("insufficient credit limit")
	}
	discountedAmount, _ := calculateDiscountedAmount(transactionAmount, merchant.DiscountPercent)
	newCreditLimit := user.AvailableCreditLimit - transactionAmount

	order, err = s.repo.HandleUserOrder(user, merchantID, transactionAmount, discountedAmount, newCreditLimit)

	if err != nil {
		return order, err
	}

	return order, nil
}

func (s simpleService) HandleUserTransaction(userID int, paymentAmount float32) (models.User, error) {
	user, err := s.repo.GetUser(userID)

	if err != nil {
		fmt.Println("Error in fetching user")
		return user, err
	}

	user.AvailableCreditLimit += paymentAmount

	err = s.repo.HandleUserPayment(&user, paymentAmount)

	return user, nil
}

func calculateDiscountedAmount(transactionAmount, merchantDiscount float32) (float32, error) {
	discountedAmount := transactionAmount * (1 - merchantDiscount/100)

	return discountedAmount, nil

}
