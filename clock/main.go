package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/dustin/go-humanize"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "clock"
	app.Version = "2.1"
	app.Usage = "a simple timekeeping utility"
	app.UsageText = "clock [-ncul] [-tz=<zone>] <fmt>\n   clock [global opts] cmd [cmdopts]"
	app.Action = clock
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "noline",
			Aliases: []string{"n"},
			Usage:   "do not print a newline (useful for pipe)",
		},
		&cli.BoolFlag{
			Name:    "copy",
			Aliases: []string{"c"},
			Usage:   "copy the output to the clipboard for easy paste",
		},
		&cli.StringFlag{
			Name:    "tz",
			Aliases: []string{"t"},
			Usage:   "specify the timezone as UTC, Local or an IANA database timezone",
			Value:   "Local",
			EnvVars: []string{"TZ"},
		},
		&cli.BoolFlag{
			Name:    "utc",
			Aliases: []string{"u"},
			Usage:   "shortcut for -tz=utc",
		},
		&cli.BoolFlag{
			Name:    "local",
			Aliases: []string{"l"},
			Usage:   "shortcut for -tz=local",
		},
	}

	// Define other commands available to the application
	app.Commands = []*cli.Command{
		{
			Name:      "after",
			Usage:     "get the date or time after the specified duration",
			UsageText: "clock [global opts] after [opts] <duration>",
			Action:    after,
		},
		{
			Name:      "until",
			Usage:     "get the amount of time until the specified date/time",
			UsageText: "clock [global opts] until [opts] datetime",
			Action:    until,
		},
		{
			Name:      "fmt",
			Usage:     "list the named format strings and describe formats",
			UsageText: "clock fmt",
			Action:    fmtHelp,
		},
	}

	// Run the CLI program
	app.Run(os.Args)
}

//===========================================================================
// CLI Commands
//===========================================================================

func clock(c *cli.Context) (err error) {
	// Get the current time in the specified location
	var loc *time.Location
	locName := c.String("tz")
	if c.Bool("local") {
		locName = "Local"
	}
	if c.Bool("utc") {
		locName = "UTC"
	}
	if loc, err = time.LoadLocation(locName); err != nil {
		return cli.Exit(fmt.Errorf("cannot parse location %q", locName), 1)
	}

	dt := time.Now().In(loc)

	// Determine how to output the time
	var layout string
	if layout, err = parseLayout(strings.Join(c.Args().Slice(), " ")); err != nil {
		return cli.Exit(err, 1)
	}

	ts := dt.Format(layout)

	if c.Bool("copy") {
		if clipboard.Unsupported {
			return cli.Exit("clipboard not supported", 1)
		}
		clipboard.WriteAll(ts)
	} else {
		if c.Bool("noline") {
			fmt.Print(ts)
		} else {
			fmt.Println(ts)
		}
	}

	return nil
}

func after(c *cli.Context) (err error) {
	return cli.Exit("this command has not been implemented yet", 1)
}

func until(c *cli.Context) (err error) {
	// Parse the input date, time or datetime
	var ts time.Time
	if ts, err = parseDatetime(strings.Join(c.Args().Slice(), " "), c.String("tz"), c.Bool("local"), c.Bool("utc")); err != nil {
		return cli.Exit(err, 1)
	}

	if c.Bool("copy") {
		if clipboard.Unsupported {
			return cli.Exit("clipboard not supported", 1)
		}
		clipboard.WriteAll(humanize.Time(ts))
	} else {
		if c.Bool("noline") {
			fmt.Print(humanize.Time(ts))
		} else {
			fmt.Println(humanize.Time(ts))
		}
	}

	return nil
}

var fmtHelpStr = `
The clock command prints out the current timestamp with a specific format so that you
can use the timestamp in a variety of applications. The most specific way to lay out a
command is to use the Go time layout by specifying the format via a reference time:

Mon Jan 2 15:04:05 MST 2006

You can also specify one of the following helper formats:

- json (default)
- code
- date
- today
- blog
- file
- ansic
- ruby
- unix
- kitchen
- rfc3339 (or rfc3339nano)
- rfc822 (or rfc822z)
- rfc850
- rfc1123 (or rfc1123z)
- stamp (or stampmilli, stampmicro, stampnano)
`

func fmtHelp(c *cli.Context) (err error) {
	fmt.Println(strings.TrimSpace(fmtHelpStr))
	return nil
}

//===========================================================================
// Helper Function
//===========================================================================

// parse the layout name or verify that the layout is valid
func parseLayout(s string) (layout string, err error) {
	name := strings.ToLower(s)
	switch name {
	case "", "json", "rfc3339":
		return time.RFC3339, nil
	case "code":
		return "Mon Jan 02 15:04:05 2006 -0700", nil
	case "date", "today":
		return "January 02, 2006", nil
	case "blog":
		return "2020-01-02 15:04:05 -0700", nil
	case "file":
		return "202001021504", nil
	case "ansic":
		return time.ANSIC, nil
	case "ruby":
		return time.RubyDate, nil
	case "unix":
		return time.UnixDate, nil
	case "kitchen":
		return time.Kitchen, nil
	case "rfc3339nano":
		return time.RFC3339Nano, nil
	case "rfc822":
		return time.RFC822, nil
	case "rfc822z":
		return time.RFC822Z, nil
	case "rfc850":
		return time.RFC850, nil
	case "rfc1123":
		return time.RFC1123, nil
	case "rfc1123z":
		return time.RFC1123Z, nil
	case "stamp":
		return time.Stamp, nil
	case "stampmilli":
		return time.StampMilli, nil
	case "stampmicro":
		return time.StampMicro, nil
	case "stampnano":
		return time.StampNano, nil
	}

	var dt time.Time
	if dt, err = time.Parse(s, time.Now().Format(s)); err != nil || dt.IsZero() {
		// Why does this not return isZero?!
		return "", fmt.Errorf("%q is not a valid layout or layout name", s)
	}
	return s, nil
}

func parseDatetime(s, tz string, local, utc bool) (dt time.Time, err error) {
	var loc *time.Location
	switch {
	case local:
		loc, err = time.LoadLocation("Local")
	case utc:
		loc, err = time.LoadLocation("UTC")
	default:
		loc, err = time.LoadLocation(tz)
	}

	if err != nil {
		return time.Time{}, err
	}

	for _, layout := range []string{"2006-01-02", "2006-01-02 15:04", "2006-01-02 15:04:05"} {
		if dt, err = time.ParseInLocation(layout, s, loc); err == nil && !dt.IsZero() {
			return dt, nil
		}
	}

	if dt, err = time.ParseInLocation("15:04", s, loc); err == nil && !dt.IsZero() {
		today := time.Now().In(loc)
		return time.Date(today.Year(), today.Month(), today.Day(), dt.Hour(), dt.Minute(), 0, 0, loc), nil
	}

	return time.Time{}, fmt.Errorf("could not parse %q into a datetime", s)
}
