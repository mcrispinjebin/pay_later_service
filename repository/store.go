package repository

import (
	"database/sql"
	"fmt"
	"pay_later_service/models"
	"pay_later_service/utils"
)

type repoStruct struct {
	DB *sql.DB
}

func New(db *sql.DB) RepoInterface {
	return &repoStruct{DB: db}
}

func (r repoStruct) AddMerchant(merchant *models.Merchant) error {
	queryStr := fmt.Sprintf("INSERT INTO merchant(merchant_name, discount_percent, created_at, updated_at) VALUES ('%s', %v, %v, %v)", merchant.MerchantName, merchant.DiscountPercent, merchant.CreatedAt, merchant.UpdatedAt)

	result, err := r.DB.Exec(queryStr)

	if err != nil {
		return err
	}
	insertedID, err := result.LastInsertId()

	if err != nil {
		return err
	}

	fmt.Println("Merchant Added with ID", insertedID)
	merchant.MerchantID = int(insertedID)

	return nil
}

func (r repoStruct) UpdateMerchant(merchant models.Merchant) error {
	queryStr := fmt.Sprintf("UPDATE merchant SET discount_percent=%v, updated_at=%v WHERE merchant_id=%d", merchant.DiscountPercent, merchant.UpdatedAt, merchant.MerchantID)

	result, err := r.DB.Exec(queryStr)

	if err != nil {
		return err
	}
	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsUpdated != 1 {
		return fmt.Errorf("no records matched to be updated")
	}

	return nil
}

func (r repoStruct) GetMerchant(merchantID int) (models.Merchant, error) {
	merchant := models.Merchant{}

	query := `SELECT * FROM merchant where merchant_id=?;`
	row := r.DB.QueryRow(query, merchantID)

	err := row.Scan(&merchant.MerchantID, &merchant.MerchantName, &merchant.DiscountPercent, &merchant.CreatedAt, &merchant.UpdatedAt)

	if err != nil {
		return merchant, err
	}
	return merchant, nil
}

func (r repoStruct) AddUser(userDetails *models.User) error {
	queryStr := fmt.Sprintf("INSERT INTO user(user_name, user_email, credit_limit_offered, available_credit_limit, created_at) VALUES ('%s', '%s', %v, %v, %v)", userDetails.UserName, userDetails.UserEmail, userDetails.CreditLimitOffered, userDetails.CreditLimitOffered, userDetails.CreatedAt)

	result, err := r.DB.Exec(queryStr)

	if err != nil {
		return err
	}
	insertedID, err := result.LastInsertId()

	if err != nil {
		return err
	}

	fmt.Println("User Added with ID", insertedID)
	userDetails.UserID = int(insertedID)

	return nil
}

func (r repoStruct) GetUser(userID int) (models.User, error) {
	fetchedUser := models.User{}

	query := `SELECT * FROM user where user_id=?;`
	row := r.DB.QueryRow(query, userID)

	err := row.Scan(&fetchedUser.UserID, &fetchedUser.UserName, &fetchedUser.CreditLimitOffered, &fetchedUser.AvailableCreditLimit, &fetchedUser.CreatedAt, &fetchedUser.UserEmail)

	if err != nil {
		fmt.Println(err.Error())
		return fetchedUser, err
	}
	return fetchedUser, nil
}

func (r repoStruct) HandleUserPayment(user *models.User, paymentAmount float32) error {
	currentTime := utils.CurrentMillis()

	tx, err := r.DB.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()
	ledgerInsertQuery := fmt.Sprintf("INSERT INTO ledger(user_id, amount, status, created_at) VALUES (%v, %v, '%s', %v)", user.UserID, paymentAmount, "success", currentTime)

	result, err := tx.Exec(ledgerInsertQuery)

	if err != nil {
		return err
	}
	insertedID, err := result.LastInsertId()

	if err != nil {
		return err
	}
	fmt.Println("Ledger Added with ID", insertedID)

	userUpdateQuery := fmt.Sprintf("UPDATE user SET available_credit_limit=%v WHERE user_id=%v", user.AvailableCreditLimit, user.UserID)

	userUpdateResult, err := tx.Exec(userUpdateQuery)

	if err != nil {
		return err
	}

	affectedRows, _ := userUpdateResult.RowsAffected()

	fmt.Println("User rows updated", affectedRows)

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error in committing transaction")
	}
	//user.UpdatedAt = ""

	return nil
}

