package tradechannelparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMessageWithMarkerCurrencyAndRange(t *testing.T) {
	mess := "OMNI 0.0222-0.0225 покупка"

	res, isOk := Parse(mess)
	assert.Equal(t, isOk, true)
	assert.NotNil(t, res.Currency)
	assert.NotNil(t, res.Range)
	assert.Equal(t, res.Range.Left, 0.0222)
	assert.Equal(t, res.Range.Right, 0.0225)
	assert.Equal(t, res.Currency.Coin, "OMNI")
}

func TestParseMessageWithMarkerCurrencyAndRangeWithComma(t *testing.T) {
	mess := "Покупка 0,00161-0,00162 BTC"

	res, isOk := Parse(mess)
	assert.Equal(t, isOk, true)
	assert.NotNil(t, res.Currency)
	assert.NotNil(t, res.Range)
	assert.Equal(t, res.Range.Left, 0.00161)
	assert.Equal(t, res.Range.Right, 0.00162)
}

func TestParseMessageWithMarkerCurrencyRangeAndExchange(t *testing.T) {
	mess := "MUSIC (Bittrex) покупка 835-840"

	res, isOk := Parse(mess)
	assert.Equal(t, isOk, true)
	assert.NotNil(t, res.Currency)
	assert.NotNil(t, res.Range)
	assert.Equal(t, res.Range.Left, 835.0)
	assert.Equal(t, res.Range.Right, 840.0)
	assert.Equal(t, res.Exchange, "Bittrex")
}

func TestParseMessageWithoutMarker(t *testing.T) {
	mess := "MUSIC (Bittrex) 835-840"

	_, isOk := Parse(mess)
	assert.Equal(t, isOk, false)
}
