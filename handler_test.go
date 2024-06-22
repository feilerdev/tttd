package main

import "testing"

func Test_Encode(t *testing.T) {
	t.Parallel()

	testCases := []TestCase{
		{
			name:        "success when content is valid",
			fileContent: `[".github/workflows/satd-detection.yaml","main.go"]`,
			want:        nil,
			err:         nil,
		}, {
			name:        "success when content is empty",
			fileContent: `[]`,
			want:        nil,
			err:         nil,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			// <RUN>: output
			_, err := Decode(tt.fileContent)
			if err != nil {
				t.Errorf("error %v", err)
			}

			// <VALIDATE>: verify results
		})
	}
}
