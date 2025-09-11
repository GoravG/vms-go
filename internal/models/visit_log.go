package models

import "time"

type VisitLog struct {
	ID           int64
	Visitor      string
	VisitDate    time.Time
	CheckInTime  time.Time
	CheckOutTime time.Time
}

/*
Go's standard library does not have a distinct Date type that represents only the year, month, and day without time information.
Instead, dates and times are handled by the time.Time struct within the time package.
While time.Time represents a specific point in time including year, month, day, hour, minute, second, and nanosecond,
you can effectively work with "dates" by focusing on the year, month, and day components of a time.Time value.
*/

func (VisitLog) CreateTableSQL() string {
	return `
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		visitor VARCHAR(50) NOT NULL,
		visit_date DATE NOT NULL,
		check_in_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		check_out_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (visitor) REFERENCES users(email)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
}
