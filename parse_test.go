package main

import (
	"fmt"
	"testing"
)

type TestCase struct {
	name        string
	fileContent string
	fileName    string
	want        []TechnicalDebt
	err         error
}

func Test_ParseString(t *testing.T) {
	t.Parallel()

	testCases := []TestCase{
		{
			name:     "success when content is valid",
			fileName: "file-lorem.go",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// TODO: td-design > improve package division.
				fmt.Println("Hello tttd!")
			}`,
			want: []TechnicalDebt{
				{
					Type:        "td-design",
					Description: "improve package division.",
					File:        "",
					Line:        8,
				},
			},
			err: nil,
		},
		{
			name:     "empty satd list when no satd exists",
			fileName: "file-lorem.go",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")
			}`,
			want: nil,
			err:  nil,
		},
		{
			name:     "empty satd list when invalid satd exists",
			fileName: "file-lorem.go",
			fileContent: `package main

			import "fmt"

			func main() {
				// TODO: example
				fmt.Println("Hello world!")
			}`,
			want: nil,
			err:  nil,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// <RUN>: output
			got, err := Parse(tt.fileContent)
			if err != nil {
				t.Errorf("error %v", err)
			}

			// <VALIDATE>: verify results
			for _, g := range got {
				for _, w := range tt.want {
					if g.Type != w.Type {
						t.Errorf("type: got %q, wanted %q", g.Type, w.Type)
					}

					if g.Description != w.Description {
						t.Errorf("description: got %q, wanted %q", g.Description, w.Description)
					}

					if g.File != w.File {
						t.Errorf("file: got %q, wanted %q", g.File, w.File)
					}

					if g.Line != w.Line {
						t.Errorf("line: got %q, wanted %q", g.Line, w.Line)
					}
				}
			}
		})
	}
}

func Test_ParseRegex(t *testing.T) {
	t.Parallel()

	testCases := []TestCase{
		{
			name: "success when content is valid",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// TODO: td-design > improve package division.
				fmt.Println("Hello tttd!")
			}`,
			want: []TechnicalDebt{
				{
					Type:        "td-design",
					Description: "improve package division.",
					File:        "",
					Line:        8,
				},
			},
			err: nil,
		},
		{
			name: "empty satd list when no satd exists",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")
			}`,
			want: nil,
			err:  nil,
		},
		{
			name: "empty satd list when invalid satd exists",
			fileContent: `package main

			import "fmt"

			func main() {
				// TODO: example
				fmt.Println("Hello world!")
			}`,
			want: nil,
			err:  nil,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// <RUN>: output
			got, err := ParseRegex(tt.fileContent, tt.fileName)
			if err != nil {
				t.Errorf("error %v", err)
			}

			// <VALIDATE>: verify results
			for _, g := range got {
				for _, w := range tt.want {
					fmt.Println(g.Author)
					fmt.Println(g.Description)
					fmt.Println(g.File)
					fmt.Println(g.Line)

					if g.Author != w.Author {
						t.Errorf("author: got %q, wanted %q", g.Author, w.Author)
					}

					if g.Type != w.Type {
						t.Errorf("type: got %q, wanted %q", g.Type, w.Type)
					}

					if g.Description != w.Description {
						t.Errorf("description: got %q, wanted %q", g.Description, w.Description)
					}

					if g.File != w.File {
						t.Errorf("file: got %q, wanted %q", g.File, w.File)
					}

					if g.Line != w.Line {
						t.Errorf("line: got %q, wanted %q", g.Line, w.Line)
					}
				}
			}
		})
	}
}
