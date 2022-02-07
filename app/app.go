package app

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"pay_later_service/handlers"
	"pay_later_service/repository"
	"pay_later_service/service"
)

func Start() {
	err := repository.InitDB("root:Global!23@tcp(localhost:3306)/simple_service")

	if err != nil {
		fmt.Println("DB Conn Error occurred: ", err.Error())
		os.Exit(1)
	}
	// TODO: requests validations needs to be added,
	// TODO: DB date to modified to sql.date type, needs to add updated_at in users db

	repo := repository.New(repository.DB)
	serv := service.NewServ(repo)

	handlers.ProcessRequest(serv)

	// Reporting DB calls, TODO: need to route these reports through service
	//discount, _ := repo.GetMerchantDiscountsReport(1)
	//fmt.Println(discount)

	//users, _ := repo.GetUsersAtLowCreditLimit(0)
	//fmt.Printf("%#v", users)

	//users, _ := repo.GetAllUsers()
	//for _, user := range users {
	//	fmt.Println(user.UserName, user.CreditLimitOffered - user.AvailableCreditLimit)
	//}
	//user, _ := repo.GetUser(2)
	//fmt.Println(user.CreditLimitOffered - user.AvailableCreditLimit)
}
