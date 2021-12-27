package statink

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core"
	"github.com/cass-dlcm/creamypeanutbutteredsalmon/core/types"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

// GetAllShifts downloads every shiftStatInk from the provided stat.ink server and saves it to a gzipped jsonlines file.
func GetAllShifts(statInkServer types.Server, client *http.Client) (errs []error) {
	var jsonLinesWriter *gzip.Writer
	file, err := os.Create(fmt.Sprintf("statink_shifts/%s_out.jl.gz", statInkServer.ShortName))
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	jsonLinesWriter = gzip.NewWriter(file)
	shift := shiftStatInk{}
	getShift := func(id int) (data []shiftStatInk, errs []error) {
		url := fmt.Sprintf("%suser-salmon", statInkServer.Address)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			errs = append(errs, err)
			return nil, errs
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", statInkServer.APIKey))
		query := req.URL.Query()
		query.Set("newer_than", fmt.Sprint(id))
		query.Set("order", "asc")
		req.URL.RawQuery = query.Encode()
		log.Println(req.URL)
		resp, err := client.Do(req)
		if err != nil {
			errs = append(errs, err)
			return nil, errs
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				errs = append(errs, err)
			}
		}()
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			errs = append(errs, err)
			return nil, errs
		}
		for i := range data {
			if _, err := os.Stat("statink_shifts"); errors.Is(err, os.ErrNotExist) {
				if err := os.Mkdir("statink_shifts", os.ModePerm); err != nil {
					errs = append(errs, err)
					return nil, errs
				}
			}
			fileText, err := json.Marshal(data[i])
			if err != nil {
				errs = append(errs, err)
				return nil, errs
			}
			if _, err := jsonLinesWriter.Write(fileText); err != nil {
				errs = append(errs, err)
				return nil, errs
			}
			if _, err := jsonLinesWriter.Write([]byte("\n")); err != nil {
				errs = append(errs, err)
				return nil, errs
			}
		}
		return data, nil
	}
	if _, err := os.Stat(fmt.Sprintf("statink_shifts/%s.jl.gz", statInkServer.ShortName)); err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(fmt.Sprintf("statink_shifts/%s.jl.gz", statInkServer.ShortName))
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
	fileIn, err := os.Open(fmt.Sprintf("statink_shifts/%s.jl.gz", statInkServer.ShortName))
	if err != nil {
		errs = append(errs, err)
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
	gzipReader, err := gzip.NewReader(fileIn)
	if err != nil && !errors.Is(err, io.EOF) {
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := fileIn.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		errs = append(errs, err)
		return errs
	}
	bufioScan := bufio.NewScanner(gzipReader)
	id := 1
	if !errors.Is(err, io.EOF) {
		for bufioScan.Scan() {
			if err := json.Unmarshal([]byte(bufioScan.Text()), &shift); err != nil {
				errs = append(errs, err)
				if err := fileIn.Close(); err != nil {
					errs = append(errs, err)
				}
				if err := jsonLinesWriter.Close(); err != nil {
					errs = append(errs, err)
				}
				if err := gzipReader.Close(); err != nil {
					errs = append(errs, err)
				}
				if err := file.Close(); err != nil {
					errs = append(errs, err)
				}
				return errs
			}
			id = shift.ID
			if _, err := jsonLinesWriter.Write([]byte(bufioScan.Text())); err != nil {
				errs = append(errs, err)
				if err := fileIn.Close(); err != nil {
					errs = append(errs, err)
				}
				if err := jsonLinesWriter.Close(); err != nil {
					errs = append(errs, err)
				}
				if err := gzipReader.Close(); err != nil {
					errs = append(errs, err)
				}
				if err := file.Close(); err != nil {
					errs = append(errs, err)
				}
				return errs
			}
			if _, err := jsonLinesWriter.Write([]byte("\n")); err != nil {
				errs = append(errs, err)
				if err := fileIn.Close(); err != nil {
					errs = append(errs, err)
				}
				if err := jsonLinesWriter.Close(); err != nil {
					errs = append(errs, err)
				}
				if err := gzipReader.Close(); err != nil {
					errs = append(errs, err)
				}
				if err := file.Close(); err != nil {
					errs = append(errs, err)
				}
				return errs
			}
		}
		if err := gzipReader.Close(); err != nil {
			errs = append(errs, err)
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := fileIn.Close(); err != nil {
				errs = append(errs, err)
			}
			return errs
		}
	}
	if err := fileIn.Close(); err != nil {
		errs = append(errs, err)
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		return errs
	}
	for true {
		tempData, errs2 := getShift(id)
		if errs2 != nil && len(errs2) > 0 {
			errs = append(errs, errs2...)
			return errs
		}
		if len(tempData) == 0 {
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
				if err := file.Close(); err != nil {
					errs = append(errs, err)
				}
				return errs
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
				return errs
			}
			if err := os.Remove(fmt.Sprintf("statink_shifts/%s.jl.gz", statInkServer.ShortName)); err != nil {
				errs = append(errs, err)
				return errs
			}
			if err := os.Rename(fmt.Sprintf("statink_shifts/%s_out.jl.gz", statInkServer.ShortName), fmt.Sprintf("statink_shifts/%s.jl.gz", statInkServer.ShortName)); err != nil {
				errs = append(errs, err)
				return errs
			}
			return nil
		}
		id = tempData[len(tempData)-1].ID
	}
	return nil
}

type shiftStatInkIterator struct {
	serverAddr string
	f          *os.File
	buffRead   *bufio.Scanner
	gzipReader *gzip.Reader
}

func (s *shiftStatInkIterator) Next() (shift core.Shift, errs []error) {
	data := shiftStatInk{}
	if s.buffRead.Scan() {
		if err := json.Unmarshal([]byte(s.buffRead.Text()), &data); err != nil {
			errs = append(errs, err)
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
	if err := s.gzipReader.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := s.f.Close(); err != nil {
		errs = append(errs, err)
	}
	errs = append(errs, &core.NoMoreShiftsError{})
	return nil, errs
}

// LoadFromFileIterator creates a core.ShiftIterator that iterates over the stat.ink jsonlimnes in the file.
func LoadFromFileIterator(server types.Server) (core.ShiftIterator, []error) {
	errs := []error{}
	returnVal := shiftStatInkIterator{serverAddr: server.Address}
	var err error
	returnVal.f, err = os.Open(fmt.Sprintf("statink_shifts/%s.jl.gz", server.ShortName))
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

func (s *shiftStatInkIterator) GetAddress() string {
	return s.serverAddr
}

func (s shiftStatInk) GetClearWave() int {
	return *s.ClearWaves
}