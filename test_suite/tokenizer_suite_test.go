package test_suite

import (
	"bytes"
	djs "github.com/Pencroff/JsonStruct"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"testing"
)

type TokenizerTestElement struct {
	idx  string
	in   []byte
	kind djs.TokenizerKind
	out  []byte
	err  error
}

func TestJStruct_Tokenizer(t *testing.T) {
	s := new(TokenizerTestSuite)
	suite.Run(t, s)
}

type TokenizerTestSuite struct {
	suite.Suite
}

func (s *TokenizerTestSuite) SetupTest() {
}

func (s *TokenizerTestSuite) TestTokenizer_Next_null() {
	tbl := []TokenizerTestElement{
		{"null:0", []byte("null"), djs.TokenNull, []byte("null"), nil},
		{"null:1", []byte("null           "), djs.TokenNull, []byte("null"), nil},
		{"null:2", []byte("null\n"), djs.TokenNull, []byte("null"), nil},
		{"null:3", []byte("null\r"), djs.TokenNull, []byte("null"), nil},
		{"null:4", []byte("null\t"), djs.TokenNull, []byte("null"), nil},
		{"null:5", []byte("null\r\n"), djs.TokenNull, []byte("null"), nil},
		{"null:6", []byte("\nnull\t\n"), djs.TokenNull, []byte("null"), nil},
		{"null:7", []byte("\nnull\t\r"), djs.TokenNull, []byte("null"), nil},
		{"null:8", []byte(" null "), djs.TokenNull, []byte("null"), nil},
		{"null:9", []byte(" null\n"), djs.TokenNull, []byte("null"), nil},
		// Invalid cases
		{"null:50", []byte(""), djs.TokenUnknown, []byte(nil), djs.InvalidJsonError{Err: io.EOF}},
		{"null:51", []byte("n"), djs.TokenUnknown, []byte("n"), djs.InvalidJsonPtrError{Pos: 1, Err: io.EOF}},
		{"null:52", []byte("   nill"), djs.TokenUnknown, []byte("ni"), djs.InvalidJsonPtrError{Pos: 4}},
		{"null:53", []byte("nnn"), djs.TokenUnknown, []byte("nn"), djs.InvalidJsonPtrError{Pos: 1}},
		{"null:54", []byte("nnnn"), djs.TokenUnknown, []byte("nn"), djs.InvalidJsonPtrError{Pos: 1}},
		{"null:55", []byte("nulle"), djs.TokenUnknown, []byte("nulle"), djs.InvalidJsonPtrError{Pos: 4}},
		{"null:56", []byte("null\t\t\tnull"), djs.TokenUnknown, []byte("null\t\t\tn"), djs.InvalidJsonPtrError{Pos: 7}},
	}
	for _, el := range tbl {
		RunTokeniserTest(el, s)
	}
}

func (s *TokenizerTestSuite) TestTokenizer_Next_bool() {
	tbl := []TokenizerTestElement{
		// False cases
		{"bool:f00", []byte("false"), djs.TokenFalse, []byte("false"), nil},
		{"bool:f01", []byte(" false "), djs.TokenFalse, []byte("false"), nil},
		// Invalid cases
		{"bool:f50", []byte(" folse "), djs.TokenUnknown, []byte("fo"), djs.InvalidJsonPtrError{Pos: 2}},
		{"bool:f51", []byte("falze"), djs.TokenUnknown, []byte("falz"), djs.InvalidJsonPtrError{Pos: 3}},
		{"bool:f52", []byte("fals"), djs.TokenUnknown, []byte("fals"), djs.InvalidJsonPtrError{Pos: 4, Err: io.EOF}},
		{"bool:f53", []byte("f "), djs.TokenUnknown, []byte("f "), djs.InvalidJsonPtrError{Pos: 1}},
		{"bool:f54", []byte("falsez"), djs.TokenUnknown, []byte("falsez"), djs.InvalidJsonPtrError{Pos: 5}},
		{"bool:f55", []byte("false\t\t\tfalse"), djs.TokenUnknown, []byte("false\t\t\tf"), djs.InvalidJsonPtrError{Pos: 8}},
		// True cases
		{"bool:t00", []byte("true"), djs.TokenTrue, []byte("true"), nil},
		{"bool:t01", []byte("\n\rtrue\n\r"), djs.TokenTrue, []byte("true"), nil},
		// Invalid cases
		{"bool:t50", []byte("truae "), djs.TokenUnknown, []byte("trua"), djs.InvalidJsonPtrError{Pos: 3}},
		{"bool:t51", []byte("trues"), djs.TokenUnknown, []byte("trues"), djs.InvalidJsonPtrError{Pos: 4}},
		{"bool:t52", []byte(" t "), djs.TokenUnknown, []byte("t "), djs.InvalidJsonPtrError{Pos: 2}},
	}
	for _, el := range tbl {
		RunTokeniserTest(el, s)
	}
}

