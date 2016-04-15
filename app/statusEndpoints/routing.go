/*
 * Copyright (c) 2015 GRNET S.A.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the
 * License. You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an "AS
 * IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language
 * governing permissions and limitations under the License.
 *
 * The views and conclusions contained in the software and
 * documentation are those of the authors and should not be
 * interpreted as representing official policies, either expressed
 * or implied, of GRNET S.A.
 *
 */

package statusEndpoints

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ARGOeu/argo-web-api/app/reports"
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// HandleSubrouter contains the different paths to follow during subrouting
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	s = respond.PrepAppRoutes(s, confhandler, appRoutesV2)
}

var appRoutesV2 = []respond.AppRoutes{
	{"status_endpoints.get", "GET", "/{report_name}/{group_type}/{group_name}/services/{service_name}/endpoints/{endpoint_name}", routeCheckGroup},
	{"status_endpoints.list", "GET", "/{report_name}/{group_type}/{group_name}/services/{service_name}/endpoints", routeCheckGroup},
	{"status_endpoints.options", "OPTIONS", "/{report_name}/{group_type}/{group_name}/services/{service_name}/endpoints/{endpoint_name}", Options},
	{"status_endpoints.options", "OPTIONS", "/{report_name}/{group_type}/{group_name}/services/{service_name}/endpoints", Options},
}

func routeCheckGroup(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("group check")
	err := error(nil)
	contentType := "application/xml"
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// Handle response format based on Accept Header
	// Default is application/xml
	format := r.Header.Get("Accept")
	if strings.EqualFold(format, "application/json") {
		contentType = "application/json"
	}

	vars := mux.Vars(r)

	// Grab Tenant DB configuration from context
	tenantcfg := context.Get(r, "tenant_conf").(config.MongoConfig)

	session, err := mongo.OpenSession(tenantcfg)
	defer mongo.CloseSession(session)
	if err != nil {
		code = http.StatusInternalServerError
		output, _ = respond.MarshalContent(respond.InternalServerErrorMessage, contentType, "", " ")
		return code, h, output, err
	}
	result := reports.MongoInterface{}
	err = mongo.FindOne(session, tenantcfg.Db, "reports", bson.M{"info.name": vars["report_name"]}, &result)

	if err != nil {
		message := "The report with the name " + vars["report_name"] + " does not exist"
		output, err := createMessageOUT(message, format) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	if vars["group_type"] != result.GetEndpointGroupType() {
		message := "The report " + vars["report_name"] + " does not define endpoint group type: " + vars["group_type"]
		output, err := createMessageOUT(message, format) //Render the response into XML or JSON
		h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
		return code, h, output, err
	}

	return ListEndpointTimelines(r, cfg)

}
