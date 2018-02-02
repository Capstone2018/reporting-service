package reports

import "testing"

func TestNewReportValidate(t *testing.T) {
	newReport := &NewReport{
		UserDescription: "This is a valid description, thanks for submitting your report",
		ReportType:      "Misleading Title",
		UserID:          1,
	}

	cases := []struct {
		name        string
		hint        string
		report      *NewReport
		expectError bool
	}{
		{
			"Valid Report",
			"This report is valid so it should validate without error",
			newReport,
			false,
		},
		{
			"Missing User Description",
			"This report has len 0 for user description",
			func(nr NewReport) *NewReport {
				nr.UserDescription = ""
				return &nr
			}(*newReport),
			true,
		},
		{
			"User Description is 300",
			"This report has a user description len == 300 so it should pass the test",
			func(nr NewReport) *NewReport {
				nr.UserDescription = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque tincidunt diam nulla, quis consectetur est sodales at. Maecenas blandit sit amet ante gravida venenatis. Fusce et felis ipsum. Proin hendrerit venenatis ultrices. Maecenas elementum sem a dignissim ornare. Pellentesque interdum ultrices ante, dictum pellentesque quam accumsan vitae. Duis bibendum purus nibh, ac pulvinar magna dapibus vel. Vivamus id consectetur dui, at sodales lacus. Donec id laoreet purus, ultricies accumsan lectus. Fusce iaculis a lorem quis eleifend. In porta tincidunt condimentum.
				Maecenas massa elit, egestas ut malesuada vitae, volutpat id justo. Sed sit amet tempor quam. Quisque sed lectus in lorem suscipit finibus eget vitae nibh. Maecenas pharetra volutpat neque id interdum. Maecenas accumsan tempor magna, suscipit porta risus fringilla a. Fusce tristique, lorem a vestibulum tempor, ex risus accumsan neque, eu blandit leo augue quis elit. Morbi dictum lobortis turpis, in faucibus nunc porttitor et. Curabitur aliquam justo libero, eget facilisis lectus semper et. Aliquam erat volutpat. Nam purus enim, pulvinar quis vestibulum ac, cursus sed eros. In eget efficitur ligula, quis bibendum orci. Sed porttitor nec lectus in sodales. In dapibus dolor at ultricies ullamcorper. In ultrices ullamcorper tempus. Aenean non urna sed odio ultrices hendrerit. Praesent ac felis sit amet ex porttitor lacinia.
				Suspendisse ultrices dignissim aliquam. Vestibulum magna quam, suscipit ac tortor vel, rhoncus suscipit risus. In nisi ante, accumsan at lacus vitae, consectetur imperdiet metus. Pellentesque sed diam tortor. Nunc mattis efficitur consequat. Fusce eros magna, cursus efficitur sagittis id, viverra ac massa. Duis eu convallis felis, ac laoreet nulla. Integer molestie dui turpis, vitae luctus mauris feugiat quis. Phasellus lacinia dolor vel consequat pellentesque. Aliquam sollicitudin auctor lacus quis condimentum. Aliquam finibus sit amet ex in tempor. Cras aliquet et lorem at viverra. Quisque ac tincidunt ante. Sed vitae semper velit, ut egestas purus. Morbi.`
				return &nr
			}(*newReport),
			false,
		},
		{
			"User Description is > 300",
			"This report has a user description len > 300",
			func(nr NewReport) *NewReport {
				nr.UserDescription = `
				Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam bibendum erat vitae massa porttitor mattis. Vivamus efficitur justo ac erat posuere vestibulum. Suspendisse egestas tincidunt mauris eget mattis. Etiam a nunc vel dolor sagittis suscipit. Phasellus consequat ornare nisi consequat efficitur. Etiam quis facilisis erat. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Sed vulputate arcu sit amet massa sagittis, a vestibulum magna lacinia. Nam elementum sem et libero rutrum, et aliquam augue imperdiet. Quisque tempus faucibus sodales. Donec id turpis posuere, semper arcu et, commodo felis.	
				Duis consequat bibendum euismod. Vivamus tristique ullamcorper ipsum, sed gravida ex sagittis at. In ultrices sem magna, quis pulvinar eros molestie facilisis. Vivamus finibus volutpat pellentesque. Fusce malesuada, leo sed fringilla posuere, risus tellus mollis erat, non bibendum nibh ligula ut massa. Sed vitae suscipit massa, in hendrerit velit. Curabitur non maximus ipsum. Nunc mattis lobortis sapien at facilisis. Quisque porta viverra ipsum, quis placerat ante molestie sit amet. Sed ac urna eget mauris posuere placerat. Donec eleifend dictum mattis.	
				Suspendisse ipsum ipsum, aliquam eget turpis id, rutrum iaculis lacus. Vestibulum sagittis ornare metus, ut maximus leo ornare eget. Sed condimentum elementum lacinia. Nam dapibus ipsum nec vulputate pellentesque. Donec sit amet nibh quis magna imperdiet molestie sed in magna. Nulla nec dolor sit amet magna mattis lobortis eget nec ipsum. Curabitur nisl justo, blandit vel ullamcorper et, iaculis vitae leo. Pellentesque scelerisque lobortis eleifend. Aenean at interdum arcu. Fusce lacinia dolor diam, at varius mi lobortis nec. Maecenas felis mauris, luctus ut posuere ut, aliquet vitae risus.		
				Maecenas a vestibulum nulla. Nunc id semper risus. Phasellus venenatis vitae justo non aliquet. Aenean iaculis neque quis nulla volutpat, eget maximus arcu sagittis. Aenean consectetur elit enim. Quisque lacinia, massa eget dignissim elementum, enim ex porttitor magna, a feugiat lacus metus.`
				return &nr
			}(*newReport),
			true,
		},
		{
			"No creator ID provided",
			"This report has no creator ID and should error",
			func(nr NewReport) *NewReport {
				nr.UserID = 0
				return &nr
			}(*newReport),
			true,
		},
		{
			"Report type is Misleading Title",
			"This report is one of the valid report types",
			func(nr NewReport) *NewReport {
				nr.ReportType = "Misleading Title"
				return &nr
			}(*newReport),
			false,
		},
		{
			"Report type is False Information",
			"This report is one of the valid report types",
			func(nr NewReport) *NewReport {
				nr.ReportType = "False Information"
				return &nr
			}(*newReport),
			false,
		},
		{
			"Report type is Other",
			"This report is one of the valid report types",
			func(nr NewReport) *NewReport {
				nr.ReportType = "Other"
				return &nr
			}(*newReport),
			false,
		},
		{
			"Report type is not valid",
			"This report is NOT one of the valid report types",
			func(nr NewReport) *NewReport {
				nr.ReportType = "WILL NOT PASS"
				return &nr
			}(*newReport),
			true,
		},
	}

	for _, c := range cases {
		err := c.report.Validate()
		if !c.expectError && err != nil {
			t.Errorf("case %s: unexpected error validating user: %v\nHINT: %s", c.name, err, c.hint)
		}
		if c.expectError && err == nil {
			t.Errorf("case %s: expected validation error but didn't get one\nHINT: %s", c.name, c.hint)
		}
	}
}

func TestNewReportToReport(t *testing.T) {
	newReport := &NewReport{
		UserDescription: "This is a valid description, thanks for submitting your report",
		ReportType:      "Misleading Title",
		UserID:          1,
	}

	cases := []struct {
		name        string
		hint        string
		report      *NewReport
		expectError bool
	}{
		{
			"Valid New Report",
			"Should not return any errors",
			newReport,
			false,
		},
	}

	for _, c := range cases {
		r, err := c.report.ToReport()
		if !c.expectError && err != nil {
			t.Errorf("case %s: unexpected error validating user: %v\nHINT: %s", c.name, err, c.hint)
		}
		if c.expectError && err == nil {
			t.Errorf("case %s: expected validation error but didn't get one\nHINT: %s", c.name, c.hint)
		}
		if r.CreatedAt.IsZero() {
			t.Errorf("case %s: expected non-zero created at field\nHINT: %s", c.name, c.hint)
		}
	}
}
