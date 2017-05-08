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

## Usage

```go
    ... 
    email := "john123.thomas123+test123@googlemail.com
    normalizedEmail := normalize_email.Normalize(input, nil, nil)
    // normalizedEmail == "john123.thomas@gmail.com"
    ...
```

```go
    ... 
    email := "john123.thomas123+test123@googlemail.com
    callback = make(chan string)
    _ := normalize_email.Normalize(input, &Option{DetectProvider: true}, callback)
    normalizedEmail := <- callback
    // normalizedEmail == "john123.thomas@gmail.com"
    ...
```

Please also refer normalize\_test.go for examples.

_ 

