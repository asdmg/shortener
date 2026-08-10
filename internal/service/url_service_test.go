package service

import "testing"

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid https url",
			url:     "https://google.com",
			wantErr: false,
		},
		{
			name:    "valid http url",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "empty url",
			url:     "",
			wantErr: true,
		},
		{
			name:    "domain without scheme",
			url:     "google.com",
			wantErr: true,
		},
		{
			name:    "invalid scheme",
			url:     "ftp://google.com",
			wantErr: true,
		},
		{
			name:    "invalid url",
			url:     "abc",
			wantErr: true,
		},
		{
			name:    "missing host",
			url:     "https://",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := validateURL(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"validateURL() error = %v, wantErr = %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
