package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
)

type CreditCardUser struct {
	gorm.Model
	Name        string
	CreditCards []CreditCard `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreditCard struct {
	gorm.Model
	Number string
	Bank   string
	UserID uint
}

func main() {
	fmt.Println("Starting go-example")
	//fmt.Println("POSTGRES_PASSWORD: ", os.Getenv("POSTGRES_PASSWORD"))
	//fmt.Println("DATABASE_PASSWORD: ", os.Getenv("DATABASE_PASSWORD"))

	//https://gorm.io/docs/connecting_to_the_database.html
	port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DATABASE_URL"), port, os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_DATABASE"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal(err)
	}

	// remove old entries
	db.Migrator().DropTable(&CreditCardUser{}, &CreditCard{})

	// create table with correct schema
	db.AutoMigrate(&CreditCardUser{})
	db.AutoMigrate(&CreditCard{})

	// insert new record
	db.Create(&CreditCardUser{Name: "mrFlux", CreditCards: []CreditCard{{Number: "1234567898", Bank: "FinFisher"}, {Number: "345657881", Bank: "MaxedOut Limited"}}})
	db.Create(&CreditCardUser{Name: "sirTuxedo", CreditCards: []CreditCard{{Number: "999999999", Bank: "FinFisher"}, {Number: "2342", Bank: "Bankxter"}}})
	db.Create(&CreditCardUser{Name: "missFraudinger", CreditCards: []CreditCard{{Number: "999999999", Bank: "FinFisher"}}})
	db.Create(&CreditCardUser{Name: "happyUser"})
	db.Create(&CreditCardUser{Name: "mrGone", CreditCards: []CreditCard{{Number: "77777777777", Bank: "BICrupt"}}})

	//////////// 1 - get all credit card records of user 'mrFlux' ////////////
	fmt.Println("---1-----------------------------------")
	creditCardsOfFlux := []CreditCardUser{}
	db.Preload("CreditCards").Where("name=?", "mrFlux").Find(&creditCardsOfFlux)
	fmt.Println("The credit cards of mrFlux are: ", creditCardsOfFlux)

	//////////// 2 - get all FinFisher Credit Card records of user 'mrFlux' ////////////
	fmt.Println("---2-----------------------------------")
	finFisherCreditCards := []CreditCard{}
	db.Joins("INNER JOIN credit_card_users ccu ON ccu.id = credit_cards.user_id").Where("ccu.name = ? AND credit_cards.bank = ?", "mrFlux", "FinFisher").Find(&finFisherCreditCards)
	fmt.Println("mrFlux's FinFisher card(s) are (request 1): ", finFisherCreditCards)

	// alternatively using preload for the same result
	mrFluxUser := CreditCardUser{}
	db.Preload("CreditCards", "bank = ?", "FinFisher").First(&mrFluxUser, "name =?", "mrFlux")
	fmt.Println("mrFlux's FinFisher card(s) are (request 2): ", mrFluxUser.CreditCards)

	//////////// 3 - update wrong creditcard number of the sirTuxedo's Bankxter card number from 2342 to 23422342 ////////////
	fmt.Println("---3-----------------------------------")
	// FIXME does not work
	op := db.Model(&CreditCard{}).Joins("INNER JOIN credit_card_users ccu ON ccu.id = credit_cards.user_id").Where("ccu.name = ? AND credit_cards.bank = ?", "sirTuxedo", "Bankxter").Update("number", "23422342")

	if op.Error != nil {
		fmt.Println("Couldn't update credit card number: ", op.Error)
	}

	//////////// 4 -  list all user(s) with a credit card from 'FinFisher' Bank ////////////
	fmt.Println("---4-----------------------------------")
	extractUserNamesFromUsers := func(users *[]CreditCardUser) string {
		s := ""
		for i := 0; i < len(*users); i++ {
			if i > 0 {
				s += ", "
			}
			s += (*users)[i].Name
		}
		return s
	}
	users := []CreditCardUser{}
	db.Joins("INNER JOIN credit_cards cc ON cc.user_id = credit_card_users.id").Where("cc.bank = ?", "FinFisher").Find(&users)
	fmt.Println(" all user(s) with a credit card from 'FinFisher' Bank: ", extractUserNamesFromUsers(&users))

	//////////// 5 - drop all fraudy creditcards from related uses where the card number is 999999999, no matter the bank name ////////////
	fmt.Println("---5-----------------------------------")
	// basically delete sirTuxedo and missFraudinger
	db.Where("number = ?", "999999999").Unscoped().Delete(&CreditCard{})

	//////////// 6 - add a creditcard to happyUser ////////////
	fmt.Println("---6-----------------------------------")
	happyUser := CreditCardUser{}
	db.Model(&CreditCardUser{}).Where("name=?", "happyUser").First(&happyUser)
	happyUser.CreditCards = []CreditCard{{Number: "666666666666", Bank: "happyBank"}}
	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&happyUser)

	creditCardsOfHappyUser := CreditCardUser{}
	db.Preload("CreditCards").Where("name=?", "happyUser").First(&creditCardsOfHappyUser)
	fmt.Println("The credit cards of HappyUser are: ", creditCardsOfHappyUser.CreditCards)

	//////////// 7 - append another entry in the the creditcard(s) of happyUser ////////////
	fmt.Println("---7-----------------------------------")
	happyUser2 := CreditCardUser{}
	db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&CreditCardUser{}).Where("name=?", "happyUser").First(&happyUser2)

		happyUser2.CreditCards = append(happyUser2.CreditCards, CreditCard{Number: "666666666666", Bank: "happyhappyBank"})
		tx.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&happyUser2)
		// return nil will commit the whole transaction
		return nil
	})

	creditCardsOfHappyUser2 := CreditCardUser{}
	db.Preload("CreditCards").Where("name=?", "happyUser").First(&creditCardsOfHappyUser2)
	fmt.Println("The credit cards of HappyUser are: ", creditCardsOfHappyUser2.CreditCards)

	//////////// 8 - delete user with associated creditcard(s) ////////////
	fmt.Println("---8-----------------------------------")
	db.Unscoped().Delete(&CreditCardUser{}, "name = ?", "mrGone")

	fmt.Println("Exiting program")
}
