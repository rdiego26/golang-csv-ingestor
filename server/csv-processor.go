package server

import (
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"time"
)

func Read(filename string) [][]string {
	logger := &Logger{}
	file, err := os.Open(filename)
	if err != nil {
		logger.Log(Fatal, fmt.Sprintf("fail to open %s file - %v", filename, err))
		return nil
	}

	defer file.Close()

	var result [][]string
	reader := csv.NewReader(file)
	for {
		data, csvErr := reader.Read()
		if csvErr == io.EOF {
			break
		} else if csvErr != nil {
			logger.Log(Fatal, fmt.Sprintf("Cannot read the content - %v", csvErr))
		}
		result = append(result, data)
	}

	return result
}

func ProcessContacts(data [][]string) []User {
	logger := &Logger{}
	var users []User
	for idx, line := range data {
		if idx == 0 { // omit CSV header
			logger.Log(Info, "Skipping the csv header line")
			continue
		}

		newID := uuid.Must(uuid.NewRandom())
		var user User
		for col, value := range line {
			switch col {
			case 0:
				// we will avoid use integer for ID instead of we will work with UUID
				user.ID = newID
			case 1:
				user.FirstName = value
			case 2:
				user.LastName = value
			case 3:
				user.Email = value
			case 4:
				// Convert millisecond string to integer
				milliseconds, err := parseMilliseconds(value)
				if err != nil {
					logger.Log(Fatal, fmt.Sprintf("Fatal parsing milliseconds from created_at value: %v", err))
				}

				// Convert milliseconds to time.Time
				timeValue := time.Unix(0, milliseconds*int64(time.Millisecond))
				user.CreatedAt = timeValue
			case 5:
				// Convert millisecond string to integer
				milliseconds, err := parseMilliseconds(value)
				if err != nil {
					logger.Log(Fatal, fmt.Sprintf("Fatal parsing milliseconds from deleted_at value: %v", err))
				}

				// Convert milliseconds to time.Time
				timeValue := time.Unix(0, milliseconds*int64(time.Millisecond))
				user.DeletedAt = timeValue
			case 6:
				// Convert millisecond string to integer
				milliseconds, err := parseMilliseconds(value)
				if err != nil {
					logger.Log(Fatal, fmt.Sprintf("Fatal parsing milliseconds from merged_at value: %v", err))
				}

				// Convert milliseconds to time.Time
				timeValue := time.Unix(0, milliseconds*int64(time.Millisecond))
				user.MergedAt = timeValue
			case 7:
				parentId, err := uuid.Parse(value)
				if err != nil {
					logger.Log(Warning, fmt.Sprintf("Fatal parsing parent_id value(%s): %v", value, err))
				}
				user.ParentId = parentId
			}
		}

		users = append(users, user)
	}

	return users
}

// parseMilliseconds converts a millisecond string to an integer.
func parseMilliseconds(millisecondString string) (int64, error) {
	milliseconds, err := time.ParseDuration(millisecondString + "ms")
	if err != nil {
		return 0, err
	}
	return int64(milliseconds), nil
}
