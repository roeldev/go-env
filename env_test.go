package env

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_Merge(t *testing.T) {
	tests := map[string]struct {
		env   Map
		merge map[string]string
		want  Map
	}{
		"append": {
			env:   Map{"foo": "bar"},
			merge: map[string]string{"qux": "xoo"},
			want:  Map{"foo": "bar", "qux": "xoo"},
		},
		"replace": {
			env:   Map{"foo": "bar", "qux": "xoo"},
			merge: map[string]string{"qux": "bar", "foo": "xoo"},
			want:  Map{"foo": "xoo", "qux": "bar"},
		},
		"merge": {
			env:   Map{"foo": "bar", "bar": "baz"},
			merge: map[string]string{"baz": "foo", "bar": "qux"},
			want:  Map{"foo": "bar", "bar": "qux", "baz": "foo"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.env.Merge(tc.merge)
			assert.Exactly(t, tc.want, tc.env)
		})
	}
}

func TestLookupEnv(t *testing.T) {
	map1 := Map{"foo": "bar", "qux": "xoo"}
	map2 := Map{"bruce": "batman", "clark": "superman"}
	se1, se2 := os.LookupEnv("GOROOT")

	tests := map[string]struct {
		maps []Map
		key  string
		want [2]interface{}
	}{
		"empty map": {
			maps: []Map{{}},
			key:  "foo",
			want: [2]interface{}{"", false},
		},
		"one map": {
			maps: []Map{map1},
			key:  "foo",
			want: [2]interface{}{"bar", true},
		},
		"one map, invalid key": {
			maps: []Map{map1},
			key:  "bar",
			want: [2]interface{}{"", false},
		},
		"two maps": {
			maps: []Map{map1, map2},
			key:  "clark",
			want: [2]interface{}{"superman", true},
		},
		"two maps, invalid key": {
			maps: []Map{map1, map2},
			key:  "peter",
			want: [2]interface{}{"", false},
		},
		"system env": {
			maps: []Map{map1, map2},
			key:  "GOROOT",
			want: [2]interface{}{se1, se2},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			a, b := LookupEnv(tc.key, tc.maps...)
			x, y := tc.want[0], tc.want[1]

			assert.Exactly(t, x, a)
			assert.Exactly(t, y, b)
			assert.Exactly(t, x, Getenv(tc.key, tc.maps...))
		})
	}
}

func TestEnviron(t *testing.T) {
	have, n := Environ()
	want := os.Environ()

	assert.Equal(t, len(want), n)

	for _, w := range want {
		key, wantVal := ParsePair(w)
		haveVal, _ := have[key]

		if !assert.Equal(t, wantVal, haveVal) {
			assert.Contains(t, wantVal, haveVal)
		}
	}
}

func TestOpen_not_existing_file(t *testing.T) {
	have := make(Map)
	n, err := Open("doesnot.exist", have)

	assert.Equal(t, 0, n, "should return 0 parsed lines")
	assert.Error(t, err, "expecting an error when trying to open a file that does not exist")
}

func TestOpen(t *testing.T) {
	have := make(Map)
	n, err := Open("test.env", have)

	assert.NoError(t, err)

	want := Map{
		"FOO": "bar",
		"bar": "baz",
		"qux": "#xoo",
	}

	assert.Equal(t, len(want), n, "return value should be the number of parsed items")
	assert.Exactly(t, want, have)
}

func TestRead(t *testing.T) {
	r := strings.NewReader(`FOO=bar
bar='baz'
qux="#xoo"
`)

	have := make(Map)
	n, err := Read(r, have)

	assert.NoError(t, err)

	want := Map{
		"FOO": "bar",
		"bar": "baz",
		"qux": "#xoo",
	}

	assert.Equal(t, len(want), n, "return value should be the number of parsed items")
	assert.Exactly(t, want, have)
}

func TestParseFlagArgs(t *testing.T) {
	tests := map[string]struct {
		flag    string
		input   []string
		wantRes []string
		wantMap Map
	}{
		"empty": {
			flag:    "e",
			input:   []string{},
			wantRes: []string{},
			wantMap: Map{},
		},
		"none": {
			flag:    "e",
			input:   []string{"-a", "-e", "-b=1", "-c", "2", "-e"},
			wantRes: []string{"-a", "-b=1", "-c", "2"},
			wantMap: Map{},
		},
		"single dash": {
			flag:    "e",
			input:   []string{"-e=foo=bar"},
			wantRes: []string{},
			wantMap: Map{"foo": "bar"},
		},
		"single dash next arg": {
			flag:    "env",
			input:   []string{"-env", "foo=bar"},
			wantRes: []string{},
			wantMap: Map{"foo": "bar"},
		},
		"single double dash": {
			flag:    "env",
			input:   []string{"--env=qux=xoo"},
			wantRes: []string{},
			wantMap: Map{"qux": "xoo"},
		},
		"single double dash next arg": {
			flag:    "e",
			input:   []string{"--e", "qux=xoo"},
			wantRes: []string{},
			wantMap: Map{"qux": "xoo"},
		},
		"mixed": {
			flag:    "e",
			input:   []string{"-e", "foo=bar", "-e=empty", "bar", "--e=qux=xoo", "-e", "bar=baz", "--skip", "-e"},
			wantRes: []string{"-e=empty", "bar", "--skip"},
			wantMap: Map{"foo": "bar", "qux": "xoo", "bar": "baz"},
		},
		"multi mixed": {
			flag:    "e",
			input:   []string{"-e", "foo=bar", "empty=", "bar", "--e=qux=xoo", "nop=nop", "-e", "bar=baz", "--skip", "-e"},
			wantRes: []string{"bar", "nop=nop", "--skip"},
			wantMap: Map{"foo": "bar", "empty": "", "qux": "xoo", "bar": "baz"},
		},
		"lookahead": {
			flag:    "e",
			input:   []string{"-e", "-t", "foo=bar", "-e", "qux=xoo", "empty=", "baz"},
			wantRes: []string{"-t", "foo=bar", "baz"},
			wantMap: Map{"qux": "xoo", "empty": ""},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			haveMap := make(Map)
			haveRes, n := ParseFlagArgs(tc.flag, tc.input, haveMap)

			assert.Exactly(t, tc.wantMap, haveMap, "destination map should include all parsed env. vars")
			assert.Exactly(t, tc.wantRes, haveRes, "return value should be without parsed env. vars")
			assert.Equal(t, len(tc.wantMap), n, "second return value should be the number of parsed items")
		})
	}
}

func TestParsePair(t *testing.T) {
	tests := map[string][2]string{
		"=::=::":         {"=::", "::"}, // legit windows entry
		"foo=bar":        {"foo", "bar"},
		"bar='baz'":      {"bar", "baz"},
		`qUx="xoo"`:      {"qUx", "xoo"},
		`PASSWD=$ecR3t`:  {"PASSWD", "$ecR3t"},
		"# some comment": {"", ""},
		"empty=":         {"empty", ""},
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			key, val := ParsePair(input)

			assert.Exactly(t, want[0], key)
			assert.Exactly(t, want[1], val)
		})
	}
}
