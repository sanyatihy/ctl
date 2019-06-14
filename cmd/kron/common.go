package kron

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"os"
	"time"
)

// For storing the location of a job for select and favorite.
type location struct {
	Contexts  []string `json:"contexts"`
	Namespace string   `json:"namespace"`
}

type selectedJob struct {
	Name     string   `json:"name"`
	Location location `json:"location"`
}

func toLocation(obj interface{}) location {
	m, ok := obj.(map[string]interface{})
	if !ok {
		fmt.Println("Failed")
		return location{} // maybe panic??
	}
	c := m["contexts"].([]string)
	n := m["namespace"].(string)
	return location{c, n}
}

func createConfig() {
	os.Mkdir(os.Getenv("HOME")+"/.kron/", 0777)
	err := viper.WriteConfigAs(os.Getenv("HOME") + "/.kron/config.yaml")
	if err != nil {
		panic(err.Error())
	}
}

func getSelected() (s selectedJob, err error) {
	err = viper.UnmarshalKey("selected", &s)
	fmt.Println(s)
	return
}

func getFavorites() (f map[string]location, err error) {
	err = viper.UnmarshalKey("favorites", &f)
	return
}

func matchesCronJobLocation(c client.CronJobDiscovery, l location) bool {
	if l.Namespace != "" && l.Namespace != c.Namespace {
		return false
	}
	if len(l.Contexts) == 0 {
		return true
	}
	for _, ctx := range l.Contexts {
		if c.Context == ctx {
			return true
		}
	}
	return false
}

func filterFromFavorites(lst []client.CronJobDiscovery) []client.CronJobDiscovery {
	f, err := getFavorites()
	if err != nil {
		panic(err.Error())
	}
	var filtered []client.CronJobDiscovery
	for _, c := range lst {
		if l, ok := f[c.Name]; ok && matchesCronJobLocation(c, l) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

type byLastRun []client.CronJobDiscovery

func (l byLastRun) Len() int      { return len(l) }
func (l byLastRun) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byLastRun) Less(i, j int) bool {
	return l[i].Status.LastScheduleTime.Time.After(l[j].Status.LastScheduleTime.Time)
}

type byNextRun []client.CronJobDiscovery

func (l byNextRun) Len() int      { return len(l) }
func (l byNextRun) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byNextRun) Less(i, j int) bool {
	a, _ := cron.ParseStandard(l[i].Spec.Schedule)
	b, _ := cron.ParseStandard(l[j].Spec.Schedule)
	now := time.Now()
	return a.Next(now).Before(b.Next(now))
}