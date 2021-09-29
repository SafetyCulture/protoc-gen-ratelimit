package genratelimit

import (
	"encoding/json"
	"testing"

	"gotest.tools/assert"
)

func TestSortLimits(t *testing.T) {
	limits := map[string]*Limit{
		"|a|a|a|a": &Limit{
			Key: "|a|a|a|a",
		},
		"c|c|a|a|a": &Limit{
			Key: "c|c|a|a|a",
		},
		"b|c|a|a|a": &Limit{
			Key: "b|c|a|a|a",
		},
		"a|c|a|a|a": &Limit{
			Key: "a|c|a|a|a",
		},
		"c|b|a|a|a": &Limit{
			Key: "c|b|a|a|a",
		},
		"b|b|a|a|a": &Limit{
			Key: "b|b|a|a|a",
		},
		"a|b|a|a|a": &Limit{
			Key: "a|b|a|a|a",
		},
		"c|a|a|a|a": &Limit{
			Key: "c|a|a|a|a",
		},
		"b|a|a|a|a": &Limit{
			Key: "b|a|a|a|a",
		},
		"a|a|a|a|a": &Limit{
			Key: "a|a|a|a|a",
		},
		"a|a|a|a|c": &Limit{
			Key: "a|a|a|a|c",
		},
		"a|a|a|a|d": &Limit{
			Key: "a|a|a|a|d",
		},
	}

	l := sortLimits(limits)

	keys := []string{}
	for _, v := range l {
		keys = append(keys, v.Key)
	}

	expected := []string{
		"a|a|a|a|a",
		"a|a|a|a|c",
		"a|a|a|a|d",
		"a|b|a|a|a",
		"a|c|a|a|a",
		"b|a|a|a|a",
		"b|b|a|a|a",
		"b|c|a|a|a",
		"c|a|a|a|a",
		"c|b|a|a|a",
		"c|c|a|a|a",
		"|a|a|a|a",
	}

	assert.DeepEqual(t, expected, keys)
}

