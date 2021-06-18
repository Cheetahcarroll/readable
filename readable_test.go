package readable

import (
	"testing"
)

type ReadableTest struct {
	Field   string     `json:"field" readable:"Generic Field" compare:"Field value changed from %v to %v"`
	Another []NestTest `json:"another" readable:"Nested Object"`
}

type NestTest struct {
	IntField     int          `json:"nestField" readable:"Nested Field" compare:"IntField value changed from %v to %v"`
	AnotherField string       `json:"nesterField" readable:"Another Nested Field" compare:"AnotherField value changed from %v to %v"`
	MoreFields   *SubNestTest `json:"current" readable:"Current Fields"`
}

type SubNestTest struct {
	Field        string    `json:"subNestField" readable:"Sub Nested Field" compare:"SubField value changed from %v to %v"`
	AnotherField string    `json:"subNesterField" readable:"Another Sub Nested Field" compare:"SubAnotherField value changed from %v to %v"`
	MoreFields   *NestTest `json:"subCurrent" readable:"Sub Current Fields"`
}

type test struct {
	testData    ReadableTest
	newTestData ReadableTest
	testValue   string
	flatValue   string
}

var stringTestCases = map[string]test{
	"Test 1": {
		testData: ReadableTest{
			Field: "Test 1",
			Another: []NestTest{
				{
					IntField:     1,
					AnotherField: "Another Tested 1",
					MoreFields: &SubNestTest{
						Field:        "Deep Nest",
						AnotherField: "Such Nest Much Wow",
					},
				},
				{
					IntField:     2,
					AnotherField: "Another Tested 1-2",
					MoreFields: &SubNestTest{
						Field:        "Deeper Nest",
						AnotherField: "Such Nesting Much Wow",
					},
				},
			},
		},
		testValue: "Generic Field: \"Test 1\"\n" +
			"Nested Object: \n" +
			"	Nested Field: 1\n" +
			"	Another Nested Field: \"Another Tested 1\"\n" +
			"	Current Fields: \n" +
			"		Sub Nested Field: \"Deep Nest\"\n" +
			"		Another Sub Nested Field: \"Such Nest Much Wow\"\n" +
			"		Sub Current Fields: \n" +
			"	Nested Field: 2\n" +
			"	Another Nested Field: \"Another Tested 1-2\"\n" +
			"	Current Fields: \n" +
			"		Sub Nested Field: \"Deeper Nest\"\n" +
			"		Another Sub Nested Field: \"Such Nesting Much Wow\"\n" +
			"		Sub Current Fields: \n",
	}, "Test 2": {
		testData: ReadableTest{
			Field: "",
			Another: []NestTest{
				{
					IntField:     2,
					AnotherField: "Another Tested 2",
				},
			},
		},
		testValue: "Generic Field: \"\"\n" +
			"Nested Object: \n" +
			"	Nested Field: 2\n" +
			"	Another Nested Field: \"Another Tested 2\"\n" +
			"	Current Fields: \n",
	},
	"Test 3": {
		testData: ReadableTest{
			Field: "Test 3",
			Another: []NestTest{
				{
					IntField:     3,
					AnotherField: "Another Tested 3",
					MoreFields: &SubNestTest{
						Field:        "Deep Nest",
						AnotherField: "Such Nest Much Wow",
						MoreFields: &NestTest{
							IntField:     4,
							AnotherField: "Another Tested 4",
						},
					},
				},
			},
		},
		testValue: "Generic Field: \"Test 3\"\n" +
			"Nested Object: \n" +
			"	Nested Field: 3\n" +
			"	Another Nested Field: \"Another Tested 3\"\n" +
			"	Current Fields: \n" +
			"		Sub Nested Field: \"Deep Nest\"\n" +
			"		Another Sub Nested Field: \"Such Nest Much Wow\"\n" +
			"		Sub Current Fields: \n" +
			"			Nested Field: 4\n" +
			"			Another Nested Field: \"Another Tested 4\"\n" +
			"			Current Fields: \n",
	},
}

