package response

import (
	"ghost-api/cmd/ghost-api/config"
	"log"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"time"

	"github.com/go-faker/faker/v4"
)

func GenerateResponse(node *config.Node) interface{} {
	if node.Value != nil {
		return node.Value
	}
	switch node.Type {
	case "date":
		now := time.Now()
		hour := int(time.Hour / 1_000_000)
		var r int
		if node.Range == nil {
			r = randIntBetween(hour*24*-400, hour*24*400) * int(time.Millisecond)
		} else {
			r = parseRange(node.Range, true)
		}
		date := now.Add(time.Duration(r))

		var layout string
		if node.Metadata == "" {
			layout = time.RFC3339
		} else {
			layout = node.Metadata
		}
		return date.Format(layout)
	case "string":
		s, _ := GenerateString(node)
		return s
	case "number":
		n, _ := GenerateNumber(node)
		return n

	case "object":
		return GenerateObject(node)

	case "array":
		return GenerateArray(node)

	default:
		return nil
	}
}

func GenerateArray(node *config.Node) []interface{} {
	l := parseRange(node.Range, false)
	items := make([]interface{}, l)
	for i := range l {
		items[i] = GenerateResponse(node.Items)
	}
	return items
}

func GenerateObject(node *config.Node) map[string]interface{} {
	items := make(map[string]interface{}, len(node.Fields))
	for k, v := range node.Fields {
		items[k] = GenerateResponse(v)
	}
	return items
}

func GenerateNumber(node *config.Node) (float64, error) {
	randLen := parseRange(node.Range, false)

	if node.Metadata != "" {
		return float64(randLen), nil
	} else {
		s := getRandNum(node.Metadata)
		return s, nil
	}
}

func GenerateString(node *config.Node) (string, error) {
	if node.Metadata == "" {
		var randLen int
		if node.Range == nil {
			randLen = randIntBetween(4, 16)
		} else {
			randLen = parseRange(node.Range, false)
		}

		finalString := ""
		averageLength := 6
		wordLength := float32(randLen) / float32(averageLength)

		for i := 0; i < int(wordLength); i++ {
			word, e := faker.GetLorem().Word(reflect.Value{})

			if e != nil {
				return "", e
			}

			finalString += " " + word.(string)
		}
		return strings.TrimSpace(finalString), nil
	} else {
		s := getRandString(node.Metadata)
		return s, nil
	}
}

func getRandNum(t string) float64 {
	switch t {
	case "latitude":
		return faker.Latitude()

	case "longitude":
		return faker.Longitude()

	case "unix_time":
		return float64(faker.UnixTime())

	default:
		return float64(rand.Intn(100))
	}
}

func getRandString(t string) string {
	switch t {
	case "date":
		return faker.Date()

	case "time_string":
		return faker.TimeString()

	case "month_name":
		return faker.MonthName()

	case "year_string":
		return faker.YearString()

	case "day_of_week":
		return faker.DayOfWeek()

	case "day_of_month":
		return faker.DayOfMonth()

	case "timestamp":
		return faker.Timestamp()

	case "century":
		return faker.Century()

	case "timezone":
		return faker.Timezone()

	case "timeperiod":
		return faker.Timeperiod()

	// Internet
	case "email":
		return faker.Email()

	case "mac_address":
		return faker.MacAddress()

	case "domain_name":
		return faker.DomainName()

	case "url":
		return faker.URL()

	case "username":
		return faker.Username()

	case "Ipv4":
		return faker.IPv4()

	case "IPv6":
		return faker.IPv6()

	case "password":
		return faker.Password()

	// Words and Sentences
	case "word":
		return faker.Word()

	case "sentence":
		return faker.Sentence()

	case "paragraph":
		return faker.Paragraph()

	// Payment
	case "CCType":
		return faker.CCType()

	case "CCNumber":
		return faker.CCNumber()

	case "currency":
		return faker.Currency()

	case "amount_with_currency":
		return faker.AmountWithCurrency()

	// Person
	case "title_male":
		return faker.TitleMale()

	case "title_female":
		return faker.TitleFemale()

	case "first_name":
		return faker.FirstName()

	case "first_name_male":
		return faker.FirstNameMale()

	case "first_name_female":
		return faker.FirstNameFemale()

	case "last_name":
		return faker.LastName()

	case "name":
		return faker.Name()

	// Phone
	case "phonenumber":
		return faker.Phonenumber()

	case "toll_free_phone_number":
		return faker.TollFreePhoneNumber()

	case "E164_phone_number":
		return faker.E164PhoneNumber()

	//  UUID
	case "UUIDHyphenated":
		return faker.UUIDHyphenated()

	case "UUIDDigit":
		return faker.UUIDDigit()

	default:
		return t
	}
}

func parseRange(value []string, isdate bool) int {
	var min, max int
	if !isdate {
		min64, err := strconv.ParseInt(value[0], 10, 32)
		if err != nil {
			min = 0
		}
		max64, err := strconv.ParseInt(value[1], 10, 32)
		if err != nil {
			max = 0
		}
		min = int(min64)
		max = int(max64)
	} else {
		min = int(parseTimeGap(value[0]))
		max = int(parseTimeGap(value[1]))
	}
	return randIntBetween(min, max)
}

func parseTimeGap(s string) int {
	isNeg := false
	if s[0] == byte('-') {
		isNeg = true
	}
	qtyPattern := `(\d*\.\d+|\d+)`
	unitPattern := `(seconds?|s|minutes?|m|hours?|h|days?|d|months?|mo|years?|y)`
	pattern := `(?<qty>` + qtyPattern + `)\s*(?<unit>` + unitPattern + `)\s*`

	re := regexp.MustCompile(pattern)

	match := re.FindStringSubmatch(s)

	if len(match) == 0 {
		log.Fatal("not found")
	}

	qty, err := strconv.ParseInt(match[1], 10, 32)
	if err != nil {
		return 0
	}
	if isNeg {
		qty = -1 * qty
	}
	unit := match[3]
	dur := time.Duration(0)
	switch unit {
	case "seconds", "s":
		dur = time.Duration(qty) * time.Second
	case "minutes", "m":
		dur = time.Duration(qty) * time.Minute
	case "hours", "h":
		dur = time.Duration(qty) * time.Hour
	case "days", "d":
		dur = time.Duration(qty) * 24 * time.Hour
	case "months", "mo", "l":
		dur = time.Duration(qty) * 30 * 24 * time.Hour
	case "years", "y":
		dur = time.Duration(qty) * 365 * 24 * time.Hour
	}

	return int(dur / time.Millisecond)
}

func randIntBetween(min, max int) int {
	return rand.Intn(max-min) + min
}