func (s *TokenizerTestSuite) TestTokenizer_Next_number() {
	tbl := []TokenizerTestElement{
		{"num:00", []byte("123"), djs.TokenIntNumber, []byte("123"), nil},
		{"num:01", []byte("0"), djs.TokenIntNumber, []byte("0"), nil},
		{"num:02", []byte("-0"), djs.TokenIntNumber, []byte("-0"), nil},
		{"num:03", []byte("1"), djs.TokenIntNumber, []byte("1"), nil},
		{"num:04", []byte("-1"), djs.TokenIntNumber, []byte("-1"), nil},
		{"num:05", []byte(" -1 "), djs.TokenIntNumber, []byte("-1"), nil},
		{"num:06", []byte("123456789"), djs.TokenIntNumber, []byte("123456789"), nil},
		{"num:07", []byte("-123456789"), djs.TokenIntNumber, []byte("-123456789"), nil},
		{"num:08", []byte("9223372036854775807"), djs.TokenIntNumber, []byte("9223372036854775807"), nil},
		{"num:09", []byte("-9223372036854775808"), djs.TokenIntNumber, []byte("-9223372036854775808"), nil},
		{"num:10", []byte("9223372036854775808"), djs.TokenIntNumber, []byte("9223372036854775808"), nil},
		{"num:11", []byte("-9223372036854775809"), djs.TokenIntNumber, []byte("-9223372036854775809"), nil},   // -9.223372036854776e+18
		{"num:12", []byte("18446744073709551615"), djs.TokenIntNumber, []byte("18446744073709551615"), nil},   // 1.8446744073709552e+19
		{"num:13", []byte("-18446744073709551615"), djs.TokenIntNumber, []byte("-18446744073709551615"), nil}, // -1.8446744073709552e+19
		{"num:14", []byte("\n9064\n\r"), djs.TokenIntNumber, []byte("9064"), nil},
		{"num:15", []byte("340282366920938463463374607431768211455"), djs.TokenIntNumber, []byte("340282366920938463463374607431768211455"), nil},
		// Num errors
		{"num:50", []byte("9 0 6 4"), djs.TokenUnknown, []byte(`9 0`), djs.InvalidJsonPtrError{Pos: 2}},
		{"num:51", []byte("-e"), djs.TokenUnknown, []byte(`-e`), djs.InvalidJsonPtrError{Pos: 1}},
		{"num:52", []byte("25$E1"), djs.TokenUnknown, []byte(`25$`), djs.InvalidJsonPtrError{Pos: 2}},
		{"num:53", []byte("123l1"), djs.TokenUnknown, []byte(`123l`), djs.InvalidJsonPtrError{Pos: 3}},
		{"num:54", []byte("1e"), djs.TokenUnknown, []byte(`1e`), djs.InvalidJsonPtrError{Pos: 1, Err: io.EOF}},
		{"num:55", []byte("1234e  "), djs.TokenUnknown, []byte(`1234e `), djs.InvalidJsonPtrError{Pos: 5}},
		{"num:56", []byte("11$!"), djs.TokenUnknown, []byte(`11$`), djs.InvalidJsonPtrError{Pos: 2}},
		{"num:57", []byte("- 123"), djs.TokenUnknown, []byte(`- `), djs.InvalidJsonPtrError{Pos: 1}},
	}
	for _, el := range tbl {
		RunTokeniserTest(el, s)
	}
}

