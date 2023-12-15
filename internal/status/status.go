package status

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var once sync.Once
var s *Status

var gGitTag = "unknown"
var gGitCommitID = "unknown"
var gBuildTime = "unknown"

type Status struct {
	lock        sync.RWMutex
	gitTag      string    // 版本号 v3.23.1
	gitCommitID string    // git commit id
	buildTime   string    // 编译时间
	startTime   time.Time // 启动时间
}

func GetStatus() *Status {
	once.Do(func() {
		s = &Status{
			lock:      sync.RWMutex{},
			gitTag:    strings.ToLower(gGitTag),
			buildTime: gBuildTime,
			startTime: time.Now(),
		}
		versionSlice := strings.Split(gGitCommitID, "+")
		if len(versionSlice) > 1 {
			s.gitCommitID = versionSlice[1]
		}
	})
	return s
}
func (s *Status) Init() error {
	return nil
}

func (s *Status) LogVersion() error {
	logrus.Infof("Version=%v", s.gitTag)
	logrus.Infof("CommitID=%v", s.gitCommitID)
	logrus.Infof("BuildTime=%v", s.buildTime)
	logrus.Infof("StartTime=%v", s.startTime.Format("2006-01-02 15:04:05"))
	return nil
}

func (s *Status) ShowVersion() {
	fmt.Printf("Version=%v\n", s.gitTag)
	fmt.Printf("CommitID=%v\n", s.gitCommitID)
	fmt.Printf("BuildTime=%v\n", s.buildTime)
}

func (s *Status) GetStartTime() time.Time {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.startTime
}
