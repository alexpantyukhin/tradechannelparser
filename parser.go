package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ParseResult struct {
	Currency Currency
	Range    Range
	Exchange string
}

type Range struct {
	Left  float64
	Right float64
}

type Currency struct {
	Coin string
	Name string
}

var currencies = []Currency{
	Currency{Coin: "AMP", Name: "Synereo"},
	Currency{Coin: "ARDR", Name: "Ardor"},
	Currency{Coin: "BCN", Name: "Bytecoin"},
	Currency{Coin: "BCY", Name: "BitCrystals"},
	Currency{Coin: "BELA", Name: "Belacoin"},
	Currency{Coin: "BLK", Name: "BlackCoin"},
	Currency{Coin: "BTC", Name: "Bitcoin"},
	Currency{Coin: "BTCD", Name: "BitcoinDark"},
	Currency{Coin: "BTM", Name: "Bitmark"},
	Currency{Coin: "BTS", Name: "BitShares"},
	Currency{Coin: "BURST", Name: "Burst"},
	Currency{Coin: "CLAM", Name: "CLAMS"},
	Currency{Coin: "DASH", Name: "Dash"},
	Currency{Coin: "DCR", Name: "Decred"},
	Currency{Coin: "DGB", Name: "DigiByte"},
	Currency{Coin: "DOGE", Name: "Dogecoin"},
	Currency{Coin: "EMC2", Name: "Einsteinium"},
	Currency{Coin: "ETC", Name: "Ethereum Classic"},
	Currency{Coin: "ETH", Name: "Ethereum"},
	Currency{Coin: "EXP", Name: "Expanse"},
	Currency{Coin: "FCT", Name: "Factom"},
	Currency{Coin: "FLDC", Name: "FoldingCoin"},
	Currency{Coin: "FLO", Name: "Florincoin"},
	Currency{Coin: "GAME", Name: "GameCredits"},
	Currency{Coin: "GNO", Name: "Gnosis"},
	Currency{Coin: "GNT", Name: "Golem"},
	Currency{Coin: "GRC", Name: "Gridcoin Research"},
	Currency{Coin: "HUC", Name: "Huntercoin"},
	Currency{Coin: "LBC", Name: "LBRY Credits"},
	Currency{Coin: "LSK", Name: "Lisk"},
	Currency{Coin: "LTC", Name: "Litecoin"},
	Currency{Coin: "MAID", Name: "MaidSafeCoin"},
	Currency{Coin: "NAUT", Name: "Nautiluscoin"},
	Currency{Coin: "NAV", Name: "NAVCoin"},
	Currency{Coin: "NEOS", Name: "Neoscoin"},
	Currency{Coin: "NMC", Name: "Namecoin"},
	Currency{Coin: "NOTE", Name: "DNotes"},
	Currency{Coin: "NXC", Name: "Nexium"},
	Currency{Coin: "NXT", Name: "NXT"},
	Currency{Coin: "OMNI", Name: "Omni"},
	Currency{Coin: "PASC", Name: "PascalCoin"},
	Currency{Coin: "PINK", Name: "Pinkcoin"},
	Currency{Coin: "POT", Name: "PotCoin"},
	Currency{Coin: "PPC", Name: "Peercoin"},
	Currency{Coin: "RADS", Name: "Radium"},
	Currency{Coin: "REP", Name: "Augur"},
	Currency{Coin: "RIC", Name: "Riecoin"},
	Currency{Coin: "SBD", Name: "Steem Dollars"},
	Currency{Coin: "SC", Name: "Siacoin"},
	Currency{Coin: "SJCX", Name: "Storjcoin X"},
	Currency{Coin: "STEEM", Name: "STEEM"},
	Currency{Coin: "STR", Name: "Stellar"},
	Currency{Coin: "STRAT", Name: "Stratis"},
	Currency{Coin: "SYS", Name: "Syscoin"},
	Currency{Coin: "USDT", Name: "Tether USD"},
	Currency{Coin: "VIA", Name: "Viacoin"},
	Currency{Coin: "VRC", Name: "VeriCoin"},
	Currency{Coin: "VTC", Name: "Vertcoin"},
	Currency{Coin: "XBC", Name: "BitcoinPlus"},
	Currency{Coin: "XCP", Name: "Counterparty"},
	Currency{Coin: "XEM", Name: "NEM"},
	Currency{Coin: "XMR", Name: "Monero"},
	Currency{Coin: "XPM", Name: "Primecoin"},
	Currency{Coin: "XRP", Name: "Ripple"},
	Currency{Coin: "XVC", Name: "Vcash"},
	Currency{Coin: "ZEC", Name: "Zcash"},
	Currency{Coin: "Music", Name: "Musicoin"},
	Currency{Coin: "GUP", Name: "GUP"},
}

