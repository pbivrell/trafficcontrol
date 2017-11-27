/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	log "github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

func TestCDNsPOST(t *testing.T) {

	for _, cdn := range testData.CDNs {

		b, err := json.Marshal(cdn)
		if err != nil {
			t.Errorf("Could not marshal data %v\n", err)
		}
		fmt.Printf("b ---> %v\n", string(b))
		// TODO: drichardson centralize these routes
		resp, err := Request(*TOSession, http.MethodGet, "/api/1.2/cdns", b)
		fmt.Printf("err ---> %v\n", err)
		fmt.Printf("resp ---> %v\n", resp)
		if err != nil {
			t.Errorf("Could not POST to cdns: %v\n", err)
		}
		defer resp.Body.Close()

		var alerts tc.Alerts
		if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
			t.Errorf("Could not decode alert response.  Error is: %v\n", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Response was not successful %v\n", err)
		}
	}

	for _, cdn := range testData.CDNs {
		DeleteCDNByName(t, cdn.Name)
	}

}

func DeleteCDNByName(t *testing.T, cdnName string) {
	route := fmt.Sprintf("/api/1.2/cdns/name/%s", cdnName)
	fmt.Printf("route ---> %v\n", route)
	resp, err := Request(*TOSession, http.MethodDelete, route, nil)
	fmt.Printf("err ---> %v\n", err)
	fmt.Printf("resp ---> %v\n", resp)
	if err != nil {
		log.Warnf("Could not DELETE cdns: %s - %v\n", cdnName, err)
		var alerts tc.Alerts
		if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
			t.Errorf("Could not decode alert response.  Error is: %v\n", err)
		}
	}
}

//TestCDNGETs compares the results of the CDN api and CDN client
func TestCDNs(t *testing.T) {
	resp, err := Request(*TOSession, http.MethodGet, "/api/1.2/cdns.json", nil)
	if err != nil {
		t.Errorf("Could not get cdns.json reponse was: %v\n", err)
	}

	defer resp.Body.Close()
	var apiCDNRes tc.CDNsResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiCDNRes); err != nil {
		t.Errorf("Could not decode CDN json.  Error is: %v\n", err)
	}
	apiCDNs := apiCDNRes.Response

	//get CDNs data from client
	clientCDNs, err := TOSession.CDNs()
	if err != nil {
		t.Errorf("Could not get CDNs from client.  Error is: %v\n", err)
	}

	if len(apiCDNs) != len(clientCDNs) {
		t.Errorf("Array lengths from client and API are different...API = %s, Client = %s\n", apiCDNs, clientCDNs)
	}

	for _, apiCDN := range apiCDNs {
		matchFound := false
		for _, clientCDN := range clientCDNs {
			if clientCDN.Name != apiCDN.Name {
				continue
			}
			matchFound = true
			//compare API results and client results
			if apiCDN.Name != clientCDN.Name {
				t.Errorf("CDN Name from client and API are different...API = %s, Client = %s\n", apiCDN.Name, clientCDN.Name)
			}
			if apiCDN.LastUpdated != clientCDN.LastUpdated {
				t.Errorf("CDN Last Updated from client and API are different...API = %s, Client = %s\n", apiCDN.Name, clientCDN.Name)
			}
		}
		if !matchFound {
			t.Errorf("A match for %s from the API was not found in the client results\n", apiCDN.Name)
		}
	}
}

//TestCDNName ensures the client returns a CDN by name
func TestCDNName(t *testing.T) {
	//Get CDNs Data from API
	resp, err := Request(*TOSession, http.MethodGet, "/api/1.2/cdns.json", nil)
	if err != nil {
		t.Errorf("Could not get cdns.json reponse was: %v\n", err)
	}

	defer resp.Body.Close()
	var apiCDNRes tc.CDNsResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiCDNRes); err != nil {
		t.Errorf("Could not decode CDN json.  Error is: %v\n", err)
	}

	apiCDNs := apiCDNRes.Response
	apiName := apiCDNs[0].Name
	apiLastUpdated := apiCDNs[0].LastUpdated

	//get CDNs data from client
	clientCDN, err := TOSession.CDNName(apiName)

	if len(clientCDN) != 1 {
		t.Errorf("The length of the client CDN response %v is greater than 1!\n", len(apiCDNs))
	}

	clientName := clientCDN[0].Name
	clientLastUpdated := clientCDN[0].LastUpdated

	//compare API results and client results
	if apiName != clientName {
		t.Errorf("CDN Name from client and API are different...API = %s, Client = %s\n", apiName, clientName)
	}
	if apiLastUpdated != clientLastUpdated {
		t.Errorf("CDN Last Updated from client and API are different...API = %s, Client = %s\n", apiName, clientName)
	}
}