func (s *TokenizerTestSuite) TestTokenizer_Next_float() {
	tbl := []TokenizerTestElement{
		{"float:00", []byte("123.45"), djs.TokenFloatNumber, []byte("123.45"), nil},
		{"float:01", []byte("0.0"), djs.TokenFloatNumber, []byte("0.0"), nil},
		{"float:02", []byte("-0.0"), djs.TokenFloatNumber, []byte("-0.0"), nil},
		{"float:03", []byte("1.0"), djs.TokenFloatNumber, []byte("1.0"), nil},
		{"float:04", []byte("-1.0"), djs.TokenFloatNumber, []byte("-1.0"), nil},
		{"float:05", []byte("3.1415"), djs.TokenFloatNumber, []byte("3.1415"), nil},
		{"float:06", []byte("-3.1415"), djs.TokenFloatNumber, []byte("-3.1415"), nil},
		{"float:07", []byte("3.141592653589793238462643383279502884197169"), djs.TokenFloatNumber, []byte("3.141592653589793238462643383279502884197169"), nil},   // 3.141592653589793
		{"float:08", []byte("-3.141592653589793238462643383279502884197169"), djs.TokenFloatNumber, []byte("-3.141592653589793238462643383279502884197169"), nil}, // -3.141592653589793
		{"float:09", []byte("3.141592653589793238462643383279502884197169e15"), djs.TokenFloatNumber, []byte("3.141592653589793238462643383279502884197169e15"), nil},
		{"float:10", []byte("-141592653589793238462643383279502884197169e+10"), djs.TokenFloatNumber, []byte("-141592653589793238462643383279502884197169e+10"), nil},
		{"float:11", []byte("3.141592653589793238462643383279502884197169e-10"), djs.TokenFloatNumber, []byte("3.141592653589793238462643383279502884197169e-10"), nil},
		{"float:12", []byte("-3.141592653589793238462643383279502884197169e-10"), djs.TokenFloatNumber, []byte("-3.141592653589793238462643383279502884197169e-10"), nil},
		{"float:13", []byte("92653589793238462643383279502884197169e-10"), djs.TokenFloatNumber, []byte("92653589793238462643383279502884197169e-10"), nil},
		{"float:14", []byte("-926535897932384626433.83279502884197169e-10"), djs.TokenFloatNumber, []byte("-926535897932384626433.83279502884197169e-10"), nil},
		{"float:15", []byte(" 3.1415E5 "), djs.TokenFloatNumber, []byte("3.1415E5"), nil},
		{"float:16", []byte("\n-3.1415E+5\n"), djs.TokenFloatNumber, []byte("-3.1415E+5"), nil},
		{"float:17", []byte("-3.1415E-5"), djs.TokenFloatNumber, []byte("-3.1415E-5"), nil},
		{"float:18", []byte("3.1415E-5"), djs.TokenFloatNumber, []byte("3.1415E-5"), nil},
		{"float:19", []byte("1.6180339887498948482045868343656381e999"), djs.TokenFloatNumber, []byte("1.6180339887498948482045868343656381e999"), nil},
		{"float:20", []byte("-1.6180339887498948482045868343656381e-999"), djs.TokenFloatNumber, []byte("-1.6180339887498948482045868343656381e-999"), nil},
		{"float:21", []byte("0.01"), djs.TokenFloatNumber, []byte("0.01"), nil},
		{"float:22", []byte("-0.01"), djs.TokenFloatNumber, []byte("-0.01"), nil},
		{"float:23", []byte(" 0.01 "), djs.TokenFloatNumber, []byte("0.01"), nil},
		{"float:24", []byte("\n\r-0.01\n\r"), djs.TokenFloatNumber, []byte("-0.01"), nil},
		{"float:25", []byte("0.1e-1"), djs.TokenFloatNumber, []byte("0.1e-1"), nil},
		// Float errors
		{"float:50", []byte("-"), djs.TokenUnknown, []byte("-"), djs.InvalidJsonPtrError{Pos: 0, Err: io.EOF}},
		{"float:51", []byte("-e"), djs.TokenUnknown, []byte(`-e`), djs.InvalidJsonPtrError{Pos: 1}},
		{"float:52", []byte("0."), djs.TokenUnknown, []byte("0."), djs.InvalidJsonPtrError{Pos: 1, Err: io.EOF}},
		{"float:53", []byte("0.e"), djs.TokenUnknown, []byte("0.e"), djs.InvalidJsonPtrError{Pos: 2}},
		{"float:54", []byte("0.e1"), djs.TokenUnknown, []byte("0.e"), djs.InvalidJsonPtrError{Pos: 2}},
		{"float:55", []byte("0.1e"), djs.TokenUnknown, []byte("0.1e"), djs.InvalidJsonPtrError{Pos: 3, Err: io.EOF}},
		{"float:56", []byte(".01"), djs.TokenUnknown, []byte(nil), djs.InvalidJsonError{}},
		{"float:57", []byte("123.4l1"), djs.TokenUnknown, []byte("123.4l"), djs.InvalidJsonPtrError{Pos: 5}},
		{"float:58", []byte("-3."), djs.TokenUnknown, []byte("-3."), djs.InvalidJsonPtrError{Pos: 2, Err: io.EOF}},
		{"float:59", []byte("-3.e"), djs.TokenUnknown, []byte("-3.e"), djs.InvalidJsonPtrError{Pos: 3}},
		{"float:60", []byte("-3.e1"), djs.TokenUnknown, []byte("-3.e"), djs.InvalidJsonPtrError{Pos: 3}},
		{"float:61", []byte("-3.1e"), djs.TokenUnknown, []byte("-3.1e"), djs.InvalidJsonPtrError{Pos: 4, Err: io.EOF}},
		{"float:62", []byte("3.1415926535.89793"), djs.TokenUnknown, []byte("3.1415926535."), djs.InvalidJsonPtrError{Pos: 12}},
		{"float:63", []byte("3.14159265Ee589793"), djs.TokenUnknown, []byte("3.14159265Ee"), djs.InvalidJsonPtrError{Pos: 11}},
		{"float:64", []byte("3.14159265E+"), djs.TokenUnknown, []byte("3.14159265E+"), djs.InvalidJsonPtrError{Pos: 11, Err: io.EOF}},
		{"float:65", []byte("3.14159265E-"), djs.TokenUnknown, []byte("3.14159265E-"), djs.InvalidJsonPtrError{Pos: 11, Err: io.EOF}},
		{"float:66", []byte("161803398.874989opq8204e28"), djs.TokenUnknown, []byte("161803398.874989o"), djs.InvalidJsonPtrError{Pos: 16}},
		{"float:67", []byte("16180.3398.874989e8204e+28"), djs.TokenUnknown, []byte("16180.3398."), djs.InvalidJsonPtrError{Pos: 10}},
	}
	for _, el := range tbl {
		RunTokeniserTest(el, s)
	}
}

