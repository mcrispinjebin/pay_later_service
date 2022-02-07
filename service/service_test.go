package service

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"pay_later_service/mocks"
	"pay_later_service/models"
	"pay_later_service/repository"
	"testing"
	"github.com/undefinedlabs/go-mpatch"
)

func TestSimpleService_HandleUserOrder(t *testing.T) {
	var monkeyPatch *mpatch.Patch

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB := mocks.NewMockRepoInterface(mockCtrl)

	type repoFields struct {
		repo repository.RepoInterface
	}

	type args struct {
		UserID            int
		MerchantID        int
		TransactionAmount float32
	}

	tests := []struct {
		name           string
		fields         repoFields
		args           args
		wantErr        string
		patchUser      func()
		patchMerchant  func()
		patchOrders    func()
		pathchDiscount func()
	}{
		{"happy flow", repoFields{mockDB}, args{1, 101, 100}, "",
			func() {
				mockDB.EXPECT().GetUser(gomock.Eq(1)).Return(models.User{UserID: 1, AvailableCreditLimit: 1000}, nil)
			},
			func() {
				mockDB.EXPECT().GetMerchant(gomock.Eq(101)).Return(models.Merchant{MerchantID: 101, DiscountPercent: 10}, nil)
			},
			func() {
				mockDB.EXPECT().HandleUserOrder(models.User{UserID: 1, AvailableCreditLimit: 1000}, gomock.Eq(101), gomock.Eq(float32(100)), gomock.Eq(float32(90)), gomock.Eq(float32(900))).Return(models.Order{}, nil)
			},
			func() {
				var err error
				monkeyPatch, err = mpatch.PatchMethod(calculateDiscountedAmount, func(amount, discPercent float32) (float32, error) {
					return 100, nil
				})
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{"insufficient Credit", repoFields{mockDB}, args{1, 101, 1100}, "insufficient credit limit",
			func() {
				mockDB.EXPECT().GetUser(gomock.Eq(1)).Return(models.User{UserID: 1, AvailableCreditLimit: 1000}, nil)
			},
			func() {
				mockDB.EXPECT().GetMerchant(gomock.Eq(101)).Return(models.Merchant{MerchantID: 101, DiscountPercent: 10}, nil)
			},
			func() {},
			func() {
				var err error
				monkeyPatch, err = mpatch.PatchMethod(calculateDiscountedAmount, func(amount, discPercent float32) (float32, error) {
					return 90, nil
				})
				if err != nil {
					t.Fatal(err)
				}
			},

		},
		{"incorrect discount processed", repoFields{mockDB}, args{1, 101, 100}, "",
			func() {
				mockDB.EXPECT().GetUser(gomock.Eq(1)).Return(models.User{UserID: 1, AvailableCreditLimit: 1000}, nil)
			},
			func() {
				mockDB.EXPECT().GetMerchant(gomock.Eq(101)).Return(models.Merchant{MerchantID: 101, DiscountPercent: 10}, nil)
			},
			func() {
				mockDB.EXPECT().HandleUserOrder(models.User{UserID: 1, AvailableCreditLimit: 1000}, gomock.Eq(101), gomock.Eq(float32(100)), gomock.Eq(float32(100)), gomock.Eq(float32(900))).Return(models.Order{}, nil)
			},
			func() {
				var err error
				monkeyPatch, err = mpatch.PatchMethod(calculateDiscountedAmount, func(amount, discPercent float32) (float32, error) {
					return 100, nil
				})
				if err != nil {
					t.Fatal(err)
				}
			},
		},

		{"user DB error", repoFields{mockDB}, args{1, 101, 100}, "Error in fetching user",
			func() {
				mockDB.EXPECT().GetUser(gomock.Eq(1)).Return(models.User{UserID: 2, AvailableCreditLimit: 1}, errors.New("Error in fetching user"))
			},
			func() {},
			func() {},
			func() {},
		},

		{"merchant DB error", repoFields{mockDB}, args{1, 101, 100}, "Error in fetching merchant",
			func() {
				mockDB.EXPECT().GetUser(gomock.Eq(1)).Return(models.User{UserID: 1, AvailableCreditLimit: 100}, nil)
			},
			func() {
				mockDB.EXPECT().GetMerchant(gomock.Eq(101)).Return(models.Merchant{MerchantID: 101, DiscountPercent: 10}, fmt.Errorf("Error in fetching merchant"))
			},
			func() {},
			func() {},
		},
	}

	s := NewServ(mockDB)

	for _, subtest := range tests {
		t.Run(subtest.name, func(t *testing.T) {
			fmt.Println("sub", subtest)
			subtest.patchUser()
			subtest.patchMerchant()
			subtest.patchOrders()
			subtest.pathchDiscount()
			fmt.Printf("%#v", monkeyPatch)
			monkeyPatch.Patch()
			if subtest.wantErr != "" {
				_, err := s.HandleUserOrder(subtest.args.UserID, subtest.args.MerchantID, subtest.args.TransactionAmount)
				if err == nil {
					t.Errorf("HandleUserOrder() with args %v, %v, %v : FAILED, expected  error %v but got error value nil", subtest.args.UserID, subtest.args.MerchantID, subtest.args.TransactionAmount, subtest.wantErr)
				} else if err.Error() != subtest.wantErr {
					t.Errorf("HandleUserOrder() with args %v, %v, %v : FAILED, expected  error %v but got error value %v", subtest.args.UserID, subtest.args.MerchantID, subtest.args.TransactionAmount, subtest.wantErr, err.Error())
				} else {
					t.Logf("HandleUserOrder() with args %v, %v, %v : PASSED, expected error %v and got error value %v", subtest.args.UserID, subtest.args.MerchantID, subtest.args.TransactionAmount, subtest.wantErr, err.Error())
				}
			} else {
				s.HandleUserOrder(subtest.args.UserID, subtest.args.MerchantID, subtest.args.TransactionAmount)
			}
			err := monkeyPatch.Unpatch()
			if err != nil {
				fmt.Println("err2", err)
			}
		})
	}

}