var jsonTestCases = map[string]test{
	"Test 1": {
		testData: ReadableTest{
			Field: "Test 1",
			Another: []NestTest{
				{
					IntField:     1,
					AnotherField: "Another Tested 1",
					MoreFields: &SubNestTest{
						Field:        "Deep Nest",
						AnotherField: "Such Nest Much Wow",
					},
				},
			},
		},
		testValue: "{\n" +
			"	\"field\": \"Test 1\",\n" +
			"	\"another\": {\n" +
			"		\"nestField\": 1,\n" +
			"		\"nesterField\": \"Another Tested 1\",\n" +
			"		\"current\": {\n" +
			"			\"subNestField\": \"Deep Nest\",\n" +
			"			\"subNesterField\": \"Such Nest Much Wow\",\n" +
			"			\"subCurrent\": {}\n" +
			"		}\n" +
			"	}\n" +
			"}",
		flatValue: "{\"field\":\"Test 1\",\"another\":{\"nestField\":1,\"nesterField\":\"Another Tested 1\",\"current\":{\"subNestField\":\"Deep Nest\",\"subNesterField\":\"Such Nest Much Wow\",\"subCurrent\":{}}}}",
	}, "Test 2": {
		testData: ReadableTest{
			Field: "",
			Another: []NestTest{
				{
					IntField:     2,
					AnotherField: "Another Tested 2",
				},
			},
		},
		testValue: "{\n" +
			"	\"field\": \"\",\n" +
			"	\"another\": {\n" +
			"		\"nestField\": 2,\n" +
			"		\"nesterField\": \"Another Tested 2\",\n" +
			"		\"current\": {}\n" +
			"	}\n" +
			"}",
		flatValue: "{\"field\":\"\",\"another\":{\"nestField\":2,\"nesterField\":\"Another Tested 2\",\"current\":{}}}",
	},
	"Test 3": {
		testData: ReadableTest{
			Field: "Test 3",
			Another: []NestTest{
				{
					IntField:     3,
					AnotherField: "Another Tested 3",
					MoreFields: &SubNestTest{
						Field:        "Deep Nest",
						AnotherField: "Such Nest Much Wow",
						MoreFields: &NestTest{
							IntField:     4,
							AnotherField: "Another Tested 4",
						},
					},
				},
			},
		},
		testValue: "{\n" +
			"	\"field\": \"Test 3\",\n" +
			"	\"another\": {\n" +
			"		\"nestField\": 3,\n" +
			"		\"nesterField\": \"Another Tested 3\",\n" +
			"		\"current\": {\n" +
			"			\"subNestField\": \"Deep Nest\",\n" +
			"			\"subNesterField\": \"Such Nest Much Wow\",\n" +
			"			\"subCurrent\": {\n" +
			"				\"nestField\": 4,\n" +
			"				\"nesterField\": \"Another Tested 4\",\n" +
			"				\"current\": {}\n" +
			"			}\n" +
			"		}\n" +
			"	}\n" +
			"}",
		flatValue: "{\"field\":\"Test 3\",\"another\":{\"nestField\":3,\"nesterField\":\"Another Tested 3\",\"current\":{\"subNestField\":\"Deep Nest\",\"subNesterField\":\"Such Nest Much Wow\",\"subCurrent\":{\"nestField\":4,\"nesterField\":\"Another Tested 4\",\"current\":{}}}}}",
	},
}

var unitTestCases = map[string][]string{
	"Test 1": {
		"{\n" +
			"	\"field\": \"Test 1\",\n" +
			"	\"another\": [{\n" +
			"		\"nestField\": 1,\n" +
			"		\"nesterField\": \"Another Tested 1\",\n" +
			"		\"current\": {\n" +
			"			\"subNestField\": \"Deep Nest\",\n" +
			"			\"subNesterField\": \"Such Nest Much Wow\"\n" +
			"		}\n" +
			"	}]\n" +
			"}",
		"ReadableTest{\n" +
			"	Field: \"Test 1\",\n" +
			"	Another: []NestTest{\n" +
			"		{\n" +
			"			IntField: 1,\n" +
			"			AnotherField: \"Another Tested 1\",\n" +
			"			MoreFields: &SubNestTest{\n" +
			"				Field: \"Deep Nest\",\n" +
			"				AnotherField: \"Such Nest Much Wow\",\n" +
			"			},\n" +
			"		},\n" +
			"	},\n" +
			"}",
	},
}

