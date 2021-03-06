package slog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bouk/monkey"
	"os"
	"strings"
	"testing"
)

func TestSetDefaultOutput(t *testing.T) {
	od := defaultOut

	SetShowLines(true)
	SetDefaultOutput(os.Stdout)
	// Crash Tests
	Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Info("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Log("Test %s %d %f %v", "huebr", 1, 10.0, true)
	LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)

	SetDefaultOutput(os.Stderr)
	// Crash Tests
	Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Info("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Log("Test %s %d %f %v", "huebr", 1, 10.0, true)
	LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)

	buff := bytes.NewBufferString("")
	SetDefaultOutput(buff)
	// Crash Tests
	Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Info("Test %s %d %f %v", "huebr", 1, 10.0, true)
	Log("Test %s %d %f %v", "huebr", 1, 10.0, true)
	LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)

	buff.Reset()

	// Test output
	Info("Test %s %d %f %v", "huebr", 1, 10.0, true)

	if strings.Index(buff.String(), "Test huebr 1 10.000000 true") == -1 {
		t.Errorf("Expected string %s in %s", "Test huebr 1 10.000000 true", buff.String())
	}

	SetDefaultOutput(od)
}

func TestArgsOnly(t *testing.T) {
	SetFieldRepresentation(NoFields)
	buff := bytes.NewBufferString("")
	i := Scope("ArgsOnly").WithCustomWriter(buff)

	i.Info("huebr", 1, 10.0, true)

	o := buff.String()
	if strings.Index(o, "huebr 1 10 true") == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", "huebr 1 10.0 true", o)
	}

	SetFieldRepresentation(JSONFields)
	buff = bytes.NewBufferString("")
	i = Scope("ArgsOnly").WithFields(map[string]interface{}{
		"a": "b",
	}).WithCustomWriter(buff)

	i.Info(555, 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, "555 1 10 true") == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", "huebr 1 10.0 true", o)
	}

	SetFieldRepresentation(KeyValueFields)
	buff = bytes.NewBufferString("")
	i = Scope("ArgsOnly").WithFields(map[string]interface{}{
		"a": "b",
	}).WithCustomWriter(buff)

	i.Info("huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, "huebr 1 10 true") == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", "huebr 1 10.0 true", o)
	}
}

