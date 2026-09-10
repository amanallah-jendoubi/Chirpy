package main

import (
	"errors"
	"fmt"
	tinytime "github.com/wagslane/go-tinytime"
	"time"
)

func costCalc() {
	var insufficientFundMessage string = "Purchase failed. Insufficient funds."
	var purchaseSuccessMessage string = "Purchase successful."
	var accountBalance float64 = 100.0
	var bulkMessageCost float64 = 75.0
	var isPremiumUser bool = true
	var discountRate float64 = 0.10
	var finalCost float64
	finalCost = bulkMessageCost
	if isPremiumUser {
		finalCost *= (1 - discountRate)
	}
	if accountBalance-finalCost >= 0 {
		accountBalance -= finalCost
		fmt.Println(purchaseSuccessMessage)
	} else {
		fmt.Println(insufficientFundMessage)
	}
	fmt.Println("Account balance:", accountBalance)
}

func userLog() {
	fname := "Dalinar"
	lname := "Kholin"
	age := 45
	messageRate := 0.5
	isSubscribed := false
	message := "Sometimes a hypocrite is nothing more than a man in the process of changing."

	userLog := fmt.Sprintf("Name: %s %s, Age: %d, Rate: %.1f, Is Subscribed: %t, Message: %s",
		fname,
		lname,
		age,
		messageRate,
		isSubscribed,
		message)
	fmt.Println(userLog)
}

func subscription() {
	const premiumPlanName = "Premium Plan"
	const basicPlanName = "Basic Plan"
	fmt.Println("plan:", premiumPlanName)
	fmt.Println("plan:", basicPlanName)
}

func monthlyBillIncrease(costPerSend, numLastMonth, numThisMonth int) int {
	var lastMonthBill int
	var thisMonthBill int
	lastMonthBill = getBillForMonth(costPerSend, numLastMonth)
	thisMonthBill = getBillForMonth(costPerSend, numThisMonth)
	return thisMonthBill - lastMonthBill
}

func getBillForMonth(costPerSend, messagesSent int) int {
	return (costPerSend * messagesSent)
}

func getProductMessage(tier string) string {
	quantityMsg, priceMsg, _ := getProductInfo(tier)
	return "You get " + quantityMsg + " for " + priceMsg + "."
}

// don't touch below this line

func getProductInfo(tier string) (string, string, string) {
	switch tier {
	case "basic":
		return "1,000 texts per month", "$30 per month", "most popular"
	case "premium":
		return "50,000 texts per month", "$60 per month", "best value"
	case "enterprise":
		return "unlimited texts per month", "$100 per month", "customizable"
	default:
		return "", "", ""
	}
}

type insuranceStatus struct {
	hasInsurance bool
	isTotaled    bool
	isDented     bool
	isBigDent    bool
}

func getInsuranceAmount(status insuranceStatus) int {
	if !status.hasInsurance {
		return 1
	}
	if status.isTotaled {
		return 10000
	}
	if !status.isDented {
		return 160
	}
	if status.isBigDent {
		return 270
	}
	return 0
}

func printReports(intro, body, outro string) {
	printCostReport(func(message string) int {
		return (len(message) * 2)
	}, intro)
	printCostReport(func(message string) int {
		return (len(message) * 3)
	}, body)
	printCostReport(func(message string) int {
		return (len(message) * 4)
	}, outro)
}

func printCostReport(costCalculator func(string) int, message string) {
	cost := costCalculator(message)
	fmt.Printf(`Message: "%s" Cost: %v cents`, message, cost)
	fmt.Println()
}

func bootup() {
	defer func() { fmt.Println("Bootup complete!") }()
	ok := connectToDB()
	if !ok {
		return
	}
	ok = connectToPaymentProvider()
	if !ok {
		return
	}
	fmt.Println("All systems ready!")
}

var shouldConnectToDB = true

func connectToDB() bool {
	fmt.Println("Connecting to database...")
	if shouldConnectToDB {
		fmt.Println("Connected!")
		return true
	}
	fmt.Println("Connection failed")
	return false
}

var shouldConnectToPaymentProvider = true

func connectToPaymentProvider() bool {
	fmt.Println("Connecting to payment provider...")
	if shouldConnectToPaymentProvider {
		fmt.Println("Connected!")
		return true
	}
	fmt.Println("Connection failed")
	return false
}

func test(dbSuccess, paymentSuccess bool) {
	shouldConnectToDB = dbSuccess
	shouldConnectToPaymentProvider = paymentSuccess
	bootup()
	fmt.Println("====================================")
}

type messageToSend struct {
	message string
	user
}

type user struct {
	name   string
	number int
	membership
}
type membership struct {
	Type             string
	messageCharLimit int
}

func newUser(name string, membershipType string) user {
	var messageCharLimit int
	if membershipType == "premium" {
		messageCharLimit = 1000
	} else {
		messageCharLimit = 100
	}
	return user{
		name: name,
		membership: membership{
			Type:             membershipType,
			messageCharLimit: messageCharLimit,
		},
	}
}

func (u user) sendMessage(message string, messageLength int) (string, bool) {
	if messageLength > 0 {
		return message, true
	}
	return "", false
}

func canSendMessage(mToSend messageToSend) bool {
	return true
}

type authenticationInfo struct {
	username string
	password string
}

func (a authenticationInfo) authenticate() {
	fmt.Printf("Authenticating user: %s:%s", a.username, a.password)
}

func sendMessage(msg message) (string, int) {
	return msg.getMessage(), len(msg.getMessage()) * 3
}

type message interface {
	getMessage() string
}

type birthdayMessage struct {
	birthdayTime  time.Time
	recipientName string
}

