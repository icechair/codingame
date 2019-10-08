package main

import (
	"bufio"
	"fmt"
	"os"
)

type Customer struct {
	item  string
	award int
}

func NewCustomer(item string, award int) *Customer {
	c := new(Customer)
	c.award = award
	c.item = item
	return c
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var numAllCustomers int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &numAllCustomers)
	customers := make([]*Customer, numAllCustomers)
	for i := 0; i < numAllCustomers; i++ {
		// customerItem: the food the customer is waiting for
		// customerAward: the number of points awarded for delivering the food
		var customerItem string
		var customerAward int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &customerItem, &customerAward)
		customers[i] = NewCustomer(customerItem, customerAward)
	}
	for i := 0; i < 7; i++ {
		scanner.Scan()
		kitchenLine := scanner.Text()
		fmt.Fprintln(os.Stderr, kitchenLine)
	}
	for {
		var turnsRemaining int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &turnsRemaining)

		var playerX, playerY int
		var playerItem string
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &playerX, &playerY, &playerItem)

		var partnerX, partnerY int
		var partnerItem string
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &partnerX, &partnerY, &partnerItem)

		// numTablesWithItems: the number of tables in the kitchen that currently hold an item
		var numTablesWithItems int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &numTablesWithItems)

		for i := 0; i < numTablesWithItems; i++ {
			var tableX, tableY int
			var item string
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &tableX, &tableY, &item)
		}
		// ovenContents: ignore until wood 1 league
		var ovenContents string
		var ovenTimer int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &ovenContents, &ovenTimer)

		// numCustomers: the number of customers currently waiting for food
		var numCustomers int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &numCustomers)

		for i := 0; i < numCustomers; i++ {
			var customerItem string
			var customerAward int
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &customerItem, &customerAward)
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")

		// MOVE x y
		// USE x y
		// WAIT
		fmt.Println("USE 5 0")
	}
}
