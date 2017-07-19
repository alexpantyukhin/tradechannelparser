// Package contains methods for parsing channel signals
package tradechannelparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ParseResult contains parsed info
type ParseResult struct {
	Currency Currency
	Range    Range
	Exchange string
}

// Range contains range left and right bounds
type Range struct {
	Left  float64
	Right float64
}

// Currency contains info about currency
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
var markerWords = []string{"покупаем", "покупай", "скупаем", "скупай", "купить", "покупка", "скупка", "закупаем", "закупка", "берем", "ордер"}

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
	withReplacedComma := strings.Replace(s, ",", ".", -1)
	f, _ := strconv.ParseFloat(withReplacedComma, 64)
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

// Parse tries to parse the string and returns (ParseResult, true) if success and (nil, false) else
func Parse(message string) (*ParseResult, bool) {
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
