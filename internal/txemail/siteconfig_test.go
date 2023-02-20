package txemail

import (
	"strconv"
	"testing"

	"github.com/hexops/autogold/v2"

	"github.com/sourcegraph/sourcegraph/schema"
)

type mockSiteConf schema.SiteConfiguration

func (m mockSiteConf) SiteConfig() schema.SiteConfiguration { return schema.SiteConfiguration(m) }

func TestValidateSiteConfigTemplates(t *testing.T) {
	for i, tt := range []struct {
		conf mockSiteConf
		want autogold.Value
	}{
		{
			conf: mockSiteConf{
				EmailTemplates: nil,
			},
			want: autogold.Expect([]string{}),
		},
		{
			conf: mockSiteConf{
				EmailTemplates: &schema.EmailTemplates{},
			},
			want: autogold.Expect([]string{}),
		},
		{
			conf: mockSiteConf{
				EmailTemplates: &schema.EmailTemplates{
					SetPassword: &schema.EmailTemplate{
						Subject: "",
						Text:    "",
						Html:    "<body>hello world from {{.Host}}</body>",
					},
				},
			},
			want: autogold.Expect([]string{"`email.templates.setPassword` is invalid: fields 'subject' and 'html' are required"}),
		},
		{
			conf: mockSiteConf{
				EmailTemplates: &schema.EmailTemplates{
					SetPassword: &schema.EmailTemplate{
						Subject: "Set up your Sourcegraph Cloud account for {{.Host}}!",
						Text:    "",
						Html:    "<body>hello world from {{.Host}}</body>",
					},
				},
			},
			want: autogold.Expect([]string{}),
		},
		{
			conf: mockSiteConf{
				EmailTemplates: &schema.EmailTemplates{
					SetPassword: &schema.EmailTemplate{
						Subject: "Set up your Sourcegraph Cloud account for {{.Host}}!",
						Text:    "hello world from {{.Hos",
						Html:    "<body>hello world from {{.Host}}</body>",
					},
				},
			},
			want: autogold.Expect([]string{"`email.templates.setPassword` is invalid: template: :1: unclosed action"}),
		},
		{
			conf: mockSiteConf{
				EmailTemplates: &schema.EmailTemplates{
					SetPassword: &schema.EmailTemplate{
						Subject: "Set up your Sourcegraph Cloud account for {{.Host}}!",
						Text:    "hello world from {{.Host}}",
						Html:    "<body>hello world from {{.Host}}</body>",
					},
				},
			},
			want: autogold.Expect([]string{}),
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			problems := validateSiteConfigTemplates(tt.conf)
			tt.want.Equal(t, problems.Messages())
		})
	}
}