func TestLimitsDescriptors(t *testing.T) {
	l := limits{
		&Limit{
			Key: "|a|a|a|a",
		},
		&Limit{
			Key: "c|c|a|a|a",
		},
		&Limit{
			Key: "b|c|a|a|a",
		},
		&Limit{
			Key: "a|c|a|a|a",
		},
		&Limit{
			Key: "c|b|a|a|a",
		},
		&Limit{
			Key: "b|b|a|a|a",
		},
		&Limit{
			Key: "a|b|a|a|a",
		},
		&Limit{
			Key: "c|a|a|a|a",
		},
		&Limit{
			Key: "b|a|a|a|a",
		},
		&Limit{
			Key: "a|a|a|a|a",
		},
		&Limit{
			Key: "a|a|a|a|c",
		},
		&Limit{
			Key: "a|a|a|a|d",
		},
	}.Descriptors([]string{"1", "2", "3", "4", "5"})

	assert.Equal(t, len(l), 4)

	b, err := json.MarshalIndent(l, "", "  ")
	assert.NilError(t, err)

	expected := `[
  {
    "Key": "1",
    "Value": "",
    "RateLimit": null,
    "Descriptors": [
      {
        "Key": "2",
        "Value": "a",
        "RateLimit": null,
        "Descriptors": [
          {
            "Key": "3",
            "Value": "a",
            "RateLimit": null,
            "Descriptors": [
              {
                "Key": "4",
                "Value": "a",
                "RateLimit": null,
                "Descriptors": [
                  {
                    "Key": "5",
                    "Value": "a",
                    "RateLimit": null,
                    "Descriptors": null
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  {
    "Key": "1",
    "Value": "c",
    "RateLimit": null,
    "Descriptors": [
      {
        "Key": "2",
        "Value": "c",
        "RateLimit": null,
        "Descriptors": [
          {
            "Key": "3",
            "Value": "a",
            "RateLimit": null,
            "Descriptors": [
              {
                "Key": "4",
                "Value": "a",
                "RateLimit": null,
                "Descriptors": [
                  {
                    "Key": "5",
                    "Value": "a",
                    "RateLimit": null,
                    "Descriptors": null
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        "Key": "2",
        "Value": "b",
        "RateLimit": null,
        "Descriptors": [
          {
            "Key": "3",
            "Value": "a",
            "RateLimit": null,
            "Descriptors": [
              {
                "Key": "4",
                "Value": "a",
                "RateLimit": null,
                "Descriptors": [
                  {
                    "Key": "5",
                    "Value": "a",
                    "RateLimit": null,
                    "Descriptors": null
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        "Key": "2",
        "Value": "a",
        "RateLimit": null,
        "Descriptors": [
          {
            "Key": "3",
            "Value": "a",
            "RateLimit": null,
            "Descriptors": [
              {
                "Key": "4",
                "Value": "a",
                "RateLimit": null,
                "Descriptors": [
                  {
                    "Key": "5",
                    "Value": "a",
                    "RateLimit": null,
                    "Descriptors": null
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  {
    "Key": "1",
    "Value": "b",
    "RateLimit": null,
    "Descriptors": [
      {
        "Key": "2",
        "Value": "c",
        "RateLimit": null,
        "Descriptors": [
          {
            "Key": "3",
            "Value": "a",
            "RateLimit": null,
            "Descriptors": [
              {
                "Key": "4",
                "Value": "a",
                "RateLimit": null,
                "Descriptors": [
                  {
                    "Key": "5",
                    "Value": "a",
                    "RateLimit": null,
                    "Descriptors": null
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        "Key": "2",
        "Value": "b",
        "RateLimit": null,
        "Descriptors": [
          {
            "Key": "3",
            "Value": "a",
            "RateLimit": null,
            "Descriptors": [
              {
                "Key": "4",
                "Value": "a",
                "RateLimit": null,
                "Descriptors": [
                  {
                    "Key": "5",
                    "Value": "a",
                    "RateLimit": null,
                    "Descriptors": null
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        "Key": "2",
        "Value": "a",
        "RateLimit": null,
        "Descriptors": [
          {
            "Key": "3",
            "Value": "a",
            "RateLimit": null,
            "Descriptors": [
              {
                "Key": "4",
                "Value": "a",
                "RateLimit": null,
                "Descriptors": [
                  {
                    "Key": "5",
                    "Value": "a",
                    "RateLimit": null,
                    "Descriptors": null
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  {
    "Key": "1",
    "Value": "a",
    "RateLimit": null,
    "Descriptors": [
      {
        "Key": "2",
        "Value": "c",
        "RateLimit": null,
        "Descriptors": [
          {
            "Key": "3",
            "Value": "a",
            "RateLimit": null,
            "Descriptors": [
              {
                "Key": "4",
                "Value": "a",
                "RateLimit": null,
                "Descriptors": [
                  {
                    "Key": "5",
                    "Value": "a",
                    "RateLimit": null,
                    "Descriptors": null
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        "Key": "2",
        "Value": "b",
        "RateLimit": null,
        "Descriptors": [
          {
            "Key": "3",
            "Value": "a",
            "RateLimit": null,
            "Descriptors": [
              {
                "Key": "4",
                "Value": "a",
                "RateLimit": null,
                "Descriptors": [
                  {
                    "Key": "5",
                    "Value": "a",
                    "RateLimit": null,
                    "Descriptors": null
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        "Key": "2",
        "Value": "a",
        "RateLimit": null,
        "Descriptors": [
          {
            "Key": "3",
            "Value": "a",
            "RateLimit": null,
            "Descriptors": [
              {
                "Key": "4",
                "Value": "a",
                "RateLimit": null,
                "Descriptors": [
                  {
                    "Key": "5",
                    "Value": "a",
                    "RateLimit": null,
                    "Descriptors": null
                  },
                  {
                    "Key": "5",
                    "Value": "c",
                    "RateLimit": null,
                    "Descriptors": null
                  },
                  {
                    "Key": "5",
                    "Value": "d",
                    "RateLimit": null,
                    "Descriptors": null
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  }
]`

	assert.Equal(t, expected, string(b))
}
