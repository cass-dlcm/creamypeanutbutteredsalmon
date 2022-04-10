package types

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var currVersion = version{0, 8, 0}

type version struct {
	Major  uint64
	Minor  uint64
	Bugfix uint64
}

func (v *version) compareVersion(v2 *version) int {
	if v.Major > v2.Major {
		return -3
	}
	if v.Major < v2.Major {
		return 3
	}
	if v.Minor > v2.Minor {
		return -2
	}
	if v.Minor < v2.Minor {
		return 2
	}
	if v.Bugfix > v2.Bugfix {
		return -1
	}
	if v.Bugfix < v2.Bugfix {
		return 1
	}
	return 0
}

func (v *version) toString() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Bugfix)
}

type releasesType []struct {
	URL             string    `json:"url"`
	HTMLURL         string    `json:"html_url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	ID              int       `json:"id"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Body            string    `json:"body"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Assets []struct {
		URL                string    `json:"url"`
		BrowserDownloadURL string    `json:"browser_download_url"`
		ID                 int       `json:"id"`
		NodeID             string    `json:"node_id"`
		Name               string    `json:"name"`
		Label              string    `json:"label"`
		State              string    `json:"state"`
		ContentType        string    `json:"content_type"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		Uploader           struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
	} `json:"assets"`
}

/*
CheckForUpdate downloads the latest version of this file and checks to see which version number is higher.
If the downloaded version is higher, execution ends.
If the binary version is higher, it warns about being a prerelease.
*/
func CheckForUpdate(client *http.Client, quiet bool) (errs []error) {
	url := "https://api.github.com/repos/cass-dlcm/creamypeanutbutteredsalmon/releases"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	resp, err := client.Do(req)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			errs = append(errs, err)
		}
	}()
	if resp.StatusCode == http.StatusNotFound {
		log.Println("Can't find the current version. Maybe none have been released yet?")
		return nil
	}
	var releases releasesType
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		errs = append(errs, err)
		return errs
	}
	versionStr := releases[0].TagName
	versionSubstrs := strings.Split(versionStr, ".")
	if versionSubstrs[0] == "" {
		return nil
	}
	major, err := strconv.ParseUint(versionSubstrs[0][1:], 10, 32)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	minor, err := strconv.ParseUint(versionSubstrs[1], 10, 32)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	bugfix, err := strconv.ParseUint(versionSubstrs[2], 10, 32)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	testVers := version{major, minor, bugfix}
	versionComparison := currVersion.compareVersion(&testVers)
	if versionComparison >= 1 {
		errs = append(errs, fmt.Errorf("A new version is available. Please update to the new version.\nCurrent Version: %s\nNew Version: %s\nExiting.", currVersion.toString(), testVers.toString()))
		return errs
	} else if versionComparison <= -1 && !quiet {
		log.Printf("You are running a unreleased version.\nLatest released version:%s\nCurrent version:%s\n", testVers.toString(), currVersion.toString())
	}
	return errs
}