func (bm birthdayMessage) getMessage() string {
	return fmt.Sprintf("Hi %s, it is your birthday on %s", bm.recipientName, bm.birthdayTime.Format(time.RFC3339))
}

type sendingReport struct {
	reportName    string
	numberOfSends int
}

func (sr sendingReport) getMessage() string {
	return fmt.Sprintf(`Your "%s" report is ready. You've sent %v messages.`, sr.reportName, sr.numberOfSends)
}

type employee interface {
	getName() string
	getSalary() int
}

type contractor struct {
	name         string
	hourlyPay    int
	hoursPerYear int
}

func (c contractor) getName() string {
	return c.name
}

func (c contractor) getSalary() int {
	return c.hourlyPay * c.hoursPerYear
}

type fullTime struct {
	name   string
	salary int
}

func (ft fullTime) getSalary() int {
	return ft.salary
}

func (ft fullTime) getName() string {
	return ft.name
}

func employeeSalaryReport(e employee) {
	fmt.Printf("%s's salary is $%d\n", e.getName(), e.getSalary())
}

func sendSMSToCouple(msgToCustomer, msgToSpouse string) (int, error) {
	cost1, err1 := sendSMS(msgToCustomer)
	if err1 != nil {
		return 0, err1
	}
	cost2, err2 := sendSMS(msgToSpouse)
	if err2 != nil {
		return 0, err2
	}
	return cost1 + cost2, nil
}

func sendSMS(message string) (int, error) {
	const maxTextLen = 25
	const costPerChar = 2
	if len(message) > maxTextLen {
		return 0, fmt.Errorf("can't send texts over %v characters", maxTextLen)
	}
	return costPerChar * len(message), nil
}

func validateStatus(status string) error {
	if len(status) == 0 {
		return errors.New("status cannot be empty")
	} else if len(status) > 140 {
		return errors.New("status cannot be longer than 140 characters")
	}
	return nil
}

// connecions handling
func countConnections(groupSize int) int {
	return (groupSize * (groupSize - 1)) / 2
}

type cost struct {
	day   int
	value float64
}

func getDayCost(costs []cost, day int) []float64 {
	dayCosts := []float64{}
	for _, c := range costs {
		if c.day == day {
			dayCosts = append(dayCosts, c.value)
		}
	}
	return dayCosts
}

type Message interface {
	Type() string
}

type TextMessage struct {
	Sender  string
	Content string
}

func (tm TextMessage) Type() string {
	return "text"
}

type MediaMessage struct {
	Sender    string
	MediaType string
	Content   string
}

func (mm MediaMessage) Type() string {
	return "media"
}

type LinkMessage struct {
	Sender  string
	URL     string
	Content string
}

func (lm LinkMessage) Type() string {
	return "link"
}

// Don't touch above this line

func filterMessages(messages []Message, filterType string) []Message {
	filtredMessages := []Message{}
	for _, message := range messages {
		if message.Type() == filterType {
			filtredMessages = append(filtredMessages, message)
		}
	}
	return filtredMessages
}

func isValidPassword(password string) bool {
	if len(password) < 5 {
		return false
	}
	hasDigit := false
	hasUppercase := false
	for i := 0; i < len(password); i++ {
		if password[i] >= '0' && password[i] <= '9' {
			hasDigit = true
		}
		if password[i] >= 'A' && password[i] <= 'Z' {
			hasUppercase = true
		}
	}
	return hasDigit && hasUppercase
}

type Analytics struct {
	MessagesTotal     int
	MessagesFailed    int
	MessagesSucceeded int
}

type msg struct {
	Recipient string
	Success   bool
}

func analyzeMessages(msg msg, analytics *Analytics) {
	if msg.Success {
		analytics.MessagesSucceeded++
	} else {
		analytics.MessagesFailed++
	}
	analytics.MessagesTotal++
}

func (e *email) setMessage(newMessage string) {
	e.message = newMessage
}

type email struct {
	message     string
	fromAddress string
	toAddress   string
}
type customer struct {
	id      int
	balance float64
}

type transactionType string

const (
	transactionDeposit    transactionType = "deposit"
	transactionWithdrawal transactionType = "withdrawal"
)

type transaction struct {
	customerID      int
	amount          float64
	transactionType transactionType
}

func updateBalance(c *customer, t transaction) error {
	if t.transactionType == transactionDeposit {
		c.balance += t.amount
		return nil
	} else if t.transactionType == transactionWithdrawal {
		if c.balance < t.amount {
			return errors.New("insufficient funds")
		}
		c.balance -= t.amount
		return nil
	}
	return errors.New("invalid transaction type")
}

func main() {
	test(true, true)
	authInfo := authenticationInfo{
		username: "AJ",
		password: "password123",
	}
	authInfo.authenticate()
	messageStart := "Happy birthday! You are now"
	age := 21
	messageEnd := "years old!"
	fmt.Println(messageStart, age, messageEnd)
	subscription()
	userLog()
	printReports(
		"Welcome to the Hotel California",
		"Such a lovely place",
		"Plenty of room at the Hotel California",
	)
	fullTimeEmployee := fullTime{
		name:   "John Doe",
		salary: 60000,
	}
	contractorEmployee := contractor{
		name:         "Jane Smith",
		hourlyPay:    50,
		hoursPerYear: 2000,
	}
	employeeSalaryReport(fullTimeEmployee)
	employeeSalaryReport(contractorEmployee)
	fmt.Println(countConnections(5))
	tt := tinytime.New(1585750374)
	tt = tt.Add(time.Hour * 48)
	fmt.Println("1585750374 converted to a tinytime is:", tt)
}
