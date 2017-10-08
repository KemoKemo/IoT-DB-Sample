package main

import (
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DateQueryParam is used to query from db.
type DateQueryParam struct {
	Start time.Time
	End   time.Time
}

// YYYYMMDD-hhmmss
const (
	layout   = "20060102-150405"
	locale   = "Local"
	duration = 3
)

// rParseQuery returns a DateQueryParam like below:
//  - start and end are empty: recent 3 days
//  - only start is specified: from start to (start + 3days)
//  - only end is specified: from (end - 3 days) to end
//  - start and end are specified: from start to end
func rParseQuery(r *http.Request) (DateQueryParam, error) {
	var data DateQueryParam
	loc, err := time.LoadLocation(locale)
	if err != nil {
		return data, err
	}

	startVal := r.FormValue("start")
	endVal := r.FormValue("end")
	var start, end time.Time

	if startVal == "" && endVal == "" {
		end = time.Now()
		start = end.AddDate(0, 0, -duration)
	} else if startVal == "" && endVal != "" {
		end, err = getTime(endVal, loc)
		if err != nil {
			return data, err
		}
		start = end.AddDate(0, 0, -duration)
	} else if startVal != "" && endVal == "" {
		start, err = getTime(startVal, loc)
		if err != nil {
			return data, err
		}
		end = start.AddDate(0, 0, duration)
	} else {
		start, err = getTime(startVal, loc)
		if err != nil {
			return data, err
		}
		end, err = getTime(endVal, loc)
		if err != nil {
			return data, err
		}
	}

	return DateQueryParam{
		Start: start,
		End:   end,
	}, nil
}

func getTime(s string, loc *time.Location) (time.Time, error) {
	var t time.Time
	if s == "" {
		return time.Time{}, nil
	}

	t, err := time.ParseInLocation(layout, s, loc)
	if err != nil {
		return t, err
	}
	return t, nil
}

func getChartData(q DateQueryParam, db *mgo.Database) (*ChartData, error) {
	var data ChartData
	err := db.C(columnName).
		Find(
			bson.M{
				"$and": []bson.M{
					bson.M{"date": bson.M{"$gte": q.Start}},
					bson.M{"date": bson.M{"$lte": q.End}},
				}}).Sort("-$natural").All(&data.DataSetList)
	if err != nil {
		return &data, err
	}
	return &data, nil
}