var exchanges = []string{"Bittrex", "Poloniex"}
var markerWords = []string{"покупаем", "скупаем", "купить", "покупка", "скупка", "закупаем", "закупка", "берем", "ордер"}

var numberPattern = "[+-]?([0-9]*[.|,])?[0-9]+"
var whitespacePattern = "[\\s]*"
var hyphen = "-"

var twoNumberPattern = fmt.Sprintf("(?P<first>%s)%s(?P<second>%s)",
	numberPattern,
	whitespacePattern+hyphen+whitespacePattern,
	numberPattern)
var oneNumberPattern = fmt.Sprintf("(?P<first>%s)", numberPattern)

var twoNumberRegexpCompiled = regexp.MustCompile(twoNumberPattern)
var oneNumberRegexpCompiled = regexp.MustCompile(oneNumberPattern)

func getGroupsValues(reg *regexp.Regexp, match []string) map[string]string {
	result := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}

	return result
}

func toFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func findRange(message string) (Range, bool) {
	match := twoNumberRegexpCompiled.FindStringSubmatch(message)
	if match != nil {
		result := getGroupsValues(twoNumberRegexpCompiled, match)
		return Range{Left: toFloat(result["first"]), Right: toFloat(result["second"])}, true
	}

	match = oneNumberRegexpCompiled.FindStringSubmatch(message)
	if match != nil {
		result := getGroupsValues(oneNumberRegexpCompiled, match)
		return Range{Left: toFloat(result["first"]), Right: -1}, true
	}

	return Range{Left: -1, Right: -1}, false
}

func findCurrency(mess string) (*Currency, bool) {
	lowerMess := strings.ToLower(mess)

	for _, element := range currencies {
		lowerCoin := strings.ToLower(element.Coin)
		lowerName := strings.ToLower(element.Name)

		if strings.Contains(lowerMess, lowerCoin) || strings.Contains(lowerMess, lowerName) {
			return &element, true
		}
	}

	return nil, false
}

func findExchange(mess string) (string, bool) {
	lowerMess := strings.ToLower(mess)

	for _, element := range exchanges {
		lowerExchange := strings.ToLower(element)

		if strings.Contains(lowerMess, lowerExchange) {
			return element, true
		}
	}

	return "", false
}

func containsMarkerWords(mess string) bool {
	lowerMess := strings.ToLower(mess)

	for _, element := range markerWords {
		lowerMarker := strings.ToLower(element)

		if strings.Contains(lowerMess, lowerMarker) {
			return true
		}
	}

	return false
}

func parse(message string) (*ParseResult, bool) {
	res := ParseResult{}

	rang, isRangeDeterminated := findRange(message)
	if isRangeDeterminated {
		res.Range = rang
	}

	currency, isCurrencyFound := findCurrency(message)
	if isCurrencyFound {
		res.Currency = *currency
	}

	exchange, isExchangeFound := findExchange(message)
	if isExchangeFound {
		res.Exchange = exchange
	}

	if containsMarkerWords(message) && isRangeDeterminated && isCurrencyFound {
		return &res, true
	}

	return nil, false
}

func testParseDiapason(mess string) {
	rang, isOk := findRange(mess)
	if isOk {
		fmt.Println(mess, ":", rang.Left, rang.Right)
	} else {
		fmt.Println(mess, ": fail")
	}
}

func testParse(mess string) {
	res, isOk := parse(mess)
	if isOk {
		fmt.Println(mess, ":", res.Currency.Name, res.Range.Left, res.Range.Right, res.Exchange)
	} else {
		fmt.Println(mess, ": fail")
	}
}

func main() {
	testParse("OMNI 0.0222-0.0225 покупка")
	testParse("Покупка 0.00161-0.00162 BTC")
	testParse("STRATIS - покупка 0.0023-0.00231")
	testParse("GUP (Bittrex) ставим ордер 6700-6800")
	testParse("MUSIC (Bittrex) покупка 835-840")
	testParse("MUSIC (Bittrex) 835-840")

	// testParseDiapason("hello, 123.45-6435.78")
	// testParseDiapason("hello, 123.45")
	// testParseDiapason("hello, 123-6435")
	// testParseDiapason("hello")
	// testParseDiapason("hello, 123.45   - 6435.78")
	// testParseDiapason("ZEC, 123.45   - 6435.78")

}
