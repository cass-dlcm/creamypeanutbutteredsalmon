Creamy Peanut Buttered Salmon
================

[![Go Reference](https://pkg.go.dev/badge/github.com/cass-dlcm/creamypeanutbutteredsalmon.svg)](https://pkg.go.dev/github.com/cass-dlcm/peanutbutteredsalmon)
[![Go Report Card](https://goreportcard.com/badge/github.com/cass-dlcm/creamypeanutbutteredsalmon)](https://goreportcard.com/report/github.com/cass-dlcm/creamypeanutbutteredsalmon)
[![DeepSource](https://deepsource.io/gh/cass-dlcm/creamypeanutbutteredsalmon.svg/?label=active+issues&show_trend=true&token=sUeiypSGmfk2eicGTR18wXkW)](https://deepsource.io/gh/cass-dlcm/creamypeanutbutteredsalmon/?ref=repository-badge)

Creamy Peanut Buttered Salmon is a program that downloads data from the SplatNet 2 app ("Nintendo Switch Online") or stat.ink to find personal bests in the *Splatoon 2* Salmon Run game mode.

## Usage

If running from source

```Shell
$ go run cmd/creamyPeanutButteredSalmon.go [-stage ""] [-event ""] [-tide ""] [-weapon ""] [-splatnet] [-statink ""] [-salmonstats ""] [-outfile ""]
```

If running a binary

```Shell
$ ./creamyPeanutButteredSalmon [-stage ""] [-event ""] [-tide ""] [-weapon ""] [-splatnet] [-statink ""] [-salmonstats ""] [-outfile ""]
```

The `-stage` flag takes in a string with the set of stages to include like such: `"spawning_grounds marooners_bay lost_outpost salmonid_smokeyard ruins_of_ark_polaris"`

The `-event` flag takes in a string with the set of events to include like such: `"water_levels rush fog goldie_seeking griller cohock_charge mothership"`

The `-tide` flag takes in a string with the set of tides to include like such: `"LT NT HT"`

The `-weapon` flag takes in a string with the set of weapon types to include like such: `"set single_random four_random random_gold"`

The `-splatnet` flag signals to download data from SplatNet 2.

The `-statink` flag signals to download data from a stat.ink instance. Use `-statink "official"` to use the https://stat.ink instance.

The `-salmonstats` flag downloads data from a salmon-stats instance. Use `-salmonstats "official"` to use the https://salmon-stats-api.yuki.games instance.

The `-outfile` flag specifies to write the JSON output to a file, and enables progress information to standard out.

### Example usage

Running `go run cmd/creamyPeanutButteredSalmon.go -splatnet` from the command line launches the program to check on Splatnet 2 for new results, save them, and find personal bests.

Running `go run cmd/creamyPeanutButteredSalmon.go --statink "official"` finds all personal bests from your data on stat.ink.

## Features

Finds personal bests for the following records

- [x] Total Golden Eggs
- [x] Total Golden Eggs 2 Nights
- [x] Total Golden Eggs 1 Nights
- [x] Total Golden Eggs No Night
- [x] All valid combinations of events and tides

Output features:

- [x] outputs as valid JSON
- [x] structured with record as top level, stage as middle level, and weapon schedule as bottom level of hierarchy
- [x] saving shift data to JSON files

Input features:
- [x] reading from SplatNet 2
- [x] reading from stat.ink
- [x] reading from local files
- [x] reading from salmon-stats.yuki.games

## Setup instructions

*These instructions are meant to be accessible and easy-to-follow for all users, and this is the recommended way to run the script. If you run into trouble, please reach out! However, an alternative [simple version](https://github.com/cass-dlcm/creamypeanutbutteredsalmon/wiki/simple-setup-instructions) is also available.*

1. Download and install Go. On Windows, download the installer from the [official website](https://www.golang.org/dl/). On macOS, install [Homebrew](https://brew.sh/) and then run `brew install go`.

2. If you're on Windows, install [Git](https://git-scm.com/download) (pre-installed on macOS).

3. Download the program from the command line (macOS: Terminal; Windows: Command Prompt/PowerShell) by running `git clone https://github.com/cass-dlcm/peanutbutteredsalmon.git`.

4. You will be prompted to enter your [language code](https://github.com/frozenpandaman/splatnet2statink/wiki/languages) (locale).

* Running the script for the first time with the `-statink` flag will fail if you haven't set an API key. To enter your stat.ink API key, put it in the `config.json` file at the `statink_servers[0].api_key` location. The API key can be found in your [profile settings](https://stat.ink/profile).
* Running the script for the first time with the `-salmonstats` flag will fail if you haven't set a Splatnet user ID. To enter your Splatnet user ID, put it in the `config.json` file at the `user_id` location. The Splatnet user ID can be found in a record downloaded from SplatNet, which would be in any of the `shifts/{job_number},json` files, at the `my_result.pid` location.
* Running the script for the first time with the `-splatnet` flag will start the cookie generation.

**NOTE: Read the "Cookie generation" section below before proceeding. [→](#cookie-generation)**

You will then be asked to navigate to a specific URL on Nintendo.com, log in, and follow simple instructions to obtain your `session_token`; this will be used to generate an `iksm_session` cookie. If you are opting against automatic cookie generation, enter "skip" for this step, at which point you will be asked to manually input your `iksm_session` cookie instead (see the [mitmproxy instructions](https://github.com/frozenpandaman/splatnet2statink/wiki/mitmproxy-instructions)).

This cookie (used to access your SplatNet battle results) along with your stat.ink API key and language will automatically be saved into `config.txt` for you. You're now ready to upload battles!

Have any questions, issues, or suggestions? Feel free to message me on [Twitter](https://twitter.com/cass-dlcm) or create an [issue](https://github.com/cass-dlcm/creamypeanutbutteredsalmon/issues) here.

### Accessing SplatNet 2 from your browser

If you wish to access SplatNet 2 from your computer rather than via the phone app, navigate to [https://app.splatoon2.nintendo.net/home](https://app.splatoon2.nintendo.net/home) (it should show a forbidden error). Use a cookie editor – such as [EditThisCookie](https://chrome.google.com/webstore/detail/editthiscookie/fngmhnnpilhplaeedifhccceomclgfbg?hl=en) for Chrome – to change `iksm_session` to the value you obtained previously (automatically or via [mitmproxy](https://github.com/frozenpandaman/splatnet2statink/wiki/mitmproxy-instructions), stored as  `cookie` in `config.txt`), and refresh the page. If you only want to access SplatNet and don't have a stat.ink API key, simply enter "skip" for this step during setup.

*Splatoon 2* stage rotation information (including Salmon Run) and current SplatNet gear are viewable at [splatoon2.ink](https://splatoon2.ink/).

---

## Cookie generation

For Creamy Peanut Buttered Salmon to work with SplatNet, a cookie known as `iksm_session` is required. This cookie may be obtained automatically, using the program, or manually via the app. Please read the following sections carefully to decide whether or not you want to use automatic cookie generation.

### Automatic

Automatic cookie generation involves making a *secure request to two non-Nintendo servers with minimal, non-identifying information*. We aim to be 100% transparent about this and provide in-depth information on security and privacy below. Users who feel uncomfortable with this may opt to manually acquire their cookie instead.

The v1.1.0 update to the Nintendo Switch Online app, released in September 2017, introduced the requirement of a [message authentication code](https://en.wikipedia.org/wiki/Message_authentication_code) (known as `f`), thereby complicating the ability to generate cookies within the script. After figuring out the [key](https://en.wikipedia.org/wiki/Key_\(cryptography\)) previously used to generate `f` tokens, the calculation method was changed in September 2018's v1.4.1 update, heavily obfuscating the new process. As a workaround, an Android server was set up to emulate the app, specifically to generate `f` tokens.

Generation now requires a [hash value](https://en.wikipedia.org/wiki/Hash_function) to further verify the authenticity of the request. The algorithm to calculate this, originally done within the app, is sensitive; to prevent sharing it publicly (i.e. distributing it in the script's source code), github user frozenpandaman created a small [API](https://en.wikipedia.org/wiki/Application_programming_interface) which generates a hash value given a valid input. This can be passed to the Android server to generate the corresponding `f` token, which is then used to retrieve an `iksm_session` cookie.

**Privacy statement:** No identifying information is ever sent to the API server. Usernames and passwords are far removed from where the API comes into play and are never readable by anyone but you. Returned hash values are never logged or stored and do not contain meaningful information. It is not possible to use either sent or stored data to identify which account/user performed a request, to view any identifying information about a user, or to gain access to an account.

See the **[API documentation wiki page](https://github.com/frozenpandaman/splatnet2statink/wiki/api-docs)** for more information.

### Manual

Users who decide against automatic cookie generation via their computer may instead generate/retrieve `iksm_session` cookies from the SplatNet app.

In this case, users must obtain their cookie from their phone by intercepting their device's web traffic and inputting it into Creamy Peanut Buttered Salmon when prompted (or manually adding it to `config.txt`). Follow the [mitmproxy instructions](https://github.com/frozenpandaman/splatnet2statink/wiki/mitmproxy-instructions) to obtain and configure your cookie manually. To opt to manually acquire your cookie, enter "skip" when prompted to enter the "Select this account" URL.

## License

[AGPLv3](https://www.gnu.org/licenses/agpl-3.0.html) or any later version