func (s *TokenizerTestSuite) TestTokenizer_Next_string() {
	testCases := []TokenizerTestElement{
		{"str:00", []byte(`""`), djs.TokenString, []byte(`""`), nil},
		{"str:01", []byte(`"abc"`), djs.TokenString, []byte(`"abc"`), nil},
		{"str:02", []byte(` "abc" `), djs.TokenString, []byte(`"abc"`), nil},
		{"str:03", []byte(`
									"abc"
								`), djs.TokenString, []byte(`"abc"`), nil},
		{"str:04", []byte(`"abc xyz"`), djs.TokenString, []byte(`"abc xyz"`), nil},
		{"str:05", []byte(`"hello world!"`), djs.TokenString, []byte(`"hello world!"`), nil},
		{"str:06", []byte(`"The quick brown fox jumps over the lazy dog"`), djs.TokenString, []byte(`"The quick brown fox jumps over the lazy dog"`), nil},
		{"str:07", []byte(`"a\"z"`), djs.TokenString, []byte{0x22, 0x61, 0x5c, 0x22, 0x7a, 0x22}, nil},
		{"str:08", []byte(`"a\\z"`), djs.TokenString, []byte(`"a\\z"`), nil},
		{"str:09", []byte(`"a\/z"`), djs.TokenString, []byte(`"a\/z"`), nil},
		{"str:10", []byte(`"a/z"`), djs.TokenString, []byte(`"a/z"`), nil},
		{"str:11", []byte(`"a\bz"`), djs.TokenString, []byte(`"a\bz"`), nil},
		{"str:12", []byte(`"a\fz"`), djs.TokenString, []byte(`"a\fz"`), nil},
		{"str:13", []byte(`"a\nz"`), djs.TokenString, []byte(`"a\nz"`), nil},
		{"str:14", []byte(`"a\rz"`), djs.TokenString, []byte(`"a\rz"`), nil},
		{"str:15", []byte(`"a\tz"`), djs.TokenString, []byte(`"a\tz"`), nil},
		{"str:16", []byte(`"abc\u00A0xyz"`), djs.TokenString, []byte(`"abc\u00A0xyz"`), nil},
		{"str:17", []byte(`"abc\u002Fxyz"`), djs.TokenString, []byte(`"abc\u002Fxyz"`), nil},
		{"str:18", []byte(`"abc\u002fxyz"`), djs.TokenString, []byte(`"abc\u002fxyz"`), nil},
		{"str:19", []byte(`"\u2070"`), djs.TokenString, []byte(`"\u2070"`), nil},
		{"str:20", []byte(`"\u0008"`), djs.TokenString, []byte(`"\u0008"`), nil},
		{"str:21", []byte(`"\u000C"`), djs.TokenString, []byte(`"\u000C"`), nil},
		{"str:22", []byte(`"\uD834\uDD1E"`), djs.TokenString, []byte(`"\uD834\uDD1E"`), nil},
		{"str:23", []byte(`"D'fhuascail ??osa, ??rmhac na h??ighe Beannaithe, p??r ??ava agus ??dhaimh"`), djs.TokenString, []byte(`"D'fhuascail ??osa, ??rmhac na h??ighe Beannaithe, p??r ??ava agus ??dhaimh"`), nil},
		{"str:24", []byte(`"????????????????????????????????????"`), djs.TokenString, []byte(`"????????????????????????????????????"`), nil},

		// String errors
		{"str:50", []byte(`"abc`), djs.TokenUnknown, []byte(`"abc`), djs.InvalidJsonPtrError{Pos: 3, Err: io.EOF}},
		{"str:51", []byte(`"abc"xyz`), djs.TokenUnknown, []byte(`"abc"x`), djs.InvalidJsonPtrError{Pos: 5}},
		{"str:52", []byte(`abc"`), djs.TokenUnknown, []byte(nil), djs.InvalidJsonError{}},
		{"str:53", []byte(`"""`), djs.TokenUnknown, []byte(`"""`), djs.InvalidJsonPtrError{Pos: 2}},
		{"str:54", []byte(`""\"`), djs.TokenUnknown, []byte(`""\`), djs.InvalidJsonPtrError{Pos: 2}},
		{"str:55", []byte(`"\u2O70"`), djs.TokenUnknown, []byte(`"\u2O`), djs.InvalidJsonPtrError{Pos: 4, Err: djs.InvalidHexNumberError}},
		{"str:56", []byte(`"\uD8Y4\uDU1E"`), djs.TokenUnknown, []byte(`"\uD8Y`), djs.InvalidJsonPtrError{Pos: 5, Err: djs.InvalidHexNumberError}},
		// Skip invalid UTF-8 characters. Will validate it on next level
		// by strconv.Unquote / strconv.Quote.
		// {"str:57", []byte(`"\uD834\q"`), djs.TokenUnknown, []byte(`"\uD834\q`), djs.InvalidJsonPtrError{Pos: 8}},
		// {"str:25", []byte(`"\uD834\n"`), djs.TokenString, []byte(`"\uD834\n`), djs.InvalidJsonPtrError{Pos: 8}},
	}
	for _, el := range testCases {
		RunTokeniserTest(el, s)
	}
}

func (s *TokenizerTestSuite) TestTokenizer_Next_time() {
	tbl := []TokenizerTestElement{
		{"time:00", []byte(`"2015-05-14T12:34:56+02:00"`), djs.TokenTime, []byte(`"2015-05-14T12:34:56+02:00"`), nil},
		{"time:01", []byte(`"2015-05-14T12:34:56.3+02:00"`), djs.TokenTime, []byte(`"2015-05-14T12:34:56.3+02:00"`), nil},
		{"time:02", []byte(`"2015-05-14T12:34:56.37+02:00"`), djs.TokenTime, []byte(`"2015-05-14T12:34:56.37+02:00"`), nil},
		{"time:03", []byte(`"2015-05-14T12:34:56.379+02:00"`), djs.TokenTime, []byte(`"2015-05-14T12:34:56.379+02:00"`), nil},
		{"time:04", []byte(`"1970-01-01T00:00:00Z"`), djs.TokenTime, []byte(`"1970-01-01T00:00:00Z"`), nil},
		{"time:05", []byte(`"0001-01-01T00:00:00Z"`), djs.TokenTime, []byte(`"0001-01-01T00:00:00Z"`), nil},
		{"time:06", []byte(`"1985-04-12T23:20:50.52Z"`), djs.TokenTime, []byte(`"1985-04-12T23:20:50.52Z"`), nil},
		{"time:07", []byte(`"1996-12-19T16:39:57-08:00"`), djs.TokenTime, []byte(`"1996-12-19T16:39:57-08:00"`), nil},
		{"time:08", []byte(`"1990-12-31T23:59:60Z"`), djs.TokenTime, []byte(`"1990-12-31T23:59:60Z"`), nil},
		{"time:09", []byte(`"1990-12-31T15:59:60-08:00"`), djs.TokenTime, []byte(`"1990-12-31T15:59:60-08:00"`), nil},
		{"time:10", []byte(`"1937-01-01T12:00:27.87+00:20"`), djs.TokenTime, []byte(`"1937-01-01T12:00:27.87+00:20"`), nil},
		{"time:11", []byte(`"2022-02-24T04:00:00+02:00"`), djs.TokenTime, []byte(`"2022-02-24T04:00:00+02:00"`), nil},
		{"time:12", []byte(`"2022-07-12T21:55:16+01:00"`), djs.TokenTime, []byte(`"2022-02-24T04:00:00+02:00"`), nil},
		{"time:13", []byte(`"2015-05-14T12:34:56.123Z"`), djs.TokenTime, []byte(`"2015-05-14T12:34:56.123Z"`), nil},
		// Invalid cases fall back to string
		{"time:50", []byte(`"2015-05-14E12:34:56.379+02:00"`), djs.TokenString, []byte(`"2015-05-14E12:34:56.379+02:00"`), nil},
		{"time:51", []byte(`"2O15-O5-14T12:34:56.379+02:00"`), djs.TokenString, []byte(`"2O15-O5-14T12:34:56.379+02:00"`), nil},
		{"time:52", []byte(`"1985-04-12T23:20:50.52ZZZZ"`), djs.TokenString, []byte(`"1985-04-12T23:20:50.52ZZZZ"`), nil},
		{"time:53", []byte(`"2022-07-12 21:55:16"`), djs.TokenString, []byte(`"2022-07-12 21:55:16"`), nil},
		{"time:54", []byte(`"20220712T215516Z"`), djs.TokenString, []byte(`"20220712T215516Z"`), nil},
		{"time:55", []byte(`"20220712T215516+01:00"`), djs.TokenString, []byte(`"20220712T215516+01:00"`), nil},
		{"time:56", []byte(`"1985-04-12T23:20:50.Z"`), djs.TokenString, []byte(`"1985-04-12T23:20:50.Z"`), nil},
		{"time:56", []byte(`"not a Timestamp"`), djs.TokenString, []byte(`"not a Timestamp"`), nil},
	}
	for _, el := range tbl {
		RunTokeniserTest(el, s)
	}
}

func RunTokeniserTest(el TokenizerTestElement, s *TokenizerTestSuite) {
	s.T().Run(el.idx, func(t *testing.T) {
		b := bytes.NewBuffer(el.in)
		sc := djs.NewJStructScanner(b)
		tk := djs.NewJSStructTokenizer(sc)
		e := tk.Next()
		v := tk.Value()
		k := tk.Kind()

		assert.Equal(t, el.out, v, "%s Value %v != %v (%v)", el.idx, v, el.out, el.in)
		assert.Equal(t, el.kind, k, "%s Kind %v != %v", el.idx, k, el.kind)
		assert.ErrorIs(t, e, el.err, "%s Next err: %v", el.idx, e)
	})
}
