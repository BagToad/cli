package list

import (
	"net/http"

	"github.com/cli/cli/v2/api"
)

type SponsorsList struct {
	User struct {
		Sponsors struct {
			Edges []struct {
				Node struct {
					Login string `json:"login,omitempty"`
				}
			}
			TotalCount int `json:"totalCount"`
		}
	}
}

func listSponsors(client *http.Client, opts *ListOptions) (SponsorsList, error) {
	query := `query SponsorsList($login: String!) {
        user(login: $login) {
            sponsors(first: 30) {
                edges {
                    node {
                        ... on User {			
                            login
                        }
                        ... on Organization {
                            login
                        }
                    }
                }
                totalCount
            }
        }
    }`

	apiClient := api.NewClientFromHTTP(client)

	var data SponsorsList

	cfg, err := opts.Config()
	if err != nil {
		return SponsorsList{}, err
	}
	hostname, _ := cfg.Authentication().DefaultHost()
	variables := map[string]interface{}{
		"login": opts.User,
	}

	err = apiClient.GraphQL(hostname, query, variables, &data)
	if err != nil {
		return SponsorsList{}, err
	}

	return data, nil
}
