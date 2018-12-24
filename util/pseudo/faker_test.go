package pseudo

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/corestoreio/pkg/util/assert"
)

type SomeStruct struct {
	Inta    int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Float32 float32
	Float64 float64

	UInta  uint
	UInt8  uint8
	UInt16 uint16
	UInt32 uint32
	UInt64 uint64

	Latitude           float32 `faker:"lat"`
	LATITUDE           float64 `faker:"lat"`
	Long               float32 `faker:"long"`
	LONG               float64 `faker:"long"`
	String             string
	CreditCardType     string `faker:"cc_type"`
	CreditCardNumber   string `faker:"cc_number"`
	Email              string `faker:"email"`
	IPV4               string `faker:"ipv4"`
	IPV6               string `faker:"ipv6"`
	Bool               bool
	SString            []string
	SInt               []int
	SInt8              []int8
	SInt16             []int16
	SInt32             []int32
	SInt64             []int64
	SFloat32           []float32
	SFloat64           []float64
	SBool              []bool
	Struct             AStruct
	Time               time.Time
	Stime              []time.Time
	Currency           string  `faker:"currency"`
	Amount             float64 `faker:"price"`
	AmountWithCurrency string  `faker:"price_currency"`
	ID                 int64   `faker:"id"`
	UUID               string  `faker:"uuid"`
	HyphenatedID       string  `faker:"uuid_string"`

	MapStringString        map[string]string
	MapStringStruct        map[string]AStruct
	MapStringStructPointer map[string]*AStruct
}
type AStruct struct {
	Number        int64
	Height        int64
	AnotherStruct CStruct
}

type BStruct struct {
	Image string
}
type CStruct struct {
	BStruct
	Name string
}

type TaggedStruct struct {
	Latitude           float32 `faker:"lat"`
	Longitude          float32 `faker:"long"`
	CreditCardNumber   string  `faker:"cc_number"`
	CreditCardType     string  `faker:"cc_type"`
	Email              string  `faker:"email"`
	IPV4               string  `faker:"ipv4"`
	IPV6               string  `faker:"ipv6"`
	Password           string  `faker:"password"`
	PhoneNumber        string  `faker:"phone_number"`
	MacAddress         string  `faker:"mac_address"`
	URL                string  `faker:"url"`
	UserName           string  `faker:"username"`
	FirstName          string  `faker:"first_name"`
	FirstNameMale      string  `faker:"male_first_name"`
	FirstNameFemale    string  `faker:"female_first_name"`
	LastName           string  `faker:"last_name"`
	Name               string  `faker:"name"`
	UnixTime           int64   `faker:"unix_time"`
	Date               string  `faker:"date"`
	Time               string  `faker:"time"`
	MonthName          string  `faker:"month_name"`
	Year               string  `faker:"year"`
	Month              string  `faker:"month"`
	DayOfWeek          string  `faker:"week_day"`
	Timestamp          string  `faker:"timestamp"`
	TimeZone           string  `faker:"timezone"`
	Word               string  `faker:"word"`
	Sentence           string  `faker:"sentence"`
	Paragraph          string  `faker:"paragraph"`
	Currency           string  `faker:"currency"`
	Amount             float64 `faker:"price"`
	AmountWithCurrency string  `faker:"price_currency"`
	ID                 []byte  `faker:"uuid"`
	HyphenatedID       string  `faker:"uuid_string"`
}

type NotTaggedStruct struct {
	Latitude         float32
	Long             float32
	CreditCardType   string
	CreditCardNumber string
	Email            string
	IPV4             string
	IPV6             string
}

func TestFakerData(t *testing.T) {
	s := NewService(0, nil)
	var a SomeStruct
	err := s.FakeData(&a)
	assert.NoError(t, err, "\n%+v", err)

	data, err := json.Marshal(&a)
	assert.NoError(t, err, "\n%+v", err)
	assert.LenBetween(t, data, 100, 140000)

	// t.Logf("SomeStruct: %+v\n", a)

	var b TaggedStruct
	err = s.FakeData(&b)
	assert.NoError(t, err, "%+v", err)

	// repr.Println(b)

	// t.Logf("TaggedStruct: %+v\n", b)
	data, err = json.Marshal(&a)
	assert.NoError(t, err, "\n%+v", err)
	assert.LenBetween(t, data, 100, 135000)

	// Example Result :
	// {Int:8906957488773767119 Int8:6 Int16:14 Int32:391219825 Int64:2374447092794071106 String:poraKzAxVbWVkMkpcZCcWlYMd Bool:false SString:[MehdV aVotHsi] SInt:[528955241289647236 7620047312653801973 2774096449863851732] SInt8:[122 -92 -92] SInt16:[15679 -19444 -30246] SInt32:[1146660378 946021799 852909987] SInt64:[6079203475736033758 6913211867841842836 3269201978513619428] SFloat32:[0.019562425 0.12729558 0.36450312] SFloat64:[0.7825838989890364 0.9732903338838912 0.8316541489234004] SBool:[true false true] Struct:{Number:7693944638490551161 Height:6513508020379591917}}

}

