package sendgrid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSuppressDiffForPendingUsers(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		old          string
		new          string
		userStatus   string
		isSSO        bool
		wantSuppress bool
	}{
		{
			name:         "pending user with empty old value",
			key:          "username",
			old:          "",
			new:          "newuser",
			userStatus:   "pending",
			isSSO:        false,
			wantSuppress: true,
		},
		{
			name:         "pending user with existing old value",
			key:          "username",
			old:          "olduser",
			new:          "newuser",
			userStatus:   "pending",
			isSSO:        false,
			wantSuppress: false,
		},
		{
			name:         "active user",
			key:          "username",
			old:          "",
			new:          "newuser",
			userStatus:   "active",
			isSSO:        false,
			wantSuppress: false,
		},
		{
			name:         "non-SSO user first_name change",
			key:          "first_name",
			old:          "John",
			new:          "Jane",
			userStatus:   "active",
			isSSO:        false,
			wantSuppress: true,
		},
		{
			name:         "non-SSO user last_name change",
			key:          "last_name",
			old:          "Doe",
			new:          "Smith",
			userStatus:   "active",
			isSSO:        false,
			wantSuppress: true,
		},
		{
			name:         "non-SSO user first_name same value",
			key:          "first_name",
			old:          "John",
			new:          "John",
			userStatus:   "active",
			isSSO:        false,
			wantSuppress: false,
		},
		{
			name:         "SSO user first_name change",
			key:          "first_name",
			old:          "John",
			new:          "Jane",
			userStatus:   "active",
			isSSO:        true,
			wantSuppress: false,
		},
		{
			name:         "SSO user last_name change",
			key:          "last_name",
			old:          "Doe",
			new:          "Smith",
			userStatus:   "active",
			isSSO:        true,
			wantSuppress: false,
		},
		{
			name:         "other field for non-SSO user",
			key:          "email",
			old:          "old@example.com",
			new:          "new@example.com",
			userStatus:   "active",
			isSSO:        false,
			wantSuppress: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock ResourceData
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"user_status": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"is_sso": {
					Type:     schema.TypeBool,
					Optional: true,
				},
			}, map[string]interface{}{
				"user_status": tt.userStatus,
				"is_sso":      tt.isSSO,
			})

			got := suppressDiffForPendingUsers(tt.key, tt.old, tt.new, d)
			if got != tt.wantSuppress {
				t.Errorf("suppressDiffForPendingUsers() = %v, want %v", got, tt.wantSuppress)
			}
		})
	}
}
