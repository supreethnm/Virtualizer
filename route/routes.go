package route

import (
	"net/http"
	"strings"

	c "virtualizer/configuration"
	"virtualizer/constants"
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
		if strings.Contains(strings.ToUpper(service.Method), constants.HTTP_METHOD_POST) {
			logrus.WithFields(logrus.Fields{}).Info("Request for POST ", service.Path)

			r := Route{service.Sname,
				constants.HTTP_METHOD_POST,
				service.Path,
				h.PostHandler(service)}

			routes = append(routes, r)
		} else if strings.Contains(strings.ToUpper(service.Method), constants.HTTP_METHOD_GET) {
			logrus.WithFields(logrus.Fields{}).Info("Request for GET ", service.Path)

			r := Route{service.Sname,
				constants.HTTP_METHOD_GET,
				service.Path,
				h.GetHandler(service)}

			routes = append(routes, r)
		} else if strings.Contains(strings.ToUpper(service.Method), constants.HTTP_METHOD_DELETE) {
			logrus.WithFields(logrus.Fields{}).Info("Request for DELETE ", service.Path)

			r := Route{service.Sname,
				constants.HTTP_METHOD_DELETE,
				service.Path,
				h.DeleteHandler(service)}

			routes = append(routes, r)
		}
	}

}