func TestUnsuportedMapStringInterface(t *testing.T) {
	s := NewService(0, nil)

	type Sample struct {
		Map map[string]interface{}
	}
	sample := new(Sample)
	if err := s.FakeData(sample); err == nil {
		t.Error("Expected Got Error. But got nil")
	}
}

func TestSetDataIfArgumentNotPtr(t *testing.T) {
	s := NewService(0, nil)
	temp := struct{}{}
	err := s.FakeData(temp)
	assert.True(t, errors.NotSupported.Match(err), "%+v", err)
}

func TestSetDataIfArgumentNotHaveReflect(t *testing.T) {
	temp := func() {}
	s := NewService(0, nil)
	err := s.FakeData(temp)
	assert.True(t, errors.NotSupported.Match(err), "%+v", err)
}

func TestSetDataErrorDataParseTagStringType(t *testing.T) {
	temp := &struct {
		Test string `faker:"test"`
	}{}
	t.Logf("%+v ", temp)
	s := NewService(0, nil)
	if err := s.FakeData(temp); err == nil {
		t.Error("Exptected error Unsupported tag, but got nil")
	}
}

func TestSetDataErrorDataParseTagIntType(t *testing.T) {
	temp := &struct {
		Test int `faker:"test"`
	}{}
	s := NewService(0, nil)
	if err := s.FakeData(temp); err == nil {
		t.Error("Exptected error Unsupported tag, but got nil")
	}
}

func TestSetDataWithTagIfFirstArgumentNotPtr(t *testing.T) {
	temp := struct{}{}
	s := NewService(0, nil)
	err := s.setDataWithTag(reflect.ValueOf(temp), "", 0)
	assert.True(t, errors.NotSupported.Match(err), "%+v", err)
}

func TestSetDataWithTagIfFirstArgumentSlice(t *testing.T) {
	temp := []int{}
	s := NewService(0, nil)
	err := s.setDataWithTag(reflect.ValueOf(&temp), "", 0)
	assert.True(t, errors.NotFound.Match(err), "%+v", err)
}

func TestSetDataWithTagIfFirstArgumentNotFound(t *testing.T) {
	temp := struct{}{}
	s := NewService(0, nil)
	err := s.setDataWithTag(reflect.ValueOf(&temp), "", 0)
	assert.True(t, errors.NotFound.Match(err), "%+v", err)
}

