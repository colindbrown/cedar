package evergreen

import (
	"net/http"
	"testing"
	"time"

	"github.com/mongodb/grip"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

func init() {
	grip.SetName("sink.evergreen.projects.test")
}

type ProjectSuite struct {
	client *http.Client
	info   *ConnectionInfo
	suite.Suite
}

func TestProjectsSuite(t *testing.T) {
	suite.Run(t, new(ProjectSuite))
}

func (s *ProjectSuite) SetupSuite() {
	s.info = &ConnectionInfo{
		RootURL: "https://evergreen.mongodb.com/",
		User:    "USER",
		Key:     "KEY",
	}
	s.client = &http.Client{}
}

// Tests getProjects(), which retrieves all projects from Evergreen.
// Authentication is needed for this route.
func (s *ProjectSuite) TestGetProjects() {
	Client := NewClient(s.client, s.info)
	Client.maxRetries = 2
	output := Client.getProjects(context.Background())
	for out := range output {
		s.NoError(out.err)
		s.NotEmpty(out.output.Identifier)
	}
}

// Tests getTaskCostsByProject(), which retrieves all task costs from Evergreen
// for the project given. Authentication is needed for this route.
func (s *ProjectSuite) TestGetTaskCostsByProject() {
	Client := NewClient(s.client, s.info)
	output := Client.getTaskCostsByProject(context.Background(), "mci", "2017-07-25T10:00:00Z", "4h")
	for out := range output {
		s.NoError(out.err)
		s.NotEmpty(out.taskcost.TimeTaken)
	}
}

// Tests GetEvergreenProjectsData(), which retrieves all task costs
// for each project in Evergreen. Authentication is needed for this route.
func (s *ProjectSuite) TestGetEvergreenProjectsData() {
	Client := NewClient(s.client, s.info)
	Client.maxRetries = 2
	starttime, _ := time.Parse(time.RFC3339, "2017-07-25T10:00:00Z")
	duration, _ := time.ParseDuration("1h")
	projectUnits, err := Client.GetEvergreenProjectsData(context.Background(), starttime, duration)
	s.NoError(err)
	for _, pu := range projectUnits {
		s.NotEmpty(pu.Name)
		s.NotEmpty(pu.Tasks)
	}
}