func (r repoStruct) HandleUserOrder(user models.User, merchantID int, orderAmount, discountedAmount, newAvailableCreditLimit float32) (models.Order, error) {
	var order models.Order

	currentTime := utils.CurrentMillis()

	tx, err := r.DB.Begin()

	if err != nil {
		return order, err
	}
	defer tx.Rollback()

	orderInsertQuery := fmt.Sprintf("INSERT INTO orders(user_id, merchant_id, order_amount, order_status, created_at) VALUES (%v, %v, %v, '%s', %v)", user.UserID, merchantID, orderAmount, "success", currentTime)

	result, err := tx.Exec(orderInsertQuery)

	if err != nil {
		return order, err
	}
	orderInsertedID, err := result.LastInsertId()

	if err != nil {
		return order, err
	}
	fmt.Println("Order Added with ID", orderInsertedID)

	payoutInsertQuery := fmt.Sprintf("INSERT INTO payout(order_id, payout_amount, payout_status, created_at) VALUES (%v, %v, '%s', %v)", orderInsertedID, discountedAmount, "success", currentTime)

	payoutResult, err := tx.Exec(payoutInsertQuery)

	if err != nil {
		return order, err
	}
	insertedPayoutID, err := payoutResult.LastInsertId()

	if err != nil {
		return order, err
	}
	fmt.Println("Payout Added with ID", insertedPayoutID)

	userUpdateQuery := fmt.Sprintf("UPDATE user SET available_credit_limit=%v WHERE user_id=%v", newAvailableCreditLimit, user.UserID)

	userUpdateResult, err := tx.Exec(userUpdateQuery)

	if err != nil {
		return order, err
	}

	affectedRows, _ := userUpdateResult.RowsAffected()

	fmt.Println("User rows updated", affectedRows)

	if err = tx.Commit(); err != nil {
		return order, fmt.Errorf("error in committing transaction")
	}

	return order, nil
}

func (r repoStruct) GetMerchantDiscountsReport(merchantID int) (float32, error) {
	var totalDiscounts float32

	query := "SELECT SUM(o.order_amount - p.payout_amount) FROM orders o JOIN payout p ON o.order_id=p.order_id WHERE o.merchant_id=?;"
	row := r.DB.QueryRow(query, merchantID)

	err := row.Scan(&totalDiscounts)

	if err != nil {
		fmt.Println(err.Error())
		return totalDiscounts, err
	}
	return totalDiscounts, nil
}

func (r repoStruct) GetUsersAtLowCreditLimit(thresholdLimit float32) ([]models.User, error) {
	users := make([]models.User, 0)

	query := "SELECT user_id, user_email, user_name, credit_limit_offered, available_credit_limit, created_at FROM user " +
		"WHERE available_credit_limit BETWEEN 0 AND ?;"
	rows, err := r.DB.Query(query, thresholdLimit)

	if err != nil {
		fmt.Println(err.Error())
		return users, err
	}

	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.UserID, &user.UserEmail, &user.UserName, &user.CreditLimitOffered, &user.AvailableCreditLimit, &user.CreatedAt); err != nil {
			fmt.Println(err.Error())
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return users, err
	}

	return users, nil
}

func (r repoStruct) GetAllUsers() ([]models.User, error) {
	users := make([]models.User, 0)
	query := "SELECT user_id, user_email, user_name, credit_limit_offered, available_credit_limit, created_at FROM user;"

	rows, err := r.DB.Query(query)

	if err != nil {
		fmt.Println(err.Error())
		return users, err
	}

	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.UserID, &user.UserEmail, &user.UserName, &user.CreditLimitOffered, &user.AvailableCreditLimit, &user.CreatedAt); err != nil {
			fmt.Println(err.Error())
			return users, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return users, err
	}

	return users, nil
}