func BenchmarkFakerDataNOTTagged(b *testing.B) {
	s := NewService(0, nil)
	for i := 0; i < b.N; i++ {
		a := NotTaggedStruct{}
		err := s.FakeData(&a)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFakerDataTagged(b *testing.B) {
	s := NewService(0, nil)
	for i := 0; i < b.N; i++ {
		a := TaggedStruct{}
		err := s.FakeData(&a)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type PointerStructA struct {
	SomeStruct *SomeStruct
}
type PointerStructB struct {
	PointA PointerStructA
}

type PointerC struct {
	TaggedStruct *TaggedStruct
}

func TestStructPointer(t *testing.T) {
	s := NewService(0, nil)
	a := new(PointerStructB)
	err := s.FakeData(a)
	if err != nil {
		t.Error("Expected Not Error, But Got: ", err)
	}
	// t.Logf("A value: %+v , Somestruct Value: %+v  ", a, a)

	tagged := new(PointerC)
	err = s.FakeData(tagged)
	if err != nil {
		t.Error("Expected Not Error, But Got: ", err)
	}
	// t.Logf(" tagged value: %+v , TaggedStruct Value: %+v  ", a, a.PointA.SomeStruct)
}

type CustomString string
type CustomInt int
type CustomMap map[string]string
type CustomPointerStruct PointerStructB
type CustomTypeStruct struct {
	CustomString        CustomString
	CustomInt           CustomInt
	CustomMap           CustomMap
	CustomPointerStruct CustomPointerStruct
}

func TestCustomType(t *testing.T) {
	a := new(CustomTypeStruct)
	s := NewService(0, nil)
	err := s.FakeData(a)
	assert.NoError(t, err)
	// t.Logf("A value: %+v , Somestruct Value: %+v  ", a, a)
}

type SampleStruct struct {
	name string
	Age  int
}

func TestUnexportedFieldStruct(t *testing.T) {
	// This test is to ensure that the faker won't panic if trying to fake data on struct that has unexported field
	a := new(SampleStruct)
	s := NewService(0, nil)
	err := s.FakeData(a)

	if err != nil {
		t.Error("Expected Not Error, But Got: ", err)
	}
	t.Logf("A value: %+v , SampleStruct Value: %+v  ", a, a)
}

func TestPointerToCustomScalar(t *testing.T) {
	// This test is to ensure that the faker won't panic if trying to fake data on struct that has field
	a := new(CustomInt)
	s := NewService(0, nil)
	err := s.FakeData(a)

	if err != nil {
		t.Error("Expected Not Error, But Got: ", err)
	}
	t.Logf("A value: %+v , Custom scalar Value: %+v  ", a, a)
}

type PointerCustomIntStruct struct {
	V *CustomInt
}

func TestPointerToCustomIntStruct(t *testing.T) {
	// This test is to ensure that the faker won't panic if trying to fake data on struct that has field
	a := new(PointerCustomIntStruct)
	s := NewService(0, nil)
	err := s.FakeData(a)

	if err != nil {
		t.Error("Expected Not Error, But Got: ", err)
	}
	t.Logf("A value: %+v , PointerCustomIntStruct scalar Value: %+v  ", a, a)
}

func TestSkipField(t *testing.T) {
	// This test is to ensure that the faker won't fill field with tag skip

	a := struct {
		ID              int
		ShouldBeSkipped int `faker:"-"`
	}{}
	s := NewService(0, nil)
	err := s.FakeData(&a)

	if err != nil {
		t.Error("Expected Not Error, But Got: ", err)
	}

	if a.ShouldBeSkipped != 0 {
		t.Error("Expected that field will be skipped")
	}

}

func TestExtend(t *testing.T) {
	// This test is to ensure that faker can be extended new providers

	a := struct {
		ID string `faker:"test"`
	}{}
	s := NewService(0, nil)
	previous := s.AddProvider("test", func(maxLen int) (interface{}, error) {
		return "test", nil
	})
	assert.Nil(t, previous)

	assert.NoError(t, s.FakeData(&a))

	if a.ID != "test" {
		t.Error("ID should be equal test value")
	}
}

func TestTagAlreadyExists(t *testing.T) {
	// This test is to ensure that existing tag cannot be rewritten
	s := NewService(0, nil)
	prev := s.AddProvider("email", func(maxLen int) (interface{}, error) { return "", nil })
	assert.NotNil(t, prev)

}

func TestSetLang(t *testing.T) {
	s := NewService(0, nil)
	err := s.SetLang("ru")
	if err != nil {
		t.Error("SetLang should successfully set lang")
	}

	err = s.SetLang("sd")
	if err == nil {
		t.Error("SetLang with nonexistent lang should return error")
	}
}

func TestFakerRuWithoutCallback(t *testing.T) {
	s := NewService(0, &Options{
		EnFallback: false,
	})
	assert.NoError(t, s.SetLang("ru"))
	brand := s.Brand()
	if brand != "" {
		t.Error("Fake call with no samples should return blank string")
	}
}

func TestFakerRuWithCallback(t *testing.T) {
	s := NewService(0, &Options{
		EnFallback: true,
	})
	assert.NoError(t, s.SetLang("ru"))
	brand := s.Brand()
	if brand == "" {
		t.Error("Fake call for name with no samples with callback should not return blank string")
	}
}

// TestConcurrentSafety runs fake methods in multiple go routines concurrently.
// This test should be run with the race detector enabled.
func TestConcurrentSafety(t *testing.T) {

	s := NewService(0, &Options{
		EnFallback: true,
	})

	var funcs = []func() string{
		s.FirstName,
		s.LastName,
		s.Gender,
		s.FullName,
		s.WeekDayShort,
		s.Country,
		s.Company,
		s.Industry,
		s.Street,
	}

	bgwork.Wait(len(funcs), func(idx int) {
		_ = funcs[idx]()
		for i := 0; i < 20; i++ {
			for _, fn := range funcs {
				_ = fn()
			}
		}
	})
}

type CoreConfigData struct {
	ConfigID      uint32       `json:"config_id,omitempty" max_len:"10"`
	Scope         string       `json:"scope,omitempty" max_len:"8"`
	ScopeID       int32        `json:"scope_id" xml:"scope_id"`
	Expires       null.Time    `json:"expires,omitempty" `
	Path          string       `json:"x_path" xml:"y_path" max_len:"255"`
	Value         null.String  `json:"value,omitempty" max_len:"65535"`
	ColDecimal100 null.Decimal `json:"col_decimal_10_0,omitempty"  max_len:"10"`
	ColBlob       []byte       `json:"col_blob,omitempty"  max_len:"65535"`
}

func TestMaxLen(t *testing.T) {
	t.Run("CoreConfigData", func(t *testing.T) {
		s := NewService(0, nil)
		ccd := new(CoreConfigData)
		assert.NoError(t, s.FakeData(ccd))
		assert.LenBetween(t, ccd.ConfigID, 0, 10)
		assert.LenBetween(t, ccd.Scope, 1, 8)
		assert.LenBetween(t, ccd.ScopeID, 0, math.MaxInt32)
		assert.LenBetween(t, ccd.Path, 1, 255)
		assert.LenBetween(t, ccd.Value.String, 1, 65535)
		assert.LenBetween(t, ccd.ColBlob, 1, 65535)
		// t.Logf("%#v", ccd.ColDecimal100)
	})
}
