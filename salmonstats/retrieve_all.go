package salmonstats

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

/*
GetAllShifts downloads every shiftSalmonStats from the provided salmon-stats/api server and saves it to a gzipped jsonlines file.
*/
func GetAllShifts(server types.Server, client *http.Client, quiet bool) (errs []error) {
	if !quiet {
		log.Println("Pulling Salmon Run data from online...")
	}
	var jsonLinesWriter *gzip.Writer
	file, err := os.Create(fmt.Sprintf("salmonstats_shifts/%s_out.jl.gz", server.ShortName))
	if err != nil {
		return []error{err}
	}
	jsonLinesWriter = gzip.NewWriter(file)
	getShifts := func(page int) (found bool, errs []error) {
		url := fmt.Sprintf("%splayers/%s/results", server.Address, viper.GetString("user_id"))
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			return false, errs
		}
		query := req.URL.Query()
		query.Set("raw", "1")
		query.Set("count", "200")
		query.Set("page", fmt.Sprint(page))
		req.URL.RawQuery = query.Encode()

		if !quiet {
			log.Println(req.URL)
		}

		resp, err := client.Do(req)
		if err != nil {
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			return false, errs
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				errs = append(errs, err)
			}
		}()
		var data shiftPage
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			return false, errs
		}

		for i := range data.Results {
			if _, err := os.Stat("salmonstats_shifts"); errors.Is(err, os.ErrNotExist) {
				err := os.Mkdir("salmonstats_shifts", os.ModePerm)
				if err != nil {
					errs = append(errs, err)
					buf := make([]byte, 1<<16)
					stackSize := runtime.Stack(buf, false)
					errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
					return false, errs
				}
			}
			fileText, err := json.Marshal(data.Results[i])
			if err != nil {
				errs = append(errs, err)
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, false)
				errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
				return false, errs
			}
			if _, err := jsonLinesWriter.Write(fileText); err != nil {
				errs = append(errs, err)
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, false)
				errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
				return false, errs
			}
			if _, err := jsonLinesWriter.Write([]byte("\n")); err != nil {
				errs = append(errs, err)
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, false)
				errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
				return false, errs
			}
		}
		return len(data.Results) > 0, nil
	}
	if _, err := os.Stat(fmt.Sprintf("salmonstats_shifts/%s.jl.gz", server.ShortName)); err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(fmt.Sprintf("salmonstats_shifts/%s.jl.gz", server.ShortName))
			if err != nil {
				errs = append(errs, err)
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, false)
				errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
				return errs
			}
			if err := f.Close(); err != nil {
				errs = append(errs, err)
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, false)
				errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
				return errs
			}
		}
	}
	f, err := os.Open(fmt.Sprintf("salmonstats_shifts/%s.jl.gz", server.ShortName))
	if err != nil {
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return errs
	}
	count := 0
	gzRead, err := gzip.NewReader(f)
	if err != nil && !errors.Is(err, io.EOF) {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		if err := f.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		return errs
	}
	bufScan := bufio.NewScanner(gzRead)
	if !errors.Is(err, io.EOF) {
		for bufScan.Scan() {
			count++
			if _, err := jsonLinesWriter.Write([]byte(bufScan.Text())); err != nil {
				errs = append(errs, err)
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, false)
				errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
				if err := f.Close(); err != nil {
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
				return errs
			}
			if _, err := jsonLinesWriter.Write([]byte("\n")); err != nil {
				errs = append(errs, err)
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, false)
				errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
				if err := f.Close(); err != nil {
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
				return errs
			}
		}
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
			if err := f.Close(); err != nil {
				errs = append(errs, err)
			}
			return errs
		}
	}
	if err := f.Close(); err != nil {
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
		return errs
	}

	page := count/200 + 1
	hasPages := true
	for hasPages {
		hasPages, errs = getShifts(page)
		if len(errs) > 0 {
			return errs
		}
		page++
	}
	if err := jsonLinesWriter.Close(); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		return errs
	}
	if err := file.Close(); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return errs
	}
	if err := os.Remove(fmt.Sprintf("salmonstats_shifts/%s.jl.gz", server.ShortName)); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return errs
	}
	if err := os.Rename(fmt.Sprintf("salmonstats_shifts/%s_out.jl.gz", server.ShortName), fmt.Sprintf("salmonstats_shifts/%s.jl.gz", server.ShortName)); err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return errs
	}
	return nil
}

type shiftSalmonStatsIterator struct {
	serverAddr string
	f          *os.File
	buffRead   *bufio.Scanner
	gzipReader *gzip.Reader
}

func (s *shiftSalmonStatsIterator) Next() (shift core.Shift, errs []error) {
	data := shiftSalmonStats{}
	if s.buffRead.Scan() {
		if err := json.Unmarshal([]byte(s.buffRead.Text()), &data); err != nil {
			errs = append(errs, err)
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, false)
			errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
			if err := s.f.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := s.gzipReader.Close(); err != nil {
				errs = append(errs, err)
			}
			return nil, errs
		}
		return &data, nil
	}
	if err := s.f.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := s.gzipReader.Close(); err != nil {
		errs = append(errs, err)
	}
	errs = append(errs, &core.NoMoreShiftsError{})
	return nil, errs
}

/*
LoadFromFileIterator creates a core.ShiftIterator that iterates over the salmon-stats/api jsonlimnes in the file.
*/
func LoadFromFileIterator(server types.Server) (core.ShiftIterator, []error) {
	errs := []error{}
	iter := &shiftSalmonStatsIterator{serverAddr: server.Address}
	var err error
	iter.f, err = os.Open(fmt.Sprintf("salmonstats_shifts/%s.jl.gz", server.ShortName))
	if err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return nil, errs
	}
	iter.gzipReader, err = gzip.NewReader(iter.f)
	if err != nil {
		errs = append(errs, err)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)
		errs = append(errs, fmt.Errorf("%s", buf[0:stackSize]))
		return nil, errs
	}
	iter.buffRead = bufio.NewScanner(iter.gzipReader)
	return iter, nil
}

func (s *shiftSalmonStatsIterator) GetAddress() string {
	return s.serverAddr
}

func (s *shiftSalmonStats) GetClearWave() int {
	return s.ClearWaves
}
