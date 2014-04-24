/*
 * Copyright (c) 2014 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
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
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package main

import (
	"api/apiCalls/availabilityProfiles"
	"api/apiCalls/ngis"
	//"api/apiCalls/profileCRUD"
	"api/apiCalls/recalculations"
	"api/apiCalls/services"
	"api/apiCalls/serviceFlavors"
	"api/apiCalls/sites"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {

	//Create the server router
	main_router := mux.NewRouter()
	get_subrouter := main_router.Methods("GET").Subrouter()                                //Routes only GET requests
	post_subrouter := main_router.Methods("POST").Headers("x-api-key", "").Subrouter()     //Routes only POST requests
	delete_subrouter := main_router.Methods("DELETE").Headers("x-api-key", "").Subrouter() //Routes only DELETE requests
	put_subrouter := main_router.Methods("PUT").Headers("x-api-key", "").Subrouter()       //Routes only PUT requests
	//All requests that modify data must provide with authentication credentials
		
	// Grouping calls. 
	// Groups are routed depending on the value of the parameter group type. FUTURE WORK: 
	// 1) Move calls to a separate subrouter. 	
	// 2) Provide with a default call informing the user of an invalid parameter
	
	get_subrouter.HandleFunc("/api/v1/group_availability_in_profile", Respond("text/xml", "utf-8", vos.VoAvailabilityInProfile)).
	Queries("group_type", "vo")
	get_subrouter.HandleFunc("/api/v1/group_availability_in_profile", Respond("text/xml", "utf-8", sites.SitesAvailabilityInProfile)).
	Queries("group_type", "site")
	
	//Basic api calls
	get_subrouter.HandleFunc("/api/v1/service_availability_in_profile", Respond("text/xml", "utf-8", services.ServiceAvailabilityInProfile))
	get_subrouter.HandleFunc("/api/v1/ngi_availability_in_profile", Respond("text/xml", "utf-8", ngis.NgiAvailabilityInProfile))
	get_subrouter.HandleFunc("/api/v1/service_flavor_availability", Respond("text/xml", "utf-8", serviceFlavors.ServiceFlavorAvailabilityInProfile)).
	Queries("group_type", "sf")
	//POEM PROFILE MANAGEMENT TO BE REMOVED/MODIFIED!!!!
	// post_subrouter.HandleFunc("/api/v1/profiles/create", Respond("text/xml", "utf-8", profileCRUD.CreateProfile))
	// get_subrouter.HandleFunc("/api/v1/profiles", Respond("text/xml", "utf-8", profileCRUD.ReadAllProfiles))
	// get_subrouter.HandleFunc("/api/v1/profiles/getone", Respond("text/xml", "utf-8", profileCRUD.ReadOneProfile))
	//SOME UPDATE METHOD MISSING
	//post_subrouter.HandleFunc("/api/v1/profiles/remove", Respond("text/xml", "utf-8", profileCRUD.DeleteProfile))

	post_subrouter.HandleFunc("/api/v1/AP", Respond("text/xml", "utf-8", availabilityProfiles.CreateProfiles))
	get_subrouter.HandleFunc("/api/v1/AP", Respond("text/xml", "utf-8", availabilityProfiles.ReadProfiles))
	put_subrouter.HandleFunc("/api/v1/AP/{id}", Respond("text/xml", "utf-8", availabilityProfiles.UpdateProfiles))
	delete_subrouter.HandleFunc("/api/v1/AP/{id}", Respond("text/xml", "utf-8", availabilityProfiles.DeleteProfiles))

	//Recalculations
	post_subrouter.HandleFunc("/api/v1/recalculate", Respond("text/xml", "utf-8", recalculations.Recalculate))
	get_subrouter.HandleFunc("/api/v1/get_recalculation_requests", Respond("text/xml", "utf-8", recalculations.GetRecalculationRequests))
	http.Handle("/", main_router)

	//Cache
	//get_subrouter.HandleFunc("/api/v1/reset_cache", Respond("text/xml", "utf-8", ResetCache))

	//Web service binds to server. Requests served over HTTPS.
	err := http.ListenAndServeTLS(cfg.Server.Bindip+":"+strconv.Itoa(cfg.Server.Port), "/etc/pki/tls/certs/localhost.crt", "/etc/pki/tls/private/localhost.key", nil)

	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
