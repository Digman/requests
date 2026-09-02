package requests

import (
	"testing"

	"github.com/Digman/tls-client/profiles"
)

func TestGetClientProfileUsesLatestBrowserProfiles(t *testing.T) {
	tests := []struct {
		name      string
		userAgent string
		want      string
	}{
		{
			name:      "chrome 150",
			userAgent: "Mozilla/5.0 Chrome/150.0.0.0 Safari/537.36",
			want:      profiles.Chrome_150_PSK.GetClientHelloStr(),
		},
		{
			name:      "chrome 151 falls back to 150",
			userAgent: "Mozilla/5.0 Chrome/151.0.0.0 Safari/537.36",
			want:      profiles.Chrome_150_PSK.GetClientHelloStr(),
		},
		{
			name:      "chrome 152",
			userAgent: "Mozilla/5.0 Chrome/152.0.0.0 Safari/537.36",
			want:      profiles.Chrome_152_PSK.GetClientHelloStr(),
		},
		{
			name:      "chrome 149 falls back to 146",
			userAgent: "Mozilla/5.0 Chrome/149.0.0.0 Safari/537.36",
			want:      profiles.Chrome_146_PSK.GetClientHelloStr(),
		},
		{
			name:      "firefox 148",
			userAgent: "Mozilla/5.0 Firefox/148.0",
			want:      profiles.Firefox_148.GetClientHelloStr(),
		},
		{
			name:      "brave 146",
			userAgent: "Mozilla/5.0 Chrome/146.0.0.0 Safari/537.36 Brave/1.80.115",
			want:      profiles.Brave_146_PSK.GetClientHelloStr(),
		},
		{
			name:      "safari macos 26",
			userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 26_0) AppleWebKit/605.1.15 Version/26.0 Safari/605.1.15",
			want:      profiles.Safari_26.GetClientHelloStr(),
		},
		{
			name:      "safari ios 26",
			userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 26_0 like Mac OS X) AppleWebKit/605.1.15 Version/26.0 Mobile/15E148 Safari/604.1",
			want:      profiles.Safari_IOS_26_0.GetClientHelloStr(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getClientProfile(tt.userAgent).GetClientHelloStr(); got != tt.want {
				t.Fatalf("getClientProfile() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestShouldRandomizeTLSExtensionOrder(t *testing.T) {
	tests := []struct {
		name    string
		profile profiles.ClientProfile
		want    bool
	}{
		{name: "chrome", profile: profiles.Chrome_152, want: true},
		{name: "brave", profile: profiles.Brave_146, want: true},
		{name: "firefox", profile: profiles.Firefox_148, want: false},
		{name: "safari", profile: profiles.Safari_IOS_26_0, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldRandomizeTLSExtensionOrder(tt.profile); got != tt.want {
				t.Fatalf("shouldRandomizeTLSExtensionOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
