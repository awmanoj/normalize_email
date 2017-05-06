package normalize_email

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	testCases := map[string]string{
		// Only supported domains
		"a.b.c+tag@example.com": "a.b.c+tag@example.com",

		// Gmail dots
		"a.b.c@gmail.com": "abc@gmail.com",
		"a.b.c@yahoo.com": "a.b.c@yahoo.com",

		// Plus
		"a.b.c+tag@gmail.com": "abc@gmail.com",
		"a.b.c+tag@yahoo.com": "a.b.c+tag@yahoo.com",

		// Non-standard TLDs
		"a.b.c+tag@something.co.uk": "a.b.c+tag@something.co.uk",

		// Yahoo!
		"a.b.c.d+tag@yahoo.com": "a.b.c.d+tag@yahoo.com",
		"a.b.c-tag@yahoo.com":   "a.b.c@yahoo.com",
		"a.b.c-tag@yahoo.co.uk": "a.b.c@yahoo.co.uk",
		"a-b.c-tag@yahoo.ro":    "a@yahoo.ro",

		// Microsoft
		"a.b.c+tag@outlook.com":   "a.b.c@outlook.com",
		"a.b.c-tag@hotmail.com":   "a.b.c-tag@hotmail.com",
		"a.b.c-tag@outlook.co.uk": "a.b.c-tag@outlook.co.uk",
		"a.b.c+d@live.com":        "a.b.c@live.com",

		// Google Apps for Work
		"a.b.c+tag@idorecall.com": "a.b.c+tag@idorecall.com",

		// FastMail
		"a.b.c+tag@fastmail.com": "a.b.c@fastmail.com",
		"a.b.c+tag@fastmail.fm":  "a.b.c@fastmail.fm",
		// http://outcoldman.com/en/archive/2014/05/08/fastmail/
		"denis+tag@outcoldman.com": "denis+tag@outcoldman.com",
	}

	// TODO: force remove dot and force remove tags
	testCases2 := map[string]string{
		// Non-standard TLDs
		"a.b.c+tag@something.co.uk": "abc@something.co.uk",
	}

	// TODO: detect provider
	testCases3 := map[string]string{
		// Google Apps for Work
		"a.b.c+tag@tokopedia.com": "abc@tokopedia.com",
		// Async test Google Apps for Work
		"a.b.c.d+tag@blueseed.com": "abcd@blueseed.com",
		// Async test no special provider
		"ad.missions+impossible@stanford.edu": "ad.missions+impossible@stanford.edu",
	}

	for input, expected := range testCases {
		result := Normalize(input, nil)
		if result != expected {
			t.Errorf("err got %s, expected %s\n", result, expected)
		}
	}

	for input, expected := range testCases2 {
		var option = Option{
			ForceRemoveDots: true,
			ForceRemoveTags: true,
		}
		result := Normalize(input, &option)
		if result != expected {
			t.Errorf("err got %s, expected %s\n", result, expected)
		}
	}

	for input, expected := range testCases3 {
		var option = Option{
			DetectProvider: true,
		}

		result := Normalize(input, &option)
		if result != expected {
			t.Errorf("err got %s, expected %s\n", result, expected)
		}
	}
}
