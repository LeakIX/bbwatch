package program_sources

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func init() {
	ProgramSources["h1"] = &HackerOneSource{}
}

type HackerOneSource struct {
	ApiKey   string
	Username string
}

func (ls *HackerOneSource) GetName() string {
	return "HackerOne"
}

func (ls *HackerOneSource) Configure(config map[string]interface{}) {
	for configKey, configValue := range config {
		switch configKey {
		case "api_key":
			ls.ApiKey = configValue.(string)
		case "username":
			ls.Username = configValue.(string)
		default:
			continue
		}
	}
	return
}

const H1_ROOT_API string = "https://api.hackerone.com/v1"

func (ls *HackerOneSource) ApiQuery(method string, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, H1_ROOT_API+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(ls.Username, ls.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "BBWatcher 0.1")
	return http.DefaultClient.Do(req)
}

func (ls *HackerOneSource) GetPrograms(bountyOnly bool) chan Program {
	programChan := make(chan Program)
	go func() {
		var endpoint = "/hackers/programs?page[size]=100"
		for {
			resp, err := ls.ApiQuery(http.MethodGet, endpoint, nil)
			if err != nil {
				log.Println(err)
				return
			}
			var programsResp H1ProgramsResp
			err = json.NewDecoder(resp.Body).Decode(&programsResp)
			if err != nil {
				log.Println(err)
				return
			}

			for _, program := range programsResp.Data {
				if bountyOnly && !program.Attributes.OffersBounties {
					continue
				}
				bbProgram := Program{
					Name:     program.Attributes.Name,
					Platform: ls.GetName(),
					Reward:   program.Attributes.OffersBounties,
				}
				log.Printf("Loading scopes for program %s", program.Attributes.Handle)
				programScope, err := ls.GetProgram(program)
				if err != nil {
					log.Println(err)
					return
				}
				for _, scope := range programScope.Relationships.StructuredScopes.Data {
					if scope.Attributes.AssetType == "URL" {
						bbProgram.Assets = append(bbProgram.Assets, Asset{
							Domain:   scope.Attributes.AssetIdentifier,
							Wildcard: false,
						})
					}
				}
				programChan <- bbProgram
			}
			if len(programsResp.Links.Next) < 1 {
				break
			}
			endpoint = programsResp.Links.GetNextEndpoint()
		}
		close(programChan)
	}()
	return programChan
}

func (ls *HackerOneSource) GetProgram(program H1Program) (*H1Program, error) {
	resp, err := ls.ApiQuery(http.MethodGet, "/hackers/programs/"+program.Attributes.Handle, nil)
	if err != nil {
		return nil, err
	}
	var programResp H1Program
	err = json.NewDecoder(resp.Body).Decode(&programResp)
	if err != nil {
		return nil, err
	}
	return &programResp, nil
}

type H1ProgramsResp struct {
	Data  []H1Program `json:"data"`
	Links H1Links
}

type H1Links struct {
	Self string `json:"self"`
	Next string `json:"next"`
}

func (h1l H1Links) GetNextEndpoint() string {
	return strings.Replace(h1l.Next, H1_ROOT_API, "", 1)
}

type H1Program struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Handle                          string     `json:"handle"`
		Name                            string     `json:"name"`
		Currency                        string     `json:"currency"`
		SubmissionState                 string     `json:"submission_state"`
		TriageActive                    bool       `json:"triage_active"`
		State                           string     `json:"state"`
		StartedAcceptingAt              *time.Time `json:"started_accepting_at,omitempty"`
		NumberOfReportsForUser          int64      `json:"number_of_reports_for_user"`
		NumberOfValidReportsForUser     int64      `json:"number_of_valid_reports_for_user"`
		BountyEarnedForUser             float64    `json:"bounty_earned_for_user"`
		LastInvitationAcceptedAtForUser *time.Time `json:"last_invitation_accepted_at_for_user,omitempty"`
		Bookmarked                      bool       `json:"bookmarked"`
		AllowsBountySplitting           bool       `json:"allows_bounty_splitting"`
		OffersBounties                  bool       `json:"offers_bounties"`
	} `json:"attributes"`
	Relationships struct {
		StructuredScopes struct {
			Data []H1Scope `json:"data"`
		} `json:"structured_scopes"`
	} `json:"relationships"`
}

type H1Scope struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		AssetType             string    `json:"asset_type"`
		AssetIdentifier       string    `json:"asset_identifier"`
		EligibleForBounty     bool      `json:"eligible_for_bounty"`
		EligibleForSubmission bool      `json:"eligible_for_submission"`
		Instruction           string    `json:"instruction"`
		MaxSeverity           string    `json:"max_severity"`
		CreatedAt             time.Time `json:"created_at"`
		UpdatedAt             time.Time `json:"updated_at"`
	}
}
