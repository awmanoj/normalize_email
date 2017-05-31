# normalize-email 

## Based on
This repo is a golang port of JS library: https://github.com/iDoRecall/email-normalize

## Introduction 
normalize-email normalizes email addresses to their canonical form.

* remove dots in GMail or email addresses hosted by similar services (Google Apps for Work, FastMail)
* remove [address tags](https://en.wikipedia.org/wiki/Email_address#Sub-addressing) starting with '+', or '-' for Yahoo!, or '=' as well if desired
* converts alias domains to the canonical one, e.g. googlemail.com to gmail

## Works for

* `gmail.com`, `googlemail.com` and `google.com` - the [only domains Google is known to use for incoming email](https://en.wikipedia.org/wiki/List_of_Google_domains)
* Any domain [hosted with Google Apps for Work](https://en.wikipedia.org/wiki/Google_Apps_for_Work#Gmail)
* FastMail and [domains hosted with FastMail](https://www.fastmail.com/help/receive/domains.html)
* Microsoft's outlook.com, hotmail.com, live.com
* Yahoo! domains

## Sample Normalizations 

```
		// Gmail dots
		"a.b.c@gmail.com": "abc@gmail.com",
		"a.b.c@yahoo.com": "a.b.c@yahoo.com",

		// Google domains
		"a.b.c@googlemail.com": "abc@gmail.com",

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
```


## Usage
sync version
```go
    ... 
    rawEmail := "john123.thomas123+test123@googlemail.com
    normEmail := normalize_email.Normalize(rawEmail, nil, nil)

    // normEmail == "john123.thomas@gmail.com"
    ...
```
async version for detect provider

```go
    ... 
    callback = make(chan string)
    ... 
    rawEmail := "john123.thomas123+test123@googlemail.com
    _ := normalize_email.Normalize(rawEmail, &Option{DetectProvider: true}, callback)
    ... 
    normEmail := <- callback

    // normEmail == "john123.thomas@gmail.com"
    ...
```

Please also refer normalize\_test.go for examples.

_ 