func TestDefaultOutput(t *testing.T) {
	i := Scope("DefaultOutput").WithCustomWriter(nil)

	// Crash Tests
	i.Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Log("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)
}

func TestSetTestMode(t *testing.T) {
	SetTestMode()
	// Test Mode should be all unset

	if DebugEnabled() {
		t.Fatalf("Debug is set to true! Should be false")
	}

	if WarningEnabled() {
		t.Fatalf("Warn is set to true! Should be false")
	}

	if WarningEnabled() {
		t.Fatalf("Error is set to true! Should be false")
	}

	if InfoEnabled() {
		t.Fatalf("Info is set to true! Should be false")
	}

	UnsetTestMode()
	// Test Mode should be all set

	if !DebugEnabled() {
		t.Fatalf("Debug is set to false! Should be true")
	}

	if !WarningEnabled() {
		t.Fatalf("Warn is set to false! Should be true")
	}

	if !ErrorEnabled() {
		t.Fatalf("Error is set to false! Should be true")
	}

	if !InfoEnabled() {
		t.Fatalf("Info is set to false! Should be true")
	}
}

func TestSubScope(t *testing.T) {
	i := Scope("ABCD").SubScope("EFGH").(*slogInstance)

	if stringSliceIndexOf("ABCD", i.scope) == -1 {
		t.Errorf("Expected ABCD in Scope")
	}
	if stringSliceIndexOf("EFGH", i.scope) == -1 {
		t.Errorf("Expected EFGH in Scope")
	}

	// No Crash tests
	i.Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Log("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)
}

func TestWithFields(t *testing.T) {
	SetFieldRepresentation(KeyValueFields)
	i := Scope("WithFields").WithFields(map[string]interface{}{
		"a": "b",
		"b": 5,
	}).(*slogInstance)

	if i.fields["a"] != "b" {
		t.Errorf("Expected field \"a\" to be \"b\"")
	}

	if i.fields["b"] != 5 {
		t.Errorf("Expected field \"b\" to be 5")
	}

	// Child should inherit parent fields and replace the existent ones
	i = i.WithFields(map[string]interface{}{
		"c": 3.14,
		"a": 9,
	}).(*slogInstance)

	if i.fields["a"] != 9 {
		t.Errorf("Expected field \"a\" to be 9")
	}

	if i.fields["b"] != 5 {
		t.Errorf("Expected field \"b\" to be 5")
	}

	if i.fields["c"] != 3.14 {
		t.Errorf("Expected field \"b\" to be 5")
	}

	// No Crash tests
	i.Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.Log("Test %s %d %f %v", "huebr", 1, 10.0, true)
	i.LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)
}

func TestWithFieldsJSON(t *testing.T) {
	SetFieldRepresentation(JSONFields)
	buff := bytes.NewBufferString("")
	i := Scope("WithFieldsJSON").WithFields(map[string]interface{}{
		"a": "b",
		"b": 5,
	}).WithCustomWriter(buff).(*slogInstance)

	jsonDataB, _ := json.Marshal(i.fields)
	jsonData := string(jsonDataB)

	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o := buff.String()
	if strings.Index(o, jsonData) == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()

	i.Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, jsonData) == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()

	i.Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, jsonData) == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()

	i.Error("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, jsonData) == -1 {
		t.Errorf("Expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()
}

func TestWithFieldsKV(t *testing.T) {
	SetFieldRepresentation(KeyValueFields)
	buff := bytes.NewBufferString("")
	i := Scope("WithFieldsKV").WithFields(map[string]interface{}{
		"a": "b",
		"b": 5,
	}).WithCustomWriter(buff).(*slogInstance)

	kvData := ""

	testFields := func(out string) { // This is nescessary since the range orders can randomly change
		for k, v := range i.fields {
			kvz := fmt.Sprintf("%s=%v,", k, v)
			if strings.Index(out, kvz) == -1 {
				t.Errorf("Expected \"%s\" in output: \"%s\"", kvData, out)
			}
		}
	}

	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o := buff.String()
	testFields(o)

	buff.Reset()

	i.Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	testFields(o)

	buff.Reset()

	i.Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	testFields(o)

	buff.Reset()

	i.Error("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	testFields(o)

	buff.Reset()
}

func TestWithFieldsNoFields(t *testing.T) {
	SetFieldRepresentation(NoFields)
	buff := bytes.NewBufferString("")
	i := Scope("WithFieldsKV").WithFields(map[string]interface{}{
		"a": "b",
		"b": 5,
	}).WithCustomWriter(buff).(*slogInstance)

	kvData := ""

	for k, v := range i.fields {
		kvData += fmt.Sprintf("%s=%v,", k, v)
	}

	jsonDataB, _ := json.Marshal(i.fields)
	jsonData := string(jsonDataB)

	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o := buff.String()
	if strings.Index(o, kvData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", kvData, o)
	}

	if strings.Index(o, jsonData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()

	i.Warn("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, kvData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", kvData, o)
	}

	if strings.Index(o, jsonData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()

	i.Debug("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, kvData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", kvData, o)
	}

	if strings.Index(o, jsonData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", jsonData, o)
	}
	buff.Reset()

	i.Error("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o = buff.String()
	if strings.Index(o, kvData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", kvData, o)
	}

	if strings.Index(o, jsonData) != -1 {
		t.Errorf("Not expected \"%s\" in output: \"%s\"", jsonData, o)
	}

	buff.Reset()
}

func TestDebug(t *testing.T) {
	Debug("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestError(t *testing.T) {
	Error("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestWarn(t *testing.T) {
	Warn("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestInfo(t *testing.T) {
	Info("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestLog(t *testing.T) {
	Log("Test %s %d %f %v", "huebr", 1, 10.0, true) // Shouldn't crash
}

func TestLogNoFormat(t *testing.T) {
	LogNoFormat("Test %s %d %f %v", "huebr", 1, 10.0, true)
}

func TestFatal(t *testing.T) {
	fakeExit := func(int) {
		panic("os.Exit called")
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	assertPanic(t, func() {
		Fatal("Test Fatal")
	}, "Fatal should os.Exit")
	assertPanic(t, func() {
		Fatal("Test %s %d %f %v", "huebr", 1, 10.0, true)
	}, "Fatal should os.Exit")
}

func TestScope(t *testing.T) {
	scoped := Scope("test-scope").(*slogInstance)
	if scoped.scope[0] != "test-scope" {
		t.Fatalf("Expected test-scope got %s", scoped.scope)
	}
}

func TestSetDebug(t *testing.T) {
	SetDebug(true)
	if !DebugEnabled() {
		t.Fatalf("Debug is set to false! Should be true")
	}
	SetDebug(false)
	if DebugEnabled() {
		t.Fatalf("Debug is set to true! Should be false")
	}
}

func TestSetError(t *testing.T) {
	SetError(true)
	if !ErrorEnabled() {
		t.Fatalf("Error is set to false! Should be true")
	}
	SetError(false)
	if ErrorEnabled() {
		t.Fatalf("Error is set to true! Should be false")
	}
}

func TestSetInfo(t *testing.T) {
	SetInfo(true)
	if !InfoEnabled() {
		t.Fatalf("Info is set to false! Should be true")
	}
	SetInfo(false)
	if InfoEnabled() {
		t.Fatalf("Info is set to true! Should be false")
	}
}

func TestSetWarn(t *testing.T) {
	SetWarning(true)
	if !WarningEnabled() {
		t.Fatalf("Warn is set to false! Should be true")
	}
	SetWarning(false)
	if WarningEnabled() {
		t.Fatalf("Warn is set to true! Should be false")
	}
}

func TestSetShowLines(t *testing.T) {
	SetShowLines(true)
	if !ShowLinesEnabled() {
		t.Fatal("ShowLines is set to false! Should be true")
	}
	SetShowLines(false)
	if ShowLinesEnabled() {
		t.Fatal("ShowLines is set to true! Should be false")
	}
}

type test struct{}

func (test) String() string {
	return "test"
}

func TestAsString(t *testing.T) {
	var tstringcast StringCast

	tstringcast = &test{}

	tests := []interface{}{
		"string",
		123456,
		123456.1,
		true,
		map[string]string{},
		[]int{1, 2, 3, 4, 5},
		complex(float32(1), float32(1)),
		complex(float64(1), float64(1)),
		fmt.Errorf("error format"),
		tstringcast,
	}

	outputs := make([]string, len(tests))

	for i, v := range tests { // Fill tests
		outputs[i] = fmt.Sprint(v) // Should be same output
	}

	for i, v := range tests {
		s := asString(v)
		if s != outputs[i] {
			t.Errorf("#%d expected %s got %s.", i, outputs[i], s)
		}
	}
}

func assertPanic(t *testing.T, f func(), message string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(message)
		}
	}()
	f()
}

func TestExtensions(t *testing.T) {
	UnsetTestMode()
	i := Scope("Extensions")

	i.Note("Test %s %d %f %v", "huebr", 1, 10.0, true)    // No Crash
	i.Await("Test %s %d %f %v", "huebr", 1, 10.0, true)   // No Crash
	i.IO("Test %s %d %f %v", "huebr", 1, 10.0, true)      // No Crash
	i.Done("Test %s %d %f %v", "huebr", 1, 10.0, true)    // No Crash
	i.Success("Test %s %d %f %v", "huebr", 1, 10.0, true) // No Crash

	i.WarnNote("Test %s %d %f %v", "huebr", 1, 10.0, true)    // No Crash
	i.WarnAwait("Test %s %d %f %v", "huebr", 1, 10.0, true)   // No Crash
	i.WarnIO("Test %s %d %f %v", "huebr", 1, 10.0, true)      // No Crash
	i.WarnDone("Test %s %d %f %v", "huebr", 1, 10.0, true)    // No Crash
	i.WarnSuccess("Test %s %d %f %v", "huebr", 1, 10.0, true) // No Crash

	i.ErrorNote("Test %s %d %f %v", "huebr", 1, 10.0, true)    // No Crash
	i.ErrorAwait("Test %s %d %f %v", "huebr", 1, 10.0, true)   // No Crash
	i.ErrorIO("Test %s %d %f %v", "huebr", 1, 10.0, true)      // No Crash
	i.ErrorDone("Test %s %d %f %v", "huebr", 1, 10.0, true)    // No Crash
	i.ErrorSuccess("Test %s %d %f %v", "huebr", 1, 10.0, true) // No Crash

	i.DebugNote("Test %s %d %f %v", "huebr", 1, 10.0, true)    // No Crash
	i.DebugAwait("Test %s %d %f %v", "huebr", 1, 10.0, true)   // No Crash
	i.DebugIO("Test %s %d %f %v", "huebr", 1, 10.0, true)      // No Crash
	i.DebugDone("Test %s %d %f %v", "huebr", 1, 10.0, true)    // No Crash
	i.DebugSuccess("Test %s %d %f %v", "huebr", 1, 10.0, true) // No Crash
}

func TestTag(t *testing.T) {
	UnsetTestMode()
	buff := bytes.NewBufferString("")
	i := Scope("Tag").WithCustomWriter(buff).Tag("MyHUETAG").(*slogInstance)
	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o := buff.String()
	if strings.Index(o, "MyHUETAG") == -1 {
		t.Errorf("Expected tag \"%s\" in output: \"%s\"", "MyHUETAG", o)
	}
}

func TestOperation(t *testing.T) {
	UnsetTestMode()
	buff := bytes.NewBufferString("")
	i := Scope("Operation").WithCustomWriter(buff).Operation(AWAIT).(*slogInstance)
	i.Info("Test %s %d %f %v", "huebr", 1, 10.0, true)

	o := buff.String()
	if strings.Index(o, string(AWAIT)) == -1 {
		t.Errorf("Expected operation \"%s\" in output: \"%s\"", string(AWAIT), o)
	}
}

func TestMultiLine(t *testing.T) {
	UnsetTestMode()
	i := Scope("MultiLine").Operation(AWAIT).(*slogInstance)
	i.Info("Test %s %d %f %v\nHUEBR\nLINE", "huebr", 1, 10.0, true)

	// TODO: Test line padding
}

func TestInstanceScope(t *testing.T) {
	i := Scope("MultiLine").Operation(AWAIT).(*slogInstance)
	i2 := i.Scope("ABCD").(*slogInstance)

	if i2.scope[0] != "ABCD" {
		t.Errorf("Expected scope[0] to be ABCD got %s", i2.scope[0])
	}

	if len(i2.scope) != 1 {
		t.Errorf("Expected scope length to be 1 got %d", len(i2.scope))
	}
}

func TestScopeLength(t *testing.T) {
	n := 15
	SetScopeLength(n)
	if scopeLength != 15 {
		t.Errorf("Expected scope length to be 15")
	}

	// TODO: Check scope padding
}

func TestToFormat(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedFormat Format
	}{
		{
			name:           "input is JSON",
			input:          "jSoN",
			expectedFormat: JSON,
		},
		{
			name:           "input is pipe",
			input:          "PIPE",
			expectedFormat: PIPE,
		},
		{
			name:           "input is empty",
			input:          "",
			expectedFormat: PIPE,
		},
		{
			name:           "input is invalid string",
			input:          "abcde",
			expectedFormat: PIPE,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToFormat(tc.input)
			if result != tc.expectedFormat {
				t.Errorf("Got %q want %q.", result, tc.expectedFormat)
			}
		})
	}
}

func TestJsonFormat(t *testing.T) {
	defer func() {
		SetLogFormat(PIPE)
		SetShowLines(false)
	}()

	SetLogFormat(JSON)
	SetShowLines(true)

	buff := bytes.NewBufferString("")
	i := Scope("UnitTest").WithCustomWriter(buff)
	i = i.SubScope("Test")
	i = i.Operation(NOTE)
	i = i.Tag("CLASSIFICATION")
	i = i.WithFields(map[string]interface{}{"customField": "123"})
	i.Info("Test message as JSON %s", "test")

	var values map[string]interface{}
	_ = json.Unmarshal(buff.Bytes(), &values)

	if values["time"] == "" {
		t.Errorf("Got empty want not empty.")
	}

	if values["scope"] != "UnitTest - Test" {
		t.Errorf("Got %q want %q.", values["scope"], "UnitTest")
	}

	if values["op"] != "NOTE" {
		t.Errorf("Got %q want %q.", values["op"], "NOTE")
	}

	if values["tag"] != "CLASSIFICATION" {
		t.Errorf("Got %q want %q.", values["tag"], "CLASSIFICATION")
	}

	if values["level"] != "info" {
		t.Errorf("Got %q want %q.", values["level"], "info")
	}

	if values["msg"] != "Test message as JSON test" {
		t.Errorf("Got %q want %q.", values["msg"], "")
	}

	if values["customField"] != "123" {
		t.Errorf("Got %q want %q.", values["customField"], "123")
	}

	if values["lines"] == "" {
		t.Errorf("Got empty want not empty.")
	}
}

func TestLogWithoutCustomValue(t *testing.T) {
	defer func() {
		SetLogFormat(PIPE)
	}()

	SetLogFormat(JSON)

	buff := bytes.NewBufferString("")
	i := Scope("UnitTest").WithCustomWriter(buff)
	i.Info("Test message as JSON %s", "test")

	var values map[string]interface{}
	_ = json.Unmarshal(buff.Bytes(), &values)

	if _, ok := values["customField"]; ok {
		t.Errorf("Got %q want ''.", values["customField"])
	}
}

func TestLogWithInvalidLogFormat(t *testing.T) {
	defer func() {
		SetLogFormat(PIPE)
	}()

	var ABC Format = "abc"
	SetLogFormat(ABC)

	buff := bytes.NewBufferString("")
	i := Scope("UnitTestInvalidLogFormat").WithCustomWriter(buff)
	i.Info("Test message as JSON %s", "test")

	if !strings.Contains(buff.String(), "Untreated log format abc") {
		t.Errorf("Got %v, want to contain 'Untreated log format abc'.", buff.String())
	}
}
