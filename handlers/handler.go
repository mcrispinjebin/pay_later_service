package handlers

import (
	"fmt"
	"os"
	"pay_later_service/service"
	"strconv"
)

func ProcessRequest(serviceInterface service.SimpleService) {
	request := os.Args[1:]
	RouteRequest(request, serviceInterface)
}

func RouteRequest(request []string, serviceInterface service.SimpleService) {
	switch request[0] {
	case "new":
		switch request[1] {
		case "user":
			creditLimit, err := strconv.ParseFloat(request[4], 32)
			if err != nil {
				fmt.Println("Error Occurred: ", err.Error())
			}
			serviceInterface.CreateUser(request[2], request[3], float32(creditLimit))

		case "merchant":
			discountPercentStr := request[3][:len(request[3])-1]
			discountPercent, err := strconv.ParseFloat(discountPercentStr, 32)
			if err != nil {
				fmt.Println("Error Occurred: ", err.Error())
			}

			serviceInterface.CreateMerchant(request[2], float32(discountPercent))

		case "txn":
			transactionAmount, err := strconv.ParseFloat(request[4], 32)
			userID, err := strconv.Atoi(request[2])
			merchantID, err := strconv.Atoi(request[3])
			if err != nil {
				fmt.Println("Error Occurred: ", err.Error())
			}

			serviceInterface.HandleUserOrder(userID, merchantID, float32(transactionAmount))

		default:
			fmt.Println("Invalid syntax")
		}

	case "update":
		if request[1] == "merchant" {
			discountPercentStr := request[3][:len(request[3])-1]
			discountPercent, err := strconv.ParseFloat(discountPercentStr, 32)
			merchantID, err := strconv.Atoi(request[2])
			if err != nil {
				fmt.Println("Error Occurred: ", err.Error())
			}

			serviceInterface.UpdateMerchantDiscount(merchantID, float32(discountPercent))
		} else {
			fmt.Println("Invalid request")
		}

	case "payback":
		paymentAmount, err := strconv.ParseFloat(request[2], 32)
		userID, err := strconv.Atoi(request[1])
		if err != nil {
			fmt.Println("Error Occurred: ", err.Error())
		}

		serviceInterface.HandleUserTransaction(userID, float32(paymentAmount))
	default:
		fmt.Println("Invalid syntax")
	}
}
