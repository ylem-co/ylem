package helpers

import (
    _"fmt"
    "math"
    "strconv"
    "strings"
)

func plural(count int, singular string) (result string) {
    if (count == 1) {
        result = strconv.Itoa(count) + " " + singular + "\n"
    } else if (count == 0) {
              result = ""
    } else {
        result = strconv.Itoa(count) + " " + singular + "s\n"
    }
    return
}

func ReleaseSecondsToHuman(input float64) (result string) {
    if (math.Abs(input) < 60 * 60 * 24) {
        result = "less than a day";
        return;
    }

    result = SecondsToHuman(input);

    return;
}

func SecondsToHuman(input float64) (result string) {
    negative := false
    if (input < 0) {
        negative = true
        input = math.Abs(input);
    }

    years := math.Floor(input / 60 / 60 / 24 / 7 / 4 / 12)

    intInput := int(input)

    seconds := intInput % (60 * 60 * 24 * 7 * 4 * 12)
    months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 4)
    seconds = intInput % (60 * 60 * 24 * 7 * 4)
    weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
    seconds = intInput % (60 * 60 * 24 * 7)
    days := math.Floor(float64(seconds) / 60 / 60 / 24)
    seconds = intInput % (60 * 60 * 24)
    hours := math.Floor(float64(seconds) / 60 / 60)
    seconds = intInput % (60 * 60)
    minutes := math.Floor(float64(seconds) / 60)
    seconds = intInput % 60

    if years > 0 {
        result = plural(int(years), "year") + plural(int(months), "month") + plural(int(weeks), "week") + plural(int(days), "day")
    } else if months > 0 {
        result = plural(int(months), "month") + plural(int(weeks), "week") + plural(int(days), "day") + plural(int(hours), "hour")
    } else if weeks > 0 {
        result = plural(int(weeks), "week") + plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute")
    } else if days > 0 {
        result = plural(int(days), "day") + plural(int(hours), "hour") + plural(int(minutes), "minute")
    } else if hours > 0 {
        result = plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
    } else if minutes > 0 {
        result = plural(int(minutes), "minute") + plural(int(seconds), "second")
    } else {
        result = plural(int(seconds), "second")
    }

    result = strings.TrimSuffix(result, "\n");

    if negative {
        result = "- " + result;
    }

    return
}
