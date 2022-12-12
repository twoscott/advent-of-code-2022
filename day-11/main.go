package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	part1Rounds = 20
	part2Rounds = 10000
)

var monkeyRegex = regexp.MustCompile(
	`(?mi)Monkey \d+:\n` +
		`\s*Starting items: ([\d,\s]+?)\n` +
		`\s*Operation: new = ([\w\s\+\-\*\/]+?)\n` +
		`\s*Test: divisible by (\d+)\n` +
		`\s*If true: throw to monkey (\d+)\n` +
		`\s*If false: throw to monkey (\d+)`,
)

type monkeyOperation struct {
	operand1 string
	operator byte
	operand2 string
}

func (op monkeyOperation) resolveOperand(
	oldItem int64, operand string) (resolved int64) {

	if operand == "old" {
		resolved = oldItem
	} else {
		resolved, _ = strconv.ParseInt(operand, 10, 64)
	}
	return
}

func (op monkeyOperation) resolveOperand1(oldItem int64) (resolved int64) {
	return op.resolveOperand(oldItem, op.operand1)
}

func (op monkeyOperation) resolveOperand2(oldItem int64) (resolved int64) {
	return op.resolveOperand(oldItem, op.operand2)
}

func (op monkeyOperation) run(oldItem int64) (result int64) {
	op1 := op.resolveOperand1(oldItem)
	op2 := op.resolveOperand2(oldItem)

	switch op.operator {
	case '+':
		result = op1 + op2
	case '-':
		result = op1 - op2
	case '*':
		result = op1 * op2
	case '/':
		result = op1 / op2
	default:
		panic("unknown operator")
	}
	return
}

type monkeyTest struct {
	divisibleBy int64
	trueTarget  int
	falseTarget int
}

func (t monkeyTest) run(troop []*monkeyThief, worry int64) (target *monkeyThief) {
	if (worry % t.divisibleBy) == 0 {
		target = troop[t.trueTarget]
	} else {
		target = troop[t.falseTarget]
	}
	return
}

type monkeyThief struct {
	items          []int64
	operation      monkeyOperation
	test           monkeyTest
	itemsInspected int64
}

func (m *monkeyThief) takeTurn(troop []*monkeyThief, relief bool, lcm int64) {
	for len(m.items) > 0 {
		item := m.items[0]
		m.items = m.items[1:]
		m.itemsInspected++

		item = m.operation.run(item)
		if relief {
			item /= 3
		} else {
			if item > lcm {
				item %= lcm
			}
		}

		target := m.test.run(troop, item)
		m.throwTo(target, item)
	}
}

func (m *monkeyThief) takeTurnWithRelief(troop []*monkeyThief) {
	m.takeTurn(troop, true, 0)
}

func (m *monkeyThief) takeTurnManaged(
	troop []*monkeyThief, lowestCommonMultiple int64) {

	m.takeTurn(troop, false, lowestCommonMultiple)
}

func (m *monkeyThief) throwTo(target *monkeyThief, item int64) {
	target.catchItem(item)
}

func (m *monkeyThief) catchItem(item int64) {
	m.items = append(m.items, item)
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	input := strings.TrimSpace(string(bytes))

	troop := make([]*monkeyThief, 0)
	monkeyMatches := monkeyRegex.FindAllStringSubmatch(input, -1)
	for _, match := range monkeyMatches {
		var (
			itemsString       = match[1]
			operationString   = match[2]
			testDivideString  = match[3]
			trueMonkeyString  = match[4]
			falseMonkeyString = match[5]
		)

		items := make([]int64, 0)
		for _, item := range strings.Split(itemsString, ",") {
			worryLevel, _ := strconv.ParseInt(strings.TrimSpace(item), 10, 64)

			items = append(items, worryLevel)
		}

		var operation monkeyOperation
		fmt.Sscanf(
			operationString,
			"%s %c %s",
			&operation.operand1,
			&operation.operator,
			&operation.operand2,
		)

		var (
			testDivide, _  = strconv.ParseInt(testDivideString, 10, 64)
			trueMonkey, _  = strconv.Atoi(trueMonkeyString)
			falseMonkey, _ = strconv.Atoi(falseMonkeyString)
		)
		monkey := monkeyThief{
			items:     items,
			operation: operation,
			test: monkeyTest{
				divisibleBy: testDivide,
				trueTarget:  trueMonkey,
				falseTarget: falseMonkey,
			},
		}

		troop = append(troop, &monkey)
	}

	troop2 := make([]*monkeyThief, len(troop))
	for i := range troop {
		c := *troop[i]
		troop2[i] = &c
	}

	part1(troop)
	part2(troop2)
}

func findTwoMostActive(troop []*monkeyThief) (twoMostActive []int64) {
	twoMostActive = make([]int64, 2)
	for _, monkey := range troop {
		for i, activity := range twoMostActive {
			if monkey.itemsInspected > activity {
				twoMostActive = append(
					twoMostActive[:i],
					append(
						[]int64{monkey.itemsInspected},
						twoMostActive[i:len(twoMostActive)-1]...,
					)...,
				)
				break
			}
		}
	}
	return
}

func part1(troop []*monkeyThief) {
	for i := 0; i < part1Rounds; i++ {
		for _, monkey := range troop {
			monkey.takeTurnWithRelief(troop)
		}
	}

	mostActive := findTwoMostActive(troop)
	monkeyBusiness := mostActive[0] * mostActive[1]
	fmt.Println(
		"The level of monkey business after 20 rounds is:",
		monkeyBusiness,
	)
}

func part2(troop []*monkeyThief) {
	lcm := troop[0].test.divisibleBy
	for _, m := range troop[1:] {
		lcm *= m.test.divisibleBy
	}

	for i := 0; i < part2Rounds; i++ {
		for _, monkey := range troop {
			monkey.takeTurnManaged(troop, lcm)
		}
	}

	mostActive := findTwoMostActive(troop)
	monkeyBusiness := mostActive[0] * mostActive[1]
	fmt.Println(
		"The level of monkey business after 10000 rounds with no relief is:",
		monkeyBusiness,
	)
}
