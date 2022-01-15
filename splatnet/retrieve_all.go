package splatnet

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core"
	"github.com/cass-dlcm/splatnetiksm"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

/*
GetAllShifts downloads every shiftSplatnet from the SplatNet server and saves it to a gzipped jsonlines file.
*/
func GetAllShifts(sessionToken, cookie, locale string, client *http.Client, quiet bool) (*string, *string, []error) {
	var errs []error
	_, timezone := time.Now().Zone()
	timezone = -timezone / 60
	appHead := http.Header{
		"Host":              []string{"app.splatoon2.nintendo.net"},
		"x-unique-id":       []string{"32449507786579989235"},
		"x-requested-with":  []string{"XMLHttpRequest"},
		"x-timezone-offset": []string{fmt.Sprint(timezone)},
		"User-Agent":        []string{"Mozilla/5.0 (Linux; Android 7.1.2; Pixel Build/NJH47D; wv) AppleWebKit/537.36 (KHTML, like Gecko) version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36"},
		"Accept":            []string{"*/*"},
		"Referer":           []string{"https://app.splatoon2.nintendo.net/home"},
		"Accept-Encoding":   []string{"gzip deflate"},
		"Accept-Language":   []string{locale},
	}

	if !quiet {
		log.Println("Pulling Salmon Run data from online...")
	}

	url := "https://app.splatoon2.nintendo.net/api/coop_results"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return &sessionToken, &cookie, errs
	}

	req.Header = appHead

	if cookie == "" {
		newSessionToken, newCookie, errs2 := splatnetiksm.GenNewCookie(locale, sessionToken, "blank", client)
		if len(errs2) > 0 {
			errs = append(errs, errs2...)
			return &sessionToken, &cookie, errs
		}
		sessionToken = *newSessionToken
		cookie = *newCookie
	}

	req.AddCookie(&http.Cookie{Name: "iksm_session", Value: cookie})

	resp, err := client.Do(req)
	if err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return &sessionToken, &cookie, errs
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			errs = append(errs, err)
		}
	}()

	var data shiftList

	var jsonLinesWriter *gzip.Writer
	if _, err := os.Stat("shifts.jl.gz"); err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create("shifts.jl.gz")
			if err != nil {
				errs = append(errs, err)
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, false)
				errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
				return &sessionToken, &cookie, errs
			}
			if err := f.Close(); err != nil {
				errs = append(errs, err)
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, false)
				errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
				return &sessionToken, &cookie, errs
			}
		}
	}
	fileIn, err := os.Open("shifts.jl.gz")
	if err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return &sessionToken, &cookie, errs
	}
	gzRead, err := gzip.NewReader(fileIn)
	eof := false
	if err != nil {
		if !errors.Is(err, io.EOF) {
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			if err := fileIn.Close(); err != nil {
				errs = append(errs, err)
			}
			return &sessionToken, &cookie, errs
		}
		eof = true
	}
	bufScan := bufio.NewScanner(gzRead)
	file, err := os.Create("shifts_out.jl.gz")
	if err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return &sessionToken, &cookie, errs
	}
	jsonLinesWriter = gzip.NewWriter(file)
	var text string

	for !eof && bufScan.Scan() {
		text = bufScan.Text()
		if _, err := jsonLinesWriter.Write([]byte(text)); err != nil {
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			if err := fileIn.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := gzRead.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return &sessionToken, &cookie, errs
		}
		if _, err := jsonLinesWriter.Write([]byte("\n")); err != nil {
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			if err := fileIn.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := gzRead.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return &sessionToken, &cookie, errs
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		newSessionToken, newCookie, errs2 := splatnetiksm.GenNewCookie(locale, sessionToken, "blank", client)
		if len(errs2) > 0 {
			errs = append(errs, errs2...)
		}
		sessionToken = *newSessionToken
		cookie = *newCookie
		if err := fileIn.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := gzRead.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			return &sessionToken, &cookie, errs
		}
		newSessionToken, newCookie, errsRec := GetAllShifts(sessionToken, cookie, locale, client, quiet)
		if len(errsRec) > 0 {
			errs = append(errs, errsRec...)
			return &sessionToken, &cookie, errs
		}
		return newSessionToken, newCookie, nil
	}

	if data.Code != nil {
		newSessionToken, newCookie, errs2 := splatnetiksm.GenNewCookie(locale, sessionToken, "blank", client)
		if len(errs2) > 0 {
			errs = append(errs, errs2...)
		}
		sessionToken = *newSessionToken
		cookie = *newCookie
		if err := fileIn.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := gzRead.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			return &sessionToken, &cookie, errs
		}
		newSessionToken, newCookie, errsRec := GetAllShifts(sessionToken, cookie, locale, client, quiet)
		if len(errsRec) > 0 {
			errs = append(errs, errsRec...)
			return &sessionToken, &cookie, errs
		}
		return newSessionToken, newCookie, nil
	}

	if err := fileIn.Close(); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if !eof {
		if err := gzRead.Close(); err != nil {
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		return &sessionToken, &cookie, errs
	}

	shift := &shiftSplatnet{}
	if json.Valid([]byte(text)) {
		if err := json.Unmarshal([]byte(text), shift); err != nil {
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return &sessionToken, &cookie, errs
		}
	}

	for i := range data.Results {
		if data.Results.inList(*shift) {
			break
		}
		fileText, err := json.Marshal(data.Results[i])
		if err != nil {
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return &sessionToken, &cookie, errs
		}
		if _, err := jsonLinesWriter.Write(fileText); err != nil {
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return &sessionToken, &cookie, errs
		}
		if _, err := jsonLinesWriter.Write([]byte("\n")); err != nil {
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return &sessionToken, &cookie, errs
		}
	}

	if err := jsonLinesWriter.Close(); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
	}
	if err := file.Close(); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
	}
	if len(errs) > 0 {
		return &sessionToken, &cookie, errs
	}

	if err := os.Remove("shifts.jl.gz"); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return &sessionToken, &cookie, errs
	}
	if err := os.Rename("shifts_out.jl.gz", "shifts.jl.gz"); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return &sessionToken, &cookie, errs
	}
	return &sessionToken, &cookie, nil
}

/*
LoadFromFileIterator creates a core.ShiftIterator that iterates over the SplatNet jsonlimnes in the file.
*/
func LoadFromFileIterator() (core.ShiftIterator, []error) {
	returnVal := shiftSplatnetIterator{}
	var errs []error
	var err error
	returnVal.f, err = os.Open("shifts.jl.gz")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return nil, errs
	}
	returnVal.gzipReader, err = gzip.NewReader(returnVal.f)
	if err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return nil, errs
	}
	returnVal.buffRead = bufio.NewScanner(returnVal.gzipReader)
	return &returnVal, nil
}