var compareTestCases = map[string]test{
	"Test 1": {
		testData: ReadableTest{
			Field: "Test 1",
			Another: []NestTest{
				{
					IntField:     1,
					AnotherField: "Another Tested 1",
					MoreFields: &SubNestTest{
						Field:        "Deep Nest",
						AnotherField: "Such Nest Much Wow",
					},
				},
				{
					IntField:     2,
					AnotherField: "Another Tested 1-2",
					MoreFields: &SubNestTest{
						Field:        "Deeper Nest",
						AnotherField: "Such Nesting Much Wow",
					},
				},
			},
		},
		newTestData: ReadableTest{
			Field: "Test 1-2",
			Another: []NestTest{
				{
					IntField:     10,
					AnotherField: "Another Tested 10",
					MoreFields: &SubNestTest{
						Field:        "Deep Nest 2",
						AnotherField: "Such Nest Much Wow 2",
					},
				},
				{
					IntField:     20,
					AnotherField: "Another Tested 1-20",
					MoreFields: &SubNestTest{
						Field:        "Deeper Nest 20",
						AnotherField: "Such Nesting Much Wow 20",
					},
				},
			},
		},
		testValue: "Field value changed from \"Test 1\" to \"Test 1-2\"\n" +
			"IntField value changed from 1 to 10\n" +
			"AnotherField value changed from \"Another Tested 1\" to \"Another Tested 10\"\n" +
			"SubField value changed from \"Deep Nest\" to \"Deep Nest 2\"\n" +
			"SubAnotherField value changed from \"Such Nest Much Wow\" to \"Such Nest Much Wow 2\"\n" +
			"IntField value changed from 2 to 20\n" +
			"AnotherField value changed from \"Another Tested 1-2\" to \"Another Tested 1-20\"\n" +
			"SubField value changed from \"Deeper Nest\" to \"Deeper Nest 20\"\n" +
			"SubAnotherField value changed from \"Such Nesting Much Wow\" to \"Such Nesting Much Wow 20\"\n",
	},
	"Test 2": {
		testData: ReadableTest{
			Field: "",
			Another: []NestTest{
				{
					IntField:     2,
					AnotherField: "Another Tested 2",
				},
			},
		},
		newTestData: ReadableTest{
			Field: "",
			Another: []NestTest{
				{
					IntField:     22,
					AnotherField: "",
				},
			},
		},
		testValue: "IntField value changed from 2 to 22\n" +
			"AnotherField value changed from \"Another Tested 2\" to \"\"\n",
	},
	"Test 3": {
		testData: ReadableTest{
			Field: "Test 3",
			Another: []NestTest{
				{
					IntField:     3,
					AnotherField: "Another Tested 3",
					MoreFields: &SubNestTest{
						Field:        "Deep Nest",
						AnotherField: "Such Nest Much Wow",
						MoreFields: &NestTest{
							IntField:     4,
							AnotherField: "Another Tested 4",
						},
					},
				},
			},
		},
		newTestData: ReadableTest{
			Field: "Test 3",
			Another: []NestTest{
				{
					IntField:     3,
					AnotherField: "Another Tested 3",
					MoreFields: &SubNestTest{
						Field:        "Deep Nest",
						AnotherField: "Such Nest Much Wow",
						MoreFields: &NestTest{
							IntField:     4,
							AnotherField: "Another Tested 4",
						},
					},
				},
			},
		},
		testValue: "",
	},
}

func Test_GetString(t *testing.T) {
	for run, testCase := range stringTestCases {
		response := GetString(testCase.testData, "readable")
		if testCase.testValue != response {
			t.Errorf("\nFailed to pass %s", run+"\n---------------------\n"+response+"\nFailed to match\n---------------\n"+testCase.testValue)
		}
	}
}

func Test_GetUnitTest(t *testing.T) {
	for run, testCase := range unitTestCases {
		response := GetUnitTest(&ReadableTest{}, testCase[0])
		if testCase[1] != response {
			t.Errorf("\nFailed to pass %s", run+"\n---------------------\n"+response+"\nFailed to match\n---------------\n"+testCase[1])
		}
	}
}

func Test_ToJSONString(t *testing.T) {
	for run, testCase := range jsonTestCases {
		response := ToJSONString(testCase.testData, 0)
		if testCase.testValue != response {
			t.Errorf("\nFailed to pass %s", run+"\n---------------------\n"+response+"\nFailed to match\n---------------\n"+testCase.testValue)
		}
	}
}

func Test_ToFlatJSONString(t *testing.T) {
	for run, testCase := range jsonTestCases {
		response := ToFlatJSONString(testCase.testData)
		if testCase.flatValue != response {
			t.Errorf("\nFailed to pass %s", run+"\n---------------------\n"+response+"\nFailed to match\n---------------\n"+testCase.testValue)
		}
	}
}

func Test_Compare(t *testing.T) {
	for run, testCase := range compareTestCases {
		response := Compare(testCase.testData, testCase.newTestData, "compare")
		if testCase.testValue != response {
			t.Errorf("\nFailed to pass %s", run+"\n---------------------\n"+response+"\nFailed to match\n---------------\n"+testCase.testValue)
		}
	}
}
