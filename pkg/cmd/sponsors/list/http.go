package list

import (
	"net/http"

	"github.com/cli/cli/v2/api"
)

type sponsorsResponse struct {
	User struct {
		Sponsors struct {
			Edges []struct {
				Node struct {
					Login string `json:"login,omitempty"`
				}
			}
		}
	}
}

type SponsorsList struct {
	Sponsors []string `json:"login"`
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
            }
        }
    }`

	apiClient := api.NewClientFromHTTP(client)

	var data sponsorsResponse

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

	listResult := SponsorsList{
		Sponsors: []string{},
	}

	for _, sponsor := range data.User.Sponsors.Edges {
		listResult.Sponsors = append(listResult.Sponsors, sponsor.Node.Login)
	}

	return listResult, nil
}
