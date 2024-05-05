package calendar

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	_ "time/tzdata"

	ps "github.com/PuloV/ics-golang"
	ics "github.com/arran4/golang-ical"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcal "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarRequestDto struct {
	Texts []string
	Semester string
	University string
}

type GoogleCalendarService struct {}

func NewGoogleCalendarService() *GoogleCalendarService {
	return &GoogleCalendarService{}
}

func (gcs *GoogleCalendarService) UploadNewCalendar(text string, calendarName string, token *oauth2.Token) (error) {
	cal, err := gcs.createCalendarInstance(text, calendarName)

	var cId string

	if err != nil {
		return err
	}

	srv, err := gcs.CreateGCalService(token)

	if err != nil {
		return err
	}

	if cId, err = gcs.insertNewGcal(cal, srv); err != nil {
		return err
	}

	if err = gcs.insertEventsToGcal(cal.GetEvents(), srv, cId); err != nil {
		return err
	
	}

	return nil
}


func (gcs *GoogleCalendarService) createCalendarInstance(input string, cname string) (*ps.Calendar, error) {
	cal, err := ics.ParseCalendar(strings.NewReader(input))

	if err != nil {
		return nil, err
	}

	cal.SetName(cname)

	ics := gcs.preProcessICS(cal.Serialize())

	parser := ps.New()
	parser.Load(ics)

	if _, err := parser.GetErrors(); err != nil {
		return nil, err
	}
	
	calendars, errCal := parser.GetCalendars()

	if errCal != nil {
		return nil, errCal
	}

	if len(calendars) == 0 {
		return nil, fmt.Errorf("không thể tạo lịch từ dữ liệu đầu vào")
	}

	return calendars[0], nil
}

func (gcs *GoogleCalendarService) preProcessICS(input string) string {
	// Change \r\n to \n
	ics := strings.ReplaceAll(input, "\r\n", "\n")

	return ics
}


func (gcs *GoogleCalendarService) insertNewGcal(calendar *ps.Calendar, srv *gcal.Service) (string, error) {
	gcal := &gcal.Calendar{
		Summary: calendar.GetName(),
		TimeZone: "Asia/Ho_Chi_Minh",
	}

	if c, err := srv.Calendars.Insert(gcal).Do(); err != nil {
		return "", err
	} else {
		return c.Id, nil
	}

}

func (gcs *GoogleCalendarService) insertEventsToGcal(events []ps.Event, srv *gcal.Service, cId string) error {
    wc := sync.WaitGroup{}
	// tz, _ := common.TimeIn("Vietnam")

    for _, event := range events {
        wc.Add(1)

        go func(event ps.Event) {
            defer wc.Done()

			// The string has malfunction newline, so we need to remove them first of all
			desc := strings.ReplaceAll(event.GetDescription(), "\n", "")
			// Then we replace the escaped newline with the real newline
			desc = strings.ReplaceAll(desc, "\\n", "\n")
			loc := strings.ReplaceAll(event.GetLocation(), "\\", "")

            gcalEvent := &gcal.Event{
                Summary:     string(event.GetSummary()),
                Location:    loc,
                Description: desc,
                Start: &gcal.EventDateTime{
                    DateTime: event.GetStart().Format("2006-01-02T15:04:05-07:00"),
                    // TimeZone: "Asia/Bangkok",
                },
                End: &gcal.EventDateTime{
                    DateTime: event.GetEnd().Format("2006-01-02T15:04:05-07:00"),
                    // TimeZone: "Asia/Bangkok",
                },
            }

            if event.GetRRule() != "" {
                rRule := "RRULE:" + gcs.setTimeToDayEnd(event.GetRRule())
                gcalEvent.Recurrence = []string{rRule}
            }

            _, err := srv.Events.Insert(cId, gcalEvent).Do()

            if err != nil {
                // Handle error appropriately, e.g., log it
                fmt.Println("Error inserting event:", err)
                return
            }
        }(event)
    }

    wc.Wait()
    return nil
}

func (gcs *GoogleCalendarService) setTimeToDayEnd(input string) string {
	// FREQ=WEEKLY;UNTIL=20240507T070000Z
	t := strings.Split(input, "T")

	chunkCount := len(t)
	// Change hour part to 23
	t[chunkCount - 1] = "23" + t[chunkCount-1][2:]

	return strings.Join(t, "T")
}


func (gcs *GoogleCalendarService) CreateGCalService(tok *oauth2.Token) (*gcal.Service, error) {
	cf, err := gcs.getConfig()

	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	client := GetClient(cf, tok)
	srv, err := gcal.NewService(ctx, option.WithHTTPClient(client))
	
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func (gcs *GoogleCalendarService) getConfig() (*oauth2.Config, error) {
	creds := []byte(os.Getenv("GOOGLE_CREDENTIALS"))
	
	return google.ConfigFromJSON(creds, "https://www.googleapis.com/auth/calendar")
}

func GetClient(config *oauth2.Config, tok *oauth2.Token) *http.Client {
	return config.Client(context.Background(), tok)
}
