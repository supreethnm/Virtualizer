package route

import (
	"net/http"
	"strings"

	c "virtualizer/configuration"
	h "virtualizer/handler"

	"github.com/Sirupsen/logrus"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes Routes

func InitializeRoutes(Services []c.Service) {
	for _, service := range Services {
		if strings.Contains(strings.ToLower(service.Method), "post") {
			logrus.WithFields(logrus.Fields{}).Info("Request for POST ", service.Path)

			r := Route{service.Sname,
				"POST",
				service.Path,
				h.PostHandler(service)}

			routes = append(routes, r)
		} else if strings.Contains(strings.ToLower(service.Method), "get") {
			logrus.WithFields(logrus.Fields{}).Info("Request for GET ", service.Path)

			r := Route{service.Sname,
				"GET",
				service.Path,
				h.GetHandler(service)}

			routes = append(routes, r)
		} else if strings.Contains(strings.ToLower(service.Method), "delete") {
			logrus.WithFields(logrus.Fields{}).Info("Request for DELETE ", service.Path)

			r := Route{service.Sname,
				"DELETE",
				service.Path,
				h.DeleteHandler(service)}

			routes = append(routes, r)
		}
	}

}
