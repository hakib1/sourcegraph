package txemail

import (
	"strconv"
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/jordan-wright/email"
	"github.com/stretchr/testify/require"

	"github.com/sourcegraph/sourcegraph/internal/txemail/txtypes"
)

func TestParseTemplate(t *testing.T) {
	type assertEmail struct {
		Subject string
		Text    string
		HTML    string
	}
	var emailData = struct {
		A string
		B string
	}{
		A: "A",
		B: `B`,
	}

	for i, tc := range []struct {
		template txtypes.Templates
		want     autogold.Value
	}{
		{
			template: txtypes.Templates{
				Subject: `{{.A}} subject {{.B}}`,
				Text: `
{{.A}} text body {{.B}}
`,
				HTML: `
{{.A}} html body <span class="{{.B}}" />
`,
			},
			want: autogold.Expect(assertEmail{
				Subject: "A subject B", Text: "A text body B",
				HTML: `A html body <span class="B" />`,
			}),
		},
		{
			template: txtypes.Templates{
				Subject: `{{.A}} subject {{.B}}`,
				Text:    "",
				HTML: `
{{.A}} html body <span class="{{.B}}" />
`,
			},
			want: autogold.Expect(assertEmail{
				Subject: "A subject B", Text: "A html body",
				HTML: `A html body <span class="B" />`,
			}),
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			pt, err := ParseTemplate(tc.template)
			require.NoError(t, err)

			var m email.Email
			err = renderTemplate(pt, emailData, &m)
			require.NoError(t, err)

			// Assert fields of interest as strings for ease of readability
			tc.want.Equal(t, assertEmail{
				Subject: m.Subject,
				HTML:    string(m.HTML),
				Text:    string(m.Text),
			})
		})
	}
}
